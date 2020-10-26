package parser

import (
	"fmt"
	"testing"

	"github.com/gocolly/colly/v2"
	"go.zoe.im/surferua"
)

func TestAvito(t *testing.T) {
	c := colly.NewCollector(
			colly.UserAgent(surferua.NewBot()),
		)
	// Find and visit all links "div.price-value"
	c.OnHTML("html", func(e *colly.HTMLElement) {
		fmt.Println( e.Text)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
		fmt.Println("Headers:", r.Headers)
	})

	err := c.Visit("https://www.avito.ru/moskva/avtomobili/volkswagen_amarok_2018_1994544125")
	if err != nil{
		t.Error(err)
	}
}


