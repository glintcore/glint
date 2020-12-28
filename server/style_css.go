package server

// styleCss returns the site css file.
func styleCss() string {
	return `
body
{
    margin-left: 1in;
    font-size: 14.5px;
}

thead {
    color: #458;
}

tbody {
    color: #458;
}

tfoot {
    color: #458;
}

h1
{
    font-family: sans-serif, "Times New Roman", "Roman", serif;
    color: #458;
    font-size: 18px;
    line-height: 1.2em;
    font-weight: 500;
    margin: 0;
    margin-top: 1.5em;
    margin-bottom: 1.5em;
}

th#sc:hover
{
    cursor: pointer;
    background: #7eb4d7;
}

th#da:hover
{
    cursor: pointer;
    background: #7eb4d7;
}

table
{
    font-family: sans-serif, "Times New Roman", "Roman", serif;
    border-collapse: collapse;
    table-layout: fixed;
    word-wrap: break-word;
}

table th 
{
    font-size: 16px;
    /* border: 1px solid #458; */
    border: 1px solid white;
    text-align: left;
    /* background-color: rgb(121, 188, 100); */
    background-color: #458;
    color: white;
    width: 100px;
    padding: 6px;
}

table td 
{
    border: 1px solid #458;
    text-align: left;
    background-color: white;
    color: black;
    font-size: 13px;
    padding: 4px;
    white-space: nowrap;
}

table tr td
{
    border-left: 0px;
    border-right: 0px;
}

table tr.row_alt td
{
    background-color: #458;
}

th#da
{
    width: 210px;
}

body .blank_row {
    height: 10px !important;
}

div
{
    padding-top: 10px;
}

input
{
    /* width: 96px; */
    padding-left: 1px;
    padding-right: 1px;
    border-radius: 0.5em;
}

#inputTable
{
    margin-top: 3px;
    margin-left: 1px;
}

#inputTable tr
{
    padding: 4px;

}

#buttontypeLogin2
{
    background-color: white;
    color: #458;
    border-radius: 5px;
    border: 1px solid #458;
    padding-top: 0px;
    padding-bottom: 3px;
    font-size: 16px;
    width: 54px;
    margin-top: 0px;
}

a {
    text-decoration: none;
}
`
}
