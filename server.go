package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
)

const Id = "id"
const Age = "age"
const Name = "name"

type xmlUser struct {
	Id        int    `xml:"id"`
	Age       int    `xml:"age"`
	FirstName string `xml:"first_name"`
	LastName  string `xml:"last_name"`
	About     string `xml:"about"`
}

type XMLUsers struct {
	Rows []xmlUser `xml:"row"`
}

func SearchServer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Header.Get("AccessToken") != GetAccessToken() {
		log.Println("Access Token Not Authorized")
		http.Error(w, "access token not authorized", http.StatusUnauthorized)
		return
	}

	file, err := os.ReadFile("dataset.xml")
	if err != nil {
		log.Printf("error reading xml file: %s", err)
		http.Error(w, "error1234", http.StatusInternalServerError)
		return
	}

	var xmlUsers XMLUsers
	err = xml.Unmarshal(file, &xmlUsers)
	if err != nil {
		log.Printf("error xml unmarshal: %s", err)
		http.Error(w, "error1234", http.StatusInternalServerError)
		return
	}

	searchRequest, err := parseURLValues(r.URL.Query())
	if err != nil {
		log.Printf("error parsing url values: %s", err)
		http.Error(w, "error1234", http.StatusInternalServerError)
		return
	}

	users := transformAndFilterUsers(&xmlUsers, searchRequest.Query)

	sortUsers(users, searchRequest.OrderField, searchRequest.OrderBy)
	users = paginateUsers(users, searchRequest.Offset, searchRequest.Limit)

	if err := json.NewEncoder(w).Encode(users); err != nil {
		log.Printf("error writing response: %s", err)
		http.Error(w, fmt.Sprintf("error encoding json: %v", err), http.StatusInternalServerError)
	}
}

func GetAccessToken() string {
	return "token"
}

func transformAndFilterUsers(users *XMLUsers, query string) []User {
	response := make([]User, 0)

	for _, xmlUser := range users.Rows {
		if matchesQuery(xmlUser, query) {
			response = append(response, User{
				Id:     xmlUser.Id,
				Name:   xmlUser.FirstName + xmlUser.LastName,
				Age:    xmlUser.Age,
				About:  xmlUser.About,
				Gender: "",
			})
		}
	}

	return response
}

func matchesQuery(user xmlUser, query string) bool {
	return strings.Contains(user.FirstName, query) ||
		strings.Contains(user.LastName, query) ||
		strings.Contains(user.About, query)
}

func paginateUsers(users []User, offset int, limit int) []User {
	if offset >= len(users) {
		return []User{}
	}

	start := offset
	end := offset + limit
	if end > len(users) {
		end = len(users)
	}

	return users[start:end]
}

func parseURLValues(values url.Values) (SearchRequest, error) {
	query := values.Get("query")

	orderField, err := parseOrderField(values.Get("order_field"))
	if err != nil {
		return SearchRequest{}, fmt.Errorf("error parsing order_field: %w", err)
	}

	orderBy, err := parseOrderBy(values.Get("order_by"))
	if err != nil {
		return SearchRequest{}, fmt.Errorf("error parsing order_by: %w", err)
	}

	offset, err := parseOffset(values.Get("offset"))
	if err != nil {
		return SearchRequest{}, fmt.Errorf("error parsing offset: %w", err)
	}

	limit, err := parseLimit(values.Get("limit"))
	if err != nil {
		return SearchRequest{}, fmt.Errorf("error parsing limit: %w", err)
	}

	return SearchRequest{
		Query:      query,
		OrderField: orderField,
		OrderBy:    orderBy,
		Offset:     offset,
		Limit:      limit,
	}, nil
}

func parseOrderField(orderField string) (string, error) {
	orderField = strings.ToLower(orderField)
	res := orderField

	if orderField == "" {
		orderField = Name
	} else if !(orderField == Id || orderField == Name || orderField == Age) {
		return "", fmt.Errorf("an unacceptable field %s for order_field", orderField)
	}

	return res, nil
}

func parseOrderBy(field string) (int, error) {
	orderBy, err := strconv.Atoi(field)
	if err != nil {
		return 0, fmt.Errorf("error parse %s to integer: %w", field, err)
	}

	if !(orderBy == OrderByAsc || orderBy == OrderByAsIs || orderBy == OrderByDesc) {
		return 0, fmt.Errorf("an unacceptable field %s for order_by", field)
	}

	return orderBy, nil
}

func parseLimit(field string) (int, error) {
	limit, err := strconv.Atoi(field)
	if err != nil {
		return 0, fmt.Errorf("error parse field %s to limit: %w", field, err)
	}

	return limit, nil
}

func parseOffset(field string) (int, error) {
	offset, err := strconv.Atoi(field)
	if err != nil {
		return 0, fmt.Errorf("error parse field %s to offset: %w", field, err)
	}

	return offset, nil
}

func sortUsers(users []User, orderField string, orderBy int) {
	if orderBy == OrderByAsIs {
		return
	}

	switch orderField {
	case Id:
		sort.Slice(users, func(i, j int) bool {
			if orderBy == OrderByAsc {
				i, j = j, i
			}

			return users[i].Id > users[j].Id
		})
	case Name:
		sort.Slice(users, func(i, j int) bool {
			if orderBy == OrderByAsc {
				i, j = j, i
			}

			return users[i].Name > users[j].Name
		})
	case Age:
		sort.Slice(users, func(i, j int) bool {
			if orderBy == OrderByAsc {
				i, j = j, i
			}

			return users[i].Age > users[j].Age
		})
	}
}
