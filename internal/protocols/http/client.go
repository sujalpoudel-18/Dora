package http

import (
	"crypto/tls"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/dora/dora/internal/configuration"
)

type Response struct {
	StatusCode    int
	Headers       map[string]string
	Body          string
	ResponseSize  int64
	ResponseTime  time.Duration
	FinalURL      string
	RedirectChain []string
	ContentType   string
	TLSVersion    string
}

type Client struct {
	httpClient *http.Client
	opts       *configuration.Options
}

func NewClient(opts *configuration.Options) *Client {
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   opts.Timeout,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:        opts.Concurrency * 2,
		MaxIdleConnsPerHost: opts.HostConcurrency,
		IdleConnTimeout:     90 * time.Second,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   opts.Timeout,
	}

	return &Client{
		httpClient: client,
		opts:       opts,
	}
}

func (c *Client) Execute(method, reqURL string, headers map[string]string, body string) (*Response, error) {
	start := time.Now()

	var bodyReader io.Reader
	if body != "" {
		bodyReader = strings.NewReader(body)
	}

	req, err := http.NewRequest(method, reqURL, bodyReader)
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	if req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", "DORA-Scanner/"+configuration.Version)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respTime := time.Since(start)

	maxSize := c.opts.MaxResponseSize
	if maxSize <= 0 {
		maxSize = 5 * 1024 * 1024 // Default 5MB
	}
	reader := io.LimitReader(resp.Body, maxSize)
	bodyBytes, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	respHeaders := make(map[string]string)
	for k, v := range resp.Header {
		if len(v) > 0 {
			respHeaders[k] = v[0]
		}
	}

	tlsVer := ""
	if resp.TLS != nil {
		switch resp.TLS.Version {
		case tls.VersionTLS10:
			tlsVer = "TLS 1.0"
		case tls.VersionTLS11:
			tlsVer = "TLS 1.1"
		case tls.VersionTLS12:
			tlsVer = "TLS 1.2"
		case tls.VersionTLS13:
			tlsVer = "TLS 1.3"
		}
	}

	return &Response{
		StatusCode:   resp.StatusCode,
		Headers:      respHeaders,
		Body:         string(bodyBytes),
		ResponseSize: int64(len(bodyBytes)),
		ResponseTime: respTime,
		FinalURL:     resp.Request.URL.String(),
		ContentType:  resp.Header.Get("Content-Type"),
		TLSVersion:   tlsVer,
	}, nil
}
