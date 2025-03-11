// repository.go
package customer

import (
	"context"
	"database/sql"
	"ecom/internal/shared/errors"
	"strings"
)

type Repository interface {
	// Customer operations
	CreateCustomer(ctx context.Context, c *Customer) error
	GetCustomerByID(ctx context.Context, id int64) (*Customer, error)
	GetCustomerByEmail(ctx context.Context, email string) (*Customer, error)
	UpdateCustomer(ctx context.Context, c *Customer) error
	DeleteCustomer(ctx context.Context, id int64) error
	ListCustomers(ctx context.Context, p *Pagination) ([]Customer, int64, error)
	
	// Group operations
	CreateGroup(ctx context.Context, g *Group) error
	GetGroupByID(ctx context.Context, id int64) (*Group, error)
	UpdateGroup(ctx context.Context, g *Group) error
	DeleteGroup(ctx context.Context, id int64) error
	ListGroups(ctx context.Context, p *Pagination) ([]Group, int64, error)
	
	// Customer-Group operations
	AssignCustomersToGroup(ctx context.Context, groupID int64, customerIDs []int64) error
	RemoveCustomersFromGroup(ctx context.Context, groupID int64, customerIDs []int64) error
	GetCustomerGroups(ctx context.Context, customerID int64) ([]Group, error)
}

type MySQLRepository struct {
	db *sql.DB
}

func NewMySQLRepository(db *sql.DB) Repository {
	return &MySQLRepository{db: db}
}

func (r *MySQLRepository) CreateCustomer(ctx context.Context, c *Customer) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
			return err
	}

	query := `
			INSERT INTO customers (name, last_name, email, phone, active, is_guest, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, NOW(), NOW())
	`
	result, err := tx.ExecContext(ctx, query,
			c.Name, c.LastName, c.Email, c.Phone, c.Active, c.IsGuest)
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

	// Asignar grupos si existen
	if len(c.Groups) > 0 {
			groupIDs := make([]int64, len(c.Groups))
			for i, g := range c.Groups {
					groupIDs[i] = g.ID
			}
			if err := r.assignCustomersToGroupTx(ctx, tx, id, groupIDs); err != nil {
					tx.Rollback()
					return err
			}
	}

	return tx.Commit()
}

func (r *MySQLRepository) GetCustomerByID(ctx context.Context, id int64) (*Customer, error) {
	query := `
			SELECT id, name, last_name, email, phone, active, is_guest, created_at, updated_at
			FROM customers
			WHERE id = ? AND deleted_at IS NULL
	`
	c := &Customer{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
			&c.ID, &c.Name, &c.LastName, &c.Email, &c.Phone,
			&c.Active, &c.IsGuest, &c.CreatedAt, &c.UpdatedAt,
	)
	if err == sql.ErrNoRows {
			return nil, errors.NewNotFoundError("Cliente no encontrado")
	}
	if err != nil {
			return nil, errors.NewMysqlError(err)
	}

	// Obtener grupos del cliente
	groups, err := r.GetCustomerGroups(ctx, id)
	if err != nil {
			return nil, err
	}
	c.Groups = groups

	return c, nil
}

func (r *MySQLRepository) GetCustomerByEmail(ctx context.Context, email string) (*Customer, error) {
	query := `
			SELECT id, name, last_name, email, phone, active, is_guest, created_at, updated_at
			FROM customers
			WHERE email = ? AND deleted_at IS NULL
	`
	c := &Customer{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
			&c.ID, &c.Name, &c.LastName, &c.Email, &c.Phone,
			&c.Active, &c.IsGuest, &c.CreatedAt, &c.UpdatedAt,
	)
	if err == sql.ErrNoRows {
			return nil, errors.NewNotFoundError("Cliente no encontrado")
	}
	if err != nil {
			return nil, errors.NewMysqlError(err)
	}

	// Obtener grupos del cliente
	groups, err := r.GetCustomerGroups(ctx, c.ID)
	if err != nil {
			return nil, err
	}
	c.Groups = groups

	return c, nil
}

func (r *MySQLRepository) UpdateCustomer(ctx context.Context, c *Customer) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
			return err
	}

	query := `
			UPDATE customers 
			SET name = ?, last_name = ?, email = ?, phone = ?, 
					active = ?, is_guest = ?, updated_at = NOW()
			WHERE id = ? AND deleted_at IS NULL
	`
	result, err := tx.ExecContext(ctx, query,
			c.Name, c.LastName, c.Email, c.Phone,
			c.Active, c.IsGuest, c.ID)
	if err != nil {
			tx.Rollback()
			return errors.NewMysqlError(err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
			tx.Rollback()
			return errors.NewInternalError("Error al actualizar cliente", err)
	}
	if rows == 0 {
			tx.Rollback()
			return errors.NewNotFoundError("Cliente no encontrado")
	}

	// Actualizar grupos si se proporcionaron
	if c.Groups != nil {
			// Eliminar grupos actuales
			if _, err := tx.ExecContext(ctx,
					"DELETE FROM n_customers_groups WHERE id_customer = ?", c.ID); err != nil {
					tx.Rollback()
					return errors.NewMysqlError(err)
			}

			// Asignar nuevos grupos
			groupIDs := make([]int64, len(c.Groups))
			for i, g := range c.Groups {
					groupIDs[i] = g.ID
			}
			if err := r.assignCustomersToGroupTx(ctx, tx, c.ID, groupIDs); err != nil {
					tx.Rollback()
					return err
			}
	}

	return tx.Commit()
}

func (r *MySQLRepository) DeleteCustomer(ctx context.Context, id int64) error {
	query := `
			UPDATE customers
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
			return errors.NewNotFoundError("Cliente no encontrado")
	}
	return nil
}

func (r *MySQLRepository) ListCustomers(ctx context.Context, p *Pagination) ([]Customer, int64, error) {
	countQuery := `
			SELECT COUNT(*) 
			FROM customers 
			WHERE deleted_at IS NULL
	`
	args := []interface{}{}

	if p.Search != "" {
			countQuery += " AND (name LIKE ? OR last_name LIKE ? OR email LIKE ?)"
			searchTerm := "%" + p.Search + "%"
			args = append(args, searchTerm, searchTerm, searchTerm)
	}

	if p.Active != nil {
			countQuery += " AND active = ?"
			args = append(args, *p.Active)
	}

	if p.IsGuest != nil {
			countQuery += " AND is_guest = ?"
			args = append(args, *p.IsGuest)
	}

	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
			return nil, 0, err
	}

	query := `
			SELECT id, name, last_name, email, phone, active, is_guest, created_at, updated_at
			FROM customers
			WHERE deleted_at IS NULL
	`

	if p.Search != "" {
			query += " AND (name LIKE ? OR last_name LIKE ? OR email LIKE ?)"
	}

	if p.Active != nil {
			query += " AND active = ?"
	}

	if p.IsGuest != nil {
			query += " AND is_guest = ?"
	}

	query += " ORDER BY created_at DESC LIMIT ? OFFSET ?"
	args = append(args, p.PerPage, (p.Page-1)*p.PerPage)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
			return nil, 0, err
	}
	defer rows.Close()

	var customers []Customer
	for rows.Next() {
			var c Customer
			err := rows.Scan(
					&c.ID, &c.Name, &c.LastName, &c.Email, &c.Phone,
					&c.Active, &c.IsGuest, &c.CreatedAt, &c.UpdatedAt,
			)
			if err != nil {
					return nil, 0, err
			}

			// Obtener grupos del cliente
			groups, err := r.GetCustomerGroups(ctx, c.ID)
			if err != nil {
					return nil, 0, err
			}
			c.Groups = groups

			customers = append(customers, c)
	}

	return customers, total, nil
}

// Group operations
func (r *MySQLRepository) CreateGroup(ctx context.Context, g *Group) error {
	query := `
			INSERT INTO customer_groups (id_price_list, name, description, created_at, updated_at)
			VALUES (?, ?, ?, NOW(), NOW())
	`
	result, err := r.db.ExecContext(ctx, query,
			g.IDPriceList, g.Name, g.Description)
	if err != nil {
			return errors.NewMysqlError(err)
	}

	id, err := result.LastInsertId()
	if err != nil {
			return errors.NewInternalError("Error al obtener el id creado", err)
	}

	g.ID = id
	return nil
}

func (r *MySQLRepository) GetGroupByID(ctx context.Context, id int64) (*Group, error) {
	query := `
			SELECT cg.id, cg.id_price_list, cg.name, cg.description, cg.created_at, cg.updated_at,
               pl.id, pl.name, pl.description, pl.is_default, pl.priority
			FROM customer_groups cg
			LEFT JOIN price_lists pl ON pl.id = cg.id_price_list
			WHERE cg.id = ? AND cg.deleted_at IS NULL
	`
	g := &Group{PriceList: &PriceList{}}
	var idPriceList sql.NullInt64
	var plID sql.NullInt64
	var plName, plDescription sql.NullString
	var plIsDefault sql.NullBool
	var plPriority sql.NullInt64

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&g.ID, &idPriceList, &g.Name, &g.Description, &g.CreatedAt, &g.UpdatedAt,
		&plID, &plName, &plDescription, &plIsDefault, &plPriority,
	)

	if err == sql.ErrNoRows {
			return nil, errors.NewNotFoundError("Grupo no encontrado")
	}
	if err != nil {
			return nil, errors.NewMysqlError(err)
	}

	if idPriceList.Valid {
		g.IDPriceList = &idPriceList.Int64
		if plID.Valid {
			g.PriceList = &PriceList{
				ID:          plID.Int64,
				Name:        plName.String,
				Description: plDescription.String,
				IsDefault:   plIsDefault.Bool,
				Priority:    int(plPriority.Int64),
			}
		}
	} else {
		g.PriceList = nil
	}

	return g, nil
}

func (r *MySQLRepository) UpdateGroup(ctx context.Context, g *Group) error {
	query := `
			UPDATE customer_groups 
			SET id_price_list = ?, name = ?, description = ?, updated_at = NOW()
			WHERE id = ? AND deleted_at IS NULL
	`
	result, err := r.db.ExecContext(ctx, query,
			g.IDPriceList, g.Name, g.Description, g.ID)
	if err != nil {
			return errors.NewMysqlError(err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
			return errors.NewInternalError("Error al actualizar grupo", err)
	}
	if rows == 0 {
			return errors.NewNotFoundError("Grupo no encontrado")
	}
	return nil
}

// Continuación de repository.go...

func (r *MySQLRepository) DeleteGroup(ctx context.Context, id int64) error {
    tx, err := r.db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }

    // Primero eliminar las relaciones
    _, err = tx.ExecContext(ctx, 
        "DELETE FROM n_customers_groups WHERE id_customer_group = ?", id)
    if err != nil {
        tx.Rollback()
        return errors.NewMysqlError(err)
    }

    // Luego soft delete del grupo
    query := `
        UPDATE customer_groups
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
        return errors.NewNotFoundError("Grupo no encontrado")
    }

    return tx.Commit()
}

func (r *MySQLRepository) ListGroups(ctx context.Context, p *Pagination) ([]Group, int64, error) {
    countQuery := `
        SELECT COUNT(*) 
        FROM customer_groups 
        WHERE deleted_at IS NULL
    `
    args := []interface{}{}

    if p.Search != "" {
        countQuery += " AND name LIKE ?"
        args = append(args, "%"+p.Search+"%")
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
        SELECT cg.id, cg.id_price_list, cg.name, cg.description, cg.created_at, cg.updated_at,
               pl.id, pl.name, pl.description, pl.is_default, pl.priority
        FROM customer_groups cg
        LEFT JOIN price_lists pl ON pl.id = cg.id_price_list
        WHERE cg.deleted_at IS NULL
    `

    if p.Search != "" {
        query += " AND cg.name LIKE ?"
    }

    query += " ORDER BY cg.created_at DESC LIMIT ? OFFSET ?"
    args = append(args, p.PerPage, (p.Page-1)*p.PerPage)

    rows, err := r.db.QueryContext(ctx, query, args...)
    if err != nil {
        return nil, 0, err
    }
    defer rows.Close()

    var groups []Group
    for rows.Next() {
         g := Group{PriceList: &PriceList{}}
        var idPriceList sql.NullInt64
        var plID sql.NullInt64
        var plName, plDescription sql.NullString
        var plIsDefault sql.NullBool
        var plPriority sql.NullInt64

        err := rows.Scan(
            &g.ID, &idPriceList, &g.Name, &g.Description, &g.CreatedAt, &g.UpdatedAt,
            &plID, &plName, &plDescription, &plIsDefault, &plPriority,
        )
        if err != nil {
            return nil, 0, err
        }

        if idPriceList.Valid {
					g.IDPriceList = &idPriceList.Int64
					if plID.Valid {
						g.PriceList = &PriceList{
							ID:          plID.Int64,
							Name:        plName.String,
							Description: plDescription.String,
							IsDefault:   plIsDefault.Bool,
							Priority:    int(plPriority.Int64),
						}
					}
        } else {
					g.PriceList = nil
        }

        groups = append(groups, g)
    }

    return groups, total, nil
}

func (r *MySQLRepository) AssignCustomersToGroup(ctx context.Context, groupID int64, customerIDs []int64) error {
    tx, err := r.db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }

    // Verificar que el grupo existe
    exists, err := r.groupExists(ctx, tx, groupID)
    if err != nil {
        tx.Rollback()
        return err
    }
    if !exists {
        tx.Rollback()
        return errors.NewNotFoundError("Grupo no encontrado")
    }

    // Insertar las relaciones
    if len(customerIDs) > 0 {
        query := "INSERT INTO n_customers_groups (id_customer, id_customer_group) VALUES "
        values := make([]string, len(customerIDs))
        args := make([]interface{}, 0, len(customerIDs)*2)

        for i, customerID := range customerIDs {
            values[i] = "(?, ?)"
            args = append(args, customerID, groupID)
        }

        query += strings.Join(values, ",")
        query += " ON DUPLICATE KEY UPDATE id_customer=id_customer" // Ignorar si ya existe

        if _, err := tx.ExecContext(ctx, query, args...); err != nil {
            tx.Rollback()
            return errors.NewMysqlError(err)
        }
    }

    return tx.Commit()
}

func (r *MySQLRepository) RemoveCustomersFromGroup(ctx context.Context, groupID int64, customerIDs []int64) error {
    query := `
        DELETE FROM n_customers_groups 
        WHERE id_customer_group = ? AND id_customer IN (?` + strings.Repeat(",?", len(customerIDs)-1) + ")"
    
    args := make([]interface{}, 0, len(customerIDs)+1)
    args = append(args, groupID)
    for _, id := range customerIDs {
        args = append(args, id)
    }

    _, err := r.db.ExecContext(ctx, query, args...)
    if err != nil {
        return errors.NewMysqlError(err)
    }

    return nil
}

func (r *MySQLRepository) GetCustomerGroups(ctx context.Context, customerID int64) ([]Group, error) {
    query := `
        SELECT cg.id, cg.id_price_list, cg.name, cg.description, cg.created_at, cg.updated_at,
               pl.id, pl.name, pl.description, pl.is_default, pl.priority
        FROM customer_groups cg
        JOIN n_customers_groups ncg ON ncg.id_customer_group = cg.id
        LEFT JOIN price_lists pl ON pl.id = cg.id_price_list
        WHERE ncg.id_customer = ? AND cg.deleted_at IS NULL
    `

    rows, err := r.db.QueryContext(ctx, query, customerID)
    if err != nil {
        return nil, errors.NewMysqlError(err)
    }
    defer rows.Close()

    var groups []Group
    for rows.Next() {
       	g  := Group{PriceList: &PriceList{}}
        var idPriceList sql.NullInt64
        var plID sql.NullInt64
        var plName, plDescription sql.NullString
        var plIsDefault sql.NullBool
        var plPriority sql.NullInt64

       	err := rows.Scan(
            &g.ID, &idPriceList, &g.Name, &g.Description, &g.CreatedAt, &g.UpdatedAt,
            &plID, &plName, &plDescription, &plIsDefault, &plPriority,
        )
        if err != nil {
            return nil, err
        }

        if idPriceList.Valid {
					g.IDPriceList = &idPriceList.Int64
					if plID.Valid {
						g.PriceList = &PriceList{
							ID:          plID.Int64,
							Name:        plName.String,
							Description: plDescription.String,
							IsDefault:   plIsDefault.Bool,
							Priority:    int(plPriority.Int64),
						}
					}
        } else {
            g.PriceList = nil
        }

        groups = append(groups, g)
    }

    return groups, nil
}

// Función auxiliar
func (r *MySQLRepository) groupExists(ctx context.Context, tx *sql.Tx, groupID int64) (bool, error) {
	var exists bool
	err := tx.QueryRowContext(ctx,
			"SELECT EXISTS(SELECT 1 FROM customer_groups WHERE id = ? AND deleted_at IS NULL)",
			groupID).Scan(&exists)
	if err != nil {
			return false, errors.NewMysqlError(err)
	}
	return exists, nil
}

func (r *MySQLRepository) assignCustomersToGroupTx(ctx context.Context, tx *sql.Tx, customerID int64, groupIDs []int64) error {
	if len(groupIDs) == 0 {
			return nil
	}

	// Construir la consulta con múltiples valores
	query := "INSERT INTO n_customers_groups (id_customer, id_customer_group) VALUES "
	values := make([]string, len(groupIDs))
	args := make([]interface{}, 0, len(groupIDs)*2)

	for i, groupID := range groupIDs {
			values[i] = "(?, ?)"
			args = append(args, customerID, groupID)
	}

	query += strings.Join(values, ",")
	query += " ON DUPLICATE KEY UPDATE id_customer=id_customer" // Ignorar si ya existe

	_, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
			return errors.NewMysqlError(err)
	}

	return nil
}