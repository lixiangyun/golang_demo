package main

import (
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"runtime/debug"
)

const (
	UPLOAD_DIR = "./uploads"

	TEMPLATE_DIR = "./views"
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		err := renderHtml(w, "upload.html", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if r.Method == "POST" {
		f, h, err := r.FormFile("image")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		filename := h.Filename
		defer f.Close()

		t, err := os.Create(UPLOAD_DIR + "/" + filename)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer t.Close()

		if _, err := io.Copy(t, f); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/view?id="+filename, http.StatusFound)
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	imageId := r.FormValue("id")
	imagePath := UPLOAD_DIR + "/" + imageId

	if exists := isExists(imagePath); !exists {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "image")
	http.ServeFile(w, r, imagePath)
}

func isExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	fileInfoArr, err := ioutil.ReadDir("./uploads")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	locals := make(map[string]interface{})
	images := []string{}

	for _, fileInfo := range fileInfoArr {
		images = append(images, fileInfo.Name())
	}

	locals["images"] = images
	err = renderHtml(w, "list.html", locals)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func renderHtml(w http.ResponseWriter, tmp1 string, locals map[string]interface{}) error {
	err := templates[tmp1].Execute(w, locals)
	return err
}

var templates = make(map[string]*template.Template)

func init() {

	fileInfoArr, err := ioutil.ReadDir(TEMPLATE_DIR)
	if err != nil {
		panic(err)
		return
	}

	var templateName, templatePath string

	for _, fileinfo := range fileInfoArr {
		templateName = fileinfo.Name()
		if ext := path.Ext(templateName); ext != ".html" {
			continue
		}

		templatePath = TEMPLATE_DIR + "/" + templateName

		log.Println("Loading template:", templatePath)

		t := template.Must(template.ParseFiles(templatePath))
		templates[templateName] = t
	}
}

func safeHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if e, ok := recover().(error); ok {
				http.Error(w, e.Error(), http.StatusInternalServerError)
				// 可以自定义50x错误页面
				log.Println("Warning: Panic in %v - %v", fn, e)
				log.Println(string(debug.Stack()))
			}
		}()

		fn(w, r)
	}
}

const (
	ListDir = 0x0001
)

func staticDirHandler(mux *http.ServeMux, prefix string, staticDir string, flags int) {
	mux.HandleFunc(prefix, func(w http.ResponseWriter, r *http.Request) {
		file := staticDir + r.URL.Path[len(prefix)-1:]
		if (flags & ListDir) == 0 {
			if exists := isExists(file); !exists {
				http.NotFound(w, r)
			}
		}
		http.ServeFile(w, r, file)
	})
}

func main() {

	mux := http.NewServeMux()
	staticDirHandler(mux, "/assets/", "./public", 0)

	mux.HandleFunc("/upload", safeHandler(uploadHandler))
	mux.HandleFunc("/view", safeHandler(viewHandler))
	mux.HandleFunc("/", safeHandler(listHandler))

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal("ListenAndServe: ", err.Error())
	}
}
