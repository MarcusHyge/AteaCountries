package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
)

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

	//Convert to json
	body, err := json.Marshal(jsonData)
	if err != nil {
		panic(err)
	}

	//Make request
	request, err := http.NewRequest("POST", "https://countries.trevorblades.com/", bytes.NewBuffer(body))
	if err != nil {
		fmt.Errorf("Error creating new request: ", err.Error())
	}

	//Set content type
	request.Header.Add("content-type", "application/json")

	//Request data
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Errorf("Response error:", err.Error())
	}
	//Close body
	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)

	//Put the json data in a struct
	var result MainData
	if err := json.Unmarshal(data, &result); err != nil {
		fmt.Errorf("Can not unmarshal JSON:", err.Error())
	}

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
	//Loop through the map and print the grouped values
	for key, element := range hashMap {
		fmt.Println("Country:", key, "=>", "Languages:", element)
	}
	return true
}

func main() {
	getCountriesByCurrency := flag.Bool("getCountriesByCurrency", false, "Test")
	flag.Parse()

	if *getCountriesByCurrency {
		GetCountriesByCurrency()
		return
	}
}
