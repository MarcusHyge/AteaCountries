package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
)

const apiurl = "https://countries.trevorblades.com/"
const apimethod = http.MethodPost

type MainData struct {
	Data struct {
		Countries []struct {
			Name      string `json:"name"`
			Languages []struct {
				Name string `json:"name"`
			} `json:"languages"`
		} `json:"countries"`
	} `json:"data"`
}

func GetCountriesByCurrency() bool {
	// Make body
	jsonData := map[string]string{
		"query": `
        {
            countries(filter: { currency: { eq: "EUR" } }) {
              name
              languages {
                name
              }
            }
          }
        `,
	}

	//Create a request with body
	request := CreateRequest(jsonData)

	//Make the request and gets the response back in a *http.Response
	response := RequestData(request)

	//Closing the body to prevent resource leak
	defer response.Body.Close()

	//Read the body
	data := ReadResponseBody(response)

	//Structure the data into MainData struct
	result := JsonDataToStruct(data)

	//Group the countries by name
	hashMap := GroupCountriesByName(result)

	//Print out the values of the hashMap
	PrintOutValuesOfMap(hashMap)
	return true
}

func ConvertToJson(jsonData map[string]string) []byte {
	body, err := json.Marshal(jsonData)
	if err != nil {
		panic(err)
	}
	return body
}

func ReadResponseBody(response *http.Response) []byte {
	data, _ := ioutil.ReadAll(response.Body)
	return data
}

func RequestData(request *http.Request) *http.Response {
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Errorf("Response error:", err.Error())
	}
	return response
}

func CreateRequest(jsonData map[string]string) *http.Request {
	request, err := http.NewRequest(http.MethodPost, apiurl, bytes.NewBuffer(ConvertToJson(jsonData)))
	if err != nil {
		fmt.Errorf("Error creating new request: ", err.Error())
	}
	//Set content type
	request.Header.Add("content-type", "application/json")
	return request
}

func JsonDataToStruct(data []byte) MainData {
	var result MainData
	if err := json.Unmarshal(data, &result); err != nil {
		fmt.Errorf("Can not unmarshal JSON:", err.Error())
	}
	return result
}

func GroupCountriesByName(result MainData) map[string][]string {
	/*
		Making map to group values.
		Key Values = Languages
		Elements = Countries
	*/
	hashMap := make(map[string][]string)
	//For each country object
	for _, Country := range result.Data.Countries {
		//For each language object
		for _, Language := range Country.Languages {
			//Setting country name to correct key value.
			CountryList := hashMap[Language.Name]
			CountryList = append(CountryList, Country.Name)
			hashMap[Language.Name] = CountryList
		}
	}
	return hashMap
}

func PrintOutValuesOfMap(hashMap map[string][]string) {
	for key, element := range hashMap {
		fmt.Println("Country:", key, "=>", "Languages:", element)
	}
}

func main() {
	getCountriesByCurrency := flag.Bool("getCountriesByCurrency", false, "Test")
	flag.Parse()

	if *getCountriesByCurrency {
		GetCountriesByCurrency()
		return
	}
}
