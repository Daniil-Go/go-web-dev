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

// объявляем структуру (новый тип)
type FileInfo struct {
	Href     string `json:"href"`
	FileLink string `json:"file"`
	FileName string `json:"name"`
}

// метод получения ссылки на скачивание: на вход получает ссылку, на выходе записывает данные в структуру
// и возвращает ошибку (*FileInfo почему по указателю?)
func getDownloadLink(reqURL string) (*FileInfo, error) {
	// мы получаем запрос doReq: на вход ссылка, метод (гет), тело (nil), авторизация (false) и записываем в body
	body, err := doReq(reqURL, "GET", nil, false)
	//обработка ошибок
	if err != nil {
		return nil, errors.Wrap(err, "get body")
	}

	// присваеваем пустую структуру переменной:
	// The new built-in function allocates memory. The first argument is a type,
	// not a value, and the value returned is a pointer to a newly
	// allocated zero value of that type.
	// func new(Type) *Type
	resp := new(FileInfo)

	// перекодируем полученную инфу json из body в структуру resp
	if err := json.Unmarshal(body, resp); err != nil {
		return nil, errors.Wrap(err, "unmarshal response")
	}

	// проверка на аналичие ссылки файла
	if len(resp.FileLink) == 0 {
		return nil, errors.New("empty link")
		// проверка на аналичие имени файла
	} else if len(resp.FileName) == 0 {
		return nil, errors.New("file name")
	}
	return resp, nil
}

// функция запроса: на вход ссылка, метод, тело запроса и авторизация -булевое значение; на выход - массив байт и ошибка
// почему isAuth bool?
func doReq(reqURL, method string, reqBody []byte, isAuth bool) ([]byte, error) {
	// присваиваем req запрос в виде масива байт
	req, err := http.NewRequest(method, reqURL, bytes.NewReader(reqBody))
	if err != nil {
		return nil, errors.Wrap(err, "creation req")
	}

	// в Header запроса добавляем токен (карта: "Authorization", значение "OAuth "+token)
	if isAuth {
		req.Header.Set("Authorization", "OAuth "+token)
	}

	// получаем ответ
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "do http req")
	}
	defer resp.Body.Close()

	// записываем в body ответ на запрос
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "read response")
	}

	// обработка ошибок, если в ответе не код 200 или 201
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, errors.Errorf("Status code not OK or Created: %v", resp.StatusCode)
	}
	return body, nil
}

const uploadFile = "https://cloud-api.yandex.net/v1/disk/resources/upload?path="

// функция сохранения файла: на вход инфо - структура, возврат ошибок
func downloadAndSaveFile(info FileInfo) error {
	// функция запроса
	body, err := doReq(info.FileLink, "GET", nil, false)
	if err != nil {
		return errors.Wrap(err, "get body by link")
	}

	// записываем файл с правами "Каждый пользователь имеет право читать и запускать на выполнение;
	//владелец может редактировать"
	if err := ioutil.WriteFile(info.FileName, body, 0755); err != nil {
		return errors.Wrap(err, "create file")
	}

	// проверка загрузки файла на сервер:
	bodyUpload, err := doReq(uploadFile+"/test/file.docx", "GET", nil, true)
	if err != nil {
		return errors.Wrap(err, "get body to upload")
	}

	// resp - новая структура
	resp := new(FileInfo)
	//перекодируем из json результат работы функции загрузки файла
	if err := json.Unmarshal(bodyUpload, resp); err != nil {
		return errors.Wrap(err, "unmarshal response")
	}

	// загружаем файл на сервер
	_, err = doReq(resp.Href, "PUT", body, false)
	if err != nil {
		return errors.Wrap(err, "send request")
	}

	return nil
}