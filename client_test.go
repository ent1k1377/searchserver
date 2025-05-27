package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
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

func TestClientDo(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := SearchClient{
		AccessToken: "",
		URL:         server.URL,
	}
	_, err := client.FindUsers(SearchRequest{})
	expectedError := "timeout for limit=1&offset=0&order_by=0&order_field=&query="
	if expectedError != err.Error() {
		t.Errorf("Expected %v, got %v", expectedError, err)
	}

	client.URL = "fakeurl"
	_, err = client.FindUsers(SearchRequest{})

	expectedError = "unknown error Get \"fakeurl?limit=1&offset=0&order_by=0&order_field=&query=\": unsupported protocol scheme \"\""
	if expectedError != err.Error() {
		t.Errorf("Expected %v, got %v", expectedError, err)
	}
}

func TestStatusUnauthorized(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer server.Close()

	client := SearchClient{
		AccessToken: "bad access token",
		URL:         server.URL,
	}

	expectedError := "Bad AccessToken"
	_, err := client.FindUsers(SearchRequest{})

	if err == nil || expectedError != err.Error() {
		t.Errorf("Expected %v, got %v", expectedError, err)
	}
}

func TestStatusInternalServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer server.Close()

	client := SearchClient{
		AccessToken: "",
		URL:         server.URL,
	}

	oldFilepath := filepath
	filepath = "fake.file"
	defer func() {
		filepath = oldFilepath
	}()

	expectedError := "SearchServer fatal error"

	_, err := client.FindUsers(SearchRequest{})
	if err == nil || expectedError != err.Error() {
		t.Errorf("Expected %v, got %v", expectedError, err)
	}
}

func TestStatusBadRequestCantUnpackJson(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, `error: bad json`, http.StatusBadRequest)
	}))
	defer server.Close()

	client := SearchClient{
		AccessToken: "",
		URL:         server.URL,
	}

	expectedError := "cant unpack error json: invalid character 'e' looking for beginning of value"
	_, err := client.FindUsers(SearchRequest{})
	if expectedError != err.Error() {
		t.Errorf("Expected %v, got %v", expectedError, err)
	}
}

func TestStatusBadRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer server.Close()

	client := SearchClient{
		AccessToken: "",
		URL:         server.URL,
	}
	testcases := []struct {
		name          string
		searchRequest SearchRequest
		expectedError string
	}{
		{
			"ErrorBadOrderField1",
			SearchRequest{
				OrderField: "unknown order field",
			},
			"OrderFeld unknown order field invalid",
		},
		{
			"Unknown bad request error",
			SearchRequest{
				OrderBy: -2,
			},
			"unknown bad request error: error parsing order_by: an unacceptable field -2 for order_by",
		},
	}

	for _, item := range testcases {
		t.Run(item.name, func(t *testing.T) {
			_, err := client.FindUsers(item.searchRequest)
			if err == nil || item.expectedError != err.Error() {
				t.Errorf("Expected %v, got %v", item.expectedError, err)
			}
		})
	}
}

func TestCantUnpackJson(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, `error: bad json`, http.StatusOK)
	}))
	defer server.Close()

	client := SearchClient{
		AccessToken: "",
		URL:         server.URL,
	}

	expectedError := "cant unpack result json: invalid character 'e' looking for beginning of value"
	_, err := client.FindUsers(SearchRequest{})
	if expectedError != err.Error() {
		t.Errorf("Expected %v, got %v", expectedError, err)
	}
}

func TestAllGood(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer server.Close()

	client := SearchClient{
		AccessToken: "",
		URL:         server.URL,
	}
	testcases := []struct {
		name             string
		searchRequest    SearchRequest
		expectedUsers    string
		expectedNextPage bool
	}{
		{
			"1 record",
			SearchRequest{
				Limit: 1,
				Query: "amet cillum",
			},
			`[{"Id":3,"Name":"EverettDillard","Age":27,"About":"Sint eu id sint irure officia amet cillum. Amet consectetur enim mollit culpa laborum ipsum adipisicing est laboris. Adipisicing fugiat esse dolore aliquip quis laborum aliquip dolore. Pariatur do elit eu nostrud occaecat.\n","Gender":""}]`,
			true,
		},
		{
			"2 Record",
			SearchRequest{
				Limit: 2,
				Query: "amet cillum",
			},
			`[{"Id":3,"Name":"EverettDillard","Age":27,"About":"Sint eu id sint irure officia amet cillum. Amet consectetur enim mollit culpa laborum ipsum adipisicing est laboris. Adipisicing fugiat esse dolore aliquip quis laborum aliquip dolore. Pariatur do elit eu nostrud occaecat.\n","Gender":""},{"Id":31,"Name":"PalmerScott","Age":37,"About":"Elit fugiat commodo laborum quis eu consequat. In velit magna sit fugiat non proident ipsum tempor eu. Consectetur exercitation labore eiusmod occaecat adipisicing irure consequat fugiat ullamco aliquip nostrud anim irure enim. Duis do amet cillum eiusmod eu sunt. Minim minim sunt sit sit enim velit sint tempor enim sint aliquip voluptate reprehenderit officia. Voluptate magna sit consequat adipisicing ut eu qui.\n","Gender":""}]`,
			false,
		},
		{
			"25 Record",
			SearchRequest{
				Limit: 26,
				Query: "amet cillum",
			},
			`[{"Id":3,"Name":"EverettDillard","Age":27,"About":"Sint eu id sint irure officia amet cillum. Amet consectetur enim mollit culpa laborum ipsum adipisicing est laboris. Adipisicing fugiat esse dolore aliquip quis laborum aliquip dolore. Pariatur do elit eu nostrud occaecat.\n","Gender":""},{"Id":31,"Name":"PalmerScott","Age":37,"About":"Elit fugiat commodo laborum quis eu consequat. In velit magna sit fugiat non proident ipsum tempor eu. Consectetur exercitation labore eiusmod occaecat adipisicing irure consequat fugiat ullamco aliquip nostrud anim irure enim. Duis do amet cillum eiusmod eu sunt. Minim minim sunt sit sit enim velit sint tempor enim sint aliquip voluptate reprehenderit officia. Voluptate magna sit consequat adipisicing ut eu qui.\n","Gender":""}]`,
			false,
		},
	}

	for _, item := range testcases {
		t.Run(item.name, func(t *testing.T) {
			expectedResponse, _ := client.FindUsers(item.searchRequest)
			if item.expectedUsers != expectedResponse.getUsers() {
				t.Errorf("Expected %v, got\n %v", item.expectedUsers, expectedResponse.getUsers())
			}
			if item.expectedNextPage != expectedResponse.getNextPage() {
				t.Errorf("Expected %v, got %v", item.expectedNextPage, expectedResponse.getNextPage())
			}
		})
	}
}

func (sr *SearchResponse) getUsers() string {
	res, _ := json.Marshal(sr.Users)

	return string(res)
}

func (sr *SearchResponse) getNextPage() bool {
	return sr.NextPage
}
