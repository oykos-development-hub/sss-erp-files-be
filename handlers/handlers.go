package handlers

import (
	"net/http"
)

type Handlers struct {
	FileHandler FileHandler
}

type FileHandler interface {
	CreateFile(w http.ResponseWriter, r *http.Request)
	DeleteFile(w http.ResponseWriter, r *http.Request)
	MultipleDeleteFile(w http.ResponseWriter, r *http.Request)
	GetFileById(w http.ResponseWriter, r *http.Request)
	GetFile(w http.ResponseWriter, r *http.Request)
	FileOverview(w http.ResponseWriter, r *http.Request)
	TemplateUpload(w http.ResponseWriter, r *http.Request)
	TemplateDownload(w http.ResponseWriter, r *http.Request)
	ReadArticles(w http.ResponseWriter, r *http.Request)
	ReadSimpleArticles(w http.ResponseWriter, r *http.Request)
}
