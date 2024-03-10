package helpers

import (
	"fmt"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"strconv"
)

func Check(err error) {
	if err != nil {
		panic(err)
	}
}

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
