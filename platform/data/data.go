package data

import (
	"database/sql"
	"fmt"
)

type Feed struct {
	DB *sql.DB
}

type Item struct {
	Name string `json:"name"`
}

func (feed *Feed) GetCountries() []Item {
	items := []Item{}

	rows, _ := feed.DB.Query("SELECT name from countries")
	var country string
	for rows.Next() {
		rows.Scan(&country)
		item := Item{
			Name: country,
		}
		items = append(items, item)
	}

	return items
}

func (feed *Feed) GetCasinos() []Item {
	allCasinos := []Item{}
	rows, _ := feed.DB.Query("Select DISTINCT(name) from websites")
	var casino string
	for rows.Next() {
		rows.Scan(&casino)
		casinoItem := Item{
			Name: casino,
		}

		allCasinos = append(allCasinos, casinoItem)
	}

	return allCasinos
}

func (feed *Feed) Get(query string) []Item {
	items := []Item{}
	stm, _ := feed.DB.Prepare(`SELECT DISTINCT(w.name) from websites as w 
								JOIN websites_countries as w_c on w_c.website_id = w.id 
								JOIN countries as c on w_c.country_id = c.id
								where c.name = ? `)
	rows, _ := stm.Query(query)
	var website string
	for rows.Next() {
		rows.Scan(&website)
		item := Item{
			Name: website,
		}
		items = append(items, item)
	}

	return items
}

func (feed *Feed) GetBlockedCountries(query string) []Item {
	items := []Item{}
	stm, _ := feed.DB.Prepare(`SELECT DISTINCT(lower(c.code) ) as code from websites as w 
								JOIN websites_countries as w_c on w_c.website_id = w.id 
								JOIN countries as c on w_c.country_id = c.id
								where  w.name like ?`)
	queryString := fmt.Sprintf("%s%%", query)
	fmt.Print(queryString)
	rows, err := stm.Query(queryString)
	if err != nil {
		fmt.Print(err)
	}
	var countryCode string
	for rows.Next() {
		rows.Scan(&countryCode)
		item := Item{
			Name: countryCode,
		}
		items = append(items, item)
	}

	return items
}

func NewConnection(db *sql.DB) *Feed {
	stm, _ := db.Prepare("CREATE TABLE IF NOT EXISTS countries ( id integer PRIMARY KEY, name text);")

	stm.Exec()

	return &Feed{
		DB: db,
	}
}
