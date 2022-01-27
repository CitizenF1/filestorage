package metabus

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
)

const (
	MB = 1 << 20
)

type Sizer interface {
	Size() int64
}

type TCM struct {
	m    *sync.Mutex
	root string
}

func CreateTCM() *TCM {
	root := "./tcm"
	os.Mkdir(root, 0755)
	return &TCM{
		m:    &sync.Mutex{},
		root: root,
	}
}

func (t *TCM) UploadHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	// case "GET":
	// 	display(w, "upload", nil)
	case "POST":
		t.UploadFile(w, r)
	}
}

func (t *TCM) UploadFile(w http.ResponseWriter, r *http.Request) {

	b, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(b)

}

// func (t *TCM) UploadFile(w http.ResponseWriter, r *http.Request) {
// 	// Maximum upload of 10 MB files
// 	r.ParseMultipartForm(100 << 20)

// 	// Get handler for filename, size and headers
// 	file, handler, err := r.FormFile("photoReport")
// 	if err != nil {
// 		fmt.Println("Error Retrieving the File")
// 		fmt.Println(err)
// 		return
// 	}

// 	defer file.Close()
// 	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
// 	fmt.Printf("File Size: %+v\n", handler.Size)
// 	fmt.Printf("MIME Header: %+v\n", handler.Header)

// 	// Create file
// 	dst, err := os.Create(filepath.Join(t.root, handler.Filename))
// 	defer dst.Close()
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	// Copy the uploaded file to the created file on the filesystem
// 	if _, err := io.Copy(dst, file); err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	path := filepath.Join(t.root, handler.Filename)
// 	fmt.Fprintf(w, "Successfully Uploaded File\n")
// 	w.Write([]byte(path))
// 	// http.ServeFile(w, r, path)
// 	// OkResponse(w)
// }

// func (t *TCM) UploadFile(w http.ResponseWriter, r *http.Request) {
// Maximum upload of 10 MB files
// b, err := io.ReadAll(r.Body)
// fmt.Println(b)
// // b, err := ioutil.ReadAll(resp.Body)  Go.1.15 and earlier
// if err != nil {
// 	log.Fatalln(err)
// }

// files := m.File["photoReport"]

// if len(files) == 0 {
// 	ErrorResponse(w, err.Error())
// }
// fmt.Println("foreach")

// for _, fileHeader := range files {
// 	file, err := fileHeader.Open()
// 	defer file.Close()
// 	if err != nil {
// 		file.Close()
// 		ErrorResponse(w, err.Error())
// 	}
// 	fullFileName := strings.ToLower(fileHeader.Filename)
// 	dst, err := os.Create(filepath.Join(t.root, fullFileName))
// 	defer dst.Close()
// 	if err != nil {
// 		ErrorResponse(w, err.Error())
// 		return
// 	}
// 	if _, err := io.Copy(dst, file); err != nil {
// 		ErrorResponse(w, err.Error())
// 		return
// 	}
// 	path := filepath.Join(t.root, fullFileName)
// 	w.Write([]byte(path))
// }

// Get handler for filename, size and headers
// file, handler, err := r.FormFile("photoReport")
// fmt.Println(file, handler)
// if err != nil {
// 	fmt.Println(file, handler, "========")
// 	fmt.Println("Error Retrieving the File")
// 	fmt.Println(err, "43")
// 	return
// }

// defer file.Close()
// fmt.Printf("Uploaded File: %+v\n", handler.Filename)
// fmt.Printf("File Size: %+v\n", handler.Size)
// fmt.Printf("MIME Header: %+v\n", handler.Header)

// // Create file
// dst, err := os.Create(filepath.Join(t.root, handler.Filename))
// defer dst.Close()
// if err != nil {
// 	http.Error(w, err.Error(), http.StatusInternalServerError)
// 	return
// }

// // Copy the uploaded file to the created file on the filesystem
// if _, err := io.Copy(dst, file); err != nil {
// 	http.Error(w, err.Error(), http.StatusInternalServerError)
// 	return
// }
// path := filepath.Join(t.root, handler.Filename)
// fmt.Fprintf(w, "Successfully Uploaded File\n")
// w.Write([]byte(path))
// http.ServeFile(w, r, path)
// OkResponse(w)
// }
