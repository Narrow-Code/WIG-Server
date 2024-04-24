// db provides functionalities for configuring connecting to the database.
package db

import (
	"WIG-Server/models"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const DefaultBorrowerUUID = "11111111-1111-1111-1111-111111111111"
const SelfBorrowerUUID = "22222222-2222-2222-2222-222222222222"
const DefaultUserUUID = "33333333-3333-3333-3333-333333333333"
const DefaultLocationUUID = "44444444-4444-4444-4444-444444444444"
const DefaultItemUUID = "55555555-5555-5555-5555-555555555555"

// DB holds the database connection instance.
var DB *gorm.DB

// Connect establishes a connection to the database.
func Connect() {
	// Load environment variables 
	godotenv.Load()

	// Set environment variables
	dbhost := os.Getenv("MYSQL_HOST")
	dbuser := os.Getenv("MYSQL_USER")
	dbpassword := os.Getenv("MYSQL_PASSWORD")
	dbname := os.Getenv("MYSQL_DBNAME")

	connection := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbuser, dbpassword, dbhost, dbname)
	
	var db *gorm.DB
	var err error
	retries := 5

	for retries > 0 {
		db, err = gorm.Open(mysql.Open(connection), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err == nil {
			break
		}
		fmt.Println("Database connection failed. Retrying in 5 seconds...")
		time.Sleep(5 * time.Second)
		retries--
	}

	if err != nil {
		panic("Database connection failed after multiple retries")
	}

	DB = db
	fmt.Println("db connected successfully")

	AutoMigrate(db)
}

/*
* AutoMigrate performs automatic migrations on the provided connection.
*
* @param connection The database connection instance on which the migrations will be applied.
 */
func AutoMigrate(connection *gorm.DB) {
	connection.Debug().AutoMigrate(
		&models.User{},
		&models.Item{},
		&models.Borrower{},
		&models.Location{},
		&models.Ownership{},
	)
	
	ensureDefaultRecords(connection)
}

// ensureDefaultRecords checks if essential tables are empty and creates default records if necessary
func ensureDefaultRecords(connection *gorm.DB) {
	ensureBorrowerRecords(connection)
	ensureUserRecords(connection)
	ensureLocationRecords(connection)
	ensureItemRecords(connection)
}

// ensureBorrowerRecords checks if Borrower table is empty and creates default records if necessary
func ensureBorrowerRecords(connection *gorm.DB) {
	var borrowerCount int64
	connection.Model(&models.Borrower{}).Count(&borrowerCount)

	if borrowerCount == 0 {
		defaultBorrower := models.Borrower{
			BorrowerName: "Default",
			BorrowerUID: uuid.MustParse(DefaultBorrowerUUID)}
		connection.Create(&defaultBorrower)

		selfBorrower := models.Borrower{
			BorrowerName: "Self",
			BorrowerUID: uuid.MustParse(SelfBorrowerUUID)}
		connection.Create(&selfBorrower)
	}
}

// ensureUserRecords checks if User table is empty and creates default records if necessary
func ensureUserRecords(connection *gorm.DB) {
	var userCount int64
	connection.Model(&models.User{}).Count(&userCount)

	if userCount == 0 {
		defaultUser := models.User{
			UserUID:  uuid.MustParse(DefaultUserUUID),
			Username: "Default User"}
		connection.Create(&defaultUser)
	}
}

// ensureLocationRecords checks if Location table is empty and creates default records if necessary
func ensureLocationRecords(connection *gorm.DB) {
	var locationCount int64
	connection.Model(&models.Location{}).Count(&locationCount)

	if locationCount == 0 {
		defaultLocation := models.Location{
			LocationUID:   uuid.MustParse(DefaultLocationUUID),
			LocationName:  "Default Location",
			LocationOwner: uuid.MustParse(DefaultUserUUID)}
		connection.Create(&defaultLocation)
	}
}

// ensureItemRecords checks if Item table is empty and creates default records if necessary
func ensureItemRecords(connection *gorm.DB) {
	var itemCount int64
	connection.Model(&models.Item{}).Count(&itemCount)

	if itemCount == 0 {
		defaultItem := models.Item{
			ItemUid: uuid.MustParse(DefaultItemUUID),
			Name: "Default Item"}
		connection.Create(&defaultItem)
	}
}

/*
* GetPort retrieves the Port to be used in the .env file.
*
* string The Port to be used.
*/
func GetPort() string {
	godotenv.Load()
	var port = os.Getenv("PORT")
	if port == "" {
		port = "80"
	}
	return port
}
