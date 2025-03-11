package database

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
    db   *sql.DB
    once sync.Once
)

type Config struct {
    Host     string
    Port     string
    User     string
    Password string
    DBName   string
}

func NewConnection(config Config) (*sql.DB, error) {
    var err error
    once.Do(func() {
        dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
            config.User,
            config.Password,
            config.Host,
            config.Port,
            config.DBName,
        )
        
        db, err = sql.Open("mysql", dsn)
        if err != nil {
            return
        }

        // Configurar el pool de conexiones
        db.SetMaxOpenConns(25)
        db.SetMaxIdleConns(5)
        db.SetConnMaxLifetime(5 * time.Minute)

        // Verificar conexión
        err = db.Ping()
    })

    if err != nil {
        return nil, fmt.Errorf("error connecting to database: %w", err)
    }

    return db, nil
}