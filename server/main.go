package main

import (
	"askGamblersApi/platform/data"
	"database/sql"
	"fmt"
	"net/http"
	"text/template"

	_ "github.com/mattn/go-sqlite3"
)

type ResponseObject struct {
	Websites  []data.Item `json:"websites"`
	Count     int         `json:"count"`
	Country   string      `json:"country"`
	PageTitle string      `json:"page_title"`
	Countries []data.Item `json:"countries"`
	Query     string      `json:"query"`
}

type HomeData struct {
	Countries []data.Item `json:"countries"`
	PageTitle string      `json:"page_title"`
}

func main() {

	db, _ := sql.Open("sqlite3", "./database/askGamblers.db")

	connection := data.NewConnection(db)
	resultTemplate := template.Must(template.ParseFiles("./templates/result.html"))
	homeTemplate := template.Must(template.ParseFiles("./templates/home.html"))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		countries := connection.GetCountries()

		responseData := HomeData{
			Countries: countries,
			PageTitle: "Casino Country",
		}
		homeTemplate.Execute(w, responseData)
	})

	http.HandleFunc("/api/data", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		items := connection.Get(r.PostForm["country"][0])
		countries := connection.GetCountries()

		responseData := ResponseObject{
			Websites:  items,
			Count:     len(items),
			Country:   r.PostForm["country"][0],
			PageTitle: "Results",
			Countries: countries,
			Query:     r.PostForm["country"][0],
		}

		//response, _ := json.Marshal(responseData)

		//w.Header().Set("Content-Type", "application/json")
		resultTemplate.Execute(w, responseData)
		//w.Write(response)
	})

	fmt.Println("Ask Gamblers Interface started !")
	fmt.Printf("API Started on port 5001 !! \n")
	http.ListenAndServe(":5001", nil)
}
