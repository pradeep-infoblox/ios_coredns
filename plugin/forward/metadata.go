// Package forward implements a forwarding proxy.
package forward

import (
	"context"

	"github.com/coredns/coredns/plugin/metadata"
	"github.com/coredns/coredns/request"
)

// Metadata implements the metadata.Provider interface.
func (f *Forward) Metadata(ctx context.Context, state request.Request) context.Context {
	// We can't know the exact upstream server at this point yet
	// because it depends on the proxy selection policy during ServeDNS.
	// But we can pre-register a metadata function that will access the upstream
	// when it's needed.
	metadata.SetValueFunc(ctx, "forward/upstream", func() string {
		// This will use the last selected proxy if available, or return a default message
		if f.Len() > 0 {
			proxies := f.List()
			if len(proxies) > 0 {
				return proxies[0].Addr()
			}
		}
		return "no upstream selected"
	})
	return ctx
}
