package server

// index returns a template for the front page.
func index() string {
	return `
<html>
  <head>
    <title>Glint</title>
  </head>
  <body>
    {{ .Root }}
  </body>
</html>
`
}
