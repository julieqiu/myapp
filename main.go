package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
)

const (
	BaldorHost   = "https://www.baldorfood.com"
	URLNewLogin  = "https://www.baldorfood.com/users/default/new-login"
	URLCart      = "https://baldorfood.com/cart"
	URLAddToCart = "https://baldorfood.com/ecommerce/shoppingcart/cart.addToCart"
)

var (
	BaldorCookieName  = os.Getenv("BALDORFOOD_COOKIE_NAME")
	BaldorCookieValue = os.Getenv("BALDORFOOD_COOKIE_VALUE")
	shoppingList      = []string{
		"product/vegetables/car10-large-carrots",
	}
)

func main() {
	c := New()
	for _, url := range shoppingList {
		// TODO: make functions work
		key := c.GetProductKey(url)
		c.AddToCart(key)
	}
	log.Println("Done! See what's in your cart at https://baldorfood.com/cart")
}

type Client struct {
	http.Client
}

func New() *Client {
	var client *Client
	if BaldorCookieName == "" || BaldorCookieValue == "" {
		client = newClientWithAuthentication()
	} else {
		client = newClientWithCookies(BaldorCookieName, BaldorCookieValue)
	}
	return client
}

func newClientWithCookies(name, val string) *Client {
	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatal(err)
	}
	u, err := url.Parse(BaldorHost)
	if err != nil {
		log.Fatal(err)
	}
	cookieJar.SetCookies(u, []*http.Cookie{{Name: name, Value: val}})
	return &Client{http.Client{Jar: cookieJar}}
}

func newClientWithAuthentication() *Client {
	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatal(err)
	}
	email := os.Getenv("BALDORFOOD_EMAIL")
	pass := os.Getenv("BALDORFOOD_PASSWORD")
	fmt.Printf("Getting cookies for baldorfood.com for: %q %q\n", email, pass)
	client := &Client{http.Client{Jar: cookieJar}}

	if _, err := client.PostForm(URLNewLogin, url.Values{
		"EmailLoginForm[email]":      {email},
		"EmailLoginForm[password]":   {pass},
		"EmailLoginForm[rememberMe]": {"1"},
		"yt0":                        {"SIGN IN"},
	}); err != nil {
		log.Fatal(err)
	}
	if _, err := client.Get(URLCart); err != nil {
		log.Fatal(err)
	}

	u, err := url.Parse(BaldorHost)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("The following cookies are needed to authenticate your request to baldorfood.com")
	var found bool
	for _, c := range client.Jar.Cookies(u) {
		if c.Name == "provided_access" || c.Name == "PHPSESSID" {
			continue
		}
		found = true
		fmt.Printf("Name: %q\n", c.Name)
		fmt.Printf("Value: %q\n", c.Value)
	}
	if !found {
		log.Fatal("Authentication failed.")
	}
	fmt.Println("Save them as the environment variables BALDORFOOD_COOKIE_NAME and BALDORFOOD_COOKIE_VALUE to skip authenticating in the future.")
	return client
}
func (c *Client) AddToCart(productKey string) {
	if _, err := c.PostForm(URLAddToCart, url.Values{
		"ShoppingCartModel[key]":      {productKey},
		"ShoppingCartModel[quantity]": {"1"},
		"ShoppingCartModel[unit]":     {"CTN"},
		"qty":                         {"1"},
	}); err != nil {
		log.Fatal(err)
	}
}

func (c *Client) GetProductKey(urlPath string) string {
	resp, err := c.Get(urlPath)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	/*
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
	*/

	// TODO: get product model key from HTML
	return ""
}
