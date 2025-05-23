// Package forward implements a forwarding proxy.
package forward

import (
	"context"
	"fmt"
	"sync"

	"github.com/coredns/coredns/plugin/metadata"
	"github.com/coredns/coredns/request"
)

// Used to track the last used proxy for a request
var (
	lastUsedProxyMutex sync.RWMutex
	lastUsedProxyMap   = make(map[string]string)
)

// SetLastUsedProxy sets the last used proxy for a given request
func SetLastUsedProxy(requestID string, proxyAddr string) {
	lastUsedProxyMutex.Lock()
	defer lastUsedProxyMutex.Unlock()
	lastUsedProxyMap[requestID] = proxyAddr
}

// GetLastUsedProxy gets the last used proxy for a given request
func GetLastUsedProxy(requestID string) string {
	lastUsedProxyMutex.RLock()
	defer lastUsedProxyMutex.RUnlock()
	if addr, ok := lastUsedProxyMap[requestID]; ok {
		// Clean up after retrieving to prevent memory leaks
		delete(lastUsedProxyMap, requestID)
		return addr
	}
	return ""
}

// Metadata implements the metadata.Provider interface.
func (f *Forward) Metadata(ctx context.Context, state request.Request) context.Context {
	// Generate a unique request ID
	requestID := fmt.Sprintf("%s-%d", state.QName(), state.Req.Id)

	// We can't know the exact upstream server at this point yet
	// because it depends on the proxy selection policy during ServeDNS.
	// But we can pre-register a metadata function that will access the upstream
	// when it's needed.
	metadata.SetValueFunc(ctx, "forward/upstream", func() string {
		// First check if we have a record of what proxy was actually used
		if addr := GetLastUsedProxy(requestID); addr != "" {
			return addr
		}

		// Fallback to the original behavior
		if f.Len() > 0 {
			proxies := f.List()
			if len(proxies) > 0 {
				return proxies[0].Addr()
			}
		}
		return "no upstream selected"
	})

	// Also add zone information to help debugging
	metadata.SetValueFunc(ctx, "forward/zone", func() string {
		return f.from
	})

	return ctx
}
