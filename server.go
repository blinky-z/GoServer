package main

import (
	"strconv"
	"encoding/json"
	"net/http"
	"fmt"
)

type Item struct {
	Name  string `json:"name"`
	Price string `json:"price"`
}

func getItems(w http.ResponseWriter, r *http.Request) {
	const defaultGoodsNumber= 5

	var userName string

	// set userName if it exists
	if r.Method == "GET" {
		if values, ok := r.URL.Query()["name"]; ok {
			userName = values[0]
		}
	} else {
		switch r.Header.Get("Content-Type") {
		case "x-www-form-urlencoded":
			if r.Form.Get("json") != "" {
				jsonBody := r.Form.Get("json")

				var body map[string]string

				json.Unmarshal([]byte(jsonBody), &body)

				if value, ok := body["name"]; ok {
					userName = value
				}
			}
		case "multipart/form-data":
			if r.FormValue("json") != "" {
				jsonBody := r.FormValue("json")

				var body map[string]string

				json.Unmarshal([]byte(jsonBody), &body)

				if value, ok := body["name"]; ok {
					userName = value
				}
			}
		}
	}

	// make response body
	responseBody := map[string]interface{}{}

	if userName != "" {
		var multiplier = 0

		items := make([]Item, len(userName))

		responseBody["nickname"] = userName
		for _, charValue := range userName {
			multiplier += int(charValue)
		}

		for currentItemNumber := 0; currentItemNumber < len(items); currentItemNumber++ {
			newItem := Item{}
			newItem.Name = userName + strconv.Itoa(currentItemNumber)
			newItem.Price = strconv.Itoa((currentItemNumber + 1) * multiplier)

			items[currentItemNumber] = newItem
		}

		responseBody["items"] = items

	} else {
		var multiplier = 30

		items := make([]Item, defaultGoodsNumber)

		for currentItemNumber := 0; currentItemNumber < len(items); currentItemNumber++ {
			newItem := Item{}
			newItem.Name = "default" + strconv.Itoa(currentItemNumber)
			newItem.Price = strconv.Itoa(currentItemNumber * multiplier)

			items[currentItemNumber] = newItem
		}

		responseBody["items"] = items
	}

	jsonResponseBody, _ := json.Marshal(responseBody)
	w.Write(jsonResponseBody)
}

func buyItems(w http.ResponseWriter, r *http.Request) {
	var itemName string

	// set itemName
	switch r.Header.Get("Content-Type") {
	case "x-www-form-urlencoded":
		jsonBody := r.Form.Get("json")

		var item Item
		json.Unmarshal([]byte(jsonBody), &item)

		itemName = item.Name
	case "multipart/form-data":
		jsonBody := r.FormValue("json")

		var item Item
		json.Unmarshal([]byte(jsonBody), &item)

		itemName = item.Name
	}

	responseBody := make(map[string]string)

	successPurchaseMessage := "success"
	failurePurchaseMessage := "failure"

	if len(itemName)%2 == 0 {
		responseBody["result"] = successPurchaseMessage
	} else {
		responseBody["result"] = failurePurchaseMessage
	}

	jsonBody, _ := json.Marshal(responseBody)
	w.Write(jsonBody)
}

func main() {
	PORT := "8080"

	mux := http.NewServeMux()

	mux.HandleFunc("/", getItems)
	mux.HandleFunc("/buy", getItems)

	fmt.Printf("listening on port %s", PORT)
	http.ListenAndServe(":" + PORT, mux)
}