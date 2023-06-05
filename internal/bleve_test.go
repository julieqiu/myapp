package internal

import (
	"testing"
)

func TestSearch(t *testing.T) {
	searcher, err := NewSearcher(nil)
	if err != nil {
		t.Fatal(err)
	}
	_, err = searcher.Search("salmon")
	if err != nil {
		t.Fatal(err)
	}
	_, err = searcher.Search("30")
	if err != nil {
		t.Fatal(err)
	}
}
