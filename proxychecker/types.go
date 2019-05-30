package proxychecker

import (
	"fmt"
	"strings"
)

type Proxy string

type ProxyResultState int

const (
	Unknown ProxyResultState = iota
	ProxyInvalid
	Timeout
	Error
	Forbidden
	Success
)

func (st ProxyResultState) String() string {
	if st == ProxyInvalid {
		return "invalid proxy"
	}
	if st == Timeout {
		return "timeout"
	}
	if st == Error {
		return "error"
	}
	if st == Forbidden {
		return "forbidden"
	}
	if st == Success {
		return "success"
	}
	return "unknown"
}

type ProxyResult struct {
	Proxy  Proxy
	State  ProxyResultState
	Error  error
	IP     string
	Status int
}

func (pr ProxyResult) SetError(err error) ProxyResult {
	pr.Error = err
	pr.State = Error
	return pr
}

func (pr ProxyResult) String() string {
	bld := strings.Builder{}
	bld.WriteString(fmt.Sprintf("%s: %s (status %d", pr.Proxy, pr.State, pr.Status))
	if pr.IP != "" {
		bld.WriteString(", ip ")
		bld.WriteString(pr.IP)
	}
	if pr.Error != nil {
		bld.WriteString(", error ")
		bld.WriteString(pr.Error.Error())
	}
	bld.WriteString(")")
	return bld.String()
}

type ProxyResults []ProxyResult

func (prs ProxyResults) Len() int {
	return len(prs)
}

func (prs ProxyResults) Less(i, j int) bool {
	a, b := prs[i], prs[j]
	if a.State == b.State {
		return a.Proxy < b.Proxy
	}
	return a.State < b.State
}

func (prs ProxyResults) Swap(i, j int) {
	prs[j], prs[i] = prs[i], prs[j]
}
