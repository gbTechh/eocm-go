package storage

import (
	"database/sql"
	"fmt"
)

const (
	mysqlMigrateInvoiceHeader = `CREATE TABLE IF NOT EXISTS invoiceHeader(
		id int auto_increment primary key,
		name varchar(25) not null,
		observations Varchar(100)
		price int not null,
		created_at Timestamp not null default now(),
		updated_at Timestamp,
	)`
)

type MysqlInvoiceHeader struct {
	db *sql.DB
}

func NewMysqlInvoiceHeader(db *sql.DB) *MysqlInvoiceHeader  {
	return &MysqlInvoiceHeader{db}
}

func (p *MysqlInvoiceHeader) Migrate() error  {
	stmt, err := p.db.Prepare(mysqlMigrateInvoiceHeader)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		return err
	}
	fmt.Println("migracion de InvoiceHeadero ejectuada correctamente")
	return nil
}