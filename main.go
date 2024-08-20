package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/projectdiscovery/gologger"
	"github.com/projectdiscovery/retryablehttp-go"
)

type Config struct {
	ApiKey string            `json:"apiKey"`
	Apis   map[string]string `json:"apis"`
}

type ResponseData struct {
	Domains []string `json:"domains"`
}

func loadConfig() (*Config, error) {
	file, err := os.Open("config.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func reverseIP(apiURL, ip string, config *Config, wg *sync.WaitGroup, mu *sync.Mutex, sem chan struct{}, domainSet map[string]struct{}, outputFile string) {
	defer wg.Done()
	defer func() { <-sem }() // Release the semaphore when the function exits

	url := fmt.Sprintf(apiURL, config.ApiKey, ip)
	client := retryablehttp.NewClient(retryablehttp.DefaultOptionsSingle)
	client.HTTPClient.Timeout = 10 * time.Second

	req, err := retryablehttp.NewRequest("GET", url, nil)
	if err != nil {
		gologger.Error().Msgf("%s: Failed to create request: %v", ip, err)
		return
	}

	req.Header.Set("User-Agent", randomUserAgent())

	resp, err := client.Do(req)
	if err != nil {
		gologger.Error().Msgf("%s: Failed to connect to the API: %v", ip, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		gologger.Error().Msgf("%s: No Data Found", ip)
		return
	}

	var data ResponseData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		gologger.Error().Msgf("%s: Failed to decode response: %v", ip, err)
		return
	}

	mu.Lock()
	defer mu.Unlock()
	file, err := os.OpenFile(outputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		gologger.Error().Msgf("Failed to open file: %v", err)
		return
	}
	defer file.Close()

	for _, domain := range data.Domains {
		if _, exists := domainSet[domain]; !exists {
			domainSet[domain] = struct{}{}
			file.WriteString(domain + "\n")
		}
	}

	gologger.Info().Msgf("%s: %d unique domains found", ip, len(data.Domains))
}

func randomUserAgent() string {
	userAgents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/64.0.3282.140 Safari/537.36 Edge/17.17134",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3359.181 Safari/537.36",
		"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:52.0) Gecko/20100101 Firefox/52.0",
	}
	return userAgents[rand.Intn(len(userAgents))]
}

func scanIPs(ips []string, config *Config, threads int, outputFile string) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	sem := make(chan struct{}, threads) // Semaphore to limit the number of concurrent goroutines
	domainSet := make(map[string]struct{})

	for _, ip := range ips {
		for _, apiURL := range config.Apis {
			wg.Add(1)
			sem <- struct{}{} // Acquire a semaphore slot
			go reverseIP(apiURL, ip, config, &wg, &mu, sem, domainSet, outputFile)
		}
	}

	wg.Wait() // Wait for all goroutines to finish
}

func showBanner() {
	banner := `
██████╗  ██████╗ ███╗   ███╗ █████╗ ██╗███╗   ██╗██████╗ ███████╗ ██████╗ ██████╗ ███╗   ██╗    
██╔══██╗██╔═══██╗████╗ ████║██╔══██╗██║████╗  ██║██╔══██╗██╔════╝██╔════╝██╔═══██╗████╗  ██║    
██║  ██║██║   ██║██╔████╔██║███████║██║██╔██╗ ██║██████╔╝█████╗  ██║     ██║   ██║██╔██╗ ██║    
██║  ██║██║   ██║██║╚██╔╝██║██╔══██║██║██║╚██╗██║██╔══██╗██╔══╝  ██║     ██║   ██║██║╚██╗██║    
██████╔╝╚██████╔╝██║ ╚═╝ ██║██║  ██║██║██║ ╚████║██║  ██║███████╗╚██████╗╚██████╔╝██║ ╚████║    
╚═════╝  ╚═════╝ ╚═╝     ╚═╝╚═╝  ╚═╝╚═╝╚═╝  ╚═══╝╚═╝  ╚═╝╚══════╝ ╚═════╝ ╚═════╝ ╚═╝  ╚═══╝    
                                                                                                   
            Reverse IP Lookup Tool         
          Powered by xReverseLabs API         
	`
	color.Cyan(banner)
	color.Cyan("==================================================")
	color.Cyan("Site : https://xreverselabs.my.id - By t.me/xxyz4")
	color.Cyan("Private & Maintained Datacenter")
	color.Cyan("==================================================")
}

func main() {
	// Set up command-line flags
	fileName := flag.String("f", "", "File containing list of IPs to scan")
	singleIP := flag.String("d", "", "Single IP address to scan")
	threads := flag.Int("t", 10, "Number of threads to use")
	outputFile := flag.String("o", "reversed.txt", "Output file to write the results")

	flag.Parse()

	if *fileName == "" && *singleIP == "" {
		flag.Usage()
		return
	}

	config, err := loadConfig()
	if err != nil {
		gologger.Fatal().Msgf("Failed to load config file: %v", err)
	}

	showBanner()

	var ips []string

	if *singleIP != "" {
		ips = append(ips, *singleIP)
	} else {
		file, err := os.Open(*fileName)
		if err != nil {
			gologger.Fatal().Msgf("Failed to open file: %v", err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			ips = append(ips, scanner.Text())
		}

		if err := scanner.Err(); err != nil {
			gologger.Fatal().Msgf("Error reading file: %v", err)
		}
	}

	scanIPs(ips, config, *threads, *outputFile)
}
