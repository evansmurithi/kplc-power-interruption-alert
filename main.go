package main

import (
	"fmt"
	"time"
)

const appName string = "evansmurithi/kplc-power-interruption-alert"
const appVersion string = "v0.0.1"

const kplcURL string = "https://kplc.co.ke/category/view/50/planned-power-interruptions"

// Date holds day, month and year information
type Date struct {
	Day   int
	Month time.Month
	Year  int
}

// PowerInterruptionNotice holds info on power interruption notices from KPLC website.
type PowerInterruptionNotice struct {
	Filename string
	Date     *Date
	URL      string
}

func main() {
	powerInterruptionNotices := scrapePowerInterruptions()

	for _, notice := range powerInterruptionNotices {
		fmt.Printf(
			"[%d/%s/%d] %s >> %s\n",
			notice.Date.Day, notice.Date.Month.String(),
			notice.Date.Year, notice.Filename, notice.URL,
		)
	}
}
