package metabus

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/gorilla/mux"
)

type Addresseble struct {
	m    *sync.Mutex
	root string
}

func CreateAddresseble() *Addresseble {
	root := "./StandaloneWindows"
	os.Mkdir(root, 0755)
	return &Addresseble{
		m:    &sync.Mutex{},
		root: root,
	}
}

func (a *Addresseble) FuncAdressebleHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		a.uploadFile(w, r)
	default:
		ErrorResponse(w, "Unknown operation")
	}
}

func (a *Addresseble) uploadFile(w http.ResponseWriter, r *http.Request) {
	// Maximum upload of 10 MB files
	r.ParseMultipartForm(10 << 20)

	// Get handler for filename, size and headers
	file, handler, err := r.FormFile("file")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}

	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	// Create file
	dst, err := os.Create(filepath.Join(a.root, handler.Filename))
	if err != nil {
		log.Println("Error creating", err)
	}
	defer dst.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Copy the uploaded file to the created file on the filesystem
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Successfully Uploaded File\n")
}

func (a *Addresseble) ServeFilesAdress(w http.ResponseWriter, r *http.Request) {
	maybePath, ok := mux.Vars(r)["path"]

	fmt.Println(maybePath, "84")
	if !ok {
		ErrorResponse(w, "Path is not defined")
		return
	}
	switch r.Method {
	case "POST":
		a.uploadFile(w, r)
	case "GET":
		path := filepath.Join(a.root, maybePath)
		fmt.Println(path, "96")
		http.ServeFile(w, r, path)
	default:
		ErrorResponse(w, "Unknown operation")
	}
}
