package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestQwe(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer ts.Close()

	sc := SearchClient{
		AccessToken: "",
		URL:         ts.URL,
	}

	sr := SearchRequest{
		Limit:      0,
		Offset:     0,
		Query:      "irure",
		OrderField: "",
		OrderBy:    0,
	}

	srq, err := sc.FindUsers(sr)
	fmt.Println(err)
	fmt.Println(srq)
}
