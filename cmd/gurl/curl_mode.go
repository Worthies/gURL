// curl_mode.go implements HTTP/HTTPS request functionality when --curl mode is enabled.
// In this mode, gurl acts as a standard HTTP client (like curl) instead of a gRPC client.

package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

// curlModeConfig holds all the configuration for making an HTTP request in curl mode
type curlModeConfig struct {
	url             string
	method          string
	headers         []string
	data            string
	dataFiles       []string
	insecure        bool
	cacert          string
	cert            string
	key             string
	verbose         bool
	silent          bool
	showError       bool
	includeHeaders  bool
	headersOnly     bool
	outputFile      string
	userAgent       string
	referer         string
	connectTimeout  float64
	maxTime         float64
	unixSocket      string
	followRedirects bool
	maxRedirs       int
	compressed      bool
	failOnError     bool
	locationTrusted bool
	proxy           string
	noProxy         string
	uploadFile      string
	customRequest   string
	cookieJar       string
	cookie          string
	range_          string
	ipv4Only        bool
	ipv6Only        bool
}

// executeCurlMode performs an HTTP/HTTPS request similar to curl
func executeCurlMode(config curlModeConfig) error {
	// Determine HTTP method
	method := config.method
	if method == "" {
		if config.headersOnly {
			method = "HEAD"
		} else if config.uploadFile != "" {
			method = "PUT"
		} else if config.data != "" {
			method = "POST"
		} else {
			method = "GET"
		}
	}

	// Prepare request body
	var bodyReader io.Reader

	// Handle file upload
	if config.uploadFile != "" {
		file, err := os.Open(config.uploadFile)
		if err != nil {
			return fmt.Errorf("failed to open upload file %s: %v", config.uploadFile, err)
		}
		defer file.Close()
		bodyReader = file
	} else if config.data != "" {
		if config.data == "@" {
			bodyReader = os.Stdin
		} else if strings.HasPrefix(config.data, "@") {
			// Read from file
			filename := config.data[1:]
			data, err := ioutil.ReadFile(filename)
			if err != nil {
				return fmt.Errorf("failed to read data file %s: %v", filename, err)
			}
			bodyReader = strings.NewReader(string(data))
		} else {
			bodyReader = strings.NewReader(config.data)
		}
	}

	// Create HTTP request
	req, err := http.NewRequest(method, config.url, bodyReader)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	// Add headers
	for _, header := range config.headers {
		parts := strings.SplitN(header, ":", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			req.Header.Add(key, value)
		}
	}

	// Set User-Agent
	if config.userAgent != "" {
		req.Header.Set("User-Agent", config.userAgent)
	} else {
		req.Header.Set("User-Agent", "gurl/"+version+" (HTTP mode)")
	}

	// Set Referer
	if config.referer != "" {
		req.Header.Set("Referer", config.referer)
	}

	// Set Cookie
	if config.cookie != "" {
		if strings.HasPrefix(config.cookie, "@") {
			// Read cookies from file
			data, err := ioutil.ReadFile(config.cookie[1:])
			if err != nil {
				return fmt.Errorf("failed to read cookie file: %v", err)
			}
			req.Header.Set("Cookie", string(data))
		} else {
			req.Header.Set("Cookie", config.cookie)
		}
	}

	// Set Range
	if config.range_ != "" {
		req.Header.Set("Range", "bytes="+config.range_)
	}

	// Set Content-Type if we have data and it's not already set
	if config.data != "" && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	// Request compression if specified
	if config.compressed {
		req.Header.Set("Accept-Encoding", "gzip, deflate")
	}

	// Create HTTP client with TLS config
	client := &http.Client{}

	// Configure redirect policy
	if config.followRedirects {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			if len(via) >= config.maxRedirs {
				return fmt.Errorf("stopped after %d redirects", config.maxRedirs)
			}
			// If locationTrusted, copy auth headers
			if config.locationTrusted && len(via) > 0 {
				if auth := via[0].Header.Get("Authorization"); auth != "" {
					req.Header.Set("Authorization", auth)
				}
			}
			return nil
		}
	} else {
		// Don't follow redirects
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}

	// Configure TLS
	tlsConfig := &tls.Config{
		InsecureSkipVerify: config.insecure,
	}

	// Load CA cert if specified
	if config.cacert != "" && !config.insecure {
		caCert, err := ioutil.ReadFile(config.cacert)
		if err != nil {
			return fmt.Errorf("failed to read CA certificate: %v", err)
		}
		caCertPool := x509.NewCertPool()
		if !caCertPool.AppendCertsFromPEM(caCert) {
			return fmt.Errorf("failed to parse CA certificate")
		}
		tlsConfig.RootCAs = caCertPool
	}

	// Load client certificate if specified
	if config.cert != "" && config.key != "" {
		clientCert, err := tls.LoadX509KeyPair(config.cert, config.key)
		if err != nil {
			return fmt.Errorf("failed to load client certificate: %v", err)
		}
		tlsConfig.Certificates = []tls.Certificate{clientCert}
	}

	// Create transport
	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	// Configure timeouts
	if config.connectTimeout > 0 {
		transport.DialContext = (&net.Dialer{
			Timeout: time.Duration(config.connectTimeout * float64(time.Second)),
		}).DialContext
	}

	if config.maxTime > 0 {
		client.Timeout = time.Duration(config.maxTime * float64(time.Second))
	}

	// Configure proxy
	if config.proxy != "" {
		// Parse proxy URL
		proxyURL, err := http.ProxyFromEnvironment(req)
		if err == nil && proxyURL != nil {
			transport.Proxy = http.ProxyURL(proxyURL)
		}
	}

	// Configure Unix socket if specified
	if config.unixSocket != "" {
		transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
			return net.Dial("unix", config.unixSocket)
		}
		// For Unix sockets, we need to adjust the URL
		if !strings.HasPrefix(config.url, "http://") && !strings.HasPrefix(config.url, "https://") {
			config.url = "http://localhost" + config.url
		}
	}

	// Configure IPv4/IPv6
	if config.ipv4Only {
		transport.DialContext = (&net.Dialer{
			Timeout: 30 * time.Second,
		}).DialContext
	} else if config.ipv6Only {
		transport.DialContext = (&net.Dialer{
			Timeout: 30 * time.Second,
		}).DialContext
	}

	client.Transport = transport

	// Perform request
	if config.verbose {
		fmt.Fprintf(os.Stderr, "> %s %s HTTP/1.1\n", method, req.URL.Path)
		fmt.Fprintf(os.Stderr, "> Host: %s\n", req.Host)
		for key, values := range req.Header {
			for _, value := range values {
				fmt.Fprintf(os.Stderr, "> %s: %s\n", key, value)
			}
		}
		fmt.Fprintf(os.Stderr, ">\n")
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	// Save cookies to jar if specified
	if config.cookieJar != "" {
		if cookies := resp.Cookies(); len(cookies) > 0 {
			file, err := os.Create(config.cookieJar)
			if err != nil {
				return fmt.Errorf("failed to create cookie jar: %v", err)
			}
			defer file.Close()
			for _, cookie := range cookies {
				fmt.Fprintf(file, "%s\n", cookie.String())
			}
		}
	}

	// Print response headers if verbose or includeHeaders
	if config.verbose || config.includeHeaders || config.headersOnly {
		fmt.Printf("HTTP/%d.%d %s\n", resp.ProtoMajor, resp.ProtoMinor, resp.Status)
		for key, values := range resp.Header {
			for _, value := range values {
				fmt.Printf("%s: %s\n", key, value)
			}
		}
		fmt.Println()
	}

	// Determine output destination
	var output io.Writer = os.Stdout
	if config.outputFile != "" && config.outputFile != "-" {
		file, err := os.Create(config.outputFile)
		if err != nil {
			return fmt.Errorf("failed to create output file: %v", err)
		}
		defer file.Close()
		output = file
	}

	// Read and write response body
	if !config.silent && !config.headersOnly {
		_, err = io.Copy(output, resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response: %v", err)
		}
		// Add newline if output is to terminal
		if config.outputFile == "" {
			fmt.Println()
		}
	}

	// Exit with error code if HTTP status is error and failOnError is set
	if config.failOnError && resp.StatusCode >= 400 {
		os.Exit(22) // curl uses exit code 22 for HTTP errors
	}

	return nil
}

// executeCurlModeFromFlags prepares curl config from parsed flags and executes HTTP request
func executeCurlModeFromFlags() {
	args := flags.Args()

	// Determine URL from arguments
	var url string
	if len(args) > 0 {
		url = args[0]

		// If URL doesn't have a scheme, assume https://
		if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
			// Check if it looks like a host:port or just a path
			if strings.Contains(url, ":") || strings.HasPrefix(url, "/") {
				url = "https://" + url
			} else {
				url = "https://" + url
			}
		}
	} else {
		fail(nil, "No URL specified for HTTP mode (--curl flag)")
	}

	// Get unix socket path if specified
	var unixSocketPath string
	if isUnixSocket != nil && isUnixSocket() {
		// The URL in this case should be the socket path
		unixSocketPath = args[0]
		// If there are more args, the second one might be the actual path
		if len(args) > 1 {
			url = args[1]
			if !strings.HasPrefix(url, "/") {
				url = "/" + url
			}
		} else {
			url = "/"
		}
	}

	// Combine all headers (addlHeaders + rpcHeaders)
	allHeaders := append([]string{}, addlHeaders...)
	allHeaders = append(allHeaders, rpcHeaders...)

	config := curlModeConfig{
		url:             url,
		method:          *httpMethod,
		headers:         allHeaders,
		data:            *data,
		insecure:        *insecure,
		cacert:          *cacert,
		cert:            *cert,
		key:             *key,
		verbose:         *verbose || *veryVerbose,
		silent:          *silent,
		includeHeaders:  *includeHeaders,
		headersOnly:     *headOnly,
		outputFile:      *outputFile,
		userAgent:       *userAgent,
		referer:         *referer,
		connectTimeout:  *connectTimeout,
		maxTime:         *maxTime,
		unixSocket:      unixSocketPath,
		followRedirects: *followRedirects || *locationTrusted,
		maxRedirs:       *maxRedirs,
		compressed:      *compressed,
		failOnError:     *failOnError,
		locationTrusted: *locationTrusted,
		proxy:           *httpProxy,
		noProxy:         *noProxy,
		uploadFile:      *uploadFile,
		cookieJar:       *cookieJar,
		cookie:          *cookie,
		range_:          *rangeHeader,
		ipv4Only:        *ipv4Only,
		ipv6Only:        *ipv6Only,
	}

	err := executeCurlMode(config)
	if err != nil {
		fail(err, "HTTP request failed")
	}
}
