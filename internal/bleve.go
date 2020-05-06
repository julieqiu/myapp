package internal

import (
	"github.com/blevesearch/bleve"
)

type Searcher struct {
	index bleve.Index
	items map[string]*Item
}

func NewSearcher(items []*Item) (*Searcher, error) {
	productIndex, err := bleve.Open(IndexDir())
	if err != nil {
		return nil, err
	}
	itemLookup := map[string]*Item{}
	for _, item := range items {
		itemLookup[item.ProductKey] = item
	}
	return &Searcher{index: productIndex, items: itemLookup}, nil
}

func NewIndex(indexName string) (*Searcher, error) {
	indexMapping := bleve.NewIndexMapping()
	productIndex, err := bleve.New(indexName, indexMapping)
	if err != nil {
		return nil, err
	}
	return &Searcher{index: productIndex}, nil
}

func (s *Searcher) Index(item *Item) {
	s.index.Index(item.ProductKey, item)
}

const numResults = 500

func (s *Searcher) Search(query string) ([]*Item, error) {
	q := bleve.NewFuzzyQuery(query)
	searchRequest := bleve.NewSearchRequestOptions(q, numResults, 0, false)
	searchResult, err := s.index.Search(searchRequest)
	if err != nil {
		return nil, err
	}

	var items []*Item
	for _, hit := range searchResult.Hits {
		if item, ok := s.items[hit.ID]; ok {
			items = append(items, item)
		}
	}
	return items, err
}

func (s *Searcher) DocCount() (uint64, error) {
	return s.index.DocCount()
}
