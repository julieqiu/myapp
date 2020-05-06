package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/julieqiu/baldorfood/internal"
)

type Page struct {
	Categories []string
	Items      []*internal.Item
}

func main() {
	allItems, err := internal.LoadItems()
	if err != nil {
		log.Fatal(err)
	}

	var items []*internal.Item
	for _, itemsInCategory := range allItems {
		items = append(items, itemsInCategory...)
	}
	searcher, err := internal.NewSearcher(items)
	if err != nil {
		log.Fatal(err)
	}

	var categories []string
	for c := range allItems {
		categories = append(categories, c)
	}
	sort.Strings(categories)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("content/static"))))
	http.HandleFunc("/add/", handleAddToCart)
	http.HandleFunc("/login/", handleLogin)
	http.HandleFunc("/search/", handleSearch(searcher, categories))
	http.HandleFunc("/", handleViewProducts(categories, allItems))
	log.Print("Listening on localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleAddToCart(w http.ResponseWriter, r *http.Request) {
	parts := strings.SplitN(strings.TrimPrefix(r.URL.Path, "/add/"), "/", 2)
	productKey := parts[0]
	unit := parts[1]
	baldorCookieName := os.Getenv("BALDORFOOD_COOKIE_NAME")
	baldorCookieValue := os.Getenv("BALDORFOOD_COOKIE_VALUE")
	s, err := internal.NewShopperWithCookies([]*http.Cookie{{Name: baldorCookieName, Value: baldorCookieValue}})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := s.AddToCart(productKey, unit); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	if os.Getenv("BALDORFOOD_COOKIE_NAME") != "" &&
		os.Getenv("BALDORFOOD_COOKIE_VALUE") != "" {
		return
	}

	email, pass, err := internal.ReadEmailAndPassword()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println()
	log.Println("Authenticating... (this may take a minute or so)")
	shopper, err := internal.NewShopperWithAuthentication(email, pass)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	cookie, err := internal.BaldorCookie(shopper.Jar)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	os.Setenv("BALDORFOOD_COOKIE_NAME", cookie.Name)
	os.Setenv("BALDORFOOD_COOKIE_VALUE", cookie.Value)
}

func handleSearch(searcher *internal.Searcher, categories []string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		values := r.URL.Query()
		query, ok := values["q"]
		p := &Page{Categories: categories}
		if ok && len(query) == 1 {
			items, err := searcher.Search(query[0])
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			p.Items = items
		}
		t, err := template.ParseFiles("content/static/html/product.html")
		if err != nil {
			log.Fatal(err)
		}
		t.Execute(w, p)
	}
}

func handleViewProducts(categories []string, allItems map[string][]*internal.Item) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := &Page{Categories: categories}
		tmpl := "index.html"
		if r.URL.Path != "/" {
			category := strings.TrimPrefix(r.URL.Path, "/")
			items, ok := allItems[category]
			if !ok || len(items) == 0 {
				http.Redirect(w, r, fmt.Sprintf("/search?q=%s", category), http.StatusFound)
				return
			}
			p.Items = items
			tmpl = "product.html"
		}
		t, err := template.ParseFiles("content/static/html/" + tmpl)
		if err != nil {
			log.Fatal(err)
		}
		t.Execute(w, p)
		return
	}
}
