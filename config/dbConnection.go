/*
* Package config provides functionalities for configuring and connecting to the database.
*
* It includes functions to establish a connection, perform auto migrations, and manage the database connection instance.
*/

package config

import (
	"WIG-Server/models"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
	"gorm.io/driver/mysql"
	"os"
	"fmt"
)

// DB holds the database connection instance.
var DB *gorm.DB

/*
* Connect establishes a connection to the database.
* 
* It loads environment variables for database configuration, creates a connection string, and initializes the database connection instance.
*/
func Connect() {
	godotenv.Load()
	dbhost := os.Getenv("MYSQL_HOST")
	dbuser := os.Getenv("MYSQL_USER")
	dbpassword := os.Getenv("MYSQL_PASSWORD")
	dbname := os.Getenv("MYSQL_DBNAME")

	connection := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbuser, dbpassword, dbhost, dbname)
	var db, err = gorm.Open(mysql.Open(connection), &gorm.Config{})

	if err != nil {
		panic("Database connection failed")
	}

	DB = db
	fmt.Println("db connected successfully")

	AutoMigrate(db)
}

/*
* AutoMigrate performs automatic migrations on the provided connection.
*/
func AutoMigrate(connection *gorm.DB) {
	connection.Debug().AutoMigrate(
		&models.User{},
		&models.Item{},
		&models.Borrower{},
		&models.Location{},
		&models.Ownership{},
	)
}

