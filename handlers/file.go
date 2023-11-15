package handlers

import (
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"gitlab.sudovi.me/erp/file-ms-api/dto"
	"gitlab.sudovi.me/erp/file-ms-api/errors"
	"gitlab.sudovi.me/erp/file-ms-api/services"

	"github.com/360EntSecGroup-Skylar/excelize"
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
	maxFileSize := int64(100 * 1024 * 1024) // Maksimalna veličina fajla je 100 MB

	err := r.ParseMultipartForm(maxFileSize)
	if err != nil {
		response := dto.FileResponse{
			Status: "failed",
		}
		_ = h.App.WriteDataResponse(w, http.StatusBadRequest, "File is not valid", response)
		return
	}

	files := r.MultipartForm.File["file"]

	var filesResponse []*dto.FileResponseDTO

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			response := dto.FileResponse{
				Status: "failed",
			}
			_ = h.App.WriteDataResponse(w, http.StatusBadRequest, "Error during fetching file", response)
			return
		}
		defer file.Close()

		uploadDir := "./files"

		fileName := generateUniqueFileName(fileHeader.Filename)

		uploadedFile, err := os.Create(filepath.Join(uploadDir, fileName))
		if err != nil {
			response := dto.FileResponse{
				Status: "failed",
			}
			_ = h.App.WriteDataResponse(w, http.StatusBadRequest, "Error during creating file", response)
			return
		}
		defer uploadedFile.Close()

		_, err = io.Copy(uploadedFile, file)
		if err != nil {
			response := dto.FileResponse{
				Status: "failed",
			}
			_ = h.App.WriteDataResponse(w, http.StatusBadRequest, "Error during uploading file", response)
			return
		}

		var input dto.FileDTO

		fileInfo, err := os.Stat(uploadedFile.Name())
		if err != nil {
			response := dto.FileResponse{
				Status: "failed",
			}
			_ = h.App.WriteDataResponse(w, http.StatusBadRequest, "Error during fetching file stats", response)
			return
		}

		ext := filepath.Ext(fileHeader.Filename)

		input.Name = fileName
		input.Size = fileInfo.Size()
		input.Type = &ext
		res, err := h.service.CreateFile(input)
		if err != nil {
			response := dto.FileResponse{
				Status: "failed",
			}
			_ = h.App.WriteDataResponse(w, http.StatusBadRequest, "Error during saving file at database", response)
			return
		}
		filesResponse = append(filesResponse, res)
	}

	response := dto.MultipleFileResponse{
		Status: "success",
		Data:   filesResponse,
	}
	_ = h.App.WriteDataResponse(w, http.StatusOK, "Files created successfully", response)
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

	response := dto.FileResponse{
		Status: "success",
	}

	_ = h.App.WriteDataResponse(w, http.StatusOK, "File deleted successfuly", response)
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

func (h *fileHandlerImpl) MultipleDeleteFile(w http.ResponseWriter, r *http.Request) {
	var input dto.MultipleDeleteFiles
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

	for _, id := range input.Files {
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
		//_ = h.App.WriteErrorResponse(w, errors.MapErrorToStatusCode(err), err)
		//return
		response := dto.FileResponse{
			Status: "failed",
		}
		_ = h.App.WriteDataResponse(w, http.StatusBadRequest, "File not found", response)
		return
	}

	filePath := "./files/" + res.Name

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

func (h *fileHandlerImpl) TemplateUpload(w http.ResponseWriter, r *http.Request) {
	maxFileSize := int64(100 * 1024 * 1024) // Maksimalna veličina fajla je 100 MB

	err := r.ParseMultipartForm(maxFileSize)
	if err != nil {
		response := dto.FileResponse{
			Status: "failed",
		}
		_ = h.App.WriteDataResponse(w, http.StatusBadRequest, "File is not valid", response)
		return
	}

	files := r.MultipartForm.File["file"]

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			response := dto.FileResponse{
				Status: "failed",
			}
			_ = h.App.WriteDataResponse(w, http.StatusBadRequest, "Error during fetching file", response)
			return
		}
		defer file.Close()

		uploadDir := "./templates"

		uploadedFile, err := os.Create(uploadDir + "/" + fileHeader.Filename)
		if err != nil {
			response := dto.FileResponse{
				Status: "failed",
			}
			_ = h.App.WriteDataResponse(w, http.StatusBadRequest, "Error during creating file", response)
			return
		}
		defer uploadedFile.Close()

		_, err = io.Copy(uploadedFile, file)
		if err != nil {
			response := dto.FileResponse{
				Status: "failed",
			}
			_ = h.App.WriteDataResponse(w, http.StatusBadRequest, "Error during uploading file", response)
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
		response := dto.FileResponse{
			Status: "failed",
		}
		_ = h.App.WriteDataResponse(w, http.StatusBadRequest, "File not exists", response)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
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

func (h *fileHandlerImpl) ReadArticles(w http.ResponseWriter, r *http.Request) {
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

	file, _, err := r.FormFile("file")
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

	procurementID := r.FormValue("public_procurement_id")

	publicProcurementID, err := strconv.Atoi(procurementID)

	if err != nil {
		response := dto.FileResponse{
			Status: "failed",
		}
		_ = h.App.WriteDataResponse(w, http.StatusBadRequest, "You must provide valid public_procurement_id", response)
		return
	}

	// Sačuvajte fajl na disku
	tempFile, err := os.CreateTemp("", "uploaded-file-")
	if err != nil {
		response := dto.FileResponse{
			Status: "failed",
		}
		_ = h.App.WriteDataResponse(w, http.StatusInternalServerError, "Error during opening file", response)
		return
	}
	defer tempFile.Close()

	_, err = io.Copy(tempFile, file)
	if err != nil {
		response := dto.FileResponse{
			Status: "failed",
		}
		_ = h.App.WriteDataResponse(w, http.StatusInternalServerError, "Error during reading file", response)
		return
	}

	// Sada možete otvoriti sačuvani fajl koristeći putanju do njega
	xlsFile, err := excelize.OpenFile(tempFile.Name())

	if err != nil {
		//http.Error(w, "Error during fetching file!", http.StatusBadRequest)
		//return
		response := dto.FileResponse{
			Status: "failed",
		}
		_ = h.App.WriteDataResponse(w, http.StatusBadRequest, "Error during opening file", response)
		return
	}

	// Prolazak kroz listu i čitanje podataka
	var articles []dto.Article

	// Pristupanje listama u Excel fajlu
	sheetMap := xlsFile.GetSheetMap()

	for _, sheetName := range sheetMap {
		if sheetName != "Stavke" {
			continue
		}

		rows, err := xlsFile.Rows(sheetName)
		if err != nil {
			response := dto.FileResponse{
				Status: "failed",
			}
			_ = h.App.WriteDataResponse(w, http.StatusBadRequest, "Error during reading file rows!", response)
			return
		}

		rowindex := 0

		for rows.Next() {
			if rowindex == 0 {
				rowindex++
				continue
			}

			cols := rows.Columns()
			if err != nil {
				response := dto.FileResponse{
					Status: "failed",
				}
				_ = h.App.WriteDataResponse(w, http.StatusBadRequest, "Error during reading column value", response)
				return
			}

			var article dto.Article
			for cellIndex, cellValue := range cols {
				value := cellValue
				switch cellIndex {
				case 0:
					article.Title = value
				case 1:
					article.Description = value
				case 2:
					if value == "" {
						break
					}

					floatValue, err := strconv.ParseFloat(value, 32)

					if err != nil {
						response := dto.FileResponse{
							Status: "failed",
						}
						_ = h.App.WriteDataResponse(w, http.StatusBadRequest, "Error during converting neto price", response)
						return
					}
					article.NetPrice = float32(floatValue)
				case 3:
					if value == "" {
						break
					}

					floatValue, err := strconv.ParseFloat(value, 32)

					if err != nil {
						response := dto.FileResponse{
							Status: "failed",
						}
						_ = h.App.WriteDataResponse(w, http.StatusBadRequest, "Error during converting neto price", response)
						return
					}

					vatPercentage := 100 * floatValue / float64(article.NetPrice)
					round := math.Round(vatPercentage)

					valueVat := strconv.Itoa(int(round))

					article.VatPercentage = valueVat
				}
			}

			article.PublicProcurementID = publicProcurementID

			if article.Title == "" || article.NetPrice == 0 || article.VatPercentage == "" {
				break
			}

			articles = append(articles, article)
		}

	}

	response := dto.ArticleResponse{
		Data:   articles,
		Status: "success",
	}

	_ = h.App.WriteDataResponse(w, http.StatusOK, "File created successfuly", response)
}
