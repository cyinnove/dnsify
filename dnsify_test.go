package dnsify

import (
	"testing"
	"github.com/miekg/dns"
	"github.com/stretchr/testify/assert"
)

// Mock DNS server address
const mockDNS = "127.0.0.1:53535"

// mockDNSHandler handles incoming DNS requests for testing
func mockDNSHandler(w dns.ResponseWriter, req *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(req)
	m.Authoritative = true

	switch req.Question[0].Qtype {
	case dns.TypeA:
		rr, _ := dns.NewRR("example.com. 3600 IN A 192.0.2.1")
		m.Answer = append(m.Answer, rr)
	case dns.TypeMX:
		rr, _ := dns.NewRR("example.com. 3600 IN MX 10 mail.example.com.")
		m.Answer = append(m.Answer, rr)
	}
	_ = w.WriteMsg(m)
}

// setupMockDNSServer starts a mock DNS server on the localhost
func setupMockDNSServer() (*dns.Server, error) {
	dns.HandleFunc(".", mockDNSHandler)

	server := &dns.Server{Addr: mockDNS, Net: "udp"}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			panic(err)
		}
	}()
	return server, nil
}

// TestResolve tests the Resolve function
func TestResolve(t *testing.T) {
	server, err := setupMockDNSServer()
	if err != nil {
		t.Fatalf("Failed to set up mock DNS server: %v", err)
	}
	defer server.Shutdown()

	client := New([]string{mockDNS}, 3)
	result, err := client.Resolve("example.com")
	assert.NoError(t, err)
	assert.NotEmpty(t, result.IPs)
	assert.Equal(t, "192.0.2.1", result.IPs[0])
	assert.Equal(t, 3600, result.TTL)
}

// TestResolveRaw tests the ResolveRaw function
func TestResolveRaw(t *testing.T) {
	server, err := setupMockDNSServer()
	if err != nil {
		t.Fatalf("Failed to set up mock DNS server: %v", err)
	}
	defer server.Shutdown()

	client := New([]string{mockDNS}, 3)
	results, raw, err := client.ResolveRaw("example.com", dns.TypeMX)
	assert.NoError(t, err)
	assert.NotEmpty(t, results)
	assert.Contains(t, raw, "mail.example.com")
	assert.Equal(t, "example.com.\t3600\tIN\tMX\t10 mail.example.com.", results[0])
}

// TestDo tests the Do function with a custom DNS message
func TestDo(t *testing.T) {
	server, err := setupMockDNSServer()
	if err != nil {
		t.Fatalf("Failed to set up mock DNS server: %v", err)
	}
	defer server.Shutdown()

	client := New([]string{mockDNS}, 3)

	msg := new(dns.Msg)
	msg.SetQuestion(dns.Fqdn("example.com"), dns.TypeA)

	resp, err := client.Do(msg)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "example.com.\t3600\tIN\tA\t192.0.2.1", resp.Answer[0].String())
}
