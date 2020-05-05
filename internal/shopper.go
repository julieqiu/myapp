package internal

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

type Shopper struct {
	http.Client
}

func NewShopperWithCookies(cookies []*http.Cookie) (*Shopper, error) {
	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	u, err := url.Parse(BaldorHost)
	if err != nil {
		return nil, err
	}
	cookieJar.SetCookies(u, cookies)
	return &Shopper{http.Client{Jar: cookieJar}}, nil
}

func NewShopperWithAuthentication(email, pass string) (*Shopper, error) {
	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	shopper := &Shopper{http.Client{Jar: cookieJar}}

	if _, err := shopper.PostForm(URLNewLogin, url.Values{
		"EmailLoginForm[email]":      {email},
		"EmailLoginForm[password]":   {pass},
		"EmailLoginForm[rememberMe]": {"1"},
		"yt0":                        {"SIGN IN"},
	}); err != nil {
		return nil, err
	}
	if _, err := shopper.Get(URLCart); err != nil {
		return nil, err
	}
	return shopper, nil
}

func (s *Shopper) AddToCart(productKey, unit string) error {
	v := url.Values{
		"ShoppingCartModel[key]":      {productKey},
		"ShoppingCartModel[unit]":     {unit},
		"ShoppingCartModel[quantity]": {"1"},
		"qty":                         {"1"},
	}
	resp, err := s.PostForm(URLAddToCart, v)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if strings.Contains(string(b), "error") {
		// The response always returns a status 200, even when there is an
		// error.
		return fmt.Errorf("Error adding to cart: %q; values = %v", string(b), v)
	}
	log.Printf("Success! %q\n", string(b))
	return nil
}

func BaldorCookies(j http.CookieJar) ([]*http.Cookie, error) {
	u, err := url.Parse(BaldorHost)
	if err != nil {
		return nil, err
	}
	return j.Cookies(u), nil
}
