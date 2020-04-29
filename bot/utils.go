package bot

import (
	"io/ioutil"
	"encoding/json"
	"net/http"
	"net/url"
	"log"
	"path/filepath"
	"bytes"
	"mime/multipart"
	"os"
	"io"
)

type File struct {
	path string
	filetype string // file type for vk upload.(doc, photo)
}

func (f *File) Init(path string, filetype string) *File {
	f.path = path
	f.filetype = filetype

	return f
}


func request(url string, params url.Values) map[string]interface {} {
	var result map[string]interface{}
	resp, err := http.PostForm(url, params)
	
	if err != nil {
		log.Fatalln(err)
	}

	json.NewDecoder(resp.Body).Decode(&result)
	
	return result
}

func UploadFile (url string, filename string, field string) map[string]interface{} {
	var result map[string]interface{}

    file, err := os.Open(filename)

    if err != nil {
        return nil
    }
    defer file.Close()

    body := &bytes.Buffer{}
    writer := multipart.NewWriter(body)
    part, err := writer.CreateFormFile(field, filepath.Base(file.Name()))

    if err != nil {
        return nil
    }

    io.Copy(part, file)
    writer.Close()
    request, err := http.NewRequest("POST", url, body)

    if err != nil {
        return nil
    }

    request.Header.Add("Content-Type", writer.FormDataContentType())
    client := &http.Client{}

    response, err := client.Do(request)

    if err != nil {
        return nil
    }
    defer response.Body.Close()

    json.NewDecoder(response.Body).Decode(&result)

    return result
}


func ReadFile(path string) []byte {
	data, err := ioutil.ReadFile(path)
    if err != nil {
        return nil
    }
    
    return data
}


func ReadFileString(path string) string {
    return string(ReadFile(path))
}


func List(dict interface {}) []interface {} {

	if (dict == nil) {
		return make([]interface {}, 1)
	}

	list := dict.([]interface {})

	return list
}

func Dict(list interface {}) map[string]interface {} {

	if (list == nil) {
		return make(map[string]interface {})
	}

	dict := list.(map[string]interface {})

	return dict
}

func HasString(list []interface {}, value string) bool {
	for i:=0; i < len(list); i++ {
		if(list[i].(string) == value) {
			return true
		}
	}

	return false
}
