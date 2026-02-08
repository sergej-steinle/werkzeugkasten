package werkzeugkasten

import (
	"errors"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const randomStrSource = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789+_"

// Werkzeug provides utility functions.
type Werkzeug struct {
	MaxFileSize      int64
	AllowedFileTypes []string
}

// RandomString returns a random string of length n using characters from randomStrSource.
func (w *Werkzeug) RandomString(n int) string {
	sb := strings.Builder{}
	sb.Grow(n)
	for i := 0; i < n; i++ {
		sb.WriteByte(randomStrSource[rand.IntN(len(randomStrSource))])
	}
	return sb.String()
}

// UploadedFile holds information about uploaded file
type UploadedFile struct {
	NewFileName      string
	OriginalFileName string
	FileSize         int64
}

// UploadOneFile is a convenience wrapper around UploadFiles that extracts and
// returns only the first file from the request. If rename is true (default),
// the file is saved with a randomly generated name while preserving its extension.
func (w *Werkzeug) UploadOneFile(r *http.Request, uploadDir string, rename ...bool) (*UploadedFile, error) {
	renameFile := true
	if len(rename) > 0 {
		renameFile = rename[0]
	}

	file, err := w.UploadFiles(r, uploadDir, renameFile)
	if err != nil {
		return nil, err
	}

	return file[0], err
}
// UploadFiles processes all files from a multipart form request and saves them
// to uploadDir. It validates each file's content type against AllowedFileTypes
// (all types are allowed when the slice is empty). If rename is true (default),
// files are saved with a randomly generated name while preserving their extension.
func (w *Werkzeug) UploadFiles(r *http.Request, uploadDir string, rename ...bool) ([]*UploadedFile, error) {
	renameFile := true
	if len(rename) > 0 {
		renameFile = rename[0]
	}

	var uploadedFiles []*UploadedFile

	if w.MaxFileSize == 0 {
		w.MaxFileSize = 1024 * 1024 * 1024
	}
	err := r.ParseMultipartForm(w.MaxFileSize)

	if err != nil {
		return nil, errors.New("the upload is to big")
	}

	for _, fHeaders := range r.MultipartForm.File {
		for _, header := range fHeaders {
			uploadedFiles, err = func(uploadedFiles []*UploadedFile) ([]*UploadedFile, error) {
				var uploadedFile UploadedFile
				infile, err := header.Open()
				if err != nil {
					return nil, err
				}
				defer infile.Close()

				buff := make([]byte, 512)
				_, err = infile.Read(buff)
				if err != nil {
					return nil, err
				}

				allowed := false
				fileType := http.DetectContentType(buff)

				if len(w.AllowedFileTypes) > 0 {
					for _, x := range w.AllowedFileTypes {
						if strings.EqualFold(fileType, x) {
							allowed = true
						}
					}
				} else {
					allowed = true
				}

				if !allowed {
					return nil, errors.New("uploaded file types are not permitted")
				}

				_, err = infile.Seek(0, 0)
				if err != nil {
					return nil, err
				}

				if renameFile {
					uploadedFile.NewFileName = fmt.Sprintf("%s%s", w.RandomString(25), filepath.Ext(header.Filename))
				} else {
					uploadedFile.NewFileName = header.Filename
				}

				uploadedFile.OriginalFileName = header.Filename

				var outfile *os.File
				defer outfile.Close()

				if outfile, err = os.Create(filepath.Join(uploadDir, uploadedFile.NewFileName)); err != nil {
					return nil, err
				} else {
					filesize, err := io.Copy(outfile, infile)
					if err != nil {
						return nil, err
					}

					uploadedFile.FileSize = filesize
				}

				uploadedFiles = append(uploadedFiles, &uploadedFile)

				return uploadedFiles, nil

			}(uploadedFiles)

			if err != nil {
				return uploadedFiles, err
			}
		}
	}

	return uploadedFiles, nil
}
