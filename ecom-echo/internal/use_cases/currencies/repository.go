// repository.go
package currency

import (
	"context"
	"database/sql"
	"ecom/internal/shared/errors"
)

type Repository interface {
	// Currency operations
	CreateCurrency(ctx context.Context, c *Currency) error
	GetCurrencyByID(ctx context.Context, id int64) (*Currency, error)
	GetCurrencyByCode(ctx context.Context, code string) (*Currency, error)
	GetBaseCurrency(ctx context.Context) (*Currency, error)
	UpdateCurrency(ctx context.Context, c *Currency) error
	DeleteCurrency(ctx context.Context, id int64) error
	ListCurrencies(ctx context.Context, p *Pagination) ([]Currency, int64, error)
	
	// Exchange rate operations
	CreateExchangeRate(ctx context.Context, e *ExchangeRate) error
	GetLatestExchangeRate(ctx context.Context, fromCurrencyID, toCurrencyID int64) (*ExchangeRate, error)
	ListExchangeRates(ctx context.Context, currencyID int64) ([]ExchangeRate, error)
	
	// Special operations
	SetBaseCurrency(ctx context.Context, currencyID int64) error
	UpdateProductPricesTx(ctx context.Context, tx *sql.Tx, oldBaseCurrencyID, newBaseCurrencyID int64, rate float64) error
}

type MySQLRepository struct {
	db *sql.DB
}

func NewMySQLRepository(db *sql.DB) Repository {
	return &MySQLRepository{db: db}
}

func (r *MySQLRepository) CreateCurrency(ctx context.Context, c *Currency) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
			return err
	}

	// Si es moneda base, desactivar la base actual
	if c.IsBase {
			_, err = tx.ExecContext(ctx, 
					"UPDATE currencies SET is_base = false WHERE is_base = true AND deleted_at IS NULL")
			if err != nil {
					tx.Rollback()
					return errors.NewMysqlError(err)
			}
	}

	query := `
			INSERT INTO currencies (name, code, symbol, is_base, active, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, NOW(), NOW())
	`
	result, err := tx.ExecContext(ctx, query,
			c.Name, c.Code, c.Symbol, c.IsBase, c.Active)
	if err != nil {
			tx.Rollback()
			return errors.NewMysqlError(err)
	}

	id, err := result.LastInsertId()
	if err != nil {
			tx.Rollback()
			return errors.NewInternalError("Error al obtener el id creado", err)
	}
	c.ID = id

	return tx.Commit()
}

func (r *MySQLRepository) GetCurrencyByID(ctx context.Context, id int64) (*Currency, error) {
	query := `
			SELECT id, name, code, symbol, is_base, active, created_at, updated_at
			FROM currencies
			WHERE id = ? AND deleted_at IS NULL
	`
	c := &Currency{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
			&c.ID, &c.Name, &c.Code, &c.Symbol, &c.IsBase, &c.Active,
			&c.CreatedAt, &c.UpdatedAt,
	)
	if err == sql.ErrNoRows {
			return nil, errors.NewNotFoundError("Moneda no encontrada")
	}
	if err != nil {
			return nil, errors.NewMysqlError(err)
	}
	return c, nil
}

func (r *MySQLRepository) GetCurrencyByCode(ctx context.Context, code string) (*Currency, error) {
	query := `
			SELECT id, name, code, symbol, is_base, active, created_at, updated_at
			FROM currencies
			WHERE code = ? AND deleted_at IS NULL
	`
	c := &Currency{}
	err := r.db.QueryRowContext(ctx, query, code).Scan(
			&c.ID, &c.Name, &c.Code, &c.Symbol, &c.IsBase, &c.Active,
			&c.CreatedAt, &c.UpdatedAt,
	)
	if err == sql.ErrNoRows {
			return nil, errors.NewNotFoundError("Moneda no encontrada")
	}
	if err != nil {
			return nil, errors.NewMysqlError(err)
	}
	return c, nil
}

func (r *MySQLRepository) GetBaseCurrency(ctx context.Context) (*Currency, error) {
	query := `
			SELECT id, name, code, symbol, is_base, active, created_at, updated_at
			FROM currencies
			WHERE is_base = true AND deleted_at IS NULL
	`
	c := &Currency{}
	err := r.db.QueryRowContext(ctx, query).Scan(
			&c.ID, &c.Name, &c.Code, &c.Symbol, &c.IsBase, &c.Active,
			&c.CreatedAt, &c.UpdatedAt,
	)
	if err == sql.ErrNoRows {
			return nil, errors.NewNotFoundError("Moneda base no encontrada")
	}
	if err != nil {
			return nil, errors.NewMysqlError(err)
	}
	return c, nil
}

func (r *MySQLRepository) UpdateCurrency(ctx context.Context, c *Currency) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
			return err
	}

	// Si estamos cambiando a moneda base
	if c.IsBase {
			_, err = tx.ExecContext(ctx, 
					"UPDATE currencies SET is_base = false WHERE is_base = true AND id != ? AND deleted_at IS NULL",
					c.ID)
			if err != nil {
					tx.Rollback()
					return errors.NewMysqlError(err)
			}
	}

	query := `
			UPDATE currencies 
			SET name = ?, code = ?, symbol = ?, is_base = ?, active = ?, updated_at = NOW()
			WHERE id = ? AND deleted_at IS NULL
	`
	result, err := tx.ExecContext(ctx, query,
			c.Name, c.Code, c.Symbol, c.IsBase, c.Active, c.ID)
	if err != nil {
			tx.Rollback()
			return errors.NewMysqlError(err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
			tx.Rollback()
			return errors.NewInternalError("Error al actualizar moneda", err)
	}
	if rows == 0 {
			tx.Rollback()
			return errors.NewNotFoundError("Moneda no encontrada")
	}

	return tx.Commit()
}

func (r *MySQLRepository) ListCurrencies(ctx context.Context, p *Pagination) ([]Currency, int64, error) {
	countQuery := `
			SELECT COUNT(*) 
			FROM currencies 
			WHERE deleted_at IS NULL
	`
	args := []interface{}{}

	if p.Search != "" {
			countQuery += " AND (name LIKE ? OR code LIKE ?)"
			searchTerm := "%" + p.Search + "%"
			args = append(args, searchTerm, searchTerm)
	}

	if p.Active != nil {
			countQuery += " AND active = ?"
			args = append(args, *p.Active)
	}

	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
			return nil, 0, err
	}

	query := `
			SELECT id, name, code, symbol, is_base, active, created_at, updated_at
			FROM currencies
			WHERE deleted_at IS NULL
	`

	if p.Search != "" {
			query += " AND (name LIKE ? OR code LIKE ?)"
	}

	if p.Active != nil {
			query += " AND active = ?"
	}

	query += " ORDER BY is_base DESC, code ASC LIMIT ? OFFSET ?"
	args = append(args, p.PerPage, (p.Page-1)*p.PerPage)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
			return nil, 0, err
	}
	defer rows.Close()

	var currencies []Currency
	for rows.Next() {
			var c Currency
			err := rows.Scan(
					&c.ID, &c.Name, &c.Code, &c.Symbol, &c.IsBase, &c.Active,
					&c.CreatedAt, &c.UpdatedAt,
			)
			if err != nil {
					return nil, 0, err
			}
			currencies = append(currencies, c)
	}

	return currencies, total, nil
}

func (r *MySQLRepository) CreateExchangeRate(ctx context.Context, e *ExchangeRate) error {
	query := `
			INSERT INTO exchange_rates (from_currency_id, to_currency_id, rate, created_at)
			VALUES (?, ?, ?, NOW())
	`
	result, err := r.db.ExecContext(ctx, query,
			e.FromCurrencyID, e.ToCurrencyID, e.Rate)
	if err != nil {
			return errors.NewMysqlError(err)
	}

	id, err := result.LastInsertId()
	if err != nil {
			return errors.NewInternalError("Error al obtener el id creado", err)
	}

	e.ID = id
	return nil
}

func (r *MySQLRepository) GetLatestExchangeRate(ctx context.Context, fromCurrencyID, toCurrencyID int64) (*ExchangeRate, error) {
	query := `
			SELECT er.id, er.from_currency_id, er.to_currency_id, er.rate, er.created_at,
							fc.id, fc.name, fc.code, fc.symbol, fc.is_base, fc.active, fc.created_at, fc.updated_at,
							tc.id, tc.name, tc.code, tc.symbol, tc.is_base, tc.active, tc.created_at, tc.updated_at
			FROM exchange_rates er
			JOIN currencies fc ON fc.id = er.from_currency_id
			JOIN currencies tc ON tc.id = er.to_currency_id
			WHERE er.from_currency_id = ? AND er.to_currency_id = ?
			ORDER BY er.created_at DESC
			LIMIT 1
	`
	
	e := &ExchangeRate{
			FromCurrency: &Currency{},
			ToCurrency:   &Currency{},
	}

	err := r.db.QueryRowContext(ctx, query, fromCurrencyID, toCurrencyID).Scan(
			&e.ID, &e.FromCurrencyID, &e.ToCurrencyID, &e.Rate, &e.CreatedAt,
			&e.FromCurrency.ID, &e.FromCurrency.Name, &e.FromCurrency.Code,
			&e.FromCurrency.Symbol, &e.FromCurrency.IsBase, &e.FromCurrency.Active,
			&e.FromCurrency.CreatedAt, &e.FromCurrency.UpdatedAt,
			&e.ToCurrency.ID, &e.ToCurrency.Name, &e.ToCurrency.Code,
			&e.ToCurrency.Symbol, &e.ToCurrency.IsBase, &e.ToCurrency.Active,
			&e.ToCurrency.CreatedAt, &e.ToCurrency.UpdatedAt,
	)

	if err == sql.ErrNoRows {
			return nil, errors.NewNotFoundError("Tipo de cambio no encontrado")
	}
	if err != nil {
			return nil, errors.NewMysqlError(err)
	}

	return e, nil
}

func (r *MySQLRepository) ListExchangeRates(ctx context.Context, currencyID int64) ([]ExchangeRate, error) {
	query := `
			SELECT er.id, er.from_currency_id, er.to_currency_id, er.rate, er.created_at,
							fc.id, fc.name, fc.code, fc.symbol, fc.is_base, fc.active, fc.created_at, fc.updated_at,
							tc.id, tc.name, tc.code, tc.symbol, tc.is_base, tc.active, tc.created_at, tc.updated_at
			FROM exchange_rates er
			JOIN currencies fc ON fc.id = er.from_currency_id
			JOIN currencies tc ON tc.id = er.to_currency_id
			WHERE (er.from_currency_id = ? OR er.to_currency_id = ?)
			ORDER BY er.created_at DESC
	`
	
	rows, err := r.db.QueryContext(ctx, query, currencyID, currencyID)
	if err != nil {
			return nil, errors.NewMysqlError(err)
	}
	defer rows.Close()

	var rates []ExchangeRate
	for rows.Next() {
			e := ExchangeRate{
					FromCurrency: &Currency{},
					ToCurrency:   &Currency{},
			}
			
			err := rows.Scan(
					&e.ID, &e.FromCurrencyID, &e.ToCurrencyID, &e.Rate, &e.CreatedAt,
					&e.FromCurrency.ID, &e.FromCurrency.Name, &e.FromCurrency.Code,
					&e.FromCurrency.Symbol, &e.FromCurrency.IsBase, &e.FromCurrency.Active,
					&e.FromCurrency.CreatedAt, &e.FromCurrency.UpdatedAt,
					&e.ToCurrency.ID, &e.ToCurrency.Name, &e.ToCurrency.Code,
					&e.ToCurrency.Symbol, &e.ToCurrency.IsBase, &e.ToCurrency.Active,
					&e.ToCurrency.CreatedAt, &e.ToCurrency.UpdatedAt,
			)
			if err != nil {
					return nil, err
			}
			
			rates = append(rates, e)
	}

	return rates, nil
}

func (r *MySQLRepository) SetBaseCurrency(ctx context.Context, currencyID int64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {		
		return err
	}

	// Obtener la moneda base actual
	var oldBaseCurrencyID int64
	err = tx.QueryRowContext(ctx,
			"SELECT id FROM currencies WHERE is_base = true AND deleted_at IS NULL").
			Scan(&oldBaseCurrencyID)
	if err != nil && err != sql.ErrNoRows {
		tx.Rollback()
		return errors.NewMysqlError(err)
	}

	// Obtener el tipo de cambio entre la antigua y nueva moneda base
	if oldBaseCurrencyID != 0 {
		var rate float64
		err = tx.QueryRowContext(ctx, `
				SELECT rate FROM exchange_rates 
				WHERE from_currency_id = ? AND to_currency_id = ?
				ORDER BY created_at DESC LIMIT 1
		`, oldBaseCurrencyID, currencyID).Scan(&rate)
		if err == nil {
			if err := r.UpdateProductPricesTx(ctx, tx, oldBaseCurrencyID, currencyID, rate); err != nil {
				tx.Rollback()
				return err
			}
		}
		
	}

	// Desactivar la moneda base actual
	_, err = tx.ExecContext(ctx,
			"UPDATE currencies SET is_base = false WHERE is_base = true AND deleted_at IS NULL")
	if err != nil {		
		tx.Rollback()
		return errors.NewMysqlError(err)
	}

	// Establecer la nueva moneda base
	_, err = tx.ExecContext(ctx,
			"UPDATE currencies SET is_base = true WHERE id = ? AND deleted_at IS NULL",
			currencyID)
	if err != nil {
			tx.Rollback()
			return errors.NewMysqlError(err)
	}

	return tx.Commit()
}

func (r *MySQLRepository) UpdateProductPricesTx(ctx context.Context, tx *sql.Tx, oldBaseCurrencyID, newBaseCurrencyID int64, rate float64) error {
	// Actualizar precios de productos que usan la moneda antigua
	_, err := tx.ExecContext(ctx, `
			UPDATE products 
			SET price = price * ?,
					currency_id = ?
			WHERE currency_id = ? AND deleted_at IS NULL
	`, rate, newBaseCurrencyID, oldBaseCurrencyID)
	
	if err != nil {
			return errors.NewMysqlError(err)
	}

	// También podríamos actualizar precios en otras tablas que manejen monedas
	// Por ejemplo: orders, invoice_items, etc.

	return nil
}

func (r *MySQLRepository) DeleteCurrency(ctx context.Context, id int64) error {
	// Verificar si es moneda base
	var isBase bool
	err := r.db.QueryRowContext(ctx,
			"SELECT is_base FROM currencies WHERE id = ? AND deleted_at IS NULL",
			id).Scan(&isBase)
	
	if err == sql.ErrNoRows {
			return errors.NewNotFoundError("Moneda no encontrada")
	}
	if err != nil {
			return errors.NewMysqlError(err)
	}

	if isBase {
		return errors.NewBadRequestError("No se puede eliminar la moneda base")
	}

	// Verificar si hay productos usando esta moneda (la tabla productos no tiene le campo currency_id)
	/*var productsCount int
	err = r.db.QueryRowContext(ctx,
			"SELECT COUNT(*) FROM products WHERE currency_id = ? AND deleted_at IS NULL",
			id).Scan(&productsCount)
	
	if err != nil {
			return errors.NewMysqlError(err)
	}

	if productsCount > 0 {
			return errors.NewBadRequestError("No se puede eliminar la moneda porque está en uso")
	}*/

	// Soft delete
	query := `
			UPDATE currencies
			SET deleted_at = NOW()
			WHERE id = ? AND deleted_at IS NULL
	`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return errors.NewMysqlError(err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.NewNotFoundError("Moneda no encontrada")
	}

	return nil
}

// Función auxiliar para verificar si una moneda existe
func (r *MySQLRepository) currencyExists(ctx context.Context, tx *sql.Tx, currencyID int64) (bool, error) {
	var exists bool
	err := tx.QueryRowContext(ctx,
			"SELECT EXISTS(SELECT 1 FROM currencies WHERE id = ? AND deleted_at IS NULL)",
			currencyID).Scan(&exists)
	if err != nil {
			return false, errors.NewMysqlError(err)
	}
	return exists, nil
}