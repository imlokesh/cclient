package cclient

import (
	"time"

	http "github.com/imlokesh/fhttp"

	"golang.org/x/net/proxy"

	utls "github.com/imlokesh/utls"
)

func NewClient(clientHello utls.ClientHelloID, proxyUrl string, allowRedirect bool, timeout time.Duration) (http.Client, error) {
	if len(proxyUrl) > 0 {
		dialer, err := newConnectDialer(proxyUrl)
		if err != nil {
			if allowRedirect {
				return http.Client{
					Timeout: timeout,
				}, err
			}
			return http.Client{
				Timeout: timeout,
				CheckRedirect: func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				},
			}, err
		}
		if allowRedirect {
			return http.Client{
				Transport: newRoundTripper(clientHello, dialer),
				Timeout:   timeout,
			}, nil
		}
		return http.Client{
			Transport: newRoundTripper(clientHello, dialer),
			Timeout:   timeout,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}, nil
	} else {
		if allowRedirect {
			return http.Client{
				Transport: newRoundTripper(clientHello, proxy.Direct),
				Timeout:   timeout,
			}, nil
		}
		return http.Client{
			Transport: newRoundTripper(clientHello, proxy.Direct),
			Timeout:   timeout,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}, nil

	}
}
