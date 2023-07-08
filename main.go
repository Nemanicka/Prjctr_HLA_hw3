package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"time"
)

var cryptoExUrl = "https://min-api.cryptocompare.com/data/price?fsym=ETH&tsyms=USD"
var gaURL = "https://www.google-analytics.com/mp/collect?api_secret=%s&measurement_id=%s"

func getUSDToETHRate() float32 {
	resp, rateErr := http.Get(cryptoExUrl)
	if rateErr != nil {
		log.Fatal("Failed to fetch currency rates", rateErr)
	}

	defer resp.Body.Close()

	data := make(map[string]float32)
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Fatal("Error decoding response body:", err)
		return 0
	}

	return data["USD"]
}

func reportRateToGA(rate float32) {

	data := map[string]interface{}{
		"client_id":            "oleksii_prjctr_hw3",
		"non_personalized_ads": false,
		"events": []interface{}{
			map[string]interface{}{
				"name": "ETH_USD_Rate",
				"params": map[string]interface{}{
					"Rate": math.Ceil(rate)
				},
			},
		},
	}

	jsonData, err := json.Marshal(data)
	fmt.Println(string(jsonData))
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	gaURLFormatted := fmt.Sprintf(gaURL, os.Getenv("GA_SECRET"), os.Getenv("MEASUREMENT_ID"))
	fmt.Println(gaURLFormatted)

	req, err := http.NewRequest("POST", gaURLFormatted, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()
	fmt.Println(time.Now(), "Response status = ", resp.StatusCode)
}

func main() {

	for {
		rate := getUSDToETHRate()

		reportRateToGA(rate)

		time.Sleep(10 * time.Second)
	}

	// fmt.Println("API_ID = ", os.Getenv("OPEN_EXCHANGE_APP_ID"))
	// client := dinero.NewClient(
	// 	os.Getenv("OPEN_EXCHANGE_APP_ID"),
	// 	"USD",
	// 	1*time.Hour,
	// )

	// client.Rates.SetBaseCurrency("USD")
	// rates, ratesError := client.Rates.List()
	// if ratesError != nil {
	// 	log.Fatalf("Failed to fetch currency rates %f \n", ratesError)
	// }

	// fmt.Println(rates)

	// r, e := client.Rates.Get("USD")
	// if e != nil {
	// 	log.Fatalf("Failed to fetch currency rate %f \n", e)
	// }

	// fmt.Println(*r)

	// Create a sample JSON payload
	// data := map[string]interface{} {
	// 	"client_id": "test123",
	// 	"non_personalized_ads": false,
	// 		"events": [
	// 			{
	//   				"name": "testEvent",
	//   				"params": {
	//     				"TestParam": 1
	//   				}
	// 			},
	// 		]
	// 	}

	// Convert data to JSON
	// jsonData, err := json.Marshal(data)
	// if err != nil {
	// 	fmt.Println("Error marshaling JSON:", err)
	// 	return
	// }

	// // Create a POST request with the JSON payload
	// req, err := http.NewRequest("POST", "http://example.com/api", bytes.NewBuffer(jsonData))
	// if err != nil {
	// 	fmt.Println("Error creating request:", err)
	// 	return
	// }

	// Set the Content-Type header to application/json
	// req.Header.Set("Content-Type", "application/json")

	// // Send the request
	// client := &http.Client{}
	// resp, err := client.Do(req)
	// if err != nil {
	// 	fmt.Println("Error sending request:", err)
	// 	return
	// }
	// defer resp.Body.Close()

	// Check the response status code
	// if resp.StatusCode != http.StatusOK {
	// 	fmt.Println("Request failed with status:", resp.StatusCode)
	// 	return
	// }

	// fmt.Println("Request successful!")
	// TODO: Process the response body here
}
