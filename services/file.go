package services

import (
	"gitlab.sudovi.me/erp/file-ms-api/data"
	"gitlab.sudovi.me/erp/file-ms-api/dto"
	newErrors "gitlab.sudovi.me/erp/file-ms-api/pkg/errors"

	"github.com/oykos-development-hub/celeritas"
)

type FileServiceImpl struct {
	App  *celeritas.Celeritas
	repo data.File
}

func NewFileServiceImpl(app *celeritas.Celeritas, repo data.File) FileService {
	return &FileServiceImpl{
		App:  app,
		repo: repo,
	}
}

func (h *FileServiceImpl) CreateFile(input dto.FileDTO) (*dto.FileResponseDTO, error) {
	data := input.ToFile()

	id, err := h.repo.Insert(*data)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo file insert")
	}

	data, err = data.Get(id)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo file get")
	}

	res := dto.ToFileResponseDTO(*data)

	return &res, nil
}

func (h *FileServiceImpl) UpdateFile(id int, input dto.FileDTO) (*dto.FileResponseDTO, error) {
	data := input.ToFile()
	data.ID = id

	err := h.repo.Update(*data)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo file update")
	}

	data, err = h.repo.Get(id)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo file get")
	}

	response := dto.ToFileResponseDTO(*data)

	return &response, nil
}

func (h *FileServiceImpl) DeleteFile(id int) error {
	err := h.repo.Delete(id)
	if err != nil {
		return newErrors.Wrap(err, "repo file delete")
	}

	return nil
}

func (h *FileServiceImpl) GetFile(id int) (*dto.FileResponseDTO, error) {
	data, err := h.repo.Get(id)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo file get")
	}
	response := dto.ToFileResponseDTO(*data)

	return &response, nil
}

func (h *FileServiceImpl) GetFileList() ([]dto.FileResponseDTO, error) {
	data, err := h.repo.GetAll(nil)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo file get all")
	}
	response := dto.ToFileListResponseDTO(data)

	return response, nil
}
