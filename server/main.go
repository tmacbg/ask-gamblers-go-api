package main

import (
	"askGamblersApi/platform/data"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"text/template"

	_ "github.com/mattn/go-sqlite3"
)

type ApiResponseObject struct {
	CountryCodes []data.Item `json:"blocked_countries"`
	Status       int32       `json:"status"`
	Message      string      `json:"message"`
}

type ResponseObject struct {
	Websites              []data.Item `json:"websites"`
	Count                 int         `json:"count"`
	Country               string      `json:"country"`
	PageTitle             string      `json:"page_title"`
	Countries             []data.Item `json:"countries"`
	Query                 string      `json:"query"`
	AvailableCasinos      []data.Item `json:"available_casinos"`
	CountAvailableCasinos int         `json:"count_available_casinos"`
}

type HomeData struct {
	Countries []data.Item `json:"countries"`
	PageTitle string      `json:"page_title"`
}

func getAvailableCasinos(allCasinos []data.Item, filteredCasinos []data.Item) []data.Item {

	for i := 0; i < len(allCasinos); {
		exist := false
		for _, b := range filteredCasinos {
			if b.Name == allCasinos[i].Name {
				exist = true
				break
			}
		}
		if exist {
			allCasinos = append(allCasinos[:i], allCasinos[i+1:]...)
		} else {
			i++
		}
	}
	return allCasinos
}

func main() {

	db, _ := sql.Open("sqlite3", "./database/askGamblers.db")

	connection := data.NewConnection(db)
	resultTemplate := template.Must(template.ParseFiles("./templates/result.html", "./templates/partials/header.html"))
	homeTemplate := template.Must(template.ParseFiles("./templates/home.html", "./templates/partials/header.html"))
	searchTemplate := template.Must(template.ParseFiles("./templates/search.html", "./templates/partials/header.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		responseData := HomeData{
			Countries: nil,
			PageTitle: "Casino tool",
		}

		homeTemplate.Execute(w, responseData)
	})

	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		countries := connection.GetCountries()

		responseData := HomeData{
			Countries: countries,
			PageTitle: "Casino Country",
		}
		searchTemplate.Execute(w, responseData)
	})

	http.HandleFunc("/api/data", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		items := connection.Get(r.PostForm["country"][0])
		countries := connection.GetCountries()
		allCasinos := connection.GetCasinos()

		availableCasinos := getAvailableCasinos(allCasinos, items)
		responseData := ResponseObject{
			Websites:              items,
			Count:                 len(items),
			Country:               r.PostForm["country"][0],
			PageTitle:             "Results",
			Countries:             countries,
			Query:                 r.PostForm["country"][0],
			AvailableCasinos:      availableCasinos,
			CountAvailableCasinos: len(availableCasinos),
		}

		resultTemplate.Execute(w, responseData)
	})

	http.HandleFunc("/api/getBlocked", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		casino := r.PostForm["casino"][0]
		items := connection.GetBlockedCountries(casino)
		jsonResponseObject := ApiResponseObject{
			CountryCodes: items,
			Status:       200,
			Message:      "Success",
		}
		jsonResponse, _ := json.Marshal(jsonResponseObject)

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)

	})

	fmt.Println("Ask Gamblers Interface started !")
	fmt.Printf("API Started on port 5001 !! \n")
	http.ListenAndServe(":5001", nil)
}
