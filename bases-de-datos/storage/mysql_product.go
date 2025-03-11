package storage

import (
	"database/sql"
	"fmt"
)

const (
	mysqlMigrateProduct = `CREATE TABLE IF NOT EXISTS products(
		id int auto_increment primary key,
		name varchar(25) not null,
		observations Varchar(100)
		price int not null,
		created_at Timestamp not null default now(),
		updated_at Timestamp,
	)`
)

type MysqlProduct struct {
	db *sql.DB
}

func NewMysqlProduct(db *sql.DB) *MysqlProduct  {
	return &MysqlProduct{db}
}

func (p *MysqlProduct) Migrate() error  {
	stmt, err := p.db.Prepare(mysqlMigrateProduct)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		return err
	}
	fmt.Println("migracion de producto ejectuada correctamente")
	return nil
}