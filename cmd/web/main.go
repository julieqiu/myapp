package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/julieqiu/baldorfood/internal"
)

var allItems = loadItems()

type Page struct {
	Categories []string
	Items      []*internal.Item
}

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("content/static"))))
	http.HandleFunc("/add/", handleAddToCart)
	http.HandleFunc("/", handler)
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

func handler(w http.ResponseWriter, r *http.Request) {
	tmpl := "index.html"
	p := &Page{Categories: categories()}
	if r.URL.Path != "/" {
		items := itemsForCategory(strings.TrimPrefix(r.URL.Path, "/"))
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

func itemsForCategory(category string) []*internal.Item {
	items, ok := allItems[category]
	if !ok {
		return nil
	}
	return items
}

func categories() []string {
	var categories []string
	for c := range allItems {
		categories = append(categories, c)
	}
	sort.Strings(categories)
	return categories
}

func loadItems() map[string][]*internal.Item {
	var files []string
	err := filepath.Walk("json", func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	allItems := map[string][]*internal.Item{}
	for _, filename := range files {
		file, err := ioutil.ReadFile(filename)
		if err != nil {
			log.Fatal(err)
		}

		var items []*internal.Item
		if err := json.Unmarshal([]byte(file), &items); err != nil {
			log.Fatal(err)
		}
		category := strings.Split(filename, "_")[1]
		allItems[category] = append(allItems[category], items...)
	}
	return allItems
}
