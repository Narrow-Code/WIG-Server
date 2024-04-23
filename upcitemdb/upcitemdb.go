// Handles the API calls with upcitemdb.com.
package upcitemdb

import (
	"WIG-Server/db"
	"WIG-Server/models"
	"compress/gzip"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"
)

/*
* Performs the GET Barcode API call with upcitemdb.com.
* If an item is retrieved it is then added to the Items table in the database.
*
* @param barcode The barcode to retrieve data for.
 */
func GetBarcode(barcode string) int {
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
		return 0
	}
	defer resp.Body.Close()

	if resp.StatusCode == 429 {
		return 429
	}

	var reader io.ReadCloser

	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			log.Println("FAIL")
			return 0
		}
		defer reader.Close()
	default:
		reader = resp.Body
	}

	var data map[string]interface{}
	decoder := json.NewDecoder(reader)
	err = decoder.Decode(&data)
	if err != nil {
		return 0
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

			newItem.ItemUid = uuid.New()

			db.DB.Create(&newItem)
		}

	}
return 0
}
