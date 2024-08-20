
# DomainRecon
![preview](https://raw.githubusercontent.com/xReverseLabs/DomainRecon/main/domainrecon.gif)

## Overview

**DomainRecon** is a powerful and efficient reverse IP lookup tool designed to help you discover domains associated with an IP address. It leverages multiple APIs to gather and map domain information efficiently, ensuring comprehensive results. 

## Features

- **Multi-API Support**: Utilizes multiple APIs (e.g., `itsuka` and `himawari`) for thorough domain discovery.
- **Concurrent Scanning**: Capable of handling multiple IPs simultaneously with configurable threading.
- **Single and Batch Scanning**: Supports scanning a single IP or multiple IPs from a file.
- **Configurable via JSON**: API keys and URLs are managed through a simple `config.json` file.
- **Duplicate Handling**: Ensures that duplicate domains from different APIs are not included in the results.
- **Customizable User-Agent**: Randomized User-Agent strings for each request to mimic different browsers.

## Installation

1. **Clone the Repository**:
    ```bash
    git clone https://github.com/YourUsername/DomainRecon.git
    cd DomainRecon
    ```

2. **Configure the Tool**:
    - Create and edit the `config.json` file:
    ```json
    {
      "apiKey": "YOUR_API_KEY_HERE",
      "apis": {
        "itsuka": "https://api.xreverselabs.my.id/itsuka?apiKey=%s&ip=%s",
        "himawari": "https://api.xreverselabs.my.id/himawari?apiKey=%s&ip=%s"
      }
    }
    ```
    - Replace `YOUR_API_KEY_HERE` with your actual API key, [REGISTER HERE TO GET YOUR API KEY](https://xreverselabs.my.id/clientarea/register).

3. **Build the Tool** (optional):
    - You can build the tool to create an executable:
    ```bash
    go build main.go
    ```

## Usage

### Single IP Scan

To scan a single IP address:
```bash
go run main.go -d 1.1.1.1 -t 10 -o output.txt
```
**Windows** : 
```bash
SubRecon-x64.exe -d 1.1.1.1 -t 10 -o output.txt
```

### Batch IP Scan

To scan multiple IPs from a file:
```bash
go run main.go -f list.txt -t 10 -o output.txt
```
**Windows** : 
```bash
SubRecon-x64.exe -f list.txt -t 10 -o output.txt
```

### Configuration Options

- `-f`: Path to the file containing a list of IPs.
- `-d`: Single IP address to scan.
- `-t`: Number of concurrent threads to use.
- `-o`: Output file where the results will be saved (default: `reversed.txt`).

## Example

```bash
go run main.go -f ip_list.txt -t 20 -o domains_found.txt
```

This command scans all IPs listed in `ip_list.txt` using 20 threads and saves the results in `domains_found.txt`.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Contact

For any inquiries, reach out to the project maintainer at [Your Email](mailto:l1nux3r69@gmail.com).
