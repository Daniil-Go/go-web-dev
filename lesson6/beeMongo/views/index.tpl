<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <title>Блог</title>
   {{template "head.tpl"}}
</head>
<body>
 {{template "header.tpl"}}
<div class="uk-card uk-card-default uk-card-body">
    <h1>Блог</h1>
    <h2>Приветствую!</h2>
    <p>Заголовок</p>
    <div class="uk-card uk-card-body">
        <ul class="uk-list">
            {{range .Posts}}
                <li>
                    <div class="uk-card uk-card-default uk-card-body">
                        <h3>{{.Title}}</h3>
                        <h4>{{.Date}}</h4>
                        <p>{{.Comment}}</p>
                        <p>{{.Link}}</p>
                        <a class="uk-button uk-button-default" href="/post/{{.Id}}">ПОДРОБНЕЕ</a>
                        <a class="uk-button uk-button-default" href="/prepare/{{.Id}}">РЕДАКТИРОВАТЬ ПОСТ</a>
                    </div>
                </li>
            {{end}}
        </ul>
    </div>
</div>
  {{template "footer.tpl"}}
</body>
</html>
