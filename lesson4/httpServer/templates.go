package main

import (
	"html/template"
	"path"
	"sync"
)

// Создаем тип
type templateName string

// Обозначаем константы (кастомные типы)
const (
	List   templateName = "list.html"
	Single templateName = "single.html"
	Edit   templateName = "edit.html"
)

// Слайс из templateName
var templates = []templateName{
	List, Single, Edit,
}

func createTemplates() map[templateName]*template.Template {
	// Создаем мапку
	out := make(map[templateName]*template.Template, len(templates))

	// Проходимся циклом по слайсу
	for _, tmplName := range templates {
		// Заполняем мапку шаблонами
		out[tmplName] = template.Must(
			template.New("MyTemplate").ParseFiles(path.Join("templates", string(tmplName))))
	}

	return out
}

// Mutex
var mu sync.Mutex

func getTemplate(mapka map[templateName]*template.Template, key templateName) *template.Template {
	// При вызове шаблона блокируем обработчик входящих соединений
	mu.Lock()
	defer mu.Unlock()

	val, ok := mapka[key]
	if !ok {
		return nil
	}

	return val
}
