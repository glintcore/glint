package server

func data() string {
	return `
<html>

<head>
    <meta http-equiv="Content-type" content="text/html; charset=utf-8">
        <link rel="stylesheet" href="/static/style.css" type="text/css"
	          media="screen" title="no title">
</head>

<body>

<h1>{{ .User }} / {{ .Path }}</h1>

<table>

{{ .Data }}

</table>

</body>

</html>
`
}
