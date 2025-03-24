package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/glebarez/sqlite"
	_ "github.com/lib/pq"
	post "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupSQLite(dbConnection string) (*gorm.DB, error) {
	if dbConnection == "" {
		dbConnection = "/docker-compose/sqlite/sqlite/assets"
	}

	// Create the sqlite file if it's not available
	if _, err := os.Stat(dbConnection); err != nil {
		if _, err = os.Create(dbConnection); err != nil {
			return nil, err
		}
	}

	db, err := gorm.Open(sqlite.Open(dbConnection), &gorm.Config{})
	return db, err
}

func setupPostgre(dbConnection string) (*gorm.DB, error) {
	if dbConnection == "" {
		dbConnection = "host=localhost user=user password=root dbname=assets_db port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	}
	sqlDb, err := sql.Open("postgres", dbConnection)
	if err != nil {
		log.Panicln("failed to connect database", err)
		return nil, err
	}

	sqlDb.SetConnMaxIdleTime(30)
	sqlDb.SetMaxOpenConns(50)
	sqlDb.SetConnMaxLifetime(2 * time.Minute)

	log.Println("pool database connection is created")

	ormDb, err := gorm.Open(post.New(post.Config{
		Conn: sqlDb,
	}), &gorm.Config{})

	if err != nil {
		log.Println("error on creating gorm connection")
		return nil, err
	}
	return ormDb, nil
}

func InitDb(env *EnviConfig) (*gorm.DB, error) {
	log.Println("create pool database connection")

	var db *gorm.DB
	var err error

	switch env.DbDriver {
	case "sqlite":
		db, err = setupSQLite(env.DbConnection)
		break
	case "postgre":
		db, err = setupPostgre(env.DbConnection)
	default:
		return nil, fmt.Errorf("No database found, set the DB env")
	}

	log.Println("gorm connection is created")

	return db, err
}
