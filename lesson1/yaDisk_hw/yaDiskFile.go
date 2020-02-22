package main

/*
Как получить OAuth-токен:
https://yandex.ru/dev/direct/doc/start/token-docpage/
*/

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/pkg/errors"
)

const (
	commonDownLoadURL = "https://cloud-api.yandex.net/v1/disk/public/resources?public_key="
	publicKey         = "https://yadi.sk/i/b0JgT1QT7OqQ-A"
)

func main() {
	// запуск функций с обработкой ошибок:
	// запуск функции, которая получает ссылку на скачивание
	fileInfo, err := getDownloadLink(commonDownLoadURL + publicKey)
	if err != nil {
		// Wrap возвращает ошибку с обозначанием вместе со следом вызовов
		// в месте где Wrap вызван, и доставляет сообщение.
		// If err is nil, Wrap returns nil.
		log.Println(errors.Wrap(err, "getDownloadLink"))
		return
	}

	// запуск функции скачивания файла которя на вход получает ссылку по указателю (почему?)
	if err := downloadAndSaveFile(*fileInfo); err != nil {
		log.Println(err)
	}
}

// структура, которя будет конвертироваться в json
type FileInfo struct {
	Href     string `json:"href"`
	FileLink string `json:"file"`
	FileName string `json:"name"`
}

// метод получения ссылки на скачивание: на вход получает ссылку, на выходе записывает данные в структуру
// и возвращает ошибку
func getDownloadLink(reqURL string) (*FileInfo, error) {
	body, err := doReq(reqURL, "GET", nil, false)
	if err != nil {
		return nil, errors.Wrap(err, "get body")
	}

	resp := new(FileInfo)
	if err := json.Unmarshal(body, resp); err != nil {
		return nil, errors.Wrap(err, "unmarshal response")
	}

	if len(resp.FileLink) == 0 {
		return nil, errors.New("empty link")
	} else if len(resp.FileName) == 0 {
		return nil, errors.New("file name")
	}
	return resp, nil
}

func doReq(reqURL, method string, reqBody []byte, isAuth bool) ([]byte, error) {
	req, err := http.NewRequest(method, reqURL, bytes.NewReader(reqBody))
	if err != nil {
		return nil, errors.Wrap(err, "creation req")
	}

	if isAuth {
		req.Header.Set("Authorization", "OAuth "+token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "do http req")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "read response")
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, errors.Errorf("Status code not OK or Created: %v", resp.StatusCode)
	}
	return body, nil
}

const uploadFile = "https://cloud-api.yandex.net/v1/disk/resources/upload?path="

func downloadAndSaveFile(info FileInfo) error {
	body, err := doReq(info.FileLink, "GET", nil, false)
	if err != nil {
		return errors.Wrap(err, "get body by link")
	}

	if err := ioutil.WriteFile(info.FileName, body, 0755); err != nil {
		return errors.Wrap(err, "create file")
	}

	bodyUpload, err := doReq(uploadFile+"/test/file.docx", "GET", nil, true)
	if err != nil {
		return errors.Wrap(err, "get body to upload")
	}

	resp := new(FileInfo)
	if err := json.Unmarshal(bodyUpload, resp); err != nil {
		return errors.Wrap(err, "unmarshal response")
	}

	_, err = doReq(resp.Href, "PUT", body, false)
	if err != nil {
		return errors.Wrap(err, "send request")
	}

	return nil
}