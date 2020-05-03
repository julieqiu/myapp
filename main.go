package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/julieqiu/baldorfood/internal"
)

const BaldorHost = "https://www.baldorfood.com"

var (
	shoppingList = []string{
		"product/vegetables/car10-large-carrots",
	}
)

func main() {
	shopper, err := newShopper()
	if err != nil {
		log.Fatal(err)
	}
	cookies, err := baldorCookies(shopper.Jar)
	if err != nil {
		log.Fatal(err)
	}
	findItems(cookies)
}

func newShopper() (*internal.Shopper, error) {
	baldorCookieName := os.Getenv("BALDORFOOD_COOKIE_NAME")
	baldorCookieValue := os.Getenv("BALDORFOOD_COOKIE_VALUE")
	if baldorCookieName != "" && baldorCookieValue != "" {
		return internal.NewShopperWithCookies([]*http.Cookie{{Name: baldorCookieName, Value: baldorCookieValue}})
	}

	shopper, err := internal.NewShopperWithAuthentication()
	if err != nil {
		return nil, err
	}

	var found bool
	cookies, err := baldorCookies(shopper.Jar)
	if err != nil {
		return nil, err
	}
	for _, c := range cookies {
		if c.Name == "provided_access" || c.Name == "PHPSESSID" {
			continue
		}
		found = true
		fmt.Println()
		fmt.Printf("export BALDORFOOD_COOKIE_NAME=%q\n", c.Name)
		fmt.Printf("export BALDORFOOD_COOKIE_VALUE=%q\n", c.Value)
	}
	fmt.Println()
	if !found {
		log.Fatal("Authentication failed.")
	} else {
		fmt.Println("Save these variables in your .bashrc to skip authenticating in the future.")
	}
	return shopper, nil
}

func baldorCookies(j http.CookieJar) ([]*http.Cookie, error) {
	u, err := url.Parse(BaldorHost)
	if err != nil {
		return nil, err
	}
	return j.Cookies(u), nil
}

func findItems(cookies []*http.Cookie) {
	file, err := os.Open("categories.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		u := scanner.Text() + "?viewall=1"
		fmt.Println("Visiting: ", u)
		c := internal.NewColly(BaldorHost, cookies)
		c.ConfirmLogin()
		items := c.GetItemsOnPage(u)
		fmt.Println(len(items))
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
