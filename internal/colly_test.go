package internal

import (
	"fmt"
	"net/http"
	"os"
	"testing"
)

func TestGetItemsOnPage(t *testing.T) {
	baldorCookieName := os.Getenv("BALDORFOOD_COOKIE_NAME")
	baldorCookieValue := os.Getenv("BALDORFOOD_COOKIE_VALUE")
	shopper, err := NewShopperWithCookies([]*http.Cookie{{Name: baldorCookieName, Value: baldorCookieValue}})
	if err != nil {
		t.Fatal(err)
	}
	cookies, err := BaldorCookies(shopper.Jar)
	if err != nil {
		t.Fatal(err)
	}
	c := NewColly(BaldorHost, cookies)
	if !c.IsLoggedIn() {
		t.Fatal("Not logged in.")
	}
	u := BaldorHost + "/products/vegetables/fresh-herbs"
	items := c.GetItemsOnPage(u)
	want := 30
	if len(items) != want {
		t.Fatalf("Unxpected number of items on %q: got = %d; want = %d", u, len(items), want)
	}
	for _, item := range items {
		fmt.Printf("%+v\n", item)
	}
}
