package main

import (
	"fmt"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

// scrapePowerInterruptions scrapes power interruption notices from KPLC's website.
func scrapePowerInterruptions() []*PowerInterruptionNotice {
	var powerInterruptionNotices []*PowerInterruptionNotice

	// Main collector
	collector := colly.NewCollector(
		colly.UserAgent(getUserAgent()),
		colly.AllowedDomains("kplc.co.ke"),
	)

	// Collector to scrape KPLC power interruptions details
	detailCollector := collector.Clone()

	// <main class="content">
	// 	<h2 class="generictitle">
	// 		<a href="..."></a>
	// 	</h2>
	// </main>
	collector.OnHTML(`main.content h2.generictitle a[href]`, func(e *colly.HTMLElement) {
		interruptionsLink := e.Attr("href")

		// Use detailCollector to scrape the power interruption files
		fmt.Println(interruptionsLink)
		detailCollector.Visit(interruptionsLink)
	})

	// <ul class="pagination">
	// 	<li>
	// 		<a href="..." rel="next"></a>
	// 	</li>
	// </ul>
	collector.OnHTML(`ul.pagination li a[rel="next"]`, func(e *colly.HTMLElement) {
		nextLink := e.Attr("href")
		fmt.Println("Next link:", nextLink)

		// Visit next page
		// e.Request.Visit(nextLink)
	})

	collector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// Extract interruptions
	detailCollector.OnHTML(`div.attachments a[class="docicon"]`, func(e *colly.HTMLElement) {
		powerInterruptionNotices = append(powerInterruptionNotices, processDownloadLink(e))
	})

	detailCollector.OnHTML(`div.genericintro div.intro ul li a[href]`, func(e *colly.HTMLElement) {
		powerInterruptionNotices = append(powerInterruptionNotices, processDownloadLink(e))
	})

	detailCollector.OnHTML(`div.genericintro a[class="download"]`, func(e *colly.HTMLElement) {
		powerInterruptionNotices = append(powerInterruptionNotices, processDownloadLink(e))
	})

	// Start scraping KPLC website with power interruptions
	collector.Visit(kplcURL)

	return powerInterruptionNotices
}

func processDownloadLink(e *colly.HTMLElement) *PowerInterruptionNotice {
	fileName := e.Text
	url := e.Attr("href")
	date := extractDateFromFilename(fileName)

	return &PowerInterruptionNotice{
		Filename: fileName,
		Date:     date,
		URL:      url,
	}
}

func extractDateFromFilename(fileName string) *Date {
	re := regexp.MustCompile(`\d{2}.\d{2}.\d{4}`)
	date := re.FindString(fileName)

	if date != "" {
		s := strings.Split(date, ".")

		day, _ := strconv.Atoi(s[0])
		month, _ := strconv.Atoi(s[1])
		year, _ := strconv.Atoi(s[2])

		return &Date{
			Day:   day,
			Month: time.Month(month),
			Year:  year,
		}
	}

	return nil
}

// getUserAgent generates User-Agent string to be used by HTTP requests
func getUserAgent() string {
	return fmt.Sprintf(
		"%s - %s (go; %s; %s-%s)",
		appName, appVersion, runtime.Version(),
		runtime.GOARCH, runtime.GOOS,
	)
}
