package controller

import (
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

// UploadHandler handles file upload into the server
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	var (
		status int
		err    error
	)

	defer func() {
		if err != nil {
			http.Error(w, err.Error(), status)
		}
	}()

	// parse request
	const _24K = (1 << 20) * 24
	if err = r.ParseMultipartForm(_24K); err != nil {
		status = http.StatusInternalServerError
		return
	}

	for _, fheaders := range r.MultipartForm.File {
		for _, hdr := range fheaders {
			// open uploaded
			var infile multipart.File
			if infile, err = hdr.Open(); err != nil {
				status = http.StatusInternalServerError
				return
			}

			// open destination
			var outfile *os.File
			if outfile, err = os.Create("/uploaded" + hdr.Filename); err != nil {
				status = http.StatusInternalServerError
				return
			}
			defer outfile.Close()

			// 32K buffer copy
			//var written int64
			if _, err = io.Copy(outfile, infile); err != nil {
				status = http.StatusInternalServerError
				return
			}

		}
	}
}
