package services
import (
	"gitlab.sudovi.me/erp/file-ms-api/dto"
)


type BaseService interface {
	RandomString(n int) string
	Encrypt(text string) (string, error)
	Decrypt(crypto string) (string, error)
}

type FileService interface {
	CreateFile(input dto.FileDTO) (*dto.FileResponseDTO, error)
	UpdateFile(id int, input dto.FileDTO) (*dto.FileResponseDTO, error)
	DeleteFile(id int) error
	GetFile(id int) (*dto.FileResponseDTO, error)
	GetFileList() ([]dto.FileResponseDTO, error)
}
