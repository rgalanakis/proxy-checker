# proxy-checker

Check whether a proxy is valid, or blocked by a server.

Pass in proxy servers to check on stdin, one per-line.

## Usage

Run `proxy-checker help` for more details on usage.
Flags as of this writing:

```
proxy-checker makes a bunch of calls to an API and reports which succeed and which fail.
  -concurrency int
    	Max number of proxies to try simultaneously. Certain proxies may rate-limit a source IP to prevent DDoS. (default 2)
  -forbidden int
    	HTTP status code returned for blocked proxies (default 403)
  -ip-parser string
    	Name of the parser to use ('body' is only choice right now) (default "body")
  -max int
    	Total number of proxies to check. Limit is here for your sanity so you don't get classified as a DDoS by testing too many proxies. (default 20)
  -method string
    	Method to call on URL (default "GET")
  -ok int
    	HTTP status code for successful proxies (default 200)
  -proxy-pass string
    	Password for the proxy service being tested
  -proxy-user string
    	Username for the proxy service being tested
  -sleep int
    	Milliseconds to sleep between proxy checks. Certain proxies may rate-limit a source IP to prevent DDoS. (default 1000)
  -trace
    	Enable http tracing to help debug what is going on
  -url string
    	URL to check (default "https://google.com")
exit status 2
```

## Simple Example

We can grab a free proxy IPs from https://free-proxy-list.net/,
and see if it's blocked by the site in question:

```sh
echo "217.23.69.146:8080" | go run main.go -trace -url='https://robg3d.com'
```

## Complex Example

For example, here is a command to check all UK VPNs offered by NordVPN,
and what is blocked by Google:

```sh
test -f "temp/ovpn.zip" || wget -O temp/ovpn.zip 'https://downloads.nordcdn.com/configs/archives/servers/ovpn.zip'
unzip -Z1 temp/ovpn.zip \ # Print just the filenames, one per-line
    | grep "ovpn_tcp" \ # Only get tcp lines; ignore udp
    | sed "s/ovpn_[a-z][a-z]p\///" \ # Strip leading "ovpn_tcp/"
    | sed "s/.[a-z][a-z]p.ovpn//" \ # Strip trailing .tcp.ovpn
    | grep "uk" \ # Limit to just "uk" servers
    | sed -n '1,20p;21q' \ # Grab the first 20 only
    | go run proxy-checker/main.go \ # Run the command
        -url="https://google.com" \
        -forbidden=403 \
        -ok=200 \
        -concurrency=1 \
        -sleep=1000 \
        -proxy-user="${PROXY_USER}" \
        -proxy-pass="${PROXY_PASS}"
```

You can throw that in a Makefile or bash file and run it with something like
`PROXY_USER=username PROXY_PASS=password make check-nord-vpns`
