package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

type ResponseError struct {
	Message string `json:"message"`
}

// query handler type
type QueryHandler struct {
	db    *sql.DB
	cache *Cache
}

// constructor
func NewQueryHandler(e *echo.Echo, db *sql.DB, cache *Cache) *QueryHandler {
	handler := &QueryHandler{
		db:    db,
		cache: cache,
	}
	e.POST("/query", handler.Query)
	return handler
}

// query route
func (q *QueryHandler) Query(c echo.Context) error {
	var query Query
	err := c.Bind(&query)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	if ok, err := query.isSQLValid(); !ok {
		log.Warn(fmt.Sprintf("query \"%s\" is not valid", query.Query))
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		return c.JSON(http.StatusBadRequest, errors.New("submitted query is not valid"))
	}

	cached, err := q.cache.Find(context.Background(), query.Query)
	if err != nil {
		panic("failed to load cache from the database")
	}

	if cached != nil {
		log.Debug(fmt.Sprintf("found result in cache for query \"%s\"", query.Query))
		return c.JSON(http.StatusOK, cached)
	}

	table, err := q.queryToJson(query.Query)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	log.Info(fmt.Sprintf("query \"%s\" executed", query.Query))

	status, err := q.cache.Create(context.Background(), query.Query, table)
	if err != nil {
		panic("failed to write result into the database")
	}
	if status {
		log.Debug(fmt.Sprintf("successfully written found result into the cache for query \"%s\"", query.Query))
	} else {
		log.Debug(fmt.Sprintf("failed to write found result into the cache for query \"%s\"", query.Query))
	}

	return c.JSON(http.StatusOK, table)
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
