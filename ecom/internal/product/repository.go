// internal/product/repository.go
package product

import (
	"context"
	"database/sql"
	"ecommerce/pkg/apperrors"
	"errors"
	"fmt"
)

type Repository struct {
    db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
    return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, product *Product) error {
    query := `
        INSERT INTO products (id, name, description, price, stock, category_id, created_at, updated_at)
        VALUES (?, ?, ?, ?, ?, ?, NOW(), NOW())
    `
    
    result, err := r.db.ExecContext(
        ctx,
        query,
        product.ID,
        product.Name,
        product.Description,
        product.Price,
        product.Stock,
        product.CategoryID,
    )
    
    if err != nil {
        if isDuplicateError(err) {
            return apperrors.NewError(
                apperrors.ErrAlreadyExists,
                "Product already exists",
            )
        }
        return fmt.Errorf("error creating product: %w", err)
    }

    rows, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("error getting rows affected: %w", err)
    }

    if rows == 0 {
        return apperrors.NewError(
            apperrors.ErrInternal,
            "No rows affected when creating product",
        )
    }

    return nil
}

func (r *Repository) Update(ctx context.Context, product *Product) error {
    query := `
        UPDATE products 
        SET name = ?, description = ?, price = ?, stock = ?, 
            category_id = ?, updated_at = NOW()
        WHERE id = ?
    `
    
    result, err := r.db.ExecContext(
        ctx,
        query,
        product.Name,
        product.Description,
        product.Price,
        product.Stock,
        product.CategoryID,
        product.ID,
    )
    
    if err != nil {
        return fmt.Errorf("error updating product: %w", err)
    }

    rows, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("error getting rows affected: %w", err)
    }

    if rows == 0 {
        return apperrors.NewError(
            apperrors.ErrNotFound,
            "Product not found",
        )
    }

    return nil
}

func (r *Repository) Delete(ctx context.Context, id string) error {
    query := "DELETE FROM products WHERE id = ?"
    
    result, err := r.db.ExecContext(ctx, query, id)
    if err != nil {
        return fmt.Errorf("error deleting product: %w", err)
    }

    rows, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("error getting rows affected: %w", err)
    }

    if rows == 0 {
        return apperrors.NewError(
            apperrors.ErrNotFound,
            "Product not found",
        )
    }

    return nil
}

func (r *Repository) GetByID(ctx context.Context, id string) (*Product, error) {
    query := `
        SELECT p.id, p.name, p.description, p.price, p.stock, 
               p.category_id, c.name as category_name,
               p.created_at, p.updated_at
        FROM products p
        LEFT JOIN categories c ON p.category_id = c.id
        WHERE p.id = ?
    `
    
    var product Product
    var category Category
    
    err := r.db.QueryRowContext(ctx, query, id).Scan(
        &product.ID,
        &product.Name,
        &product.Description,
        &product.Price,
        &product.Stock,
        &product.CategoryID,
        &category.Name,
        &product.CreatedAt,
        &product.UpdatedAt,
    )
    
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, apperrors.NewError(
                apperrors.ErrNotFound,
                "Product not found",
            )
        }
        return nil, fmt.Errorf("error getting product: %w", err)
    }

    category.ID = product.CategoryID
    product.Category = &category

    return &product, nil
}

func (r *Repository) List(ctx context.Context, page, limit int) ([]*Product, error) {
    offset := (page - 1) * limit
    
    query := `
        SELECT p.id, p.name, p.description, p.price, p.stock, 
               p.category_id, c.name as category_name,
               p.created_at, p.updated_at
        FROM products p
        LEFT JOIN categories c ON p.category_id = c.id
        LIMIT ? OFFSET ?
    `
    
    rows, err := r.db.QueryContext(ctx, query, limit, offset)
    if err != nil {
        return nil, fmt.Errorf("error listing products: %w", err)
    }
    defer rows.Close()

    var products []*Product
    
    for rows.Next() {
        var product Product
        var category Category
        
        err := rows.Scan(
            &product.ID,
            &product.Name,
            &product.Description,
            &product.Price,
            &product.Stock,
            &product.CategoryID,
            &category.Name,
            &product.CreatedAt,
            &product.UpdatedAt,
        )
        
        if err != nil {
            return nil, fmt.Errorf("error scanning product: %w", err)
        }

        category.ID = product.CategoryID
        product.Category = &category
        products = append(products, &product)
    }

    return products, nil
}

func (r *Repository) ListCategories(ctx context.Context) ([]*Category, error) {
    query := "SELECT id, name FROM categories"
    
    rows, err := r.db.QueryContext(ctx, query)
    if err != nil {
        return nil, fmt.Errorf("error listing categories: %w", err)
    }
    defer rows.Close()

    var categories []*Category
    
    for rows.Next() {
        var category Category
        if err := rows.Scan(&category.ID, &category.Name); err != nil {
            return nil, fmt.Errorf("error scanning category: %w", err)
        }
        categories = append(categories, &category)
    }

    return categories, nil
}

func isDuplicateError(err error) bool {
    // Implementar lógica para detectar errores de duplicado en MySQL
    return false
}