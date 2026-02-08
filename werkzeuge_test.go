package werkzeugkasten

import (
	"fmt"
	"image"
	"image/png"
	"io"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
)

func TestWerkzeuge_RandomString(t *testing.T) {
	var werkzeug Werkzeug
	length := 10
	s := werkzeug.RandomString(length)

	if len(s) != length {
		t.Errorf("Want length of %d but got length of %d\n", length, len(s))
	}
}

var uploadTests = []struct {
	name          string
	allowedTypes  []string
	renameFile    bool
	errorExpected bool
}{
	{name: "allowed no rename", allowedTypes: []string{"image/jpeg", "image/png"}, renameFile: false, errorExpected: false},
	{name: "allowed  rename", allowedTypes: []string{"image/jpeg", "image/png"}, renameFile: true, errorExpected: false},
	{name: "not allowed", allowedTypes: []string{"image/jpeg"}, renameFile: false, errorExpected: true},
}

func TestWerkzeuge_UploadFiles(t *testing.T) {
	for _, entry := range uploadTests {
		//set up a pipe to avoid buffering
		pr, pw := io.Pipe()
		writer := multipart.NewWriter(pw)
		wg := sync.WaitGroup{}
		wg.Add(1)

		go func() {
			defer writer.Close()
			defer wg.Done()

			//create form data field >>file<<
			part, err := writer.CreateFormFile("file", "./testdata/image.png")
			if err != nil {
				t.Error(err)
			}

			f, err := os.Open("./testdata/image.png")
			if err != nil {
				t.Error(err)
			}
			defer f.Close()

			img, _, err := image.Decode(f)
			if err != nil {
				t.Error("error decoding image", err)
			}

			err = png.Encode(part, img)
			if err != nil {
				t.Error(err)
			}
		}()

		//read from the pipe witch receives data
		request := httptest.NewRequest("POST", "/", pr)
		request.Header.Add("Content-Type", writer.FormDataContentType())

		var werkzeug Werkzeug
		werkzeug.AllowedFileTypes = entry.allowedTypes

		uploadedFiles, err := werkzeug.UploadFiles(request, "./testdata/uploads/", entry.renameFile)
		if err != nil && !entry.errorExpected {
			t.Error(err)
		}

		if !entry.errorExpected {
			if _, err := os.Stat(fmt.Sprintf("./testdata/uploads/%s", uploadedFiles[0].NewFileName)); os.IsNotExist(err) {
				t.Errorf("%s: expected file to exist: %s", entry.name, err.Error())
			}

			//Clean up
			_ = os.Remove(fmt.Sprintf("./testdata/uploads/%s", uploadedFiles[0].NewFileName))
		}

		if !entry.errorExpected && err != nil {
			t.Errorf("%s error expected but non received", entry.name)
		}

		wg.Wait()
	}
}

func TestWerkzeug_UploadOneFile(t *testing.T) {
	//set up a pipe to avoid buffering
	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)

	go func() {
		defer writer.Close()

		//create form data field >>file<<
		part, err := writer.CreateFormFile("file", "./testdata/image.png")
		if err != nil {
			t.Error(err)
		}

		f, err := os.Open("./testdata/image.png")
		if err != nil {
			t.Error(err)
		}
		defer f.Close()

		img, _, err := image.Decode(f)
		if err != nil {
			t.Error("error decoding image", err)
		}

		err = png.Encode(part, img)
		if err != nil {
			t.Error(err)
		}
	}()

	//read from the pipe witch receives data
	request := httptest.NewRequest("POST", "/", pr)
	request.Header.Add("Content-Type", writer.FormDataContentType())

	var werkzeug Werkzeug

	uploadedFile, err := werkzeug.UploadOneFile(request, "./testdata/uploads/", true)
	if err != nil {
		t.Error(err)
	}

	if _, err := os.Stat(fmt.Sprintf("./testdata/uploads/%s", uploadedFile.NewFileName)); os.IsNotExist(err) {
		t.Errorf("expected file to exist: %s", err.Error())
	}

	//Clean up
	_ = os.Remove(fmt.Sprintf("./testdata/uploads/%s", uploadedFile.NewFileName))

}
