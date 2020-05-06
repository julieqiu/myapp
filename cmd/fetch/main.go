// Command fetch recreates the data in the products directory.
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/julieqiu/baldorfood/internal"
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

	email, pass, err := readEmailAndPassword()
	if err != nil {
		return nil, err
	}
	log.Printf("Getting cookies for baldorfood.com for: %q \n", email)
	shopper, err := internal.NewShopperWithAuthentication(email, pass)
	if err != nil {
		return nil, err
	}
	printCookies(shopper.Jar)
	return shopper, nil
}

func readEmailAndPassword() (string, string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Enter email: ")
	email, err := reader.ReadString('\n')
	if err != nil {
		return "", "", err
	}
	fmt.Printf("Enter password: ")
	pass, err := reader.ReadString('\n')
	if err != nil {
		return "", "", err
	}
	return strings.TrimSuffix(email, "\n"), strings.TrimSuffix(pass, "\n"), nil
}

func printCookies(j http.CookieJar) error {
	var found bool
	cookies, err := internal.BaldorCookies(j)
	if err != nil {
		return err
	}
	for _, c := range cookies {
		if c.Name == "provided_access" || c.Name == "PHPSESSID" {
			continue
		}
		found = true
		log.Println()
		log.Printf("export BALDORFOOD_COOKIE_NAME=%q\n", c.Name)
		log.Printf("export BALDORFOOD_COOKIE_VALUE=%q\n", c.Value)
	}
	if !found {
		return fmt.Errorf("Authentication failed.")
	}
	log.Println("Save these variables in your .bashrc to skip authenticating in the future.")
	return nil
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
	filename := "data/" + strings.ReplaceAll(
		strings.TrimSuffix(strings.TrimPrefix(u, "https://www.baldorfood.com/products"), "?viewall=1"),
		"/", "_") + ".data"

	jsonString, err := json.Marshal(items)
	if err != nil {
		return err
	}
	ioutil.WriteFile(filename, jsonString, os.ModePerm)
	fmt.Printf("Wrote %d items to %q\n", len(items), filename)
	return nil
}
