// Command index indexes data from the products directory into a searcher index.
package main

import (
	"log"

	"github.com/julieqiu/myapp/internal"
)

func main() {
	searcher, err := internal.NewIndex(internal.IndexDir())
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Loading items from products...")
	items, err := internal.LoadItems()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Items loaded...")
	for category, itemsInCategory := range items {
		log.Println("indexing category: ", category)
		for _, item := range itemsInCategory {
			log.Println("indexing: ", item.Title)
			searcher.Index(item)
		}
	}
	count, err := searcher.DocCount()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Items indexed: %d\n", count)
}
