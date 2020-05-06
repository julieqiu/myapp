package internal

import (
	"fmt"

	"github.com/blevesearch/bleve"
)

type Searcher struct {
	index bleve.Index
}

func NewSearcher() (*Searcher, error) {
	productIndex, err := bleve.Open(IndexDir())
	if err != nil {
		return nil, err
	}
	return &Searcher{index: productIndex}, nil
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

func (s *Searcher) Search(q string) (*bleve.SearchResult, error) {
	query := bleve.NewQueryStringQuery(q)
	searchRequest := bleve.NewSearchRequest(query)
	searchResult, err := s.index.Search(searchRequest)
	if err != nil {
		return nil, err
	}
	fmt.Println(searchResult)
	return searchResult, err
}

func (s *Searcher) DocCount() (uint64, error) {
	return s.index.DocCount()
}
