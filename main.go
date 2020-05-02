package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
)

var shoppingList = []string{
	"product/vegetables/car10-large-carrots",
}

func main() {
	c := New()
	for _, url := range shoppingList {
		key := c.GetProductKey(url)
		c.AddToCart(key)
	}
	log.Println("Done! See what's in your cart at https://baldorfood.com/cart")
}

type Client struct {
	http.Client
}

func New() *Client {
	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatal(err)
	}
	email := os.Getenv("BALDORFOOD_EMAIL")
	pass := os.Getenv("BALDORFOOD_PASSWORD")
	fmt.Printf("Getting cookies for baldorfood.com for: %q %q\n", email, pass)
	client := &Client{http.Client{Jar: cookieJar}}

	doRequest("https://baldorfood.com/users/default/new-login", func(u string) (*http.Response, error) {
		return client.PostForm(u, url.Values{
			"EmailLoginForm[email]":    {email},
			"EmailLoginForm[password]": {pass},
		})
	})
	doRequest("https://baldorfood.com/cart", func(u string) (*http.Response, error) {
		return client.Get(u)
	})

	return client
}

func doRequest(url string, fn func(url string) (*http.Response, error)) {
	resp, err := fn(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	name := os.Getenv("BALDORFOOD_NAME")
	if !strings.Contains(string(b), name) {
		log.Fatalf("Authentication failed: %q not found on page %q.", name, url)
	}
}

func (c *Client) AddToCart(productKey string) {
	if _, err := c.PostForm("https://baldorfood.com/ecommerce/shoppingcart/cart.addToCart", url.Values{
		"ShoppingCartModel[key]":      {productKey},
		"ShoppingCartModel[quantity]": {"1"},
		"ShoppingCartModel[unit]":     {"CTN"},
		"qty":                         {"1"},
	}); err != nil {
		log.Fatal(err)
	}
}

func (c *Client) GetProductKey(urlPath string) string {
	_, err := c.Get(urlPath)
	if err != nil {
		log.Fatal(err)
	}
	// TODO: get product model key from HTML
	return ""
}
