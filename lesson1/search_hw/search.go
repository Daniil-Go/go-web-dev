package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
)

func main() {
	str := "html"
	sites := []string{
		"https://yandex.ru",
		"https://golang.org",
		"https://google.com",
		"https://github.com",
	}

	var (
		result = make([]string, 0, len(sites)) // создаем срез типа строка с кол-вом 0 элементов, но подготовленный
		// для увелечения кол-ва элементов как в массиве sites
		errs int
	)

	// you can add flag here
	if !true {
		result, errs = search(str, sites) // запуск функции поиска v1
	} else {
		result, errs = searchConcurrency(str, sites) // запуск функции поиска v2 - первой
	}

	if errs > 0 { // подсчет и вывод ошибок
		log.Printf("There are %v errors during request", errs)
	}

	if len(result) == 0 { // вывод сообщения об отсутствии результатов поиска
		log.Println("empty result")
		return
	}

	log.Printf("Sites: %v", result) // вывод успешных результатов поиска
}

// the first version of search()
// реализация функции поиска v1
func search(str string, sites []string) ([]string, int) { // принимает на вход строку запроса и массив сайтов, выводит
	//срез стринг и число (ошибок)
	out := make([]string, 0, 1) // подготовка среза для вывода результата
	errs := 0                   // счетчик ошибок

	for _, site := range sites { // запускается цикл, который итерируется по массиву сайтов
		res, err := getReq(site) // запуск функции GET запроса. возвращет массив байтов Body
		if err != nil {          // обработка ошибок
			errs++
			log.Print(err) // вывод ошибок
			continue
		}

		if strings.Contains(string(res), str) { //е сли в строке результата запроса содержится поисковый запрос,
			// то в срез out добавляется сайт, на котором этот запрос был найден
			out = append(out, site)
		}
	}

	return out, errs // возврат сайтов, на которых был найден поисковый запрос и количества ошибок
}

func getReq(reqURL string) ([]byte, error) { // функция GET запроса принимает на вход адрес сайта
	// и возвращает тело запроса в виде массива байтов и количество ошибок
	resp, err := http.Get(reqURL) // GET запрос
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body) // чтение из тела запроса
	if err != nil {
		return nil, err
	}

	return body, nil // возвращаем тело запроса и ничего
}

type chunk struct { // создаем структуру
	site string
	err  error
}

// the second version of search
func searchConcurrency(str string, sites []string) ([]string, int) { // тоже самое - на вход строка, массив сайтов,
	// а на выход срез и кол-во ошибок
	wg := sync.WaitGroup{}      // присваем тип WaitGroup (будем ждать завершения всех горутин)
	results := make(chan chunk) // создаем канал, который будет посылать данные в структуру

	for _, site := range sites { // итерируемся по массиву сайтов
		wg.Add(1)              // добавляем 1 в группу ожидания горутин
		go func(site string) { // запуск горутины
			defer wg.Done()          // по завершении данной горутины удаляем ее из группы ожидания
			res, err := getReq(site) // запуск функции ГЕТ запроса
			if err != nil {
				results <- chunk{site: site, err: err} // отправляем ошибки в структуру
				return
			}

			if strings.Contains(string(res), str) {
				results <- chunk{site: site} // отправляем найденное значение в структуру
			}
		}(site) // анонимная функция с замыканием; передаем в нее значение site
	}

	go func() {
		wg.Wait() // ждем, пока закончат работу все горутины
		log.Println("channel closed")
		close(results) // закрываем канал
	}()

	out := make([]string, 0, len(sites)) // готовим срез для вывода рез-татов
	errs := 0
	for chunk := range results { // results - канал, и мы из него читаем чанки, здесь chunk как элемент канала
		if chunk.err != nil {
			errs++
			// log.Printf("Error %v for site %v", chunk.err, chunk.site)
			continue
		}

		// log.Println(chunk.site)
		out = append(out, chunk.site)
	}

	return out, errs
}
