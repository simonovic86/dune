package main

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/labstack/echo"
	"net/http"
)

type ResponseError struct {
	Message string `json:"message"`
}

type QueryHandler struct {
	db *sql.DB
}

func NewQueryHandler(e *echo.Echo, db *sql.DB) *QueryHandler {
	handler := &QueryHandler{
		db:db,
	}
	e.POST("/query", handler.Query)
	return handler
}

func (q *QueryHandler) Query(c echo.Context) error {
	var query Query
	err := c.Bind(&query)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	if ok, err := query.isSQLValid(); !ok {
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		return c.JSON(http.StatusBadRequest, errors.New("submitted query is not valid"))
	}

	table, err := q.queryToJson(query.Query)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	fmt.Println(table)
	return c.JSON(http.StatusCreated, table)
}

// convert result to JSON table
func (q *QueryHandler) queryToJson(query string, args ...interface{}) ([]map[string]interface{}, error) {
	rows, err := q.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	count := len(columns)
	tableData := make([]map[string]interface{}, 0)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	for rows.Next() {
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}
		err = rows.Scan(valuePtrs...)
		if err != nil {
			return nil, err
		}
		entry := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			entry[col] = v
		}
		tableData = append(tableData, entry)
	}
	return tableData, nil
}