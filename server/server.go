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

	LatFl := helpers.StrToFloat(latStr)
	LonFl := helpers.StrToFloat(lonStr)

	q := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"filter": map[string]interface{}{
					"geo_distance": map[string]interface{}{
						"distance": "0.25km",
						"location": map[string]interface{}{
							"lat": LatFl,
							"lon": LonFl,
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
	for i := id; i < id+10; i++ {
		PlacesL = append(PlacesL, *GetData(i, es, w))
	}
	page := docs.Page{Title: "Places", Places: PlacesL, NextPage: id + 1, CurrPage: id, PrevPage: id - 1}
	return &page
}

func ComposeSearchedPage(Rsp esfunc.RspSearch) *docs.PageSearchResult {
	return &docs.PageSearchResult{Title: "Places", RspRes: Rsp, Total: len(Rsp.Hits.Hits)}
}

func GetData(id int, es *elasticsearch.Client, w http.ResponseWriter) *esfunc.Place {
	var Indexes []string
	var CurPlace esfunc.Place

	Indexes = append(Indexes, fmt.Sprintf("places/_doc/%d", id))
	request := esapi.IndicesGetRequest{
		Index:  Indexes,
		Pretty: true,
		Human:  true,
	}
	result, err := request.Do(context.Background(), es)

	if result.IsError() {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return &CurPlace
	}

	CheckErrorServer(err, w)

	data := strings.Split(result.String(), "\"_source\" : ")

	err = json.Unmarshal([]byte(data[1][:len(data[1])-2]), &CurPlace)

	CheckErrorServer(err, w)

	return &CurPlace
}

func CheckErrorServer(err error, w http.ResponseWriter) {
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
	}
}

func GetPage(w http.ResponseWriter, r *http.Request) int {
	PageNum, err := strconv.Atoi(r.URL.Query().Get("page"))
	if PageNum < 1 {
		http.Error(w, "Bad request", http.StatusBadRequest)
	}
	CheckErrorServer(err, w)
	return PageNum
}

func LaunchHttpServer(es *elasticsearch.Client) {
	var err error
	http.HandleFunc("/places/show", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/html")

		PageNum := GetPage(w, r)

		templates := template.New("places")
		_, err = templates.Parse(docs.Doc)

		CheckErrorServer(err, w)

		page := *CreateHtmlPage(PageNum, es, w)
		err = templates.Execute(w, page)

		CheckErrorServer(err, w)
	})

	http.HandleFunc("/places/closest", func(w http.ResponseWriter, r *http.Request) {
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

	})

	http.HandleFunc("/places/add", func(w http.ResponseWriter, r *http.Request) {
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
	})

	http.HandleFunc("/places/search/name", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/html")

		name := r.URL.Query().Get("name")

		templates := template.New("searched")
		_, err = templates.Parse(docs.DocSearchResult)

		CheckErrorServer(err, w)

		err = templates.Execute(w, SearchByName(name, es))

		CheckErrorServer(err, w)

	})

	err = http.ListenAndServe(":8080", nil)
	helpers.Check(err)
}
