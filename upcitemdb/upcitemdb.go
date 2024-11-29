// upcitemdb handles the API calls with upcitemdb.com.
package upcitemdb

import (
	"WIG-Server/db"
	"WIG-Server/models"
	"WIG-Server/utils"
	"compress/gzip"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
)

/*
 * GetBarcode performs the GET Barcode API call with upcitemdb.com.
 * If an item is retrieved, it is then added to the Items table in the database.
 *
 * @param barcode The barcode to retrieve data for.
 * @return An integer indicating the status of the operation.
 */
func GetBarcode(barcode string) int {
    // Construct URL
    utils.Log("began call")
    url := constructURL(barcode)

    // Send request and process response
    data, err := fetchData(url)
    if err != nil {
	utils.Log("Error is returned")
        return 0
    }

    // Process retrieved items
    if items, exists := data["items"]; exists {
        createItems(barcode, items)
    }

    return 0
}

/*
 * constructURL constructs the API URL for the given barcode.
 *
 * @param barcode The barcode to construct the URL for.
 * @return The constructed URL as a string.
 */
func constructURL(barcode string) string {
	upcItemCheck := os.Getenv("UPC_ITEM_DB")

	if upcItemCheck == "trial" {
		return "https://api.upcitemdb.com/prod/trial/lookup?upc=" + barcode
	} else {
		return "https://api.upcitemdb.com/prod/v1/lookup?upc=" + barcode
	}
}

/*
 * fetchData sends an HTTP request and returns the decoded JSON data.
 *
 * @param url The URL to send the request to.
 * @return A map containing the decoded JSON data.
 * @return An error if the request fails or decoding is unsuccessful.
 */
func fetchData(url string) (map[string]interface{}, error) {
    utils.Log("fetching data from upcitemdb")
    userKey := os.Getenv("UPC_ITEM_DB")
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        log.Fatal(err)
    }

    // Set request headers
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Accept", "application/json")
    req.Header.Set("Accept-Encoding", "gzip,deflate")

    // If not trial add token to header
    if userKey != "trial" {
    	req.Header.Set("user_key", userKey)
	req.Header.Set("key_type", "3scale")
    }

    // Send request
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    utils.Log("response from upcitemdb successful")

    // Check for rate limit
    if resp.StatusCode == 429 {
        return nil, errors.New("rate limit exceeded")
    }

    // Decode response based on Content-Encoding
    var reader io.ReadCloser
    switch resp.Header.Get("Content-Encoding") {
    case "gzip":
        reader, err = gzip.NewReader(resp.Body)
        if err != nil {
            utils.Log("failed to decompress response") 
            return nil, err
        }
        defer reader.Close()
    default:
        reader = resp.Body
    }

    // Decode JSON response
    utils.Log("decoding json response")
    var data map[string]interface{}
    decoder := json.NewDecoder(reader)
    err = decoder.Decode(&data)
    if err != nil {
        return nil, err
    }
    utils.Log("decoding successful")

    return data, nil
}

/*
 * createItems processes the retrieved items and adds them to the database.
 *
 * @param barcode The barcode for which the items are retrieved.
 * @param items The items retrieved from the API response.
 */
func createItems(barcode string, items interface{}) {
    utils.Log("began call")
    for _, item := range items.([]interface{}) {
        itemData := item.(map[string]interface{})
        var newItem models.Item

        // Set barcode
        newItem.Barcode = barcode
	utils.Log(newItem.Barcode + " set as barcode")

        // Set name if available
        if title, exists := itemData["title"]; exists {
            newItem.Name = title.(string)
	    utils.Log(newItem.Name + " set as name")
        }

        // Set brand if available
        if brand, exists := itemData["brand"]; exists {
            newItem.Brand = brand.(string)
	    utils.Log(newItem.Brand + " set as brand")
        }

        // Set image if available
        if images, exists := itemData["images"]; exists && len(images.([]interface{})) > 0 {
            newItem.Image = images.([]interface{})[0].(string)
	    utils.Log(newItem.Image + " set as image")
        }

        // Generate unique identifier
        newItem.ItemUid = uuid.New()
	utils.Log("UUID set")

        // Create item in database
        db.DB.Create(&newItem)
	utils.Log(newItem.Name + " has been successfully added to database")
    }
}

