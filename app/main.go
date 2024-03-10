package main

import (
	"Pet_Elastic/config"
	esfunc "Pet_Elastic/elasticsearch"
	"Pet_Elastic/helpers"
	"Pet_Elastic/server"
	"github.com/elastic/go-elasticsearch/v8"
)

/*
type Page struct {
	Title    string
	Places   []Place
	CurrPage int
	NextPage int
	PrevPage int
}

type PageSearchResult struct {
	Title  string
	Total  int
	RspRes RspSearch
}
*/

/*
type Place struct {
	Id      int      `json:"id"`
	Name    string   `json:"name"`
	Address string   `json:"address"`
	Phone   string   `json:"phone"`
	Loc     Location `json:"location"`
}

type Location struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type RspSearch struct {
	Hits struct {
		Hits []struct {
			Source struct {
				ID       int    `json:"id"`
				Name     string `json:"name"`
				Address  string `json:"address"`
				Phone    string `json:"phone"`
				Location struct {
					Lat float64 `json:"lat"`
					Lon float64 `json:"lon"`
				} `json:"location"`
			} `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}
*/
/*
func CreateHtmlPage(id int, es *elasticsearch.Client, w http.ResponseWriter) *Page {
	var PlacesL []Place
	for i := id; i < id+10; i++ {
		PlacesL = append(PlacesL, *GetData(i, es, w))
	}
	fmt.Println(es.Indices)
	page := Page{Title: "Places", Places: PlacesL, NextPage: id + 1, CurrPage: id, PrevPage: id - 1}
	return &page
}


*/
/*
	func Check(err error) {
		if err != nil {
			panic(err)
		}
	}

	func JsonParse(filename string) []Place {
		data, err := os.ReadFile(filename)

		Check(err)

		var Places []Place

		err = json.Unmarshal(data, &Places)

		Check(err)

		return Places
	}
*/

/*
func InsertData(es *elasticsearch.Client, Places []Place) {

	var data strings.Builder

	for i := range Places {
		tmp, err := json.Marshal(Places[i])
		Check(err)
		data.Write([]byte(fmt.Sprintf(
			"{ \"index\" : { \"_index\" : \"places\" ,  \"_id\" : \"%d\"} }\n", Places[i].Id)))
		data.Write(tmp)
		data.WriteString("\n")
	}
	blk, err := es.Bulk(strings.NewReader(data.String()), es.Bulk.WithIndex("places"))

	Check(err)

	defer blk.Body.Close()

	CheckResponse(blk, "Index request success", "Error request")
}
*/

/*

func MakeRequest(es *elasticsearch.Client, RequestFunc func(es *elasticsearch.Client) (*esapi.Response, error)) {
	res, err := RequestFunc(es)
	Check(err)

	CheckResponse(res, "Request success", "Error request")

	defer res.Body.Close()
}

func DefaultInit(es *elasticsearch.Client) {

	data := JsonParse("data/initial_data.json")

	MakeRequest(es, func(es *elasticsearch.Client) (*esapi.Response, error) {
		var Indexes []string
		Indexes = append(Indexes, "places")
		request := esapi.IndicesDeleteRequest{
			Index: Indexes,
		}
		return request.Do(context.Background(), es)
	})

	MakeRequest(es, func(es *elasticsearch.Client) (*esapi.Response, error) {
		mappings :=
			`{
			"mappings": {
				"properties": {
					"name" : { "type": "text" },
					"address" : { "type": "text" },
					"phone" : { "type": "text" },
					"location" : { "type": "geo_point" }
				}
			}
		}`
		request := esapi.IndicesCreateRequest{
			Index: "places",
			Body:  strings.NewReader(mappings),
		}
		return request.Do(context.Background(), es)
	})

	InsertData(es, data)

	MakeRequest(es, func(es *elasticsearch.Client) (*esapi.Response, error) {
		var Indexes []string
		Indexes = append(Indexes, "places/_doc/1")
		request := esapi.IndicesGetRequest{
			Index: Indexes,
		}
		return request.Do(context.Background(), es)
	})
}
*/

/*
func GetData(id int, es *elasticsearch.Client, w http.ResponseWriter) *Place {
	var Indexes []string
	var CurPlace Place

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

*/

/*

func CheckErrorServer(err error, w http.ResponseWriter) {
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
	}
}
*/
/*
func GetPage(w http.ResponseWriter, r *http.Request) int {
	pagenum, err := strconv.Atoi(r.URL.Query().Get("page"))
	CheckPagenum(pagenum, w)
	CheckErrorServer(err, w)
	return pagenum
}
*/

/*
func StrToFloat(str string) float64 {
	flt, err := strconv.ParseFloat(str, 64)
	Check(err)
	return flt
}

func CheckResponse(response *esapi.Response, GoodResponse string, BadResponse string) {
	if response.IsError() {
		fmt.Printf(BadResponse)
		fmt.Println()
	} else {
		fmt.Printf(GoodResponse)
		fmt.Println()
	}
}
*/

/*
func SearchClosest(lonStr string, latStr string, es *elasticsearch.Client) *PageSearchResult {

		LatFl := StrToFloat(latStr)
		LonFl := StrToFloat(lonStr)

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
*/

/*
	func DoSearchRequest(q map[string]interface{}, es *elasticsearch.Client) *PageSearchResult {
		var b bytes.Buffer

		err := json.NewEncoder(&b).Encode(q)

		Check(err)

		res, err := es.Search(
			es.Search.WithIndex("places"),
			es.Search.WithBody(&b),
			es.Search.WithContext(context.Background()),
		)

		Check(err)

		var Rsp RspSearch

		err = json.NewDecoder(res.Body).Decode(&Rsp)

		Check(err)

		return ComposeSearchedPage(Rsp)
	}

func SearchByName(name string, es *elasticsearch.Client) *PageSearchResult {

		q := map[string]interface{}{
			"query": map[string]interface{}{
				"match": map[string]interface{}{
					"name": name,
				},
			},
		}

		return DoSearchRequest(q, es)
	}

	func LaunchHttpServer(es *elasticsearch.Client) {
		var err error
		http.HandleFunc("/places/show", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "text/html")

			pagenum := GetPage(w, r)

			templates := template.New("places")
			_, err = templates.Parse(docs.Doc)

			CheckErrorServer(err, w)

			page := *CreateHtmlPage(pagenum, es, w)
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

			var place Place

			err := json.NewDecoder(r.Body).Decode(&place)

			if err != nil {
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}

			var places []Place

			places = append(places, place)

			InsertData(es, places)
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

		http.ListenAndServe(":8080", nil)
	}
*/
func main() {
	es, err := elasticsearch.NewClient(*config.NewCfg())

	helpers.Check(err)

	esfunc.DefaultInit(es)

	server.LaunchHttpServer(es)
}
