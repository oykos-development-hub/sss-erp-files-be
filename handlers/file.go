package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"gitlab.sudovi.me/erp/file-ms-api/dto"
	newErrors "gitlab.sudovi.me/erp/file-ms-api/pkg/errors"
	"gitlab.sudovi.me/erp/file-ms-api/services"

	"github.com/go-chi/chi/v5"
	"github.com/oykos-development-hub/celeritas"
)

// FileHandler is a concrete type that implements FileHandler
type fileHandlerImpl struct {
	App     *celeritas.Celeritas
	service services.FileService
}

// NewFileHandler initializes a new FileHandler with its dependencies
func NewFileHandler(app *celeritas.Celeritas, fileService services.FileService) FileHandler {
	return &fileHandlerImpl{
		App:     app,
		service: fileService,
	}
}

func (h *fileHandlerImpl) handleError(w http.ResponseWriter, err error, statusCode int) {
	h.App.ErrorLog.Print(err)
	w.WriteHeader(statusCode)
	_ = MarshalAndWriteJSON(w, dto.ErrorResponse{Message: err.Error()})
}

func MarshalAndWriteJSON(w http.ResponseWriter, obj interface{}) error {
	jsonResponse, err := json.Marshal(obj)
	if err != nil {
		//http.Error(w, "Error during JSON marshaling", http.StatusInternalServerError)
		return newErrors.Wrap(err, "json marshal")
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonResponse)

	if err != nil {
		return newErrors.Wrap(err, "w write")
	}

	return nil
}

func (h *fileHandlerImpl) CreateFile(w http.ResponseWriter, r *http.Request) {
	maxFileSize := int64(100 * 1024 * 1024) // Maksimalna veličina fajla je 100 MB

	err := r.ParseMultipartForm(maxFileSize)
	if err != nil {
		h.handleError(w, err, http.StatusBadRequest)
		return
	}

	files := r.MultipartForm.File["file"]

	var filesResponse []*dto.FileResponseDTO

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			h.handleError(w, err, http.StatusBadRequest)
			return
		}
		defer file.Close()

		uploadDir := "./files"

		fileName := generateUniqueFileName(fileHeader.Filename)

		uploadedFile, err := os.Create(filepath.Join(uploadDir, fileName))
		if err != nil {
			h.handleError(w, err, http.StatusBadRequest)
			return
		}
		defer uploadedFile.Close()

		_, err = io.Copy(uploadedFile, file)
		if err != nil {
			h.handleError(w, err, http.StatusBadRequest)
			return
		}

		var input dto.FileDTO

		fileInfo, err := os.Stat(uploadedFile.Name())
		if err != nil {
			h.handleError(w, err, http.StatusBadRequest)
			return
		}

		ext := filepath.Ext(fileHeader.Filename)

		input.Name = fileName
		input.Size = fileInfo.Size()
		input.Type = &ext
		res, err := h.service.CreateFile(input)
		if err != nil {
			h.handleError(w, err, http.StatusInternalServerError)
			return
		}
		filesResponse = append(filesResponse, res)
	}

	response := dto.MultipleFileResponse{
		Data:   filesResponse,
		Status: "success",
	}

	_ = MarshalAndWriteJSON(w, response)
}

func (h *fileHandlerImpl) DeleteFile(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	res, err := h.service.GetFile(id)
	if err != nil {
		h.handleError(w, err, http.StatusBadRequest)
		return
	}

	err = os.Remove("./files/" + res.Name)

	if err != nil {
		h.handleError(w, err, http.StatusBadRequest)
		return
	}

	err = h.service.DeleteFile(id)
	if err != nil {
		h.handleError(w, err, http.StatusBadRequest)
		return
	}

	response := dto.FileResponse{
		Status: "success",
	}

	_ = h.App.WriteDataResponse(w, http.StatusOK, "File deleted successfuly", response)
}

func (h *fileHandlerImpl) GetFileById(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	res, err := h.service.GetFile(id)
	if err != nil {
		h.handleError(w, err, http.StatusBadRequest)
		return
	}

	response := dto.FileResponse{
		Data:   res,
		Status: "success",
	}

	_ = h.App.WriteDataResponse(w, http.StatusOK, "", response)
}

func generateUniqueFileName(filePath string) string {
	source := rand.NewSource(time.Now().UnixNano())
	generator := rand.New(source)

	randomNum := generator.Int31()

	baseName := filepath.Base(filePath)
	ext := filepath.Ext(baseName)
	fileNameWithoutExt := baseName[:len(baseName)-len(ext)]

	uniqueFileName := fmt.Sprintf("%s_%d%s", fileNameWithoutExt, randomNum, ext)

	return uniqueFileName
}

func (h *fileHandlerImpl) MultipleDeleteFile(w http.ResponseWriter, r *http.Request) {
	var input dto.MultipleDeleteFiles
	err := h.App.ReadJSON(w, r, &input)
	if err != nil {
		h.handleError(w, err, http.StatusBadRequest)
		return
	}

	validator := h.App.Validator().ValidateStruct(&input)
	if !validator.Valid() {
		h.handleError(w, err, http.StatusBadRequest)
		return
	}

	for _, id := range input.Files {
		res, err := h.service.GetFile(id)
		if err != nil {
			h.handleError(w, err, http.StatusBadRequest)
			return
		}

		err = os.Remove("./files/" + res.Name)

		if err != nil {
			h.handleError(w, err, http.StatusBadRequest)
			return
		}

		err = h.service.DeleteFile(id)
		if err != nil {
			h.handleError(w, err, http.StatusBadRequest)
			return
		}

	}

	response := dto.FileResponse{
		Status: "success",
	}

	_ = h.App.WriteDataResponse(w, http.StatusOK, "Files deleted successfuly", response)
}

func (h *fileHandlerImpl) GetFile(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	res, err := h.service.GetFile(id)
	if err != nil {
		h.handleError(w, err, http.StatusBadRequest)
		return
	}

	filePath := "./files/" + res.Name

	// Proverite da li fajl postoji
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		h.handleError(w, err, http.StatusBadRequest)
		return
	}

	http.ServeFile(w, r, filePath)
}

func (h *fileHandlerImpl) FileOverview(w http.ResponseWriter, r *http.Request) {
	fileId := chi.URLParam(r, "id")

	id, err := strconv.Atoi(fileId)
	if err != nil {
		h.handleError(w, err, http.StatusBadRequest)
		return
	}

	data, err := h.service.GetFile(id)
	if err != nil {
		h.handleError(w, err, http.StatusBadRequest)
		return
	}

	filePath := "./files/" + data.Name

	file, err := os.Open(filePath)
	if err != nil {
		h.handleError(w, err, http.StatusBadRequest)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Disposition", "attachment; filename="+data.Name)
	w.Header().Set("Content-Type", "application/octet-stream")

	_, err = io.Copy(w, file)
	if err != nil {
		h.handleError(w, err, http.StatusBadRequest)
		return
	}
}

func (h *fileHandlerImpl) TemplateUpload(w http.ResponseWriter, r *http.Request) {
	maxFileSize := int64(100 * 1024 * 1024) // Maksimalna veličina fajla je 100 MB

	err := r.ParseMultipartForm(maxFileSize)
	if err != nil {
		h.handleError(w, err, http.StatusBadRequest)
		return
	}

	files := r.MultipartForm.File["file"]

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			h.handleError(w, err, http.StatusBadRequest)
			return
		}
		defer file.Close()

		uploadDir := "./templates"

		uploadedFile, err := os.Create(uploadDir + "/" + fileHeader.Filename)
		if err != nil {
			h.handleError(w, err, http.StatusBadRequest)
			return
		}
		defer uploadedFile.Close()

		_, err = io.Copy(uploadedFile, file)
		if err != nil {
			h.handleError(w, err, http.StatusBadRequest)
			return
		}
	}

	response := dto.MultipleFileResponse{
		Status: "success",
	}
	_ = h.App.WriteDataResponse(w, http.StatusOK, "Files created successfully", response)
}

func (h *fileHandlerImpl) TemplateDownload(w http.ResponseWriter, r *http.Request) {
	fileName := chi.URLParam(r, "*")

	filePath := "./templates/" + fileName

	file, err := os.Open(filePath)
	if err != nil {
		h.handleError(w, err, http.StatusBadRequest)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.Header().Set("Content-Type", "application/octet-stream")

	_, err = io.Copy(w, file)
	if err != nil {
		h.handleError(w, err, http.StatusBadRequest)
		return
	}
}
