package internal

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gocolly/colly/v2"
)

type Colly struct {
	Host  string
	colly *colly.Collector
}

func NewColly(hostURL string, cookies []*http.Cookie) *Colly {
	// Instantiate default collector
	c := colly.NewCollector(
		colly.Async(true),
	)
	c.SetCookies(hostURL, cookies)
	c.AllowURLRevisit = false

	// Limit the maximum parallelism to 2
	// This is necessary if the goroutines are dynamically
	// created to control the limit of simultaneous requests.
	//
	// Parallelism can be controlled also by spawning fixed
	// number of go routines.
	c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 5})
	return &Colly{colly: c, Host: hostURL}
}

type Item struct {
	Brand      string
	Label      string
	Title      string
	Link       string
	Sku        string
	Price      string
	Unit       string
	Img        string
	ProductKey string
}

func (c *Colly) ConfirmLogin() {
	// Check that cookies are set properly so that the page is a logged-in
	// page.
	c.colly.OnHTML(".support-container", func(e *colly.HTMLElement) {
		if !strings.Contains(e.ChildText(".user-support"), "Julie Qiu") {
			log.Fatalf("Not logged in.")
		}
	})
	c.colly.Visit(c.Host)
}

func (c *Colly) GetItemsOnPage(url string) []*Item {
	var items []*Item
	c.colly.OnHTML(".items", func(e *colly.HTMLElement) {
		e.ForEach(".table-cover-back", func(_ int, elem *colly.HTMLElement) {
			productKey, _ := elem.DOM.Find(".add-cart-wrap").Find("input.jq-increase").Attr("id")
			productKey = strings.TrimPrefix(productKey, "increase-item-")
			item := &Item{
				Brand:      elem.ChildText(".pct-farm"),
				Label:      elem.ChildText(".special-category-label"),
				Title:      elem.ChildText("h3"),
				Link:       elem.ChildText("a[href]"),
				Sku:        elem.ChildText(".product-sku-inline"),
				Price:      elem.ChildText(".price"),
				Unit:       elem.ChildText(".unit"),
				Img:        c.Host + elem.ChildAttr("img", "src"),
				ProductKey: productKey,
			}
			fmt.Println(item.Title, item.Price)
			items = append(items, item)
		})
	})
	c.colly.Visit(url)
	c.colly.Wait()
	return items
}
