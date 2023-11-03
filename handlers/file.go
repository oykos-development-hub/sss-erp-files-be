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

func (h *fileHandlerImpl) CreateFile(w http.ResponseWriter, r *http.Request) {
	maxFileSize := int64(100 * 1024 * 1024) // file maximum 100 MB

	err := r.ParseMultipartForm(maxFileSize)
	if err != nil {
		//http.Error(w, "File is not valid!", http.StatusBadRequest)
		//return
		response := dto.FileResponse{
			Status: "failed",
		}
		_ = h.App.WriteDataResponse(w, http.StatusBadRequest, "File is not valid", response)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		//http.Error(w, "Error during fetching file!", http.StatusBadRequest)
		//return
		response := dto.FileResponse{
			Status: "failed",
		}
		_ = h.App.WriteDataResponse(w, http.StatusBadRequest, "Error during fetching file", response)
		return
	}
	defer file.Close()

	uploadDir := "./files"

	fileName := generateUniqueFileName(header.Filename)

	uploadedFile, err := os.Create(filepath.Join(uploadDir, fileName))
	if err != nil {
		//http.Error(w, "Error during creating file!", http.StatusBadRequest)
		//return
		response := dto.FileResponse{
			Status: "failed",
		}
		_ = h.App.WriteDataResponse(w, http.StatusBadRequest, "Error during creating file", response)
		return
	}
	defer uploadedFile.Close()

	_, err = io.Copy(uploadedFile, file)
	if err != nil {
		//http.Error(w, "Error during uploading file!", http.StatusBadRequest)
		//return
		response := dto.FileResponse{
			Status: "failed",
		}
		_ = h.App.WriteDataResponse(w, http.StatusBadRequest, "Error during uploading file", response)
		return
	}

	var input dto.FileDTO

	fileInfo, err := os.Stat(uploadedFile.Name())
	if err != nil {
		//http.Error(w, "Error during fetching file stats!", http.StatusBadRequest)
		//return
		response := dto.FileResponse{
			Status: "failed",
		}
		_ = h.App.WriteDataResponse(w, http.StatusBadRequest, "Error during fetching file stats", response)
		return
	}

	ext := filepath.Ext(header.Filename)

	input.Name = fileName
	input.Size = fileInfo.Size()
	input.Type = &ext
	res, err := h.service.CreateFile(input)
	if err != nil {
		//_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		//return
		response := dto.FileResponse{
			Status: "failed",
		}
		_ = h.App.WriteDataResponse(w, http.StatusBadRequest, "Error during saving file at database", response)
		return
	}

	response := dto.FileResponse{
		Data:   res,
		Status: "success",
	}

	_ = h.App.WriteDataResponse(w, http.StatusOK, "File created successfuly", response)
}

func (h *fileHandlerImpl) DeleteFile(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	res, err := h.service.GetFile(id)
	if err != nil {
		//_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		//return

		response := dto.FileResponse{
			Status: "failed",
		}
		_ = h.App.WriteDataResponse(w, http.StatusBadRequest, "File not found", response)
		return
	}

	err = os.Remove("./files/" + res.Name)

	if err != nil {
		//_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		//return
		response := dto.FileResponse{
			Status: "failed",
		}
		_ = h.App.WriteDataResponse(w, http.StatusBadRequest, "Error during deleting file", response)
		return
	}

	err = h.service.DeleteFile(id)
	if err != nil {
		//_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		//return
		response := dto.FileResponse{
			Status: "failed",
		}
		_ = h.App.WriteDataResponse(w, http.StatusBadRequest, "Error during deleting file from database", response)
		return
	}

	_ = h.App.WriteSuccessResponse(w, http.StatusOK, "File deleted successfuly")
}

func (h *fileHandlerImpl) GetFileById(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	res, err := h.service.GetFile(id)
	if err != nil {
		//_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		//return
		response := dto.FileResponse{
			Status: "failed",
		}
		_ = h.App.WriteDataResponse(w, http.StatusBadRequest, "File not found", response)
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

func (h *fileHandlerImpl) GetFile(w http.ResponseWriter, r *http.Request) {
	filePath := "./files" + r.URL.Path[len("/api/download"):]

	// Proverite da li fajl postoji
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		//http.NotFound(w, r)
		//return
		response := dto.FileResponse{
			Status: "failed",
		}
		_ = h.App.WriteDataResponse(w, http.StatusBadRequest, "File not found", response)
		return
	}

	http.ServeFile(w, r, filePath)
}

func (h *fileHandlerImpl) FileOverview(w http.ResponseWriter, r *http.Request) {
	fileId := chi.URLParam(r, "id")

	id, err := strconv.Atoi(fileId)
	if err != nil {
		response := dto.FileResponse{
			Status: "failed",
		}
		_ = h.App.WriteDataResponse(w, http.StatusBadRequest, "ID is not valid number", response)
		return
	}

	data, err := h.service.GetFile(id)
	if err != nil {
		response := dto.FileResponse{
			Status: "failed",
		}
		_ = h.App.WriteDataResponse(w, http.StatusBadRequest, "File not found", response)
		return
	}

	filePath := "./files/" + data.Name

	file, err := os.Open(filePath)
	if err != nil {
		response := dto.FileResponse{
			Status: "failed",
		}
		_ = h.App.WriteDataResponse(w, http.StatusBadRequest, "File not exists", response)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Disposition", "attachment; filename="+data.Name)
	w.Header().Set("Content-Type", "application/octet-stream")

	_, err = io.Copy(w, file)
	if err != nil {
		response := dto.FileResponse{
			Status: "failed",
		}
		_ = h.App.WriteDataResponse(w, http.StatusBadRequest, "Error during reading file", response)
		return
	}
}
