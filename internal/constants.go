package internal

import (
	"encoding/json"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	BaldorHost   = "https://www.baldorfood.com"
	URLNewLogin  = "https://www.baldorfood.com/users/default/new-login"
	URLCart      = "https://www.baldorfood.com/cart"
	URLAddToCart = "https://www.baldorfood.com/ecommerce/shoppingcart/cart.addToCart"
)

func LoadItems() (map[string][]*Item, error) {
	var files []string
	err := filepath.Walk(ProductsDir(), func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	allItems := map[string][]*Item{}
	for _, filename := range files {
		file, err := os.ReadFile(filename)
		if err != nil {
			return nil, err
		}

		var items []*Item
		if err := json.Unmarshal([]byte(file), &items); err != nil {
			return nil, err
		}
		category := strings.Split(filename, "_")[1]
		allItems[category] = append(allItems[category], items...)
	}
	return allItems, nil
}

func ProductsDir() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("No caller information")
	}
	return filepath.Join(path.Dir(path.Dir(filename)), "products")
}

func IndexDir() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("No caller information")
	}
	return filepath.Join(path.Dir(path.Dir(filename)), "product_index")
}
