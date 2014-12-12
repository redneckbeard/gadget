package query

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/jmoiron/sqlx"
)

func Query(db *sqlx.DB, model interface{}) *Q {
	var (
		tableName string
		columns   []string
	)
	t := reflect.TypeOf(model)
	for k, _ := range db.Mapper.TypeMap(t) {
		if k != "id" {
			columns = append(columns, k)
		}
	}
	tableName = strings.ToLower(t.Name()) + "s"
	return &Q{TableName: tableName, Columns: columns}
}

type Q struct {
	TableName string
	Columns   []string
}

func (q *Q) Insert() string {
	template := "INSERT INTO %s (%s) VALUES (%s) RETURNING *"
	columns := strings.Join(q.Columns, ", ")
	values := strings.Join(stringMap(q.Columns, func(s string) string {
		return ":" + s
	}), ", ")
	return fmt.Sprintf(template, q.TableName, columns, values)
}

func (q *Q) Update() string {
	template := "UPDATE %s SET %s WHERE id=:id"
	set := strings.Join(stringMap(q.Columns, func(s string) string {
		return fmt.Sprintf(":%s=%s", s, s)
	}), ", ")
	return fmt.Sprintf(template, q.TableName, set)
}

func (q *Q) Delete() string {
	template := "DELETE FROM %s WHERE id=:id"
	return fmt.Sprintf(template, q.TableName)
}

func (q *Q) Get() string {
	return fmt.Sprintf("SELECT * FROM %s %%s LIMIT 1", q.TableName)
}

func (q *Q) Select() string {
	return fmt.Sprintf("SELECT * FROM %s %%s", q.TableName)
}

type mapper func(string) string

func stringMap(in []string, m mapper) (out []string) {
	for _, s := range in {
		out = append(out, m(s))
	}
	return out
}
