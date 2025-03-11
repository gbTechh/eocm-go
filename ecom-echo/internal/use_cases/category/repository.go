package category

import (
	"context"
	"database/sql"
	"ecom/internal/shared/errors"
	"fmt"
)

type Repository interface {
	Create(ctx context.Context, c *Category) error
	GetByID(ctx context.Context, id int64) (*Category, error)
	Update(ctx context.Context, c *Category) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, p *Pagination) ([]Category, int64, error)
	HasChildren(ctx context.Context, id int64) (bool, error)
}

type MySQLRepository struct {
	db *sql.DB
}

func NewMySQLRepository(db *sql.DB) Repository {
	return &MySQLRepository{db: db}
}

const (
	queryInsertCategory = `
			INSERT INTO categories (name, slug, description, parent_id, id_media, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, NOW(), NOW())
	`

	queryGetCategoryWithMedia = `
			SELECT 
					c.id, c.name, c.slug, c.description, c.parent_id, c.id_media,
					m.id, m.file_name, m.file_path, m.mime_type, m.size, m.created_at, m.updated_at,
					c.created_at, c.updated_at
			FROM categories c
			LEFT JOIN media m ON c.id_media = m.id
			WHERE c.id = ? AND c.deleted_at IS NULL
	`

	queryUpdateCategory = `
			UPDATE categories 
			SET name = ?, slug = ?, description = ?, parent_id = ?, id_media = ?, updated_at = NOW()
			WHERE id = ? AND deleted_at IS NULL
	`

	queryDeleteCategory = `
			UPDATE categories
			SET deleted_at = NOW()
			WHERE id = ? AND deleted_at IS NULL
	`

	queryHasChildren = `
			SELECT EXISTS(
					SELECT 1 FROM categories 
					WHERE parent_id = ? AND deleted_at IS NULL
			)
	`
)

func (r *MySQLRepository) Create(ctx context.Context, c *Category) error {
	result, err := r.db.ExecContext(ctx, queryInsertCategory,
			c.Name, c.Slug, c.Description, c.ParentID, c.IDMedia)
	if err != nil {
		return errors.NewMysqlError(err)	
	}

	id, err := result.LastInsertId()
	if err != nil {
		return errors.NewInternalError("Error al obtener el id creado", err)
	}

	c.ID = id
	return nil
}

func (r *MySQLRepository) GetByID(ctx context.Context, id int64) (*Category, error) {
	c := &Category{}

	// Usamos NullString y NullInt64 para manejar campos NULL
	var description, mediaFileName, mediaFilePath, mediaMimeType sql.NullString
	var mediaSize sql.NullFloat64
	var parentID, mediaID sql.NullInt64
	var mediaID2 sql.NullInt64
	var mediaCreatedAt, mediaUpdatedAt sql.NullTime

	err := r.db.QueryRowContext(ctx, queryGetCategoryWithMedia, id).Scan(
			&c.ID, &c.Name, &c.Slug, &description, &parentID, &mediaID,
			&mediaID2, &mediaFileName, &mediaFilePath, &mediaMimeType, &mediaSize,
			&mediaCreatedAt, &mediaUpdatedAt,
			&c.CreatedAt, &c.UpdatedAt,
	)

	if err == sql.ErrNoRows {
			return nil, errors.NewNotFoundError("Categoría no encontrada")
	}
	if err != nil {
			return nil, err
	}

	// Asignar valores opcionales
	if description.Valid {
			c.Description = description.String
	}
	if parentID.Valid {
			c.ParentID = &parentID.Int64
	}
	if mediaID.Valid {
			c.IDMedia = &mediaID.Int64
			c.Media = &Media{
					ID:        mediaID.Int64,
					FileName:  mediaFileName.String,
					FilePath:  mediaFilePath.String,
					MimeType:  mediaMimeType.String,
					Size:      mediaSize.Float64,
					CreatedAt: mediaCreatedAt.Time,
					UpdatedAt: mediaUpdatedAt.Time,
			}
	}

	return c, nil
}

func (r *MySQLRepository) Update(ctx context.Context, c *Category) error {
	result, err := r.db.ExecContext(ctx, queryUpdateCategory,
			c.Name, c.Slug, c.Description, c.ParentID, c.IDMedia, c.ID)
	if err != nil {
		return errors.NewMysqlError(err)	
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return errors.NewInternalError("Error al obtener la categoría", err)
	}
	if rows == 0 {
		return errors.NewNotFoundError("Categoría no encontrada")
	}
	return nil
}

func (r *MySQLRepository) Delete(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, queryDeleteCategory, id)
	if err != nil {
			return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
			return err
	}
	if rows == 0 {
			return errors.NewNotFoundError("Categoría no encontrada")
	}
	return nil
}

func (r *MySQLRepository) HasChildren(ctx context.Context, id int64) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, queryHasChildren, id).Scan(&exists)
	return exists, err
}

func (r *MySQLRepository) List(ctx context.Context, p *Pagination) ([]Category, int64, error) {
	// Query base para contar
	countQuery := `
			SELECT COUNT(*) 
			FROM categories 
			WHERE deleted_at IS NULL
	`
	
	// Query base para listar con JOIN a media
	listQuery := `
			SELECT 
					c.id, c.name, c.slug, c.description, c.parent_id, c.id_media,
					m.id, m.file_name, m.file_path, m.mime_type, m.size, m.created_at, m.updated_at,
					c.created_at, c.updated_at
			FROM categories c
			LEFT JOIN media m ON c.id_media = m.id
			WHERE c.deleted_at IS NULL
	`

	args := []interface{}{}

	// Agregar condiciones de búsqueda
	if p.Search != "" {
			countQuery += " AND (name LIKE ? OR slug LIKE ?)"
			listQuery += " AND (c.name LIKE ? OR c.slug LIKE ?)"
			searchTerm := "%" + p.Search + "%"
			args = append(args, searchTerm, searchTerm)
	}

	if p.ParentID != nil {
			countQuery += " AND parent_id = ?"
			listQuery += " AND c.parent_id = ?"
			args = append(args, *p.ParentID)
	}

	// Contar total
	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
			return nil, 0, err
	}

	// Agregar paginación
	listQuery += " ORDER BY c.created_at DESC LIMIT ? OFFSET ?"
	args = append(args, p.PerPage, (p.Page-1)*p.PerPage)

	 // Debug: imprimir query y argumentos
	fmt.Printf("Query: %s\n", listQuery)
	fmt.Printf("Args: %v\n", args)

	// Ejecutar query principal
	rows, err := r.db.QueryContext(ctx, listQuery, args...)
	if err != nil {
			return nil, 0, err
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
			c := Category{}

			var description, mediaFileName, mediaFilePath, mediaMimeType sql.NullString
			var mediaSize sql.NullFloat64
			var parentID, mediaID sql.NullInt64
			var mediaID2 sql.NullInt64
			var mediaCreatedAt, mediaUpdatedAt sql.NullTime

			err := rows.Scan(
					&c.ID, &c.Name, &c.Slug, &description, &parentID, &mediaID,
					&mediaID2, &mediaFileName, &mediaFilePath, &mediaMimeType, &mediaSize,
					&mediaCreatedAt, &mediaUpdatedAt,
					&c.CreatedAt, &c.UpdatedAt,
			)
			if err != nil {
					return nil, 0, err
			}

			if description.Valid {
					c.Description = description.String
			}
			if parentID.Valid {
					c.ParentID = &parentID.Int64
			}
			if mediaID.Valid {
					c.IDMedia = &mediaID.Int64
					c.Media = &Media{
							ID:        mediaID.Int64,
							FileName:  mediaFileName.String,
							FilePath:  mediaFilePath.String,
							MimeType:  mediaMimeType.String,
							Size:      mediaSize.Float64,
							CreatedAt: mediaCreatedAt.Time,
							UpdatedAt: mediaUpdatedAt.Time,
					}
			}

			categories = append(categories, c)
	}

	if err = rows.Err(); err != nil {
			return nil, 0, err
	}

	return categories, total, nil
}