<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>{{.Title}}</title>
   {{template "head.tpl"}}
</head>
<body>
 {{template "header.tpl"}}
<div class="container">
    <h3>Post</h3>
    <form method="POST" action="/post/{{.Post.Id}}">
        <div class="form-group">
            <label>Id</label>
            <input type="id" name="id" class="form-control" value="">
        </div>
        <div class="form-group">
            <label>Title</label>
            <input type="title" name="title" class="form-control" value="">
        </div>
        <div class="form-group">
            <label>Date</label>
            <input type="text" name="date" class="form-control" value="">
        </div>
        <div class="form-group">
            <label>Link</label>
            <input type="text" name="link" class="form-control" value="">
        </div>
        <div class="form-group">
            <label>Comment</label>
            <textarea name="comment" class="form-control"></textarea>
        </div>
        <input class="btn btn-primary" type="submit" value="submit">
        <a class="btn btn-outline-primary" href="/">Back</a>
    </form>
    </div>
  {{template "footer.tpl"}}
</body>
</html>