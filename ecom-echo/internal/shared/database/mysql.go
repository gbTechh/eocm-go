// internal/shared/database/mysql.go
package database

import (
	"database/sql"
	"ecom/config"
	"fmt"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	db *sql.DB
	once sync.Once
)

func NewMySQLConnection(config *config.Config) (*sql.DB, error) {
	var err error

	once.Do(func() {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
			config.Database.User,
			config.Database.Password,
			config.Database.Host,
			config.Database.Port,
			config.Database.Name,
		)

		db, err = sql.Open("mysql", dsn)
		if err != nil {
				return
		}

		if err = db.Ping(); err != nil {
				return
		}

		db.SetMaxOpenConns(25)
		db.SetMaxIdleConns(25)
		db.SetConnMaxLifetime(5 * time.Minute)

		fmt.Println("Conectado a MySQL")
	})

	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	return db, nil
}

func Pool() *sql.DB  {
	return db;
}