package main

import (
	"fmt"
	"regexp"
	"time"

	"github.com/gocolly/colly/v2"
)

func Scrape() ([]string, error) {
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.99 Safari/537.36"),
	)

	url := "https://parade.com/1027247/jessicasager/best-breakup-songs/"
	domain := "parade.com"

	var youtubeLinks []string

	// Headers c/p'd from the browser to bypass 403
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Host", domain)
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:123.0) Gecko/20100101 Firefox/123.0")
		r.Headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
		r.Headers.Set("Accept-Language", "en-US,en;q=0.5")
		r.Headers.Set("Accept-Encoding", "gzip, deflate, br")
		r.Headers.Set("DNT", "1")
		r.Headers.Set("Sec-GPC", "1")
		r.Headers.Set("Connection", "keep-alive")
		r.Headers.Set("Cookie", "datadome=5MaJ5K2ybinK6cMjWM9VsArkyRtWe5QPX56Xfq22r1y4LS16uz8FY~DogjQmgu79lZmEXQwr2vJNWBjWsqPRIeg8UiDwHsYKjsCybk8pU8~Lt7BihwFZciTduvSdqEXE; _aren_ab=g=30/e:pxp-472-aff,s:a,ex:1709907823; muid=JxZbxFlD1Z5FFLzpmtBYWA; _lc2_fpi=1081db3850f9--01hqxbraynf5tc2edmfdhch13p; _lc2_fpi_meta=%7B%22w%22%3A1709308652502%7D; _sp_id.1b15=ae731019-1a6e-41a1-9658-fabc7053d660.1709308655.10.1709906030.1709864190.d7e19081-7515-425f-9b88-6ae4f5521c6f.d349f1e3-75f0-40dc-9fb5-989ecdb77047.e1fe478c-f307-4178-9055-b88e14920236.1709906025159.2; _lr_env_src_ats=false; ArenaGeo=eyJjb3VudHJ5Q29kZSI6IlVTIiwicmVnaW9uQ29kZSI6IlBBIiwiaW5FRUEiOmZhbHNlfQ==; _li_dcdm_c=.parade.com; _lr_sampling_rate=0; _ig=JxZbxFlD1Z5FFLzpmtBYWA; sp_debug_s=mz48r52ibsimyoazdqjrvl; 3a39163a163a30272a=l296lsguoikzw4exurw899; _sp_ses.1b15=*; _lr_retry_request=true; ArenaGeo=eyJjb3VudHJ5Q29kZSI6IlVTIiwicmVnaW9uQ29kZSI6IlBBIiwiaW5FRUEiOmZhbHNlfQ==; Cookie_1=value; _aren_ab=g=30/e:pxp-472-aff,s:a,ex:1709908043; datadome=2kbwoxKtVYtFDfwapZRRCQWJi7Y4PnDVwucAFVHVLWUQOiHlAjQDQGj9dB2NCQozJ0O90gnJWYSBB5smKFpelGRLjEOBZoQM1vbkJutVNB0PdxpLQPgmD82ZKl4i08V4")
		r.Headers.Set("Upgrade-Insecure-Requests", "1")
		r.Headers.Set("Sec-Fetch-Dest", "document")
		r.Headers.Set("Sec-Fetch-Mode", "navigate")
		r.Headers.Set("Sec-Fetch-Site", "cross-site")
		r.Headers.Set("Pragma", "no-cache")
		r.Headers.Set("Cache-Control", "no-cache")
		fmt.Println("Visiting", r.URL.String())
	})
	// grab those embeds
	c.OnHTML("phoenix-iframe[src]", func(e *colly.HTMLElement) {
		link := e.Attr("src")
		// regex to match YouTube links no I do not have this regex memorized
		re := regexp.MustCompile(`(?:https?:\/\/)?(?:www\.)?youtube-nocookie\.com\/embed\/([a-zA-Z0-9_-]+)`)
		matches := re.FindStringSubmatch(link)
		if len(matches) > 1 {
			// convert the embed link to an api-usable one and append to array
			slug := matches[1]
			youtubeLinks = append(youtubeLinks, slug)
		}
	})

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		RandomDelay: 5 * time.Second,
	})
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "HTTP STATUS: ", r.StatusCode)
	})

	// Start scraping
	err := c.Visit(url)
	if err != nil {
		fmt.Println("ERROR: ", err)
	}
	fmt.Println("YouTube slugs:")

	for _, link := range youtubeLinks {
		fmt.Println(link)
	}

	return youtubeLinks, nil
}
