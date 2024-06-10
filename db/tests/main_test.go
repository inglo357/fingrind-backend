package db_tests

import (
	"database/sql"
	"fmt"
	db "github/inglo357/fingrind_backend/db/sqlc"
	"github/inglo357/fingrind_backend/utils"
	"log"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

var testQuery *db.Store

const testDBName = "fingrind_test"
const sslmode = "?sslmode=disable"

func TestMain(m *testing.M){
	config, err := utils.LoadConfig("../../")

	if err != nil{
		log.Fatal("Could not load configuration from env ", err)
	}

	conn, err := sql.Open(config.DB_driver, config.DB_source+sslmode)
	
	if err != nil{
		log.Fatalf("Could not connect to %s server %v", config.DB_driver, err)
	}

	_, err = conn.Exec(fmt.Sprintf("CREATE DATABASE %v;", testDBName))

	if err != nil{
		log.Fatalf("Could not create database %v", err)
	}
	
	tconn, err := sql.Open(config.DB_driver, config.DB_source+testDBName+sslmode)

	if err != nil{
		teardown(conn)
		log.Fatalf("Could not create database %v", err)
	}

	driver, err := postgres.WithInstance(tconn, &postgres.Config{})

	if err != nil{
		teardown(conn)
		log.Fatalf("Could not create migration driver %v", err)	
	}

	mig, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", "../migrations"),
		config.DB_driver, driver)	
	
	if err != nil{
		teardown(conn)
		log.Fatalf("Could not create migration instance %v", err)
	}

	if err = mig.Up(); err != nil && err != migrate.ErrNoChange{
		teardown(conn)
		log.Fatalf("Could not run migrations %v", err)
	}

	testQuery = db.NewStore(tconn)

	code := m.Run()

	tconn.Close()

	teardown(conn)

	os.Exit(code)
}

func teardown(conn *sql.DB){
	_, err := conn.Exec(fmt.Sprintf("DROP DATABASE %s WITH (FORCE);", testDBName))
	if err != nil{
		log.Fatalf("Could not drop database %v", err)
	}

	conn.Close()
}