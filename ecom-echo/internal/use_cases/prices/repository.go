// repository.go
package prices

import (
	"context"
	"database/sql"
	"ecom/internal/shared/errors"
)

type Repository interface {
	// PriceList operations
	CreatePriceList(ctx context.Context, pl *PriceList) error
	GetPriceListByID(ctx context.Context, id int64) (*PriceList, error)
	UpdatePriceList(ctx context.Context, pl *PriceList) error
	DeletePriceList(ctx context.Context, id int64) error
	ListPriceLists(ctx context.Context, p *Pagination) ([]PriceList, int64, error)
	GetDefaultPriceList(ctx context.Context) (*PriceList, error)
	CountPriceLists(ctx context.Context) (int64, error)
	
	// Price operations
	CreatePrice(ctx context.Context, p *Price) error
	GetPriceByID(ctx context.Context, id int64) (*Price, error)
	UpdatePrice(ctx context.Context, p *Price) error
	DeletePrice(ctx context.Context, id int64) error
	ListPrices(ctx context.Context, priceListID int64, p *Pagination) ([]Price, int64, error)
	
	// ProductVariantPrice operations
	AssignPrice(ctx context.Context, pvp *ProductVariantPrice) error
	UnassignPrice(ctx context.Context, productVariantID, priceID int64) error
	GetActivePrice(ctx context.Context, productVariantID int64) (*Price, error)
	GetAllPrices(ctx context.Context, productVariantID int64) ([]ProductVariantPrice, error)
}

type MySQLRepository struct {
	db *sql.DB
}

func NewMySQLRepository(db *sql.DB) Repository {
	return &MySQLRepository{db: db}
}

// PriceList Methods
func (r *MySQLRepository) CreatePriceList(ctx context.Context, pl *PriceList) error {
	if pl.IsDefault {
			// Si es default, deshabilitar cualquier otro default existente
			query := `
					UPDATE price_lists 
					SET is_default = false 
					WHERE is_default = true AND deleted_at IS NULL
			`
			if _, err := r.db.ExecContext(ctx, query); err != nil {
					return errors.NewMysqlError(err)
			}
	}

	query := `
			INSERT INTO price_lists (name, description, is_default, priority, status, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, NOW(), NOW())
	`
	result, err := r.db.ExecContext(ctx, query,
			pl.Name, pl.Description, pl.IsDefault, pl.Priority, pl.Status)
	if err != nil {
			return errors.NewMysqlError(err)
	}

	id, err := result.LastInsertId()
	if err != nil {
			return errors.NewInternalError("Error al obtener el id creado", err)
	}

	pl.ID = id
	return nil
}

func (r *MySQLRepository) GetPriceListByID(ctx context.Context, id int64) (*PriceList, error) {
	query := `
			SELECT id, name, description, is_default, priority, status, created_at, updated_at
			FROM price_lists
			WHERE id = ? AND deleted_at IS NULL
	`
	pl := &PriceList{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
			&pl.ID, &pl.Name, &pl.Description, &pl.IsDefault,
			&pl.Priority, &pl.Status, &pl.CreatedAt, &pl.UpdatedAt,
	)
	if err == sql.ErrNoRows {
			return nil, errors.NewNotFoundError("Lista de precios no encontrada")
	}
	if err != nil {
			return nil, errors.NewMysqlError(err)
	}
	return pl, nil
}

func (r *MySQLRepository) UpdatePriceList(ctx context.Context, pl *PriceList) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
			return err
	}

	if pl.IsDefault {
			// Si es default, deshabilitar cualquier otro default existente
			query := `
					UPDATE price_lists 
					SET is_default = false 
					WHERE is_default = true AND id != ? AND deleted_at IS NULL
			`
			if _, err := tx.ExecContext(ctx, query, pl.ID); err != nil {
					tx.Rollback()
					return errors.NewMysqlError(err)
			}
	}

	query := `
			UPDATE price_lists 
			SET name = ?, description = ?, is_default = ?, priority = ?, status = ?, updated_at = NOW()
			WHERE id = ? AND deleted_at IS NULL
	`
	result, err := tx.ExecContext(ctx, query,
			pl.Name, pl.Description, pl.IsDefault, pl.Priority, pl.Status, pl.ID)
	if err != nil {
			tx.Rollback()
			return errors.NewMysqlError(err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
			tx.Rollback()
			return errors.NewInternalError("Error al actualizar lista de precios", err)
	}
	if rows == 0 {
			tx.Rollback()
			return errors.NewNotFoundError("Lista de precios no encontrada")
	}

	return tx.Commit()
}

func (r *MySQLRepository) DeletePriceList(ctx context.Context, id int64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
			return err
	}

	// Verificar si es la lista por defecto
	var isDefault bool
	err = tx.QueryRowContext(ctx, 
			"SELECT is_default FROM price_lists WHERE id = ? AND deleted_at IS NULL", 
			id).Scan(&isDefault)
	if err != nil {
			tx.Rollback()
			return errors.NewMysqlError(err)
	}

	if isDefault {
			tx.Rollback()
			return errors.NewBadRequestError("No se puede eliminar la lista de precios por defecto")
	}

	query := `
			UPDATE price_lists
			SET deleted_at = NOW()
			WHERE id = ? AND deleted_at IS NULL
	`
	result, err := tx.ExecContext(ctx, query, id)
	if err != nil {
			tx.Rollback()
			return errors.NewMysqlError(err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
			tx.Rollback()
			return err
	}
	if rows == 0 {
			tx.Rollback()
			return errors.NewNotFoundError("Lista de precios no encontrada")
	}

	return tx.Commit()
}

func (r *MySQLRepository) ListPriceLists(ctx context.Context, p *Pagination) ([]PriceList, int64, error) {
	countQuery := `
			SELECT COUNT(*) 
			FROM price_lists 
			WHERE deleted_at IS NULL
	`
	args := []interface{}{}

	if p.Search != "" {
			countQuery += " AND name LIKE ?"
			args = append(args, "%"+p.Search+"%")
	}

	if p.Status != nil {
			countQuery += " AND status = ?"
			args = append(args, *p.Status)
	}

	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
			return nil, 0, err
	}

	query := `
			SELECT id, name, description, is_default, priority, status, created_at, updated_at
			FROM price_lists
			WHERE deleted_at IS NULL
	`

	if p.Search != "" {
			query += " AND name LIKE ?"
	}

	if p.Status != nil {
			query += " AND status = ?"
	}

	query += " ORDER BY priority DESC, created_at DESC LIMIT ? OFFSET ?"
	args = append(args, p.PerPage, (p.Page-1)*p.PerPage)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
			return nil, 0, err
	}
	defer rows.Close()

	var priceLists []PriceList
	for rows.Next() {
			var pl PriceList
			err := rows.Scan(
					&pl.ID, &pl.Name, &pl.Description, &pl.IsDefault,
					&pl.Priority, &pl.Status, &pl.CreatedAt, &pl.UpdatedAt,
			)
			if err != nil {
					return nil, 0, err
			}
			priceLists = append(priceLists, pl)
	}

	return priceLists, total, nil
}

func (r *MySQLRepository) GetDefaultPriceList(ctx context.Context) (*PriceList, error) {
	query := `
			SELECT id, name, description, is_default, priority, status, created_at, updated_at
			FROM price_lists
			WHERE is_default = true AND deleted_at IS NULL
	`
	pl := &PriceList{}
	err := r.db.QueryRowContext(ctx, query).Scan(
			&pl.ID, &pl.Name, &pl.Description, &pl.IsDefault,
			&pl.Priority, &pl.Status, &pl.CreatedAt, &pl.UpdatedAt,
	)
	if err == sql.ErrNoRows {
			return nil, errors.NewNotFoundError("Lista de precios por defecto no encontrada")
	}
	if err != nil {
			return nil, errors.NewMysqlError(err)
	}
	return pl, nil
}

// Price Methods
func (r *MySQLRepository) CreatePrice(ctx context.Context, p *Price) error {
	query := `
			INSERT INTO prices (id_price_list, amount, starts_at, ends_at, created_at, updated_at)
			VALUES (?, ?, ?, ?, NOW(), NOW())
	`
	result, err := r.db.ExecContext(ctx, query,
			p.IDPriceList, p.Amount, p.StartsAt, p.EndsAt)
	if err != nil {
			return errors.NewMysqlError(err)
	}

	id, err := result.LastInsertId()
	if err != nil {
			return errors.NewInternalError("Error al obtener el id creado", err)
	}

	p.ID = id
	return nil
}

func (r *MySQLRepository) GetPriceByID(ctx context.Context, id int64) (*Price, error) {
	query := `
			SELECT p.id, p.id_price_list, p.amount, p.starts_at, p.ends_at, p.created_at, p.updated_at,
							pl.id, pl.name, pl.description, pl.is_default, pl.priority, pl.status,
							pl.created_at, pl.updated_at
			FROM prices p
			JOIN price_lists pl ON pl.id = p.id_price_list
			WHERE p.id = ? AND p.deleted_at IS NULL
	`
	var price Price
	price.PriceList = &PriceList{}
	
	err := r.db.QueryRowContext(ctx, query, id).Scan(
			&price.ID, &price.IDPriceList, &price.Amount, &price.StartsAt, &price.EndsAt,
			&price.CreatedAt, &price.UpdatedAt,
			&price.PriceList.ID, &price.PriceList.Name, &price.PriceList.Description,
			&price.PriceList.IsDefault, &price.PriceList.Priority, &price.PriceList.Status,
			&price.PriceList.CreatedAt, &price.PriceList.UpdatedAt,
	)
	if err == sql.ErrNoRows {
			return nil, errors.NewNotFoundError("Precio no encontrado")
	}
	if err != nil {
			return nil, errors.NewMysqlError(err)
	}
	return &price, nil
}

func (r *MySQLRepository) UpdatePrice(ctx context.Context, p *Price) error {
	query := `
			UPDATE prices 
			SET amount = ?, starts_at = ?, ends_at = ?, updated_at = NOW()
			WHERE id = ? AND deleted_at IS NULL
	`
	result, err := r.db.ExecContext(ctx, query,
			p.Amount, p.StartsAt, p.EndsAt, p.ID)
	if err != nil {
			return errors.NewMysqlError(err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
			return errors.NewInternalError("Error al actualizar precio", err)
	}
	if rows == 0 {
			return errors.NewNotFoundError("Precio no encontrado")
	}
	return nil
}

func (r *MySQLRepository) DeletePrice(ctx context.Context, id int64) error {
	query := `
			UPDATE prices
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
			return errors.NewNotFoundError("Precio no encontrado")
	}
	return nil
}
func (r *MySQLRepository) ListPrices(ctx context.Context, priceListID int64, p *Pagination) ([]Price, int64, error) {
	countQuery := `
			SELECT COUNT(*) 
			FROM prices 
			WHERE deleted_at IS NULL AND id_price_list = ?
	`
	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, priceListID).Scan(&total)
	if err != nil {
			return nil, 0, err
	}

	query := `
			SELECT p.id, p.id_price_list, p.amount, p.starts_at, p.ends_at, p.created_at, p.updated_at,
							pl.id, pl.name, pl.description, pl.is_default, pl.priority, pl.status,
							pl.created_at, pl.updated_at
			FROM prices p
			JOIN price_lists pl ON pl.id = p.id_price_list
			WHERE p.deleted_at IS NULL AND p.id_price_list = ?
			ORDER BY p.created_at DESC LIMIT ? OFFSET ?
	`
	rows, err := r.db.QueryContext(ctx, query, priceListID, p.PerPage, (p.Page-1)*p.PerPage)
	if err != nil {
			return nil, 0, err
	}
	defer rows.Close()

	var prices []Price
	for rows.Next() {
			var price Price
			price.PriceList = &PriceList{}
			
			err := rows.Scan(
					&price.ID, &price.IDPriceList, &price.Amount, &price.StartsAt, &price.EndsAt,
					&price.CreatedAt, &price.UpdatedAt,
					&price.PriceList.ID, &price.PriceList.Name, &price.PriceList.Description,
					&price.PriceList.IsDefault, &price.PriceList.Priority, &price.PriceList.Status,
					&price.PriceList.CreatedAt, &price.PriceList.UpdatedAt,
			)
			if err != nil {
					return nil, 0, err
			}
			prices = append(prices, price)
	}

	return prices, total, nil
}

func (r *MySQLRepository) AssignPrice(ctx context.Context, pvp *ProductVariantPrice) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
			return err
	}

	// Verificar si ya existe algún precio activo para esta variante en la misma lista de precios
	query := `
			SELECT COUNT(*) 
			FROM n_product_variant_prices pvp
			JOIN prices p ON p.id = pvp.id_price
			WHERE pvp.id_product_variant = ? AND p.id_price_list = (
					SELECT id_price_list FROM prices WHERE id = ?
			) AND pvp.is_active = true
	`
	var count int
	err = tx.QueryRowContext(ctx, query, pvp.IDProductVariant, pvp.IDPrice).Scan(&count)
	if err != nil {
			tx.Rollback()
			return errors.NewMysqlError(err)
	}

	// Si ya existe un precio activo, desactivarlo
	if count > 0 {
			_, err = tx.ExecContext(ctx, `
					UPDATE n_product_variant_prices pvp
					JOIN prices p ON p.id = pvp.id_price
					SET pvp.is_active = false
					WHERE pvp.id_product_variant = ? AND p.id_price_list = (
							SELECT id_price_list FROM prices WHERE id = ?
					)
			`, pvp.IDProductVariant, pvp.IDPrice)
			if err != nil {
					tx.Rollback()
					return errors.NewMysqlError(err)
			}
	}

	// Insertar el nuevo precio
	query = `
			INSERT INTO n_product_variant_prices 
			(id_product_variant, id_price, is_active, created_at, updated_at)
			VALUES (?, ?, ?, NOW(), NOW())
	`
	result, err := tx.ExecContext(ctx, query, pvp.IDProductVariant, pvp.IDPrice, pvp.IsActive)
	if err != nil {
			tx.Rollback()
			return errors.NewMysqlError(err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		return errors.NewInternalError("Error al verificar la inserción del precio", err)
	}
	if rows == 0 {
		tx.Rollback()
		return errors.NewInternalError("No se pudo asignar el precio", nil)
	}
	return tx.Commit()
}

func (r *MySQLRepository) UnassignPrice(ctx context.Context, productVariantID, priceID int64) error {
	// En lugar de eliminar, actualizamos is_active a false
	query := `
			DELETE FROM n_product_variant_prices
			WHERE id_product_variant = ? AND id_price = ?
	`
	result, err := r.db.ExecContext(ctx, query, productVariantID, priceID)
	if err != nil {
			return errors.NewMysqlError(err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
			return err
	}
	if rows == 0 {
			return errors.NewNotFoundError("Asignación de precio no encontrada")
	}

	return nil
}

func (r *MySQLRepository) GetActivePrice(ctx context.Context, productVariantID int64) (*Price, error) {
	// Obtener el precio activo con mayor prioridad y más reciente
	query := `
			SELECT p.id, p.id_price_list, p.amount, p.starts_at, p.ends_at, p.created_at, p.updated_at,
							pl.id, pl.name, pl.description, pl.is_default, pl.priority, pl.status,
							pl.created_at, pl.updated_at
			FROM prices p
			JOIN price_lists pl ON pl.id = p.id_price_list
			JOIN n_product_variant_prices pvp ON pvp.id_price = p.id
			WHERE pvp.id_product_variant = ? 
			AND pvp.is_active = true
			AND p.deleted_at IS NULL
			AND pl.deleted_at IS NULL
			AND p.starts_at <= NOW()
			AND (p.ends_at IS NULL OR p.ends_at > NOW())
			ORDER BY pl.priority DESC, p.created_at DESC
			LIMIT 1
	`
	
	price := &Price{PriceList: &PriceList{}}
	err := r.db.QueryRowContext(ctx, query, productVariantID).Scan(
			&price.ID, &price.IDPriceList, &price.Amount, &price.StartsAt, &price.EndsAt,
			&price.CreatedAt, &price.UpdatedAt,
			&price.PriceList.ID, &price.PriceList.Name, &price.PriceList.Description,
			&price.PriceList.IsDefault, &price.PriceList.Priority, &price.PriceList.Status,
			&price.PriceList.CreatedAt, &price.PriceList.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
			return nil, errors.NewNotFoundError("No hay precio activo para esta variante")
	}
	if err != nil {
			return nil, errors.NewMysqlError(err)
	}

	return price, nil
}

func (r *MySQLRepository) GetAllPrices(ctx context.Context, productVariantID int64) ([]ProductVariantPrice, error) {
	query := `
			SELECT pvp.id_product_variant, pvp.id_price, pvp.is_active, 
							pvp.created_at, pvp.updated_at,
							p.id, p.id_price_list, p.amount, p.starts_at, p.ends_at, 
							p.created_at, p.updated_at,
							pl.id, pl.name, pl.description, pl.is_default, pl.priority, 
							pl.status, pl.created_at, pl.updated_at
			FROM n_product_variant_prices pvp
			JOIN prices p ON p.id = pvp.id_price
			JOIN price_lists pl ON pl.id = p.id_price_list
			WHERE pvp.id_product_variant = ? 
			AND p.deleted_at IS NULL
			AND pl.deleted_at IS NULL
			ORDER BY pl.priority DESC, p.created_at DESC
	`
	
	rows, err := r.db.QueryContext(ctx, query, productVariantID)
	if err != nil {
			return nil, errors.NewMysqlError(err)
	}
	defer rows.Close()

	var prices []ProductVariantPrice
	for rows.Next() {
			var pvp ProductVariantPrice
			pvp.Price = &Price{PriceList: &PriceList{}}
			
			err := rows.Scan(
					&pvp.IDProductVariant, &pvp.IDPrice, &pvp.IsActive,
					&pvp.CreatedAt, &pvp.UpdatedAt,
					&pvp.Price.ID, &pvp.Price.IDPriceList, &pvp.Price.Amount,
					&pvp.Price.StartsAt, &pvp.Price.EndsAt,
					&pvp.Price.CreatedAt, &pvp.Price.UpdatedAt,
					&pvp.Price.PriceList.ID, &pvp.Price.PriceList.Name,
					&pvp.Price.PriceList.Description, &pvp.Price.PriceList.IsDefault,
					&pvp.Price.PriceList.Priority, &pvp.Price.PriceList.Status,
					&pvp.Price.PriceList.CreatedAt, &pvp.Price.PriceList.UpdatedAt,
			)
			if err != nil {
					return nil, errors.NewMysqlError(err)
			}
			prices = append(prices, pvp)
	}

	return prices, nil
}

func (r *MySQLRepository) CountPriceLists(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.QueryRowContext(ctx, 
			"SELECT COUNT(*) FROM price_lists WHERE deleted_at IS NULL").
			Scan(&count)
	if err != nil {
			return 0, errors.NewMysqlError(err)
	}
	return count, nil
}