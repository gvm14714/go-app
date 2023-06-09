package main

import (
	"database/sql"
	"log"
	"time"
	mysqldriver "github.com/go-sql-driver/mysql"
)

var _connection *sql.DB

func init() {
	mysqlConfig := &mysqldriver.Config{
		User:                 "root",
		Passwd:               "1234",
		Net:                  "tcp",
		Addr:                 "mysql:3306",
		AllowNativePasswords: true,
		ParseTime:            true,
	}

	for i := 0; i < 30; i++ {
		connection, err := sql.Open("mysql", mysqlConfig.FormatDSN())
		if err != nil {
			log.Println("Failed to connect to MySQL. Retrying in 1 second...")
			time.Sleep(1 * time.Second)
			continue
		}

		if err = connection.Ping(); err != nil {
			log.Println("Failed to ping MySQL. Retrying in 1 second...")
			connection.Close()
			time.Sleep(1 * time.Second)
			continue
		}

		_connection = connection
		log.Println("Connected to MySQL")
		break
	}

	if _connection == nil {
		log.Fatal("Failed to connect to MySQL")
	}

	connection, err := sql.Open("mysql", mysqlConfig.FormatDSN())
	if err != nil {
		panic(err)
	}

	_, err = connection.Exec(`CREATE DATABASE IF NOT EXISTS internship`)
	if err != nil {
		panic(err)
	}

	mysqlConfig.DBName = `internship`

	if err = connection.Close(); err != nil {
		panic(err)
	}

	connection, err = sql.Open("mysql", mysqlConfig.FormatDSN())
	if err != nil {
		panic(err)
	}

	_connection = connection

	schema()
}

func schema() {
	_, err := _connection.Exec(`CREATE TABLE IF NOT EXISTS stuff (
		id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
		created_at DATETIME NOT NULL
	)`)

	if err != nil {
		panic(err)
	}
}
