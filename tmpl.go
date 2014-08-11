// tmpl.go
package main

const header = `
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>set</title>
	<link rel="icon" href="pub/set.ico" />
	<link href="pub/bootstrap.min.css" rel="stylesheet">
	<style type="text/css">
	.mainput {
		outline: none;
		border: 0;
		border-left: 9px solid #008cba;
		margin: 10px 0 20px 0;
		padding: 0 0 0 5px;
		font-size: 20px;
		width: 100%;
	}
	img {
		max-width: 100%;
	}
	h3 {
		color: #008cba;
	}
	</style>
</head>
`
const root = header + `
<body>
<div class="container-fluid">
	<div class="row">
		<div class="col-sm-12">

		<form action="/" method="post">
			<input type="text" class="mainput" name="search" value="{{.Search}}" placeholder="neither impatient nor common." autofocus="autofocus" autocomplete="off">
		</form>
		<hr>

		{{range .Blog}}
		<article>
			<h3>{{.Title}}</h3>
			<div>{{.Content}}</div>
			<form action="/modify" method="post">
				<input type="hidden" name="id" value="{{.ID}}">
				<input type="submit" class="btn btn-xs btn-danger" name="button" value="Trash">
				<input type="submit" class="btn btn-xs btn-warning" name="button" value="Modify">
			</form>
			<hr>
		</article>
		{{end}}

		</div>
	</div>
</div>
</body>
</html>
`
const editor = header + `
<body>
<div class="container-fluid">
	<div class="row">
		<div class="col-sm-12">

		<form action="/put" method="post">
			<input type="hidden" name="id" value="{{.ID}}">
			<input type="text" class="mainput" name="title" value="{{.Title}}" placeholder="Title" autofocus="autofocus" autocomplete="off">
			<textarea id="editor" name="content">{{.Content}}</textarea>
			<span style="float:right;"><input type="submit" class="btn btn-primary" value="Save"></span>
		</form>

		<form enctype="multipart/form-data" method="post" name="fileinfo">
			<input id="file" type="file" name="file" required multiple />
			<input type="submit" class="btn btn-default" value="Insert" />
		</form>
		<progress id="uploadprogress" min="0" max="100" value="0">0</progress>

		</div>
	</div>
</div>
<script src="pub/ckeditor/ckeditor.js"></script>
<script src="pub/set.js"></script>
</body>
</html>
`
