package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

// Определяем структуру
type Server struct {
	// db - обращение к базе данных
	db        *sql.DB
	Title     string
	// Templates - шаблоны которые будут выводить содержание
	Templates map[templateName]*template.Template
}

func main() {
	// Подключаемся к БД
	db, err := sql.Open("mysql", DSN)
	if err != nil {
		// В данном случае log.Fatal т.к. при неудачном подключении нет смысла продолжать
		log.Fatal(err)
	}
	defer db.Close()

	// Проверяем, установлено-ли соединение с БД
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	// Заполняем структуру
	server := Server{
		Title:     "Posts from Habr",
		Templates: createTemplates(),
		db:        db,
	}

	// Заполняем БД таблицами с постами (применяется метод insertDefault)
	server.insertDefault()

	// router := http.NewServeMux()
	// Настраиваем mux роутер
	router := mux.NewRouter()
	router.HandleFunc("/", server.handlePostsList)
	router.HandleFunc("/post/{id:[0-9]+}", server.handleSinglePost)
	router.HandleFunc("/edit/{id:[0-9]+}", server.handleEditPost)
	router.HandleFunc("/results", server.handleResults)

	port := "8080"
	log.Printf("start server on port: %v", port)

	go func() {
		_ = http.ListenAndServe(":"+port, router)
	}()

	// Создаем канал, где будем ждать сигнал прерывания
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// При поступлени сигналя прерывания мы очищаем таблицы в БД
	<-interrupt
	server.truncate()
}


// Отображаем все посты
func (server *Server) handlePostsList(wr http.ResponseWriter, req *http.Request) {
	// Загружаем нужный шаблон
	tmpl := getTemplate(server.Templates, List)
	// Обрабатываем ошибку
	if tmpl == nil {
		err := errors.New("empty template")
		http.Error(wr, err.Error(), http.StatusInternalServerError)
		log.Print(err)
		return
	}

	// Получаем в переменную посты (содержание строк таблицы из БД)
	posts, err := getPosts(server.db)
	if err != nil {
		log.Println(err)
		return
	}

	// Отображаем нужный шаблон с содержимым (постами)
	if err := tmpl.ExecuteTemplate(wr, "page", posts); err != nil {
		err = errors.Wrap(err, "execute template")
		http.Error(wr, err.Error(), http.StatusInternalServerError)
		log.Print(err)
		return
	}
}

// Отображаем один пост
func (server *Server) handleSinglePost(wr http.ResponseWriter, req *http.Request) {
	// Готовим шаблон
	tmpl := getTemplate(server.Templates, Single)
	if tmpl == nil {
		err := errors.New("empty template")
		http.Error(wr, err.Error(), http.StatusInternalServerError)
		log.Print(err)
		return
	}

	// Получаем переменную (ID поста)
	vars := mux.Vars(req)

	// Проверяем наличие ID
	id := vars["id"]
	if len(id) == 0 {
		log.Println(errors.New("empty id"))
		return
	}

	// Извлекаем из БД пост под номером..
	post, err := getPost(server.db, id)
	if err != nil {
		err := errors.Wrap(err, "empty post")
		http.Error(wr, err.Error(), http.StatusInternalServerError)
		log.Print(err)
		return
	}

	// Запускаем шаблон
	if err := tmpl.ExecuteTemplate(wr, "page", post); err != nil {
		err := errors.Wrap(err, "execute template")
		http.Error(wr, err.Error(), http.StatusInternalServerError)
		log.Print(err)
		return
	}
}


// Редактируем пост
func (server *Server) handleEditPost(wr http.ResponseWriter, req *http.Request) {
	// Готовим шаблон
	tmpl := getTemplate(server.Templates, Edit)
	if tmpl == nil {
		err := errors.New("empty template")
		http.Error(wr, err.Error(), http.StatusInternalServerError)
		log.Print(err)
		return
	}

	// Получаем АйДи
	vars := mux.Vars(req)

	id := vars["id"]
	if len(id) == 0 {
		log.Println("edit: empty id")
		return
	}

	// Извлекаем пост
	post, err := getPost(server.db, id)
	if err != nil {
		err := errors.Wrap(err, "empty post")
		http.Error(wr, err.Error(), http.StatusInternalServerError)
		log.Print(err)
		return
	}

	// Запускаем шаблон
	if err := tmpl.ExecuteTemplate(wr, "page", post); err != nil {
		err = errors.Wrap(err, "execute template")
		http.Error(wr, err.Error(), http.StatusInternalServerError)
		log.Print(err)
		return
	}
}


// Записываем данные в БД
func (server *Server) handleResults(wr http.ResponseWriter, req *http.Request) {
	// Получаем данные из формы
	idVal := req.FormValue("id")
	if len(idVal) == 0 {
		log.Print("results: empty id")
		return
	}

	// Конвертируем АйДи(стринг) в INT
	id, err := strconv.Atoi(idVal)
	if err != nil {
		err := errors.Wrapf(err, "id from form value: %v", idVal)
		http.Error(wr, err.Error(), http.StatusInternalServerError)
		log.Print(err)
		return
	}

	// Заполняем структуру данными из формы
	post := Post{
		Id:      id,
		Title:   req.FormValue("title"),
		Date:    req.FormValue("date"),
		Link:    req.FormValue("link"),
		Comment: req.FormValue("comment"),
	}

	// Записываем в БД
	if err := editPost(server.db, post, idVal); err != nil {
		log.Print(err)
		return
	}

	// Перенаправляем после успешной записи
	http.Redirect(wr, req, "/", http.StatusFound)
}
