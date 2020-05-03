package internal

import (
	"bufio"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
)

const (
	BaldorHost   = "https://www.baldorfood.com"
	URLNewLogin  = "https://www.baldorfood.com/users/default/new-login"
	URLCart      = "https://baldorfood.com/cart"
	URLAddToCart = "https://baldorfood.com/ecommerce/shoppingcart/cart.addToCart"
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

func NewShopperWithAuthentication() (*Shopper, error) {
	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Enter email: ")
	email, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	fmt.Printf("Enter password: ")
	pass, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	email = strings.TrimSuffix(email, "\n")
	pass = strings.TrimSuffix(pass, "\n")

	fmt.Printf("Getting cookies for baldorfood.com for: %q \n", email)
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

func (s *Shopper) AddToCart(productKey string) error {
	if _, err := s.PostForm(URLAddToCart, url.Values{
		"ShoppingCartModel[key]":      {productKey},
		"ShoppingCartModel[quantity]": {"1"},
		"ShoppingCartModel[unit]":     {"CTN"},
		"qty":                         {"1"},
	}); err != nil {
		return err
	}
	return nil
}
