package proxychecker

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

type Config struct {
	URL             string
	Method          string
	ForbiddenStatus int
	OkStatus        int
	Trace           bool
	ProxyUser       string
	ProxyPass       string
	ParseIP         func(*http.Response, string) (ip string)
	Realtime        bool
	Concurrency     int
	Proxies         []string
	Sleep           time.Duration
}

func ParseConfig(flag *flag.FlagSet) (Config, error) {
	flag.Usage = func() {
		_, _ = fmt.Fprintln(
			flag.Output(),
			"proxy-checker makes a bunch of calls to an API and reports which succeed and which fail.")
		flag.PrintDefaults()
	}
	checkUrl := flag.String("url", "https://google.com", "URL to check")
	method := flag.String("method", "GET", "Method to call on URL")
	forbidden := flag.Int("forbidden", 403, "HTTP status code returned for blocked proxies")
	ok := flag.Int("ok", 200, "HTTP status code for successful proxies")
	trace := flag.Bool("trace", false, "Enable http tracing to help debug what is going on")
	user := flag.String("proxy-user", "", "Username for the proxy service being tested")
	pass := flag.String("proxy-pass", "", "Password for the proxy service being tested")
	parserName := flag.String("ip-parser", "body",
		"Name of the parser to use ('body' is only choice right now)")
	realtime := flag.Bool("realtime", false,
		"Print results as they come in, rather than sorted at the end")
	concurrecy := flag.Int(
		"concurrency", 2,
		"Max number of proxies to try simultaneously. "+
			"Certain proxies may rate-limit a source IP to prevent DDoS.")
	max := flag.Int("max", 20,
		"Total number of proxies to check. Limit is here for your sanity so "+
			"you don't get classified as a DDoS by testing too many proxies.")
	sleep := flag.Int(
		"sleep", 1000,
		"Milliseconds to sleep between proxy checks. "+
			"Certain proxies may rate-limit a source IP to prevent DDoS.")

	err := flag.Parse(os.Args[1:])
	if err != nil {
		return Config{}, err
	}

	cfg := Config{
		URL:             *checkUrl,
		Method:          *method,
		ForbiddenStatus: *forbidden,
		OkStatus:        *ok,
		Trace:           *trace,
		ProxyUser:       *user,
		ProxyPass:       *pass,
		Concurrency:     *concurrecy,
		Proxies:         readProxiesFromStdin(),
		Sleep:           time.Millisecond * time.Duration(*sleep),
		Realtime:        *realtime,
	}

	if len(cfg.Proxies) > *max {
		return cfg, fmt.Errorf("%d exceeds max checks, limit your input or increase max", len(cfg.Proxies))
	}

	switch *parserName {
	case "body":
		cfg.ParseIP = parseIPFromBody
	default:
		return cfg, fmt.Errorf("unknown parser: %v", *parserName)
	}
	return cfg, err
}

var ipRegexp = regexp.MustCompile(`(\d\d?\d?\.\d\d?\d?\.\d\d?\d?\.\d\d?\d?)`)

func parseIPFromBody(_ *http.Response, body string) string {
	return ipRegexp.FindString(body)
}

// See https://flaviocopes.com/go-shell-pipes/
func readProxiesFromStdin() []string {
	info, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}

	if info.Mode()&os.ModeCharDevice != 0 || info.Size() <= 0 {
		fmt.Println("The command is intended to work with pipes.")
		fmt.Println("Send in the list of proxy services, one per line.")
		fmt.Println("cat proxies.txt | proxy-checker")
		os.Exit(2)
	}

	result := make([]string, 0, 512)
	reader := bufio.NewReader(os.Stdin)
	for {
		s, err := reader.ReadString('\n')
		if err != nil && err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Error reading stdin")
			os.Exit(1)
		}
		result = append(result, strings.TrimSpace(s))
	}
	return result
}
