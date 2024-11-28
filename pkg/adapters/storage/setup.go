// package storage

// import (

// 	"gorm.io/driver/postgres"
// 	"gorm.io/gorm"
// 	"Questify/config"
// )

// // SetupDatabase initializes the database connection.
// func SetupDatabase(cfg config.DatabaseConfig) (*gorm.DB, error) {
// 	// Build the DSN (Data Source Name)
// 	dsn := "host=" + cfg.Host + " port=" + string(cfg.Port) + " user=" + cfg.User +
// 		" password=" + cfg.Pass + " dbname=" + cfg.DBName + " sslmode=disable"

// 	// Open the database connection
// 	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Ensure UUID extension is enabled
// 	sqlDB, err := db.DB() // Get the native *sql.DB object
// 	if err != nil {
// 		return nil, err
// 	}
// 	if _, err := sqlDB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\""); err != nil {
// 		log.Fatalf("Failed to enable uuid-ossp extension: %v", err)
// 	}

// 	log.Println("Database connection established and uuid-ossp extension enabled.")
// 	return db, nil
// }


package storage

import (
	"log"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"Questify/config"
)

// SetupDatabase initializes the database connection.
func SetupDatabase(cfg config.DatabaseConfig) (*gorm.DB, error) {
	// Safely format the port as a string
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Pass, cfg.DBName)

	// Open the database connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Ensure UUID extension is enabled
	sqlDB, err := db.DB() // Get the native *sql.DB object
	if err != nil {
		return nil, err
	}
	if _, err := sqlDB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\""); err != nil {
		log.Fatalf("Failed to enable uuid-ossp extension: %v", err)
	}

	log.Println("Database connection established and uuid-ossp extension enabled.")
	return db, nil

	return db, nil
}
