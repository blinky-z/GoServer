package main

import (
	"strconv"
	"encoding/json"
	"net/http"
	"fmt"
	"time"
)

type Item struct {
	Name  string `json:"name"`
	Price string `json:"price"`
}

func getItems(w http.ResponseWriter, r *http.Request) {
	const defaultGoodsNumber = 5

	var userName string

	if r.Method == "GET" {
		if values, ok := r.URL.Query()["name"]; ok {
			userName = values[0]
		}
	} else {
		switch r.Header.Get("Content-Type") {
		case "application/x-www-form-urlencoded":
			r.ParseForm()

			jsonBody := r.Form.Get("json")

			var body map[string]string

			json.Unmarshal([]byte(jsonBody), &body)

			if value, ok := body["name"]; ok {
				userName = value
			}
		default:
			r.ParseMultipartForm(2048)

			jsonBody := r.PostFormValue("json")

			var body map[string]string

			json.Unmarshal([]byte(jsonBody), &body)

			if value, ok := body["name"]; ok {
				userName = value
			}
		}
	}

	response := map[string]interface{}{}

	if userName != "" {
		var multiplier = 0

		items := make([]Item, len(userName))

		response["nickname"] = userName
		for _, charValue := range userName {
			multiplier += int(charValue)
		}

		for currentItemNumber := 0; currentItemNumber < len(items); currentItemNumber++ {
			newItem := Item{}
			newItem.Name = userName + strconv.Itoa(currentItemNumber)
			newItem.Price = strconv.Itoa((currentItemNumber + 1) * multiplier)

			items[currentItemNumber] = newItem
		}

		response["items"] = items

	} else {
		var multiplier = 30

		items := make([]Item, defaultGoodsNumber)

		for currentItemNumber := 0; currentItemNumber < len(items); currentItemNumber++ {
			newItem := Item{}
			newItem.Name = "default" + strconv.Itoa(currentItemNumber)
			newItem.Price = strconv.Itoa(currentItemNumber * multiplier)

			items[currentItemNumber] = newItem
		}

		response["items"] = items
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func buyItems(w http.ResponseWriter, r *http.Request) {
	var itemName string

	// set itemName
	switch r.Header.Get("Content-Type") {
	case "application/x-www-form-urlencoded":
		r.ParseForm()

		jsonBody := r.PostForm.Get("json")

		var item Item
		json.Unmarshal([]byte(jsonBody), &item)

		itemName = item.Name
	default:
		r.ParseMultipartForm(1024)

		jsonBody := r.PostFormValue("json")

		var item Item

		json.Unmarshal([]byte(jsonBody), &item)

		itemName = item.Name
	}

	response := make(map[string]string)

	successPurchaseMessage := "success"
	failurePurchaseMessage := "failure"

	if len(itemName)%2 == 0 {
		response["result"] = successPurchaseMessage
	} else {
		response["result"] = failurePurchaseMessage
	}

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func main() {
	mux := http.NewServeMux()

	server := &http.Server{ReadTimeout: 5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout: 120 * time.Second,
		Handler: mux,
	Addr:"localhost:8080"}

	mux.HandleFunc("/", getItems)
	mux.HandleFunc("/buy", buyItems)

	fmt.Printf("listening on address %s", server.Addr)
	server.ListenAndServe()
}
