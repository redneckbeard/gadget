package query

import (
	"reflect"
	"strings"

	"github.com/jmoiron/sqlx"
)

var (
	DB     *sqlx.DB
	tables = make(map[string]*qwery)
	queue  = make(chan interface{})
)

// Connect sets the query package's *sqlx.DB. Fallaciously, it does not actually establish the connection.
func Connect(db *sqlx.DB) {
	DB = db
	go func() {
		for m := range queue {
			addTable(m)
		}
	}()
}

// Register adds models (aka structs) to the reflection cache. Registration
// does not actually take place until Connect has been called.
func Register(model interface{}) {
	go func() {
		queue <- model
	}()
}

func addTable(model interface{}) {
	var (
		tableName string
		columns   []string
	)
	t := reflect.TypeOf(model)
	for k, _ := range DB.Mapper.TypeMap(t) {
		if k != "id" {
			columns = append(columns, k)
		}
	}
	tableName = strings.ToLower(t.Name()) + "s"
	tables[tableName] = &qwery{tableName: tableName, columns: columns}
}
