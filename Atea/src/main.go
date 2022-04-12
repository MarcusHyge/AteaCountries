package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
)

const APIURL = "https://countries.trevorblades.com/"
const APIMETHOD = http.MethodPost

type Country struct {
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
	request := createRequest(jsonData)

	//Make the request and gets the response back in a *http.Response
	response := requestData(request)

	//Closing the body to prevent resource leak
	defer response.Body.Close()

	//Read the body
	data := readResponseBody(response)

	//Structure the data into Country struct
	result := jsonDataToStruct(data)

	//Group the countries by name
	hashMap := groupCountriesByName(result)

	//Print out the values of the hashMap
	printOutValuesOfMap(hashMap)
	return true
}

func convertToJson(jsonData map[string]string) []byte {
	body, err := json.Marshal(jsonData)
	if err != nil {
		panic(err)
	}
	return body
}

func readResponseBody(response *http.Response) []byte {
	data, _ := ioutil.ReadAll(response.Body)
	return data
}

func requestData(request *http.Request) *http.Response {
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Errorf("Response error:", err.Error())
	}
	return response
}

func createRequest(jsonData map[string]string) *http.Request {
	request, err := http.NewRequest(APIMETHOD, APIURL, bytes.NewBuffer(convertToJson(jsonData)))
	if err != nil {
		fmt.Errorf("Error creating new request: ", err.Error())
	}
	//Set content type
	request.Header.Add("content-type", "application/json")
	return request
}

func jsonDataToStruct(data []byte) Country {
	var result Country
	if err := json.Unmarshal(data, &result); err != nil {
		fmt.Errorf("Can not unmarshal JSON:", err.Error())
	}
	return result
}

func groupCountriesByName(result Country) map[string][]string {
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

func printOutValuesOfMap(hashMap map[string][]string) {
	for key, element := range hashMap {
		fmt.Println("Language:", key, "=>", "Countries:", element)
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
