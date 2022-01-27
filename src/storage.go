package metabus

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/gorilla/mux"
)

type Storage struct {
	m    *sync.Mutex
	root string
}

func CreateStorage() *Storage {
	root := "./storage"
	os.Mkdir(root, 0755)
	return &Storage{
		m:    &sync.Mutex{},
		root: root,
	}
}

func (this *Storage) FuncHandler(w http.ResponseWriter, r *http.Request) {
	operation := r.FormValue("operation")
	path := r.FormValue("path")
	newPath := r.FormValue("newPath")

	switch operation {
	case "show":
		this.ListFolder(w, path)
	case "mkdir":
		this.CreateFolder(w, path)
	case "remove":
		this.Remove(w, path)
	case "add":
		this.AddFile(w, r, path)
	case "move":
		this.MoveFile(w, path, newPath)
	default:
		ErrorResponse(w, "Unknown operation")
	}
}

func ErrorResponse(w http.ResponseWriter, text string) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(text))
}

func JsonResponse(w http.ResponseWriter, json []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func OkResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (this *Storage) ServeFiles(w http.ResponseWriter, r *http.Request) {
	maybePath, ok := mux.Vars(r)["path"]
	if !ok {
		ErrorResponse(w, "Path is not defined")
		return
	}
	path := strings.ToLower(maybePath)
	path = filepath.Join(this.root, path)
	// path := filepath.Join(this.root, maybePath)
	http.ServeFile(w, r, path)
}

func (this *Storage) CreateFolder(w http.ResponseWriter, folderName string) {
	this.m.Lock()
	defer this.m.Unlock()
	_, err := os.Stat(folderName)
	if os.IsNotExist(err) {
		errDir := os.Mkdir(filepath.Join(this.root, folderName), 0755)
		if errDir != nil {
			ErrorResponse(w, errDir.Error())
		}
	}
}

func (this *Storage) AddFile(w http.ResponseWriter, r *http.Request, path string) {
	this.m.Lock()
	defer this.m.Unlock()

	err := r.ParseMultipartForm(100 << 20)
	if err != nil {
		ErrorResponse(w, err.Error())
	}

	m := r.MultipartForm
	files := m.File["file"]

	if len(files) == 0 {
		ErrorResponse(w, err.Error())
	}
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		defer file.Close()
		if err != nil {
			file.Close()
			ErrorResponse(w, err.Error())
		}
		fullFileName := strings.ToLower(fileHeader.Filename)
		dst, err := os.Create(filepath.Join(this.root, path, fullFileName))
		defer dst.Close()
		if err != nil {
			ErrorResponse(w, err.Error())
			return
		}
		if _, err := io.Copy(dst, file); err != nil {
			ErrorResponse(w, err.Error())
			return
		}
	}
}

func (this *Storage) MoveFile(w http.ResponseWriter, oldPath string, newPath string) {
	this.m.Lock()
	defer this.m.Unlock()
	err := os.Rename(filepath.Join(this.root, oldPath), filepath.Join(this.root, newPath))
	if err != nil {
		ErrorResponse(w, err.Error())
	}
}

func (this *Storage) ListFolder(w http.ResponseWriter, path string) {
	type Element struct {
		Path     string `json:"path"`
		Name     string `json:"name"`
		IsFolder bool   `json:"isfolder"`
	}
	elements := []Element{}

	fileInfo, err := ioutil.ReadDir(filepath.Join(this.root, path))
	if err != nil {
		ErrorResponse(w, err.Error())
		return
	}
	for i := range fileInfo {
		elements = append(elements, Element{
			Path:     filepath.Join(path, fileInfo[i].Name()),
			Name:     fileInfo[i].Name(),
			IsFolder: fileInfo[i].IsDir(),
		})
	}

	data, err := json.Marshal(elements)
	if err != nil {
		ErrorResponse(w, err.Error())
		return
	}

	JsonResponse(w, data)

}

func (this *Storage) Remove(w http.ResponseWriter, target string) {
	this.m.Lock()
	defer this.m.Unlock()
	err := os.RemoveAll(filepath.Join(this.root, target))
	if err != nil {
		ErrorResponse(w, err.Error())
		return
	}
	OkResponse(w)
}
