package handlers
import (
	"net/http"
)


type Handlers struct {
	FileHandler FileHandler
	}

type FileHandler interface {
	CreateFile(w http.ResponseWriter, r *http.Request)
	UpdateFile(w http.ResponseWriter, r *http.Request)
	DeleteFile(w http.ResponseWriter, r *http.Request)
	GetFileById(w http.ResponseWriter, r *http.Request)
	GetFileList(w http.ResponseWriter, r *http.Request)
}
