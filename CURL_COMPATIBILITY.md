# cURL Options Compatibility Status in gURL

This document summarizes the compatibility status of cURL options in gURL's `--curl` mode.

## ✅ SUPPORTED Options

### Basic Operations
- `-d, --data <data>` - HTTP POST data
- `--data-raw <data>` - HTTP POST data, '@' allowed
- `--data-binary <data>` - HTTP POST binary data
- `--data-ascii <data>` - HTTP POST ASCII data
- `-X, --request <method>` - Specify HTTP request method
- `-G, --get` - Use GET method
- `-I, --head` - HEAD request only
- `-T, --upload-file <file>` - Upload file to destination

### Headers & Authentication
- `-H, --header <header>` - Pass custom headers to server
- `-A, --user-agent <name>` - Send User-Agent to server
- `-u, --user <user:password>` - Server user and password (converted to Authorization header)
- `-e, --referer <URL>` - Set Referer header

### Output Options
- `-o, --output <file>` - Write output to file instead of stdout
- `-i, --include` - Include protocol response headers in output
- `-v, --verbose` - Make the operation more talkative
- `-s, --silent` - Silent mode
- `-S, --show-error` - Show error even when -s is used
- `-f, --fail` - Fail silently on HTTP errors

### Network & Connection
- `--connect-timeout <seconds>` - Maximum time allowed for connection
- `-m, --max-time <seconds>` - Maximum time allowed for the transfer
- `--keepalive-time <seconds>` - Interval time for keepalive probes
- `--unix-socket <path>` - Connect through Unix domain socket
- `-4, --ipv4` - Resolve names to IPv4 addresses only
- `-6, --ipv6` - Resolve names to IPv6 addresses only

### TLS/SSL Options
- `-k, --insecure` - Allow insecure server connections when using SSL
- `--cacert <file>` - CA certificate to verify peer against
- `-E, --cert <certificate>` - Client certificate file
- `--key <key>` - Private key file name

### Redirects & Cookies
- `-L, --location` - Follow redirects
- `--max-redirs <num>` - Maximum number of redirects allowed
- `--location-trusted` - Like --location, and send auth to other hosts
- `-b, --cookie <data>` - Send cookies from string/file
- `-c, --cookie-jar <filename>` - Write cookies to file after operation

### Compression & Transfer
- `--compressed` - Request compressed response
- `-r, --range <range>` - Retrieve only the bytes within RANGE

### Proxy
- `-x, --proxy <[protocol://]host[:port]>` - Use this proxy
- `--noproxy <no-proxy-list>` - List of hosts which do not use proxy

### Misc
- `-h, --help` - Show help text
- `-V, --version` - Show version number
- `--url <url>` - URL to work with

---

## ❌ NOT SUPPORTED (and why)

### Protocol-Specific (FTP, SMTP, Telnet, TFTP)
These are not HTTP/HTTPS protocols and don't apply to gURL's use case:
- `--ftp-*` - All FTP-specific options
- `--mail-*` - SMTP email options
- `-t, --telnet-option` - Telnet options
- `--tftp-*` - TFTP options
- `-P, --ftp-port` - FTP PORT command
- `-Q, --quote` - Send command to server before transfer

### Advanced Authentication (Not HTTP Basic)
- `--anyauth` - Pick any authentication method
- `--basic` - Use HTTP Basic Authentication (implicit in -u)
- `--digest` - Use HTTP Digest Authentication
- `--negotiate` - Use HTTP Negotiate (SPNEGO) authentication
- `--ntlm` - Use HTTP NTLM authentication
- `--ntlm-wb` - Use HTTP NTLM authentication with winbind
- `--oauth2-bearer <token>` - OAuth 2 Bearer Token
- `--delegation` - GSS-API delegation permission
- `--krb <level>` - Enable Kerberos

### Advanced Proxy Features
- `--proxy-*` - All proxy-specific authentication and TLS options
- `--preproxy` - Use this proxy first
- `-p, --proxytunnel` - Operate through HTTP proxy tunnel
- `--proxy1.0` - Use HTTP/1.0 proxy
- `--socks4/5*` - SOCKS proxy options
- `--proxy-anyauth, --proxy-basic, --proxy-digest, --proxy-ntlm` - Proxy auth methods

### Alternative Protocols & Features
- `--alt-svc` - Enable alt-svc with cache file
- `--http0.9` - Allow HTTP 0.9 responses
- `-0, --http1.0` - Force HTTP 1.0
- `--http1.1` - Force HTTP 1.1
- `--http2` - Use HTTP 2
- `--http2-prior-knowledge` - Use HTTP 2 without HTTP/1.1 Upgrade
- `--http3` - Use HTTP v3
- `--haproxy-protocol` - Send HAProxy PROXY protocol v1 header

### Form Data & Multipart
- `-F, --form <name=content>` - Specify multipart MIME data
- `--form-string <name=string>` - Specify multipart MIME data
- `--data-urlencode <data>` - HTTP POST data url encoded

### Advanced Transfer Options
- `-C, --continue-at <offset>` - Resumed transfer offset
- `-a, --append` - Append to target file when uploading
- `-O, --remote-name` - Write output to a file named as the remote file
- `-J, --remote-header-name` - Use the header-provided filename
- `--remote-name-all` - Use the remote file name for all URLs
- `-R, --remote-time` - Set the remote file's time on the local output

### Rate Limiting & Progress
- `--limit-rate <speed>` - Limit transfer speed to RATE
- `-Y, --speed-limit <speed>` - Stop transfers slower than this
- `-y, --speed-time <seconds>` - Trigger 'speed-limit' abort after this time
- `-#, --progress-bar` - Display transfer progress as a bar
- `--no-progress-meter` - Do not show the progress meter
- `-Z, --parallel` - Perform transfers in parallel
- `--parallel-immediate, --parallel-max` - Parallel transfer options

### Network Interface Control
- `--interface <name>` - Use network INTERFACE (or address)
- `--local-port <num/range>` - Force use of RANGE for local port numbers
- `--dns-interface, --dns-ipv4-addr, --dns-ipv6-addr, --dns-servers` - DNS control
- `--doh-url` - Resolve host names over DOH
- `--resolve <host:port:address>` - Resolve the host+port to this address
- `--connect-to <HOST1:PORT1:HOST2:PORT2>` - Connect to host

### Retry & Error Handling
- `--retry <num>` - Retry request if transient problems occur
- `--retry-connrefused` - Retry on connection refused
- `--retry-delay <seconds>` - Wait time between retries
- `--retry-max-time <seconds>` - Retry only within this period
- `--fail-early` - Fail on first transfer error

### Advanced TLS/SSL
- `--cert-status` - Verify the status of the server certificate
- `--cert-type <type>` - Certificate file type (DER/PEM/ENG)
- `--key-type <type>` - Private key file type
- `--ciphers <list>` - SSL ciphers to use
- `--tls-max, --tls13-ciphers` - TLS version and cipher control
- `--tlsv1, --tlsv1.0, --tlsv1.1, --tlsv1.2, --tlsv1.3` - Force TLS version
- `-2, --sslv2, -3, --sslv3` - Use SSLv2/v3 (deprecated)
- `--ssl, --ssl-reqd` - Require SSL/TLS
- `--ssl-allow-beast` - Allow security flaw to improve interop
- `--ssl-no-revoke` - Disable cert revocation checks
- `--no-alpn, --no-npn` - Disable ALPN/NPN TLS extensions
- `--no-sessionid` - Disable SSL session-ID reusing
- `--pinnedpubkey <hashes>` - Public key to verify peer against
- `--crlfile <file>` - Get a CRL list in PEM format
- `--capath <dir>` - CA directory to verify peer against
- `--pass <phrase>` - Pass phrase for the private key
- `--tlsauthtype, --tlspassword, --tlsuser` - TLS authentication options

### Netrc & Config Files
- `-n, --netrc` - Must read .netrc for user name and password
- `--netrc-file <filename>` - Specify FILE for netrc
- `--netrc-optional` - Use either .netrc or URL
- `-K, --config <file>` - Read config from a file
- `-q, --disable` - Disable .curlrc

### Output Formatting & Debugging
- `-w, --write-out <format>` - Use output FORMAT after completion
- `--trace <file>` - Write a debug trace to FILE
- `--trace-ascii <file>` - Like --trace, but without hex output
- `--trace-time` - Add time stamps to trace/verbose output
- `--libcurl <file>` - Dump libcurl equivalent code
- `--stderr` - Where to redirect stderr
- `--styled-output` - Enable styled output for HTTP headers
- `-N, --no-buffer` - Disable buffering of the output stream

### Misc Advanced Options
- `-:, --next` - Make next URL use its separate set of options
- `-g, --globoff` - Disable URL sequences and ranges using {} and []
- `--ignore-content-length` - Ignore the size of the remote resource
- `--raw` - Do HTTP "raw"; no transfer decoding
- `--tr-encoding` - Request compressed transfer encoding
- `--post301, --post302, --post303` - Do not switch to GET after redirect
- `--request-target` - Specify the target for this request
- `-z, --time-cond <time>` - Transfer based on a time condition
- `--max-filesize <bytes>` - Maximum file size to download
- `-l, --list-only` - List only mode
- `-B, --use-ascii` - Use ASCII/text transfer
- `-M, --manual` - Display the full manual
- `--abstract-unix-socket` - Connect via abstract Unix domain socket
- `--no-keepalive` - Disable TCP keepalive on the connection
- `--tcp-fastopen` - Use TCP Fast Open
- `--tcp-nodelay` - Use the TCP_NODELAY option
- `--happy-eyeballs-timeout-ms` - IPv6 to IPv4 fallback timeout
- `--suppress-connect-headers` - Suppress proxy CONNECT response headers
- `--xattr` - Store metadata in extended file attributes
- `--egd-file, --random-file` - Random data sources
- `--engine` - Crypto engine to use
- `--etag-save, --etag-compare` - ETag handling
- `--expect100-timeout` - How long to wait for 100-continue
- `--false-start` - Enable TLS False Start
- `--disallow-username-in-url` - Disallow username in url
- `--disable-eprt, --disable-epsv` - Disable EPRT/EPSV
- `--create-dirs` - Create necessary local directory hierarchy
- `--crlf` - Convert LF to CRLF in upload
- `--metalink` - Process given URLs as metalink XML file
- `--proto, --proto-default, --proto-redir` - Protocol control
- `--service-name` - SPNEGO service name
- `--sasl-authzid, --sasl-ir` - SASL authentication options
- `-j, --junk-session-cookies` - Ignore session cookies read from file
- `--path-as-is` - Do not squash .. sequences in URL path
- `--pubkey, --hostpubmd5` - SSH public key options
- `--compressed-ssh` - Enable SSH compression

---

## 📊 Summary Statistics

- **Total cURL options**: ~230+
- **Supported in gURL**: ~45
- **Coverage**: ~20% (focused on HTTP/HTTPS essentials)

## 🎯 Design Philosophy

gURL's curl mode focuses on:
1. **HTTP/HTTPS essentials** - Core web request functionality
2. **Common use cases** - Options used in 80% of curl commands
3. **Security basics** - TLS, certificates, authentication
4. **Developer workflow** - Headers, data, redirects, output

**Not included**: Protocol-specific features (FTP, SMTP, etc.), advanced proxy configurations, obscure TLS options, and features rarely used in modern HTTP workflows.

## 📝 Notes

- The `--curl` flag switches gURL to HTTP mode (from gRPC mode)
- The `--compat` flag (planned) will allow curl-style options to work in gRPC mode by translating them
- gURL prioritizes simplicity and the most common curl use cases
- For advanced curl features, users should use the original curl tool
