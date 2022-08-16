# CLI Sitemap generator

## Description
Command line tool that can generate a sitemap for a given url passes as an argument

## Usage
```
./sitemap -help
./sitemap <URL> -parallel=4 -max-depth=2  -timeout=10 -max-retry=2 -retry-delay=10 -output-file=./sitemapgoogle
```
EG: `./sitemap https://google.com -parallel=4 -timeout=10 -max-depth=2`

## Flags to be used
- help
 	- Show help

- max-depth int
 	- Maximum depth of URL navigation recursion (default 3)

- max-retry int
 	- Number request retries to make to an URL (default 3)

- output-file string
	- Sitemap output filepath without extension (default "output-sitemap")

- parallel int
 	- Number of parallel workers to navigate throught site (default 1)

- retry-delay int
 	- Number of seconds to wait before making a new request to an URL that failed (default 30)

- timeout int
 	- Number of seconds to wait for response from server (default 30)

## Features
- run every request in parallel
- no external dependencies
- extracts URLs from \<a> HTML tags
- can have retries on URLs that are failing to respond
