package db

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

// DB is the global database connection
var DB *sql.DB

func InitDBFile() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	dbPath := filepath.Join(homeDir, ".quazaar", "quazaar.db")

	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	return dbPath, nil
}

func Init() error {
	path, err := InitDBFile()
	if err != nil {
		return err
	}

	var openErr error
	DB, openErr = sql.Open("sqlite3", path)
	if openErr != nil {
		return openErr
	}

	if err = DB.Ping(); err != nil {
		log.Println("❌ Error connecting to database:", err)
		return err
	}

	log.Println("✅ Database connected at", path)

	// Create all tables
	if err = createTables(); err != nil {
		log.Println("❌ Error creating tables:", err)
		return err
	}

	log.Println("✅ Database tables ready")
	return nil
}

// createTables creates all necessary tables in the database
// Optimized for single-user local server
func createTables() error {
	// Table 1: Single user (only one user for local server)
	// Using CHECK (id = 1) ensures only one user can exist
	userTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY CHECK (id = 1),
		name TEXT NOT NULL UNIQUE,
		pass TEXT NOT NULL,
		username TEXT NOT NULL UNIQUE
	);`

	// Table 2: Tokens for different services/devices
	// Each service or device gets its own token
	tokensTable := `
	CREATE TABLE IF NOT EXISTS tokens (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		tokenOf TEXT NOT NULL,
		tokenType TEXT NOT NULL,
		token TEXT NOT NULL UNIQUE,
		deviceId TEXT,
		expiry NUMERIC
	);`
	fileShareDeviceTokenTable := `
	CREATE TABLE IF NOT EXISTS file_share_device_tokens (
		id INTEGER PRIMARY KEY AUTOINCREMENT,  
		token TEXT NOT NULL UNIQUE,
		deviceId TEXT,
		expiry NUMERIC
	);`

	// Execute users table creation
	if _, err := DB.Exec(userTable); err != nil {
		log.Println("❌ Failed to create users table:", err)
		return err
	}
	log.Println("✅ Users table ready")

	// Execute tokens table creation
	if _, err := DB.Exec(tokensTable); err != nil {
		log.Println("❌ Failed to create tokens table:", err)
		return err
	}
	log.Println("✅ Tokens table ready")
	// Execute file share device tokens table creation
	if _, err := DB.Exec(fileShareDeviceTokenTable); err != nil {
		log.Println("❌ Failed to create file_share_device_tokens table:", err)
		return err
	}
	log.Println("✅ File Share Device Tokens table ready")

	// Create indexes for faster searches
	indexSQL := `
	CREATE INDEX IF NOT EXISTS idx_tokens_token ON tokens(token);
	CREATE INDEX IF NOT EXISTS idx_tokens_type ON tokens(tokenType);
	`

	if _, err := DB.Exec(indexSQL); err != nil {
		log.Println("⚠️ Failed to create indexes:", err)
	}
	log.Println("✅ Indexes created")

	return nil
}

// CloseDB closes the database connection gracefully
func CloseDB() error {
	if DB != nil {
		log.Println("🔒 Closing database connection...")
		log.Println("🔒 Closing database connection...")
		log.Println("🔒 Closing database connection...")
		log.Println("🔒 Closing database connection...")
		return DB.Close()
	}
	return nil
}

func GetDB() *sql.DB {
	return DB
}
