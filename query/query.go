package query

import (
	"fmt"
	"strings"
)

// Insert returns a templated INSERT query for the given table name based on
// the field names in the corresponding struct using the SQLx
// Named(Query|ExecStmt) bindvar syntax.
func Insert(table string) string { return tables[table].insert() }

// Insert returns a templated UPDATE query for the given table name based on
// the field names in the corresponding struct using the SQLx
// Named(Query|ExecStmt) bindvar syntax.
func Update(table string) string { return tables[table].update() }

// Insert returns a templated DELETE query for the given table name based on
// the field names in the corresponding struct using the SQLx
// Named(Query|ExecStmt) bindvar syntax.
func Delete(table string) string { return tables[table].delete() }

// Insert returns a templated SELECT ... LIMIT 1 query for the given table name
// based on the field names in the corresponding struct using the SQLx
// Named(Query|ExecStmt) bindvar syntax.
func Get(table string) string { return tables[table].get() }

// Insert returns a templated SELECT query for the given table name based on
// the field names in the corresponding struct using the SQLx
// Named(Query|ExecStmt) bindvar syntax.
func Select(table string) string { return tables[table]._select() }

// Columns returns a slice of strings listing the names of columns expected
// given the struct registered with name table.
func Columns(table string) []string { return tables[table].columns }

type qwery struct {
	tableName string
	columns   []string
}

func (q *qwery) insert() string {
	template := "INSERT INTO %s (%s) VALUES (%s) RETURNING *"
	columns := strings.Join(q.columns, ", ")
	values := strings.Join(stringMap(q.columns, func(s string) string {
		return ":" + s
	}), ", ")
	return fmt.Sprintf(template, q.tableName, columns, values)
}

func (q *qwery) update() string {
	template := "UPDATE %s SET %s WHERE id=:id"
	set := strings.Join(stringMap(q.columns, func(s string) string {
		return fmt.Sprintf(":%s=%s", s, s)
	}), ", ")
	return fmt.Sprintf(template, q.tableName, set)
}

func (q *qwery) delete() string {
	template := "DELETE FROM %s WHERE id=:id"
	return fmt.Sprintf(template, q.tableName)
}

func (q *qwery) get() string {
	return fmt.Sprintf("SELECT * FROM %s %%s LIMIT 1", q.tableName)
}

func (q *qwery) _select() string {
	return fmt.Sprintf("SELECT * FROM %s %%s", q.tableName)
}

type mapper func(string) string

func stringMap(in []string, m mapper) (out []string) {
	for _, s := range in {
		out = append(out, m(s))
	}
	return out
}
