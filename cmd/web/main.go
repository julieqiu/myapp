package main

import (
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

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("content/static"))))
	http.HandleFunc("/add/", handleAddToCart)
	http.HandleFunc("/", handleViewProducts(allItems))
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
	}
	if err := s.AddToCart(productKey, unit); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleViewProducts(allItems map[string][]*internal.Item) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl := "index.html"
		p := &Page{Categories: categories(allItems)}
		if r.URL.Path != "/" {
			items := itemsForCategory(strings.TrimPrefix(r.URL.Path, "/"), allItems)
			if len(items) == 0 {
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			}
			p.Items = items
			tmpl = "product.html"
		}
		t, err := template.ParseFiles("content/static/html/" + tmpl)
		if err != nil {
			log.Fatal(err)
		}
		t.Execute(w, p)
	}
}

func itemsForCategory(category string, allItems map[string][]*internal.Item) []*internal.Item {
	items, ok := allItems[category]
	if !ok {
		return nil
	}
	return items
}

func categories(allItems map[string][]*internal.Item) []string {
	var categories []string
	for c := range allItems {
		categories = append(categories, c)
	}
	sort.Strings(categories)
	return categories
}
