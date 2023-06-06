// Command fetch recreates the data in the products directory.
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/julieqiu/myapp/internal"
)

func main() {
	shopper, err := newShopper()
	if err != nil {
		log.Fatal(err)
	}
	cookie, err := internal.BaldorCookie(shopper.Jar)
	if err != nil {
		log.Fatal(err)
	}
	findItems([]*http.Cookie{cookie})
}

func newShopper() (*internal.Shopper, error) {
	baldorCookieName := os.Getenv("BALDORFOOD_COOKIE_NAME")
	baldorCookieValue := os.Getenv("BALDORFOOD_COOKIE_VALUE")
	if baldorCookieName != "" && baldorCookieValue != "" {
		return internal.NewShopperWithCookies([]*http.Cookie{{Name: baldorCookieName, Value: baldorCookieValue}})
	}

	email, pass, err := internal.ReadEmailAndPassword()
	if err != nil {
		return nil, err
	}
	log.Printf("Getting cookies for myapp.com for: %q \n", email)
	shopper, err := internal.NewShopperWithAuthentication(email, pass)
	if err != nil {
		return nil, err
	}
	return shopper, nil
}

func findItems(cookies []*http.Cookie) {
	file, err := os.Open("categories.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		findItemsAtURL(scanner.Text()+"?viewall=1", cookies)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func findItemsAtURL(u string, cookies []*http.Cookie) error {
	c := internal.NewColly(internal.BaldorHost, cookies)
	fmt.Printf("Visiting %q\n", u)
	items := c.GetItemsOnPage(u)
	if len(items) > 0 {
		if err := writeItemsToFile(u, items); err != nil {
			return err
		}
		return nil
	}
	fmt.Printf("No items found on %q\n", u)
	return nil
}

func writeItemsToFile(u string, items []*internal.Item) error {
	filename := internal.ProductsDir() + "/" + strings.ReplaceAll(
		strings.TrimSuffix(strings.TrimPrefix(u, "https://www.myapp.com/products"), "?viewall=1"),
		"/", "_") + ".json"

	jsonString, err := json.Marshal(items)
	if err != nil {
		return err
	}
	os.WriteFile(filename, jsonString, os.ModePerm)
	fmt.Printf("Wrote %d items to %q\n", len(items), filename)
	return nil
}
