Here's a `README.md` file for your `dnsify` package that explains its purpose, usage, and provides code examples:

---

# DNSify



`dnsify` is a lightweight DNS resolver client built in Go that allows you to resolve hostnames and retrieve DNS records (e.g., A, AAAA, NS, CNAME, etc.) efficiently. It supports retries, random resolver selection, and retrieving raw DNS responses.

## Features

- Random resolver selection for load balancing
- Configurable retry mechanism
- Resolve specific DNS record types (A, AAAA, NS, etc.)
- Fetch raw DNS responses for debugging or inspection
- Concurrency-safe with optimized read and write locks

## Installation

To install the `dnsify` package, run:

```bash
go get github.com/cyinnove/dnsify
```

Then, import it into your project:

```go
import "github.com/cyinnove/dnsify"
```

## Usage

### Basic Example: Resolving A Records

Here’s how you can use `dnsify` to resolve A records (IPv4 addresses) for a given hostname:

```go
package main

import (
	"fmt"
	"log"

	"github.com/cyinnove/dnsify"
)

func main() {
	resolvers := []string{"8.8.8.8:53", "8.8.4.4:53"}
	client := dnsify.New(resolvers, 3)  // 3 retries

	result, err := client.Resolve("example.com")
	if err != nil {
		log.Fatalf("Failed to resolve: %v", err)
	}

	fmt.Printf("Resolved IPs: %v, TTL: %d\n", result.IPs, result.TTL)
}
```

### Fetching Raw DNS Records

If you want to retrieve raw DNS records for debugging or advanced use cases, you can use the `ResolveRaw` method:

```go
package main

import (
	"fmt"
	"log"

	"github.com/cyinnove/dnsify"
	"github.com/miekg/dns"
)

func main() {
	resolvers := []string{"8.8.8.8:53", "1.1.1.1:53"}
	client := dnsify.New(resolvers, 3)  // 3 retries

	results, raw, err := client.ResolveRaw("example.com", dns.TypeMX)  // Resolve MX records
	if err != nil {
		log.Fatalf("Failed to resolve: %v", err)
	}

	fmt.Printf("Raw response: \n%s\n", raw)
	fmt.Printf("Parsed MX records: %v\n", results)
}
```

### Sending Custom DNS Requests

For advanced DNS queries, you can create a custom `dns.Msg` and send it using the `Do` method:

```go
package main

import (
	"fmt"
	"log"

	"github.com/cyinnove/dnsify"
	"github.com/miekg/dns"
)

func main() {
	resolvers := []string{"8.8.8.8:53", "1.1.1.1:53"}
	client := dnsify.New(resolvers, 3)  // 3 retries

	msg := new(dns.Msg)
	msg.SetQuestion(dns.Fqdn("example.com"), dns.TypeA)

	resp, err := client.Do(msg)
	if err != nil {
		log.Fatalf("Failed to send custom DNS request: %v", err)
	}

	fmt.Printf("Custom DNS Response: \n%s\n", resp)
}
```

## API Reference

### `New`

```go
func New(baseResolvers []string, maxRetries int) *Client
```

- **baseResolvers**: A list of DNS resolver addresses (e.g., `8.8.8.8:53`).
- **maxRetries**: The number of retry attempts when resolving fails.

### `Resolve`

```go
func (c *Client) Resolve(host string) (Result, error)
```

- **host**: The hostname to resolve (e.g., `example.com`).
- **Result**: A struct containing the IPs and TTL.

### `ResolveRaw`

```go
func (c *Client) ResolveRaw(host string, requestType uint16) ([]string, string, error)
```

- **host**: The hostname to resolve.
- **requestType**: The DNS record type (e.g., `dns.TypeA`, `dns.TypeMX`, etc.).
- **Returns**: A list of parsed records, the raw response string, and any error encountered.

### `Do`

```go
func (c *Client) Do(msg *dns.Msg) (*dns.Msg, error)
```

- **msg**: A custom `dns.Msg` message to send.
- **Returns**: The raw DNS response and any error encountered.

### `Result`

```go
type Result struct {
	IPs []string
	TTL int
}
```

The `Result` struct contains:
- **IPs**: The list of resolved IP addresses.
- **TTL**: The time-to-live (TTL) of the DNS response.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for more details.

---

### Contributions

Feel free to submit issues, fork the repository, and open pull requests to help improve `dnsify`.

