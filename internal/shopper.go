package internal

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"

	"golang.org/x/crypto/ssh/terminal"
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

func (s *Shopper) AddToCart(productKey, unit string) (int, error) {
	v := url.Values{
		"ShoppingCartModel[key]":      {productKey},
		"ShoppingCartModel[unit]":     {unit},
		"ShoppingCartModel[quantity]": {"1"},
		"qty":                         {"1"},
	}
	resp, err := s.PostForm(URLAddToCart, v)
	if err != nil {
		return resp.StatusCode, err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if strings.Contains(string(b), "error") {
		// The response always returns a status 200, even when there is an
		// error.
		return http.StatusInternalServerError, fmt.Errorf("Error adding to cart: %q; values = %v", string(b), v)
	}
	log.Printf("Success! %q\n", string(b))
	return http.StatusOK, nil
}

func BaldorCookie(j http.CookieJar) (*http.Cookie, error) {
	u, err := url.Parse(BaldorHost)
	if err != nil {
		return nil, err
	}
	cookie := &http.Cookie{}
	for _, c := range j.Cookies(u) {
		if c.Name == "provided_access" || c.Name == "PHPSESSID" || c.Name == "SESSION" {
			continue
		}
		log.Printf("BALDORFOOD_COOKIE_NAME=%q\n", c.Name)
		log.Printf("BALDORFOOD_COOKIE_VALUE=%q\n", c.Value)
		cookie.Name = c.Name
		cookie.Value = c.Value
	}
	if cookie.Name == "" {
		return nil, fmt.Errorf("Authentication cookie not found")
	}
	return cookie, nil
}

func ReadEmailAndPassword() (string, string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Enter email: ")
	email, err := reader.ReadString('\n')
	if err != nil {
		return "", "", err
	}
	fmt.Printf("Enter password: ")
	pass, err := terminal.ReadPassword(0)
	if err != nil {
		log.Fatal(err)
	}
	return strings.TrimSuffix(email, "\n"), strings.TrimSuffix(string(pass), "\n"), nil
}
