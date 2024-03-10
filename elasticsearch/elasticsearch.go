package elasticsearch

import (
	"Pet_Elastic/helpers"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"os"
	"strings"
)

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

func JsonParse(filename string) []Place {
	data, err := os.ReadFile(filename)

	helpers.Check(err)

	var Places []Place

	err = json.Unmarshal(data, &Places)

	helpers.Check(err)

	return Places
}

func MakeRequest(es *elasticsearch.Client, RequestFunc func(es *elasticsearch.Client) (*esapi.Response, error)) {
	res, err := RequestFunc(es)
	helpers.Check(err)

	helpers.CheckResponse(res, "Request success", "Error request")

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

func InsertData(es *elasticsearch.Client, Places []Place) {

	var data strings.Builder

	for i := range Places {
		tmp, err := json.Marshal(Places[i])
		helpers.Check(err)
		data.Write([]byte(fmt.Sprintf(
			"{ \"index\" : { \"_index\" : \"places\" ,  \"_id\" : \"%d\"} }\n", Places[i].Id)))
		data.Write(tmp)
		data.WriteString("\n")
	}
	blk, err := es.Bulk(strings.NewReader(data.String()), es.Bulk.WithIndex("places"))

	helpers.Check(err)

	defer blk.Body.Close()

	helpers.CheckResponse(blk, "Index request success", "Error request")
}
