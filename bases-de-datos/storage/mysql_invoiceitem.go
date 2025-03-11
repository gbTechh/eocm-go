package storage

import (
	"database/sql"
	"fmt"
)

const (
	mysqlMigrateInvoiceItem = `CREATE TABLE IF NOT EXISTS invoiceItem(
		id int auto_increment primary key,
		name varchar(25) not null,
		observations Varchar(100)
		price int not null,
		created_at Timestamp not null default now(),
		updated_at Timestamp,
	)`
)

type MysqlInvoiceItem struct {
	db *sql.DB
}

func NewMysqlInvoiceItem(db *sql.DB) *MysqlInvoiceItem  {
	return &MysqlInvoiceItem{db}
}

func (p *MysqlInvoiceItem) Migrate() error  {
	stmt, err := p.db.Prepare(mysqlMigrateInvoiceItem)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		return err
	}
	fmt.Println("migracion de InvoiceItem ejectuada correctamente")
	return nil
}