package internal

import (
	"net/http"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

type Colly struct {
	Host    string
	Cookies []*http.Cookie
	colly   *colly.Collector
}

func NewColly(hostURL string, cookies []*http.Cookie) *Colly {
	// Instantiate default collector
	c := colly.NewCollector()
	c.SetCookies(hostURL, cookies)
	c.AllowURLRevisit = false
	c.SetRequestTimeout(60 * time.Second)
	return &Colly{colly: c, Host: hostURL, Cookies: cookies}
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
	Category   string
}

func (c *Colly) IsLoggedIn() bool {
	loggedIn := true
	c.colly.OnHTML(".support-container", func(e *colly.HTMLElement) {
		if !strings.Contains(e.ChildText(".user-support"), "Julie Qiu") {
			loggedIn = false
		}
	})
	c.colly.Visit(BaldorHost)
	return loggedIn
}

func (c *Colly) GetItemsOnPage(url string) []*Item {

	var items []*Item
	c.colly.OnHTML(".items", func(e *colly.HTMLElement) {
		e.ForEach(".table-cover-back", func(_ int, elem *colly.HTMLElement) {
			productKey, _ := elem.DOM.Find(".add-cart-wrap").Find("input.jq-increase").Attr("id")
			productKey = strings.TrimPrefix(productKey, "increase-item-")
			unit, _ := elem.DOM.Find("input.ShoppingCartModel_unit").Attr("value")
			link, _ := elem.DOM.Find(".product-title-and-sku").Find("h3").Find("a").Attr("href")
			category := strings.Split(strings.TrimPrefix(link, "/product/"), "/")[0]
			item := &Item{
				Brand:      elem.ChildText(".pct-farm"),
				Label:      elem.ChildText(".special-category-label"),
				Title:      elem.DOM.Find(".product-title-and-sku").Find("h3").Text(),
				Link:       c.Host + link,
				Sku:        elem.ChildText(".product-sku-inline"),
				Price:      elem.ChildText(".price"),
				Unit:       unit,
				Img:        elem.ChildAttr("img", "src"),
				ProductKey: productKey,
				Category:   category,
			}
			items = append(items, item)
		})
	})
	c.colly.Visit(url)
	return items
}
