package storage

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

var (
	db *sql.DB
	once sync.Once
)

func NewPostgresDB()  {
	once.Do(func(){
		//se ejectua una sola vez
		var err error
		db, err := sql.Open("mysql", "enkit:123@tcp(localhost:3306)/go_ecom")
		if err != nil {
			log.Fatalf("can't open db: %v", err)
		}
		defer db.Close()

		if err := db.Ping(); err != nil {
			log.Fatalf("can't do ping: %v", err)
		}

		fmt.Println("conectado a postgres")
	})
}
func NewMysqlDB()  {
	once.Do(func(){
		//se ejectua una sola vez
		var err error
		db, err := sql.Open("mysql", "enkit:123@tcp(localhost:3306)/go_ecom")
		if err != nil {
			log.Fatalf("can't open db: %v", err)
		}
		defer db.Close()

		if err := db.Ping(); err != nil {
			log.Fatalf("can't do ping: %v", err)
		}

		fmt.Println("conectado a mysql")
	})
}

//retorna una unica instancia de db
func Pool() *sql.DB  {
	return db;
}