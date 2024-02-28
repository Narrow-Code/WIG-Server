package upcitemdb

import (
	"WIG-Server/db"
	"testing"
	"os"
	"github.com/stretchr/testify/assert"
)

func TestGetBarcodePossitive(t *testing.T) {	
	os.Setenv("MYSQL_DBNAME", "wig")
	os.Setenv("MYSQL_USER", "wig")
	os.Setenv("MYSQL_PASSWORD", "wigsecret")
	os.Setenv("MYSQL_HOST", "localhost:3306")
	
	db.Connect()
	value := GetBarcode("305212345001")
	//value := GetBarcode("1234")
	assert.Equal(t, value, 0)
}
