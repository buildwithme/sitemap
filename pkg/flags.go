package pkg

import (
	"flag"
	"os"
	"runtime"
	"sync"
)

type ParameterOptions struct {
	OutputFile string
	Parallel   int
	MaxDepth   int
	Timeout    int
	RetryDetay int
	MaxRetry   int
}

var options *ParameterOptions
var once sync.Once

const (
	flags = "FLAGS"
)

func GetParameterOptions() *ParameterOptions {
	once.Do(initSitemapFlags)
	return options
}

func initSitemapFlags() {
	help := flag.Bool("help", false, "Show help")
	outputFile := flag.String("output-file", "output-sitemap", "Sitemap output filepath without extension")
	parallel := flag.Int("parallel", 1, "Number of parallel workers to navigate throught site")
	maxDepth := flag.Int("max-depth", 3, "Maximum depth of URL navigation recursion")
	timeout := flag.Int("timeout", 30, "Number of seconds to wait for response from server")
	maxRetry := flag.Int("max-retry", 3, "Number request retries to make to an URL")
	retryDetail := flag.Int("retry-delay", 30, "Number of seconds to wait before making a new request to an URL that failed")

	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	runtime.GOMAXPROCS(*parallel)

	options = &ParameterOptions{
		OutputFile: *outputFile,
		Parallel:   *parallel,
		MaxDepth:   *maxDepth,
		Timeout:    *timeout,
		RetryDetay: *retryDetail,
		MaxRetry:   *maxRetry,
	}
}
