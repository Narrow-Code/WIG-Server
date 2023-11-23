package upcitemdb

import (
	"WIG-Server/db"
	"WIG-Server/models"
	"compress/gzip"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func GetBarcode(barcode string) {
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
		return
	}
	defer resp.Body.Close()

	var reader io.ReadCloser

	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			return
		}
		defer reader.Close()
	default:
		reader = resp.Body
	}

	var data map[string]interface{}
	decoder := json.NewDecoder(reader)
	err = decoder.Decode(&data)
	if err != nil {
		return
	}

	if items, exists := data["items"]; exists {
		for _, item := range items.([]interface{}) {
			itemData := item.(map[string]interface{})
			var newItem models.Item

			newItem.Barcode = barcode

			if title, exists := itemData["title"]; exists {
				newItem.Name = title.(string)
			}

			if brand, exists := itemData["brand"]; exists {
				newItem.Brand = brand.(string)
			}

			if images, exists := itemData["images"]; exists && len(images.([]interface{})) > 0 {
				newItem.Image = images.([]interface{})[0].(string)
			}
			db.DB.Create(&newItem)
		}

	}

}
