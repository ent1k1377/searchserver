package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
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

	file, err := os.ReadFile("dataset.xml")
	if err != nil {
		log.Printf("error reading xml file: %s", err)
		http.Error(w, "error1234", http.StatusInternalServerError)
		return
	}

	var users XMLUsers
	err = xml.Unmarshal(file, &users)
	if err != nil {
		log.Printf("error xml unmarshal: %s", err)
		return
	}

	values := r.URL.Query()
	query := values.Get("query")
	clientUsers := f(&users, query)

	orderField, err := parseOrderField(values.Get("order_field"))
	if err != nil {
		log.Printf("%s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	orderBy, err := parseOrderBy(values.Get("order_by"))
	if err != nil {
		log.Printf("%s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sortUsers(clientUsers, orderField, orderBy)

	offset, err := parseOffset(values.Get("offset"))
	if err != nil {
		log.Printf("%s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	limit, err := parseLimit(values.Get("limit"))
	if err != nil {
		log.Printf("%s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(clientUsers[offset : offset+limit])
	if err != nil {
		log.Println("error json marshal: %w")
	}

	w.Write(response)

	// if err := json.NewEncoder(w).Encode(searchResponse.Users); err != nil {
	// 	http.Error(w, fmt.Sprintf("error encoding json: %v", err), http.StatusInternalServerError)
	// }
}

func f(users *XMLUsers, query string) []User {
	response := make([]User, 0)

	for _, user := range users.Rows {
		if strings.Contains(user.FirstName, query) ||
			strings.Contains(user.LastName, query) ||
			strings.Contains(user.About, query) {
			response = append(response, User{
				user.Id,
				user.FirstName + user.LastName,
				user.Age,
				user.About,
				"",
			})
		}
	}

	return response
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
	case "id":
		sort.Slice(users, func(i, j int) bool {
			if orderBy == OrderByAsc {
				return users[i].Id < users[j].Id
			} else if orderBy == OrderByDesc {
				return users[i].Id > users[j].Id
			}

			return false
		})
	case "name":
		sort.Slice(users, func(i, j int) bool {
			if orderBy == OrderByAsc {
				return users[i].Name < users[j].Name
			} else if orderBy == OrderByDesc {
				return users[i].Name > users[j].Name
			}

			return false
		})
	case "age":
		sort.Slice(users, func(i, j int) bool {
			if orderBy == OrderByAsc {
				return users[i].Age < users[j].Age
			} else if orderBy == OrderByDesc {
				return users[i].Age > users[j].Age
			}

			return false
		})
	}
}
