// internal/media/service.go
package media

import (
	"context"
	"ecom/internal/shared/errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

type Service struct {
	repo   Repository
	upload UploadStrategy
}

// UploadStrategy define la interfaz para estrategias de almacenamiento
type UploadStrategy interface {
	SaveFile(fileName string, content []byte) (string, error)
}

// LocalUploadStrategy implementa el almacenamiento local de archivos
type LocalUploadStrategy struct {
	uploadDir string
}

// NewLocalUploadStrategy crea una nueva instancia de LocalUploadStrategy
func NewLocalUploadStrategy(uploadDir string) UploadStrategy {
	return &LocalUploadStrategy{uploadDir: uploadDir}
}

// SaveFile guarda el archivo en el sistema de archivos local
func (s *LocalUploadStrategy) SaveFile(fileName string, content []byte) (string, error) {
	// Generar nombre de archivo único
	ext := filepath.Ext(fileName)
	id, _ := gonanoid.New(10)
	uniqueFileName := fmt.Sprintf("%s%s",id, ext)
	filePath := filepath.Join(s.uploadDir, uniqueFileName)

	// Asegurarse de que el directorio exista
	err := os.MkdirAll(s.uploadDir, os.ModePerm)
	if err != nil {
		return "", err
	}

	// Guardar archivo
	err = os.WriteFile(filePath, content, 0644)
	if err != nil {
		return "", err
	}

	return filePath, nil
}

// NewService crea una nueva instancia del servicio de medios
func NewService(repo Repository, uploadStrategy UploadStrategy) *Service {
	return &Service{
		repo:   repo,
		upload: uploadStrategy,
	}
}

// Upload maneja la subida de un nuevo archivo
func (s *Service) Upload(ctx context.Context, req *CreateMediaRequest) (*Media, error) {
	// Validar tamaño máximo del archivo (por ejemplo, 10MB)
	const maxFileSize = 10 * 1024 * 1024 // 10MB
	if len(req.File) > maxFileSize {
		return nil, errors.NewValidationError("Archivo demasiado GRANDE. Máximo 10MB")
	}

	// Guardar archivo usando la estrategia de almacenamiento
	filePath, err := s.upload.SaveFile(req.FileName, req.File)
	if err != nil {
		return nil, errors.NewInternalError("Error al guardar el archivo", err)
	}

	// Crear registro de media
	media := &Media{
		FileName:   req.FileName,
		FilePath:   filePath,
		MimeType:   req.MimeType,
		Size:       float64(len(req.File)) / 1024.0, // Tamaño en KB
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Guardar en base de datos
	if err := s.repo.Upload(ctx, media); err != nil {
		// Eliminar archivo si falla la inserción en base de datos
		os.Remove(filePath)
		return nil, errors.NewInternalError("Error al guardar los metadatos", err)
	}

	return media, nil
}

// GetByID recupera un archivo por su ID
func (s *Service) GetByID(ctx context.Context, id int64) (*Media, error) {
	media, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if err == ErrNotFound {
			return nil, errors.NewNotFoundError("Archivo no encontrado")
		}
		return nil, errors.NewInternalError("Error al buscar el archivo", err)
	}
	return media, nil
}

// Update actualiza los metadatos de un archivo
func (s *Service) Update(ctx context.Context, id int64, req *UpdateMediaRequest) (*Media, error) {
    // Obtener el media actual
    oldMedia, err := s.repo.GetByID(ctx, id)
    if err != nil {
        return nil, err
    }

    // Actualizar si hay nuevo archivo
    if req.File != nil {
			const maxFileSize = 10 * 1024 * 1024 // 10MB
			if len(req.File) > maxFileSize {
				return nil, errors.NewValidationError("Archivo demasiado GRANDE. Máximo 10MB")
			}
        // Guardar el nuevo archivo
        filePath, err := s.upload.SaveFile(*req.FileName, req.File)
        if err != nil {
            return nil, errors.NewInternalError("Error al guardar el archivo", err)
        }

        // Eliminar el archivo antiguo
        if err := os.Remove(oldMedia.FilePath); err != nil {
            fmt.Printf("Error eliminando archivo antiguo %s: %v\n", oldMedia.FilePath, err)
        } else {
            fmt.Printf("Archivo antiguo eliminado exitosamente: %s\n", oldMedia.FilePath)
        }

        // Actualizar información del archivo
        oldMedia.FilePath = filePath
        oldMedia.FileName = *req.FileName
        oldMedia.MimeType = *req.MimeType
        oldMedia.Size = float64(len(req.File)) / 1024.0
    }

    // Actualizar solo el nombre si se proporciona sin archivo
    if req.FileName != nil && req.File == nil {
        oldMedia.FileName = *req.FileName
    }

    // Actualizar timestamp
    oldMedia.UpdatedAt = time.Now()

    // Guardar cambios en la base de datos
    if err := s.repo.Update(ctx, oldMedia); err != nil {
        return nil, errors.NewInternalError("Error al actualizar el archivo", err)
    }

    return oldMedia, nil
}


// Delete elimina un archivo
func (s *Service) Delete(ctx context.Context, id int64) error {
	// Primero, obtener el archivo para verificar su existencia
	media, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if err == ErrNotFound {
			return errors.NewNotFoundError("Archivo no encontrado")
		}
		return errors.NewInternalError("Error al buscar el archivo", err)
	}

	// Eliminar registro de base de datos
	if err := s.repo.Delete(ctx, id); err != nil {
		return errors.NewInternalError("Error al eliminar el archivo", err)
	}

	// Intentar eliminar archivo físico (soft delete en BD)
	if err := os.Remove(media.FilePath); err != nil && !os.IsNotExist(err) {
		// Log error pero no retornar error para no afectar la eliminación del registro
		fmt.Printf("Error al eliminar archivo físico: %v\n", err)
	}

	return nil
}

// List recupera una lista de archivos con filtrado y paginación
func (s *Service) List(ctx context.Context, p *Pagination) ([]Media, int64, error) {
   // Validar paginación
   if p.Page < 1 {
       p.Page = 1
   }
   if p.PerPage < 1 || p.PerPage > 100 {
       p.PerPage = 10
   }

   // Obtener datos del repositorio
   media, total, err := s.repo.List(ctx, p)
   if err != nil {
       return nil, 0, errors.NewInternalError("Error al listar archivos", err)
   }

   return media, total, nil
}

// Esta función es la misma que está en el handler, la duplicamos aquí para mantener la separación de responsabilidades
func isValidMimeType(mimeType string) bool {
	validTypes := map[string]bool{
		"image/jpeg":        true,
		"image/png":         true,
		"image/gif":         true,
		"application/pdf":   true,
	}
	return validTypes[mimeType]
}