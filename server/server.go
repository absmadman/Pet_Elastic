package server

import (
	docs "Pet_Elastic/doc"
	esfunc "Pet_Elastic/elasticsearch"
	"Pet_Elastic/helpers"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

func SearchClosest(lonStr string, latStr string, es *elasticsearch.Client) *docs.PageSearchResult {

	latFl := helpers.StrToFloat(latStr)
	lonFl := helpers.StrToFloat(lonStr)

	q := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"filter": map[string]interface{}{
					"geo_distance": map[string]interface{}{
						"distance": "0.25km",
						"location": map[string]interface{}{
							"lat": latFl,
							"lon": lonFl,
						},
					},
				},
			},
		},
	}

	return DoSearchRequest(q, es)
}

func DoSearchRequest(q map[string]interface{}, es *elasticsearch.Client) *docs.PageSearchResult {
	var b bytes.Buffer

	err := json.NewEncoder(&b).Encode(q)

	helpers.Check(err)

	res, err := es.Search(
		es.Search.WithIndex("places"),
		es.Search.WithBody(&b),
		es.Search.WithContext(context.Background()),
	)

	helpers.Check(err)

	var Rsp esfunc.RspSearch

	err = json.NewDecoder(res.Body).Decode(&Rsp)

	helpers.Check(err)

	return ComposeSearchedPage(Rsp)
}

func SearchByName(name string, es *elasticsearch.Client) *docs.PageSearchResult {

	q := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"name": name,
			},
		},
	}

	return DoSearchRequest(q, es)
}
func CreateHtmlPage(id int, es *elasticsearch.Client, w http.ResponseWriter) *docs.Page {
	var PlacesL []esfunc.Place
	for i := id * 10; i < id*10+10; i++ {
		if i == 0 {
			continue
		}
		PlacesL = append(PlacesL, *GetData(i, es, w))
	}
	page := docs.Page{Title: "Places", Places: PlacesL, NextPage: id + 1, CurrPage: id, PrevPage: id - 1}
	return &page
}

func ComposeSearchedPage(Rsp esfunc.RspSearch) *docs.PageSearchResult {
	return &docs.PageSearchResult{Title: "Places", RspRes: Rsp, Total: len(Rsp.Hits.Hits)}
}

func GetData(id int, es *elasticsearch.Client, w http.ResponseWriter) *esfunc.Place {
	var indexes []string
	var curPlace esfunc.Place

	indexes = append(indexes, fmt.Sprintf("places/_doc/%d", id))
	request := esapi.IndicesGetRequest{
		Index:  indexes,
		Pretty: true,
		Human:  true,
	}
	result, err := request.Do(context.Background(), es)

	if result.IsError() {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return &curPlace
	}

	CheckErrorServer(err, w)

	data := strings.Split(result.String(), "\"_source\" : ")

	err = json.Unmarshal([]byte(data[1][:len(data[1])-2]), &curPlace)

	CheckErrorServer(err, w)

	return &curPlace
}

func CheckErrorServer(err error, w http.ResponseWriter) {
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
	}
}

func GetPage(w http.ResponseWriter, r *http.Request) int {

	pageNum, err := strconv.Atoi(r.URL.Query().Get("page"))
	if pageNum < 0 {
		http.Error(w, "Bad request", http.StatusBadRequest)
	}
	CheckErrorServer(err, w)
	return pageNum
}

func ShowHandler(w http.ResponseWriter, r *http.Request, es *elasticsearch.Client) {
	var err error

	w.Header().Add("Content-Type", "text/html")

	pageNum := GetPage(w, r)

	templates := template.New("places")
	_, err = templates.Parse(docs.Doc)

	CheckErrorServer(err, w)

	page := *CreateHtmlPage(pageNum, es, w)
	err = templates.Execute(w, page)

	CheckErrorServer(err, w)
}

func ClosestHandler(w http.ResponseWriter, r *http.Request, es *elasticsearch.Client) {
	var err error

	w.Header().Add("Content-Type", "text/html")

	latStr := r.URL.Query().Get("lat")
	lonStr := r.URL.Query().Get("lon")

	if latStr == "" || lonStr == "" {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	templates := template.New("searched")
	_, err = templates.Parse(docs.DocSearchResult)

	CheckErrorServer(err, w)

	err = templates.Execute(w, SearchClosest(lonStr, latStr, es))

	CheckErrorServer(err, w)
}

func AddHandler(w http.ResponseWriter, r *http.Request, es *elasticsearch.Client) {
	w.Header().Add("Content-Type", "text/html")

	var place esfunc.Place

	err := json.NewDecoder(r.Body).Decode(&place)

	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	var places []esfunc.Place

	places = append(places, place)

	esfunc.InsertData(es, places)
}

func SearchNameHandler(w http.ResponseWriter, r *http.Request, es *elasticsearch.Client) {
	var err error

	w.Header().Add("Content-Type", "text/html")

	name := r.URL.Query().Get("name")

	templates := template.New("searched")
	_, err = templates.Parse(docs.DocSearchResult)

	CheckErrorServer(err, w)

	err = templates.Execute(w, SearchByName(name, es))

	CheckErrorServer(err, w)

}

func LaunchHttpServer(es *elasticsearch.Client) {

	http.HandleFunc("/places/show", func(w http.ResponseWriter, r *http.Request) {
		ShowHandler(w, r, es)
	})

	http.HandleFunc("/places/closest", func(w http.ResponseWriter, r *http.Request) {
		ClosestHandler(w, r, es)
	})

	http.HandleFunc("/places/add", func(w http.ResponseWriter, r *http.Request) {
		AddHandler(w, r, es)
	})

	http.HandleFunc("/places/search/name", func(w http.ResponseWriter, r *http.Request) {
		SearchNameHandler(w, r, es)
	})

	err := http.ListenAndServe(":8080", nil)

	helpers.Check(err)
}
