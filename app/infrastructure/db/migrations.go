package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	c "usermanager/app/config"
	"usermanager/app/domain"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
)

// migrate database method, returns error if ocurred.
func MigrateDb() error {
	// first we need to validate db, check if exist
	// on server, if not we need to create new one
	if err := validateDb(); err != nil {
		return err
	}

	db, err := OpenDb()
	if err != nil {
		return err
	}

	// simple db migrations from domain models
	return db.AutoMigrate(&domain.User{})
}

// validate does database exist method, returns error if ocurred
func validateDb() error {
	dbExist, err := dbExist()
	if err != nil {
		return err
	}
	if !dbExist {
		if err = createDb(); err != nil {
			return fmt.Errorf("create database failed: %v", err)
		}
	}

	return nil
}

// create new database method, returns error if ocurred.
func createDb() error {
	db, err := sql.Open("postgres", server())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", c.EnvConfig.DbName))
	if err != nil {
		return err
	}

	return nil
}

// Checks if the database exists method.
func dbExist() (bool, error) {
	db, err := sql.Open("postgres", server())
	if err != nil {
		return false, err
	}
	defer db.Close()

	sqlQuery := `SELECT EXISTS(SELECT 1 FROM pg_catalog.pg_database 
		WHERE lower(datname) = lower($1));`

	ctx := context.Background()
	row := db.QueryRowContext(ctx, sqlQuery, c.EnvConfig.DbName)

	var dbExists bool
	err = row.Scan(&dbExists)

	return dbExists, err
}

// DB connection string from env variables
func dbUrl() string {
	return fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=%v",
		c.EnvConfig.DbUser,
		c.EnvConfig.DbPass,
		c.EnvConfig.DbHost,
		c.EnvConfig.DbPort,
		c.EnvConfig.DbName,
		c.EnvConfig.SslMode)
}

// Connection string without db name. We need this one so
// we can check first does database exist on the server.
func server() string {
	return fmt.Sprintf("postgres://%v:%v@%v:%v/?sslmode=%v",
		c.EnvConfig.DbUser,
		c.EnvConfig.DbPass,
		c.EnvConfig.DbHost,
		c.EnvConfig.DbPort,
		c.EnvConfig.SslMode)
}
