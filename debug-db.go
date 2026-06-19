package main

import (
"log"
"github.com/omorojames5-prog/glass-guitar/pkg/database"
)

func main() {
database.ConnectDB()

db := database.DB

// Get the current database name
var dbName string
db.Raw("SELECT current_database()").Scan(&dbName)
log.Printf("Connected to database: %s", dbName)

// Check if the users table exists
var tableExists bool
db.Raw("SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'users')").Scan(&tableExists)
log.Printf("Users table exists: %v", tableExists)

// Check the actual columns
var columns []string
db.Raw("SELECT column_name FROM information_schema.columns WHERE table_name = 'users' ORDER BY ordinal_position").Scan(&columns)
log.Printf("Columns in users table: %v", columns)
}
