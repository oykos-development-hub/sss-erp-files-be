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
	GetFileById(w http.ResponseWriter, r *http.Request)
	GetFile(w http.ResponseWriter, r *http.Request)
	FileOverview(w http.ResponseWriter, r *http.Request)
}
