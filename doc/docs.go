package docs

import (
	esfunc "Pet_Elastic/elasticsearch"
)

type Page struct {
	Title    string
	Places   []esfunc.Place
	CurrPage int
	NextPage int
	PrevPage int
}

type PageSearchResult struct {
	Title  string
	Total  int
	RspRes esfunc.RspSearch
}

const Doc = `
		<!doctype html>
		<html>
		<head>
			<meta charset="utf-8">
			<title>Places</title>
		</head>
		<body>
		<h5>Total</h5>
		<ul>
			{{range .Places}}
				<li>
					<div>{{.Name}}</div>
					<div>{{.Address}}</div>
					<div>{{.Phone}}</div>
				</li>
			{{end}}
		</ul>
		<a href="/places/show?page={{.PrevPage}}">Previous</a>
		<a href="/places/show?page={{.NextPage}}">Next</a>
		</body>
		</html>
	`

const DocSearchResult = `
		<!doctype html>
		<html>
		<head>
			<meta charset="utf-8">
			<title>Search result :</title>
		</head>
		<body>
		<h5>Total : {{.Total}}</h5>
		<ul>
			{{range .RspRes.Hits.Hits}}
				<li>
					<div>{{.Source.Name}}</div>
					<div>{{.Source.Address}}</div>
					<div>{{.Source.Phone}}</div><br />
				</li>
			{{end}}
		</ul>
		</body>
		</html>
	`
