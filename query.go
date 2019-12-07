package main

import (
	"errors"
	"github.com/xwb1989/sqlparser"
	"strings"
)

// query struct
type Query struct {
	Query string `json:"query"`
}

// validate query string
func (q *Query) isSQLValid() (bool, error) {
	_, err := sqlparser.Parse(q.Query)
	if err != nil {
		return false, err
	}
	trimmedLower := strings.ToLower(strings.TrimSpace(q.Query))
	if !strings.HasPrefix(trimmedLower, "select") {
		return false, errors.New("only SELECT queries permitted")
	}
	return true, nil
}