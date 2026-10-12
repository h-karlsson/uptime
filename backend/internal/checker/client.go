package checker

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/netip"
	"syscall"
	"time"
)

const maxRedirects = 5

var errBlockedAddress = errors.New("address is not publicly routable")

// newClient returns an HTTP client that refuses to connect to private, loopback and other
// non-public addresses, so users can't point monitors at the server's own network (SSRF).
// The check runs on the resolved IP of every connection, including redirects, so DNS tricks
// can't get around it.
func newClient() *http.Client {
	dialer := &net.Dialer{
		Timeout: 10 * time.Second,
		Control: func(_, address string, _ syscall.RawConn) error {
			addrPort, err := netip.ParseAddrPort(address)
			if err != nil || !isPublic(addrPort.Addr()) {
				return errBlockedAddress
			}
			return nil
		},
	}
	transport := &http.Transport{
		Proxy:               nil,
		DialContext:         dialer.DialContext,
		TLSHandshakeTimeout: 10 * time.Second,
		MaxIdleConnsPerHost: 2,
		IdleConnTimeout:     90 * time.Second,
	}
	return &http.Client{
		Transport: transport,
		CheckRedirect: func(_ *http.Request, via []*http.Request) error {
			if len(via) >= maxRedirects {
				return fmt.Errorf("stopped after %d redirects", maxRedirects)
			}
			return nil
		},
	}
}

func isPublic(addr netip.Addr) bool {
	addr = addr.Unmap()
	return addr.IsGlobalUnicast() &&
		!addr.IsPrivate() &&
		!addr.IsLoopback() &&
		!addr.IsLinkLocalUnicast() &&
		!sharedAddressSpace.Contains(addr)
}

// Carrier-grade NAT range (RFC 6598); not covered by IsPrivate.
var sharedAddressSpace = netip.MustParsePrefix("100.64.0.0/10")
