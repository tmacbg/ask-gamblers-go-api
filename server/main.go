package main

import (
	"askGamblersApi/platform/data"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type ResponseObject struct {
	Websites []data.Item `json:"websites"`
	Count    int         `json:"count"`
	Country  string      `json:"country"`
}

func main() {

	db, _ := sql.Open("sqlite3", "./database/askGamblers.db")

	connection := data.NewConnection(db)

	http.HandleFunc("/api/data", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		items := connection.Get(r.PostForm["country"][0])

		responseData := ResponseObject{
			Websites: items,
			Count:    len(items),
			Country:  r.PostForm["country"][0],
		}

		response, _ := json.Marshal(responseData)

		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	})

	fmt.Println("Ask Gamblers Interface started !")
	fmt.Printf("API Started on port 5000 !! \n")
	http.ListenAndServe(":5001", nil)
}
