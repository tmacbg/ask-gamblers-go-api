package data

import (
	"database/sql"
)

type Feed struct {
	DB *sql.DB
}

type Item struct {
	Name string `json:"name"`
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

func NewConnection(db *sql.DB) *Feed {
	stm, _ := db.Prepare("CREATE TABLE IF NOT EXISTS countries ( id integer PRIMARY KEY, name text);")

	stm.Exec()

	return &Feed{
		DB: db,
	}
}
