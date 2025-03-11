// internal/media/handler.go
package media

import (
	"ecom/internal/shared/errors"
	"ecom/internal/shared/response"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/labstack/echo/v4"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

type Handler struct {
   service *Service
}

func NewHandler(e *echo.Group, s *Service) {
	h := &Handler{service: s}

	// Rutas
	e.POST("", h.Upload)  
	e.GET("/:id", h.GetByID)
	e.PUT("/:id", h.Update)
	e.DELETE("/:id", h.Delete)
	e.GET("", h.List)

  

}

func (h *Handler) Upload(c echo.Context) error {
	// Obtener el archivo
	file, err := c.FormFile("file")
	if err != nil {
			return c.JSON(http.StatusBadRequest, 
					response.Error("Error al procesar el archivo", err.Error()))
	}

	// Validar el tipo de archivo
	contentType := file.Header.Get("Content-Type")
	if !isValidMimeType(contentType) {
			return c.JSON(http.StatusBadRequest, 
					response.Error("Tipo de archivo no permitido", nil))
	}

	// Leer el archivo
	src, err := file.Open()
	if err != nil {
			return h.handleError(c, err)
	}
	defer src.Close()

	fileBytes := make([]byte, file.Size)
	_, err = src.Read(fileBytes)
	if err != nil {
			return h.handleError(c, err)
	}

	// Obtener extensión y generar nombre único
	ext := filepath.Ext(file.Filename)
	id, _ := gonanoid.New(10)
	newFileName := fmt.Sprintf("%s%s", id, ext)

	// Crear request
	req := CreateMediaRequest{
			File:     fileBytes,
			FileName: newFileName,
			MimeType: contentType,
	}

	// Validar request
	if err := c.Validate(&req); err != nil {
			return c.JSON(http.StatusBadRequest, 
					response.Error("Error de validación", err.Error()))
	}

	// Llamar al servicio
	media, err := h.service.Upload(c.Request().Context(), &req)
	if err != nil {
			return h.handleError(c, err)
	}

	// Crear respuesta
	resp := MediaResponse{
			ID:        media.ID,
			FileName:  media.FileName,
			FilePath:  media.FilePath,
			MimeType:  media.MimeType,
			Size:      media.Size,
			CreatedAt: media.CreatedAt,
			UpdatedAt: media.UpdatedAt,
	}

	return c.JSON(http.StatusCreated, 
			response.Success("Archivo subido exitosamente", resp))
}

func (h *Handler) GetByID(c echo.Context) error {
   id := c.Param("id")
   if id == "" {
       return c.JSON(http.StatusBadRequest, 
           response.Error("ID inválido", nil))
   }

   mediaID, err := strconv.ParseInt(id, 10, 64)
   if err != nil {
       return c.JSON(http.StatusBadRequest, 
           response.Error("ID inválido", err.Error()))
   }

   media, err := h.service.GetByID(c.Request().Context(), mediaID)
   if err != nil {
       return h.handleError(c, err)
   }

   resp := MediaResponse{
       ID:        media.ID,
       FileName:  media.FileName,
       FilePath:  media.FilePath,
       MimeType:  media.MimeType,
       Size:      media.Size,
       CreatedAt: media.CreatedAt,
       UpdatedAt: media.UpdatedAt,
   }

   return c.JSON(http.StatusOK, 
       response.Success("Archivo obtenido exitosamente", resp))
}

func (h *Handler) List(c echo.Context) error {
	p := Pagination{Page: 1, PerPage: 10}
	if err := c.Bind(&p); err != nil {
			return c.JSON(http.StatusBadRequest, response.Error("Error en paginación", err.Error()))
	}

	medias, total, err := h.service.List(c.Request().Context(), &p)
	if err != nil {
			return h.handleError(c, err)
	}

	data := map[string]interface{}{
			"items": medias,
			"total": total,
			"page": p.Page,
			"per_page": p.PerPage,
	}

	return c.JSON(http.StatusOK, 
        response.Success("Archivos obtenidos exitosamente", data))
}

func (h *Handler) Delete(c echo.Context) error {
   id := c.Param("id")
   if id == "" {
       return c.JSON(http.StatusBadRequest, 
           response.Error("ID inválido", nil))
   }

   mediaID, err := strconv.ParseInt(id, 10, 64)
   if err != nil {
       return c.JSON(http.StatusBadRequest, 
           response.Error("ID inválido", err.Error()))
   }

   if err := h.service.Delete(c.Request().Context(), mediaID); err != nil {
       return h.handleError(c, err)
   }

   return c.JSON(http.StatusOK, 
       response.Success("Archivo eliminado exitosamente", nil))
}

func (h *Handler) Update(c echo.Context) error {
    id := c.Param("id")
    if id == "" {
        return c.JSON(http.StatusBadRequest, response.Error("ID inválido", nil))
    }

    mediaID, err := strconv.ParseInt(id, 10, 64)
    if err != nil {
        return c.JSON(http.StatusBadRequest, response.Error("ID inválido", err.Error()))
    }
	

     // Decodificar el cuerpo de la solicitud (JSON)
    req := new(UpdateMediaRequest)
    if err := c.Bind(req); err != nil {
        return c.JSON(http.StatusBadRequest, response.Error("JSON inválido", err.Error()))
    }
		// Validar si hay un archivo en el formulario (opcional)
    file, err := c.FormFile("file")
		// Validar el tipo de archivo
		contentType := file.Header.Get("Content-Type")
		if !isValidMimeType(contentType) {
				return c.JSON(http.StatusBadRequest, 
						response.Error("Tipo de archivo no PERMITIDO", nil))
		}
    if err == nil {
        src, err := file.Open()
        if err != nil {
            return h.handleError(c, err)
        }
        defer src.Close()

        fileBytes := make([]byte, file.Size)
        _, err = src.Read(fileBytes)
        if err != nil {
            return h.handleError(c, err)
        }

        fileName := file.Filename
        mimeType := file.Header.Get("Content-Type")

        req.File = fileBytes
        req.FileName = &fileName
        req.MimeType = &mimeType
    }

    // Verificar si se envió un nuevo nombre de archivo
    newFileName := c.FormValue("fileName")
    if newFileName != "" {
        req.FileName = &newFileName
    }

   // Actualizar media
    media, err := h.service.Update(c.Request().Context(), mediaID, req)
    if err != nil {
        return h.handleError(c, err)
    }

    resp := MediaResponse{
        ID:        media.ID,
        FileName:  media.FileName,
        FilePath:  media.FilePath,
        MimeType:  media.MimeType,
        Size:      media.Size,
        CreatedAt: media.CreatedAt,
        UpdatedAt: media.UpdatedAt,
    }


    return c.JSON(http.StatusOK, response.Success("Archivo actualizado exitosamente", resp))
}

func (h *Handler) handleError(c echo.Context, err error) error {
   switch e := err.(type) {
   case *errors.AppError:
       return c.JSON(e.Code, response.Error(e.Message, nil))
   default:
       if errors.IsNotFound(err) {
           return c.JSON(http.StatusNotFound, 
               response.Error("Archivo no encontrado", nil))
       }
       return c.JSON(http.StatusInternalServerError, 
           response.Error("Error interno del servidor", nil))
   }
}

