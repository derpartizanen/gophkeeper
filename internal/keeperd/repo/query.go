package repo

import (
	"fmt"
	"strings"
)

// QueryBuilder assists in creation of DB query with variadic parameters,
// e.g. update of complex select with filters.
// This is a very low level builder to suit modest needs.
type queryBuilder struct {
	query string

	// lastDef stands for last opened command or clause in SQL query.
	lastDef string

	values []any
}

// newQueryBuilder creates new builder from the provided starting query.
func newQueryBuilder(query string) *queryBuilder {
	return &queryBuilder{query: query}
}

func (qb *queryBuilder) setLastDef(name string) *queryBuilder {
	qb.lastDef = name
	qb.query += " " + qb.lastDef

	return qb
}

// Set adds SET command to query.
func (qb *queryBuilder) Set() *queryBuilder {
	return qb.setLastDef("SET")
}

// Where adds WHERE clause to query.
func (qb *queryBuilder) Where() *queryBuilder {
	return qb.setLastDef("WHERE")
}

// And adds AND clause to query.
func (qb *queryBuilder) And() *queryBuilder {
	return qb.setLastDef("AND")
}

// Append adds new condition to query.
func (qb *queryBuilder) Append(name, operator string, value any) *queryBuilder {
	if qb.lastDef == "SET" && !strings.HasSuffix(qb.query, qb.lastDef) {
		qb.query += ","
	}

	qb.values = append(qb.values, value)
	qb.query += fmt.Sprintf(" %s %s $%d", name, operator, len(qb.values))

	return qb
}

// Query returns full query with values placeholders.
func (qb *queryBuilder) Query() string {
	return qb.query
}

// Values returns stored values.
func (qb *queryBuilder) Values() []any {
	return qb.values
}
