package proxychecker

import (
	"errors"
	"fmt"
	"golang.org/x/net/http/httpproxy"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
	"sync/atomic"
	"time"
)

type Checker struct {
	Cfg       Config
	mutex     *sync.Mutex
	doneCount *int64
	waitGroup *sync.WaitGroup
}

func NewChecker(cfg Config) *Checker {
	c := &Checker{Cfg: cfg}
	c.mutex = &sync.Mutex{}
	i := int64(0)
	c.doneCount = &i
	c.waitGroup = &sync.WaitGroup{}
	return c
}

func (c Checker) CheckProxies() (ProxyResults, error) {
	proxies := c.Cfg.Proxies

	results := make([]ProxyResult, len(proxies))

	c.waitGroup.Add(len(proxies))
	semaphore := make(chan struct{}, c.Cfg.Concurrency)
	for i := range proxies {
		go func(i int) {
			defer c.done()
			semaphore <- struct{}{}
			time.Sleep(c.Cfg.Sleep)
			results[i] = c.CheckProxy(Proxy(proxies[i]))
			<-semaphore
		}(i)
	}
	c.waitGroup.Wait()
	fmt.Print("\r")
	return results, nil
}

func (c Checker) done() {
	c.waitGroup.Done()
	nowDone := atomic.AddInt64(c.doneCount, 1)
	if nowDone%6 == 5 {
		c.mutex.Lock()
		defer c.mutex.Unlock()
		fmt.Printf("Processed %d of %d proxies\r", nowDone, len(c.Cfg.Proxies))
	}
}

func (c Checker) CheckProxy(proxy Proxy) (result ProxyResult) {
	result.Proxy = proxy

	req, err := http.NewRequest(c.Cfg.Method, c.Cfg.URL, nil)
	if err != nil {
		return result.SetError(err)
	}
	if c.Cfg.Trace {
		req = AddTrace(req)
	}

	resp, err := c.newClient(proxy).Do(req)
	// If the proxy is bad, we get a 307, but we don't have a strong reference to the error.
	// So we need to check its error string...
	if uerr, ok := err.(*url.Error); ok && uerr.Err.Error() == http.StatusText(http.StatusProxyAuthRequired) {
		result.State = ProxyInvalid
		return
	}
	if err != nil {
		return result.SetError(err)
	}

	result.Status = resp.StatusCode
	if result.Status == c.Cfg.OkStatus {
		result.State = Success
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return result.SetError(err)
	}

	if result.Status != c.Cfg.ForbiddenStatus {
		return result.SetError(errors.New("unexpected status, body: " + string(body)))
	}

	result.State = Forbidden
	result.IP = c.Cfg.ParseIP(resp, string(body))
	if result.IP == "" {
		return result.SetError(errors.New("could not parse IP from 403 body"))
	}
	return
}

func (c Checker) newClient(proxy Proxy) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			Proxy: func(req *http.Request) (*url.URL, error) {
				config := &httpproxy.Config{HTTPProxy: string(proxy)}
				ur, err := config.ProxyFunc()(req.URL)
				if err != nil {
					return nil, err
				}
				if c.Cfg.ProxyUser != "" {
					ur.User = url.UserPassword(c.Cfg.ProxyUser, c.Cfg.ProxyPass)
				}
				return ur, err
			},
			DisableKeepAlives: false,
		},
	}
}
