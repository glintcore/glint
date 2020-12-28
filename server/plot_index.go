package server

// plotIndex returns static html for the example plottin service.
func plotIndex() string {
	return `
<html>
  <head>
    <meta http-equiv="Content-type" content="text/html; charset=utf-8">
        <link rel="stylesheet" href="/static/style.css" type="text/css"
	          media="screen" title="no title">
    <title>Time Series Plotting Service</title>
  </head>
  <body>
    <p>This plotting service can plot any time series data containing a
    <b>dc:date</b> or <b>yamz:h1317</b> column.</p>
    <form>
      <div>Data set URL: <input type="text" name="dataurl" size="80"></div>
      <div><input type="submit" value="Submit"></div>
    </form>
  </body>
</html>
`
}
