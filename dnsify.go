package dnsify

import (
	"errors"
	"math/rand"
	"sync"
	"time"

	"github.com/miekg/dns"
)

const defaultPort = "53"

// Client is a DNS resolver client to resolve hostnames.
type Client struct {
	resolvers  []string
	maxRetries int
	rand       *rand.Rand
	mutex      sync.RWMutex
}

// Result contains the results from a DNS resolution.
type Result struct {
	IPs []string
	TTL int
}

// New creates a new DNS client.
func New(baseResolvers []string, maxRetries int) *Client {
	return &Client{
		rand:       rand.New(rand.NewSource(time.Now().UnixNano())),
		maxRetries: maxRetries,
		resolvers:  baseResolvers,
	}
}

// Resolve resolves a hostname and retrieves its A record IPs and TTL.
func (c *Client) Resolve(host string) (Result, error) {
	msg := buildDNSMessage(host, dns.TypeA)

	var result Result
	for i := 0; i < c.maxRetries; i++ {
		resolver := c.getRandomResolver()

		answer, err := dns.Exchange(msg, resolver)
		if err != nil || answer == nil || answer.Rcode != dns.RcodeSuccess {
			continue
		}

		return parseARecords(answer), nil
	}

	return result, errors.New("failed to resolve after max retries")
}

// ResolveRaw resolves a hostname and retrieves raw DNS records of a specific type.
func (c *Client) ResolveRaw(host string, requestType uint16) ([]string, string, error) {
	msg := buildDNSMessage(host, requestType)

	for i := 0; i < c.maxRetries; i++ {
		resolver := c.getRandomResolver()

		answer, err := dns.Exchange(msg, resolver)
		if err != nil || answer == nil || answer.Rcode != dns.RcodeSuccess {
			continue
		}

		raw := answer.String()
		return parseRecordsByType(answer, requestType), raw, nil
	}

	return nil, "", errors.New("failed to resolve after max retries")
}

// Do sends a DNS request and returns the raw DNS response.
func (c *Client) Do(msg *dns.Msg) (*dns.Msg, error) {
	for i := 0; i < c.maxRetries; i++ {
		resolver := c.getRandomResolver()

		resp, err := dns.Exchange(msg, resolver)
		if err == nil && resp != nil {
			return resp, nil
		}
	}

	return nil, errors.New("failed to send DNS request after max retries")
}

// getRandomResolver selects a random DNS resolver from the list.
func (c *Client) getRandomResolver() string {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.resolvers[c.rand.Intn(len(c.resolvers))]
}

// buildDNSMessage constructs a DNS message with the given host and request type.
func buildDNSMessage(host string, qtype uint16) *dns.Msg {
	msg := new(dns.Msg)
	msg.Id = dns.Id()
	msg.RecursionDesired = true
	msg.Question = []dns.Question{
		{
			Name:   dns.Fqdn(host),
			Qtype:  qtype,
			Qclass: dns.ClassINET,
		},
	}
	return msg
}

// parseARecords extracts A records from the DNS response.
func parseARecords(answer *dns.Msg) Result {
	var result Result
	for _, record := range answer.Answer {
		if a, ok := record.(*dns.A); ok {
			result.IPs = append(result.IPs, a.A.String())
			result.TTL = int(a.Header().Ttl)
		}
	}
	return result
}

// parseRecordsByType extracts records of the requested type from the DNS response.
func parseRecordsByType(answer *dns.Msg, requestType uint16) []string {
	var results []string
	for _, record := range answer.Answer {
		switch requestType {
		case dns.TypeA:
			if t, ok := record.(*dns.A); ok {
				results = append(results, t.A.String())
			}
		case dns.TypeNS:
			if t, ok := record.(*dns.NS); ok {
				results = append(results, t.Ns)
			}
		case dns.TypeCNAME:
			if t, ok := record.(*dns.CNAME); ok {
				results = append(results, t.Target)
			}
		case dns.TypeSOA:
			if t, ok := record.(*dns.SOA); ok {
				results = append(results, t.String())
			}
		case dns.TypePTR:
			if t, ok := record.(*dns.PTR); ok {
				results = append(results, t.Ptr)
			}
		case dns.TypeMX:
			if t, ok := record.(*dns.MX); ok {
				results = append(results, t.String())
			}
		case dns.TypeTXT:
			if t, ok := record.(*dns.TXT); ok {
				results = append(results, t.String())
			}
		case dns.TypeAAAA:
			if t, ok := record.(*dns.AAAA); ok {
				results = append(results, t.AAAA.String())
			}
		}
	}
	return results
}
