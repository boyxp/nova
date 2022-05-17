package main

import (
    "github.com/golang-migrate/migrate/v4"
    "github.com/golang-migrate/migrate/v4/database/mysql"
    _ "github.com/golang-migrate/migrate/v4/source/file"
)

import "os"
import "log"
import "github.com/boyxp/nova/database"
import _ "github.com/joho/godotenv/autoload"

func main() {
    database.Register("database", os.Getenv("database.dbname"), os.Getenv("database.dsn"))
    db        := database.Open("database")
    driver, _ := mysql.WithInstance(db, &mysql.Config{})
    m, err    := migrate.NewWithDatabaseInstance("file://migrations", "mysql", driver)

    if err != nil {
      log.Fatal(err)
    }

    err = m.Up()
    if err != nil {
      log.Fatal(err)
    }

    m.Steps(1)
}
