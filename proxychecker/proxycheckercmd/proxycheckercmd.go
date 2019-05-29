package proxycheckercmd

import (
	"flag"
	"fmt"
	"github.com/rgalanakis/proxy-checker/proxychecker"
	"os"
	"sort"
)

func Main() {
	cfg, err := proxychecker.ParseConfig(flag.CommandLine)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(2)
	}

	checker := proxychecker.NewChecker(cfg)

	fmt.Printf("Testing %d proxies\n", len(cfg.Proxies))
	results, err := checker.CheckProxies()
	if err != nil {
		fmt.Println("Fatal error checking proxies:", err)
		os.Exit(1)
	}

	sort.Sort(results)
	for _, r := range results {
		fmt.Println(r)
	}
	os.Exit(0)
}
