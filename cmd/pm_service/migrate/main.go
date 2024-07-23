package main

import (
	"log"
	"os"

	"github.com/4lerman/pm_service/internal/config"
	"github.com/4lerman/pm_service/pkg/db"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	db, err := db.NewPSQLStorage(&db.DbConfig{
		Host:     config.Envs.DBAddress,
		User:     config.Envs.DBUser,
		Port:     config.Envs.DBPort,
		Dbname:   config.Envs.DBName,
		Password: config.Envs.DBPassword,
	})

	if err != nil {
		log.Fatal("Db init error", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://cmd/pm_service/migrate/migrations",
		"postgres",
		driver,
	)

	if err != nil {
		log.Fatal("Migration err: ", err)
	}

	cmd := os.Args[len(os.Args)-1]
	if cmd == "up" {
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	}

	if cmd == "down" {
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	}
}
