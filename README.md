# Time-Adjusted Market Index (TAMI)

A universal mechanism for calculating the estimated value of a collection of assets, written in Go. 

## Usage

```go
package main

import (
	"fmt"
	"time"

	"github.com/linkpoolio/go-tami"
)

var (
	now            = time.Now()
	twoDaysAgo     = now.AddDate(0, 0, -2)
	threeDaysAgo   = now.AddDate(0, 0, -3)
	oneMonthAgo    = now.AddDate(0, -1, 0)
	sixWeeksAgo    = now.AddDate(0, 0, -6*7)
)

func main() {
	transactions := []tami.Transaction{
		{ItemID: "Mars", Price: 612, Timestamp: sixWeeksAgo},
		{ItemID: "Hyacinth", Price: 700, Timestamp: oneMonthAgo},
		{ItemID: "Hyacinth", Price: 400, Timestamp: threeDaysAgo},
		{ItemID: "Mars", Price: 1200, Timestamp: twoDaysAgo},
	}
	timeAdjustedMarketValue := tami.TAMI(transactions)
	fmt.Println(timeAdjustedMarketValue) // 1832.411067193676
}

```

## Motivation

See the [JavaScript implementation](https://github.com/Mimicry-Protocol/TAMI) for background.