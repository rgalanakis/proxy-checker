package proxychecker

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/httptrace"
	"net/textproto"
)

func AddTrace(r *http.Request) *http.Request {
	trace := &httptrace.ClientTrace{
		GetConn: func(hostPort string) {
			fmt.Println("GetConn", hostPort)
		},
		GotConn: func(t httptrace.GotConnInfo) {
			fmt.Println("GotConn", t)
		},
		PutIdleConn: func(err error) {
			fmt.Println("PutIdleConn", err)
		},
		GotFirstResponseByte: func() {
			fmt.Println("GotFirstResponseByte")
		},
		Got100Continue: func() {
			fmt.Println("Got100Continue")
		},
		Got1xxResponse: func(code int, header textproto.MIMEHeader) error {
			fmt.Println("Got1xxResponse", code, header)
			return nil
		},
		DNSStart: func(dns httptrace.DNSStartInfo) {
			fmt.Println("DNSStart")
		},
		DNSDone: func(dns httptrace.DNSDoneInfo) {
			fmt.Println("DNSDone", dns)
		},
		ConnectStart: func(network, addr string) {
			fmt.Println("ConnectStart", network, addr)
		},
		ConnectDone: func(network, addr string, err error) {
			fmt.Println("ConnectDone", network, addr, err)
		},
		TLSHandshakeStart: func() {
			fmt.Println("TLSHandshakeStart")
		},
		TLSHandshakeDone: func(t tls.ConnectionState, e error) {
			fmt.Println("TLSHandshakeDone", t, e)
		},
		WroteHeaderField: func(key string, value []string) {
			fmt.Println("WroteHeaderField", key, value)
		},
		WroteHeaders: func() {
			fmt.Println("WroteHeaders")
		},
		Wait100Continue: func() {
			fmt.Println("Wait100Continue")
		},
		WroteRequest: func(t httptrace.WroteRequestInfo) {
			fmt.Println("WroteRequest", t)
		},
	}
	return r.WithContext(httptrace.WithClientTrace(r.Context(), trace))
}
