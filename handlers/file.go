package handlers

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"gitlab.sudovi.me/erp/file-ms-api/dto"
	"gitlab.sudovi.me/erp/file-ms-api/errors"
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

func (h *fileHandlerImpl) ShowFile(w http.ResponseWriter, r *http.Request) {
	fs := http.FileServer(http.Dir("./files"))
	http.StripPrefix("/file", fs).ServeHTTP(w, r)
}

func (h *fileHandlerImpl) CreateFile(w http.ResponseWriter, r *http.Request) {
	maxFileSize := int64(100 * 1024 * 1024) // file maximum 100 MB

	err := r.ParseMultipartForm(maxFileSize)
	if err != nil {
		http.Error(w, "File is not valid!", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error during fetching file!", http.StatusBadRequest)
		return
	}
	defer file.Close()

	uploadDir := "./files"

	fileName := generateUniqueFileName(filepath.Join(uploadDir, "file"))

	uploadedFile, err := os.Create(filepath.Join(uploadDir, fileName))
	if err != nil {
		http.Error(w, "Error during creating file!", http.StatusBadRequest)
		return
	}
	defer uploadedFile.Close()

	_, err = io.Copy(uploadedFile, file)
	if err != nil {
		http.Error(w, "Error during uploading file!", http.StatusBadRequest)
		return
	}

	var input dto.FileDTO

	fileInfo, err := os.Stat(uploadedFile.Name())
	if err != nil {
		http.Error(w, "Error during fetching file stats!", http.StatusBadRequest)
		return
	}

	input.Name = fileName
	input.Size = fileInfo.Size()
	//input.Type???????

	res, err := h.service.CreateFile(input)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponse(w, http.StatusOK, "File created successfuly", res)
}

func (h *fileHandlerImpl) UpdateFile(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	var input dto.FileDTO
	err := h.App.ReadJSON(w, r, &input)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	validator := h.App.Validator().ValidateStruct(&input)
	if !validator.Valid() {
		_ = h.App.WriteErrorResponseWithData(w, errors.MapErrorToStatusCode(errors.ErrBadRequest), errors.ErrBadRequest, validator.Errors)
		return
	}

	res, err := h.service.UpdateFile(id, input)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponse(w, http.StatusOK, "File updated successfuly", res)
}

func (h *fileHandlerImpl) DeleteFile(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	res, err := h.service.GetFile(id)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	err = os.Remove("./files/" + res.Name)

	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	err = h.service.DeleteFile(id)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteSuccessResponse(w, http.StatusOK, "File deleted successfuly")
}

func (h *fileHandlerImpl) GetFileById(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	res, err := h.service.GetFile(id)
	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponse(w, http.StatusOK, "", res)
}

func (h *fileHandlerImpl) GetFileList(w http.ResponseWriter, r *http.Request) {
	res, err := h.service.GetFileList()
	if err != nil {
		_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		return
	}

	_ = h.App.WriteDataResponse(w, http.StatusOK, "", res)
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
