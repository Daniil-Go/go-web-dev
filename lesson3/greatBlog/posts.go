package main

import (
	"errors"
	"strconv"
)

func (p Posts) getID(queryID string) (Post, error) {
	mu.Lock()
	defer mu.Unlock()

	id, err := strconv.Atoi(queryID)
	if err != nil {
		return Post{}, err
	}

	post, ok := p[id]
	if !ok {
		return Post{}, errors.New("empty map for ID: ")
	}

	return post, nil

}

func createPosts() Posts {
	return map[int]Post{
		0: {
			Id:      0,
			Title:   "GitHub опубликовал отчет о блокировках и удалении контента пользователей за 2019 год",
			Date:    "21 февраля 2020 в 20:42",
			Link:    "https://habr.com/ru/news/t/489448/",
			Comment: "Отличная статья про гит",
		},
		1: {
			Id:      1,
			Title:   "Рефакторинг — мощь сокрытая в качественном коде",
			Date:    "15 августа 2016 в 13:55",
			Link:    "https://habr.com/ru/post/307762/",
			Comment: "Стоит обратить на рефакторинг",
		},
		2: {
			Id:      2,
			Title:   "20 приёмов работы в командной строке Linux, которые сэкономят уйму времени",
			Date:    "11 октября 2017 в 12:07",
			Link:    "https://habr.com/ru/company/ruvds/blog/339820/",
			Comment: "Это и правда эконочит уйму времени. Читайте!",
		},
	}
}
