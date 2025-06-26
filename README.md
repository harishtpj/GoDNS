# GoDNS

A lightweight DNS server implementation written in Go from scratch. GoDNS provides a simple, fast, and customizable DNS server that can handle A record queries with configurable domain-to-IP mappings.

## Features

- **Pure Go Implementation**: Built from scratch without external DNS libraries
- **UDP Protocol Support**: Handles DNS queries over UDP on port 53
- **A Record Resolution**: Supports IPv4 address resolution for domain names
- **Configurable Records**: Load DNS records from JSON files or use default configuration
- **NXDOMAIN Response**: Properly handles queries for non-existent domains
- **Concurrent Request Handling**: Uses goroutines for handling multiple simultaneous queries
- **Lightweight**: Minimal resource usage with fast response times

## Installation

### Prerequisites
- Go 1.24.4 or later

### Build from Source
```bash
git clone https://github.com/harishtpj/GoDNS.git
cd GoDNS
go build -o godns
```

## Usage

### Basic Usage (Default Configuration)
Run the DNS server with a default record:
```bash
sudo ./godns
```
This starts the server with `example.com` pointing to `127.0.0.1`.

### Custom Configuration
Create a JSON file with your DNS records:
```json
{
    "example.com": "192.168.1.100",
    "test.local": "10.0.0.1",
    "myapp.dev": "172.16.0.10"
}
```

Run the server with your custom records:
```bash
sudo ./godns records.json
```

### Testing the DNS Server
Once running, you can test the DNS server using tools like `dig` or `nslookup`:

```bash
# Using dig
dig @localhost example.com

# Using nslookup
nslookup example.com localhost
```

## Configuration

### DNS Records Format
DNS records should be provided in JSON format with the following structure:
```json
{
    "domain.name": "ip.address",
    "another.domain": "another.ip"
}
```

### Default Configuration
When run without arguments, GoDNS uses the following default configuration:
- `example.com` → `127.0.0.1`

## Architecture

### Components

- **DNS Header Parsing**: Handles standard DNS message headers with ID, flags, and section counts
- **Question Parsing**: Extracts domain names, query types, and classes from DNS questions
- **Response Building**: Constructs proper DNS responses with A records or NXDOMAIN responses
- **UDP Server**: Manages incoming DNS queries and sends responses back to clients

### Protocol Support
- **Query Types**: Currently supports A record queries (Type 1)
- **Response Codes**: 
  - `NOERROR` (0x8180) for successful queries
  - `NXDOMAIN` (0x8183) for non-existent domains
- **TTL**: Default TTL of 60 seconds for all records

## Development

### Project Structure
```
GoDNS/
├── main.go          # Main server logic and UDP handling
├── dnsutils.go      # DNS protocol parsing and response building
├── go.mod           # Go module definition
├── LICENSE          # MIT License
└── .gitignore       # Git ignore rules
```

### Building
```bash
go build -o godns
```

### Running Tests
```bash
go test ./...
```

## Performance

GoDNS is designed for efficiency:
- Concurrent request handling using goroutines
- Minimal memory allocation during query processing
- Fast DNS message parsing and response generation
- UDP-based communication for low latency

## Limitations

- Currently supports only A record queries (IPv4)
- No support for AAAA records (IPv6)
- No recursive resolution (authoritative only)
- No caching mechanism
- No support for other DNS record types (CNAME, MX, etc.)

## Security Considerations

- **Root Privileges**: Requires root/administrator privileges to bind to port 53
- **Input Validation**: Basic validation is performed on incoming DNS queries
- **No Authentication**: No built-in authentication mechanism

## Contributing

Contributions are welcome! Please feel free to submit issues, feature requests, or pull requests.

### Development Setup
1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Author

**Harish Kumar** - [harishtpj](https://github.com/harishtpj)

## Acknowledgments

- Built with Go's standard library networking capabilities
- Implements RFC 1035 DNS message format specifications

---

⚠️ **Note**: This DNS server is intended for development, testing, and educational purposes. For production use, consider additional security hardening and feature implementation.
