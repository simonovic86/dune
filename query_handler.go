package main

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/labstack/echo"
	"net/http"
	"reflect"
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
	// an array of JSON objects
	// the map key is the field name
	var objects []map[string]interface{}

	rows, err := q.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		// figure out what columns were returned
		// the column names will be the JSON object field keys
		columns, err := rows.ColumnTypes()
		if err != nil {
			return nil, err
		}

		// Scan needs an array of pointers to the values it is setting
		// This creates the object and sets the values correctly
		values := make([]interface{}, len(columns))
		object := map[string]interface{}{}
		for i, column := range columns {
			object[column.Name()] = reflect.New(column.ScanType()).Interface()
			values[i] = object[column.Name()]
		}

		err = rows.Scan(values...)
		if err != nil {
			return nil, err
		}

		objects = append(objects, object)
	}
	return objects, nil
}