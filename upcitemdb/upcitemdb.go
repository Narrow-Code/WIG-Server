package upcitemdb

import (
	db "WIG-Server/config"
	"WIG-Server/models"
	"fmt"
	"compress/gzip"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func GetBarcode(barcode string) {
	fmt.Println("Getting barcode")
	url := "https://api.upcitemdb.com/prod/trial/lookup?upc=" + barcode

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Encoding", "gzip,deflate")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var reader io.ReadCloser

	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		defer reader.Close()
	default:
		reader = resp.Body
	}

	fmt.Println("api reached")

	var data map[string]interface{}
	decoder := json.NewDecoder(reader)
	err = decoder.Decode(&data)
	if err != nil {
		log.Fatal(err)
	}
	
	items := data["items"].([]interface{})
	for _, item := range items {
		itemData := item.(map[string]interface{})
		newItem := models.Item{
			Barcode: barcode,
			Name: itemData["title"].(string),
			Brand: itemData["brand"].(string),
			Image: itemData["images"].([]interface{})[0].(string),
			ItemDesc: itemData["description"].(string),
		}
		db.DB.Create(&newItem)
	}


}
