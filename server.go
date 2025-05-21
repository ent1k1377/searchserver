package main

import (
	"encoding/json"
	"encoding/xml"
	"log"
	"net/http"
	"os"
	"strings"
)

type UserQ struct {
	Id        int    `xml:"id"`
	Age       int    `xml:"age"`
	FirstName string `xml:"first_name"`
	LastName  string `xml:"last_name"`
	About     string `xml:"about"`
}

type Users struct {
	Rows []UserQ `xml:"row"`
}

func SearchServer(w http.ResponseWriter, r *http.Request) {
	file, err := os.ReadFile("dataset.xml")
	if err != nil {
		log.Println("error reading xml file: %w", err)
		return
	}

	var users Users
	err = xml.Unmarshal(file, &users)
	if err != nil {
		log.Println("error xml unmarshal: %w", err)
		return
	}

	query := r.URL.Query().Get("query")

	searchResponse := SearchResponse{
		Users: make([]User, 0),
	}

	for _, row := range users.Rows {
		if strings.Contains(row.FirstName, query) || strings.Contains(row.LastName, query) || strings.Contains(row.About, query) {
			searchResponse.Users = append(searchResponse.Users, User{
				row.Id,
				row.FirstName + row.LastName,
				row.Age,
				row.About,
				"",
			})
		}
	}

	response, err := json.Marshal(searchResponse.Users)
	if err != nil {
		log.Println("error json marshal: %w")
	}

	w.Write(response)

	// w.Header().Set("Content-Type", "application/json")
	// if err := json.NewEncoder(w).Encode(searchResponse.Users); err != nil {
	// 	http.Error(w, fmt.Sprintf("error encoding json: %v", err), http.StatusInternalServerError)
	// }
}
