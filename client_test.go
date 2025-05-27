package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type TestCase struct {
	name     string
	request  SearchRequest
	response *SearchResponse
	err      error
}

func TestQwe(t *testing.T) {
	cases := []TestCase{
		{
			"req.Limit < 0",
			SearchRequest{
				Limit: -1,
			},
			nil,
			fmt.Errorf("limit must be > 0"),
		},
		{
			"req.Offset < 0",
			SearchRequest{
				Offset: -1,
			},
			nil,
			fmt.Errorf("offset must be > 0"),
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer ts.Close()

	sc := SearchClient{
		AccessToken: "",
		URL:         ts.URL,
	}

	for _, item := range cases {
		t.Run(item.name, func(t *testing.T) {
			response, err := sc.FindUsers(item.request)

			if item.err.Error() != err.Error() {
				t.Errorf("Expected %v, got %v", item.err, err)
				return
			}

			if item.response == nil {
				return
			}

			if item.response != response {
				t.Errorf("Expected %v, got %v", item.response, response)
			} else if item.response.NextPage != response.NextPage {
				t.Errorf("Expected %v, got %v", item.response.NextPage, response.NextPage)
			}
			// ...
		})
	}
}
