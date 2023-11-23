/*
* Package config provides functionalities for configuring and connecting to the database.
*
* It includes functions to establish a connection, perform auto migrations, and manage the database connection instance.
 */
package db

import (
	"WIG-Server/models"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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
	var db, err = gorm.Open(mysql.Open(connection), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		panic("Database connection failed")
	}

	DB = db
	fmt.Println("db connected successfully")

	AutoMigrate(db)
}

/*
* AutoMigrate performs automatic migrations on the provided connection.
*
* @param connection *gorm.DB - The database connection instance on which the migrations will be applied.
 */
func AutoMigrate(connection *gorm.DB) {
	connection.Debug().AutoMigrate(
		&models.User{},
		&models.Item{},
		&models.Borrower{},
		&models.Location{},
		&models.Ownership{},
	)

	// Check if Borrower table is empty
	var borrowerCount int64
	connection.Model(&models.Borrower{}).Count(&borrowerCount)

	if borrowerCount == 0 {
		// Create a default Borrower record
		defaultBorrower := models.Borrower{
			BorrowerUID:  1,
			BorrowerName: "Default Borrower"}
		connection.Create(&defaultBorrower)
	}

	// Check if User table is empty
	var userCount int64
	connection.Model(&models.User{}).Count(&userCount)

	if userCount == 0 {
		// Create a default User record
		defaultUser := models.User{
			UserUID:  1,
			Username: "Default User"}
		connection.Create(&defaultUser)
	}

	// Check if Location table is empty
	var locationCount int64
	connection.Model(&models.Location{}).Count(&locationCount)

	if locationCount == 0 {
		// Create a default Location record
		defaultLocation := models.Location{
			LocationUID:   1,
			LocationName:  "Default Location",
			LocationOwner: 1}
		connection.Create(&defaultLocation)
	}

}

func GetPort() string {
	godotenv.Load()
	var port = os.Getenv("PORT")
	if port == "" {
		port = "80"
	}
	return port
}
