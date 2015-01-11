package share

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/featen/gas/service/auth"
	log "github.com/featen/gas/utils/log"
)

func Register() {
	log.Info("share service registered")

	http.HandleFunc("/service/uploadphoto", uploadPhotoHandler)
}

func uploadPhotoHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		//parse the multipart form in the request
		err := r.ParseMultipartForm(100000)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//get a ref to the parsed multipart form
		m := r.MultipartForm

		//get the *fileheaders
		if m == nil {
			w.WriteHeader(http.StatusNotAcceptable)
			return
		}
		files := m.File["files"]
		urls := make([]string, 0, len(files))
		for i, _ := range files {
			//for each fileheader, get a handle to the actual file
			file, err := files[i].Open()
			defer file.Close()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			b, err := ioutil.ReadAll(file)
			h := md5.New()
			h.Write(b)
			filename := fmt.Sprintf("%x", h.Sum(nil))
			//create destination file making sure the path is writeable.
			//dst, err := os.Create("data/upload/" + files[i].Filename)
			fileurl := fmt.Sprintf("/upload/%s", filename)
			dst, err := os.Create("webapp/upload/" + filename)
			defer dst.Close()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			//copy the uploaded file to the destination file
			//if _, err := io.Copy(dst, file); err != nil {
			if _, err := dst.Write(b); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			urls = append(urls, fileurl)
		}

		s, err := auth.CookieStore.Get(r, "gas-session")
		if err != nil {
			w.WriteHeader(http.StatusNotAcceptable)
			return
		}
		if s.Values["id"] == nil || s.Values["time"] == nil || s.Values["magic"] == nil {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		fmt.Fprintf(w, "%s", strings.Join(urls, ";"))
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	return
}
