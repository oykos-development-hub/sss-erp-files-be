package dto

import (
	"time"

	"gitlab.sudovi.me/erp/file-ms-api/data"
)

type FileDTO struct {
	ParentID    *int    `json:"parent_id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	Size        int     `json:"size"`
	Type        string  `json:"type"`
}

type FileResponseDTO struct {
	ID          int       `json:"id"`
	ParentID    *int      `json:"parent_id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	Size        int       `json:"size"`
	Type        string    `json:"type"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (dto FileDTO) ToFile() *data.File {
	return &data.File{
		ParentID:    dto.ParentID,
		Name:        dto.Name,
		Description: dto.Description,
		Size:        dto.Size,
		Type:        dto.Type,
	}
}

func ToFileResponseDTO(data data.File) FileResponseDTO {
	return FileResponseDTO{
		ID:          data.ID,
		ParentID:    data.ParentID,
		Name:        data.Name,
		Description: data.Description,
		Size:        data.Size,
		Type:        data.Type,
		CreatedAt:   data.CreatedAt,
		UpdatedAt:   data.UpdatedAt,
	}
}

func ToFileListResponseDTO(files []*data.File) []FileResponseDTO {
	dtoList := make([]FileResponseDTO, len(files))
	for i, x := range files {
		dtoList[i] = ToFileResponseDTO(*x)
	}
	return dtoList
}
