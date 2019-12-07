package main

import (
	"database/sql"
	"dune/internal/middleware"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"

	"github.com/labstack/echo"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

// initialize app
func init() {
	viper.SetConfigFile(`config.json`)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	Formatter := new(log.TextFormatter)
	Formatter.TimestampFormat = "02-01-2006 15:04:05"
	Formatter.FullTimestamp = true
	log.SetFormatter(Formatter)
	log.SetOutput(os.Stdout)

	if viper.GetBool(`debug`) {
		fmt.Println("Service RUN on DEBUG mode")
		log.SetLevel(log.DebugLevel)
	}
}

// main entry point
func main() {
	dbHost := viper.GetString(`database.host`)
	dbUser := viper.GetString(`database.user`)
	dbPass := viper.GetString(`database.pass`)
	dbName := viper.GetString(`database.name`)
	connection := fmt.Sprintf("dbname=%s user=%s password=%s host=%s sslmode=disable", dbName, dbUser, dbPass, dbHost)
	dbConn, err := sql.Open(`postgres`, connection)
	if err != nil {
		if viper.GetBool("debug") {
			log.Fatal(fmt.Sprintf("failed to connect to connect to database %s", dbName))
		}
		os.Exit(1)
	}

	err = dbConn.Ping()
	if err != nil {
		if viper.GetBool("debug") {
			log.Fatal(fmt.Sprintf("failed to ping database %s", dbName))
		}
		os.Exit(1)
	}

	defer func() {
		err := dbConn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	e := echo.New()
	middleware := middleware.InitMiddleware()
	e.Use(middleware.CORS)
	log.Fatal(e.Start(viper.GetString("server.address")))
}