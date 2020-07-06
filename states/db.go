package states

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
)

func initDB(path string) (*sql.DB, error) {
	// Verify the DB file exists
	if info, err := os.Stat(path); os.IsNotExist(err) {
		// Create if the file does not exists
		log.Println("Creating db", path)
		file, err := os.Create(path) // Create SQLite file
		if err != nil {
			log.Fatal(err.Error())
		}
		file.Close()
		log.Println(path, "created")
	} else if info.IsDir() {
		log.Fatal(path, " exists but is a direcroty")
	}

	// Create table
	sqliteDatabase, _ := sql.Open("sqlite3", path) // Open the created SQLite File
	createTables(sqliteDatabase)                   // Create Database Tables
	return sqliteDatabase, nil                     // Defer Closing the database
}

func createTables(db *sql.DB) {
	createStateTableSQL := `CREATE TABLE IF NOT EXISTS state (
    "id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
    "name" TEXT,
    "checksum" TEXT,
    "previous" TEXT,
    CONSTRAINT state_unique UNIQUE (name, checksum)
  );` // SQL Statement for Create Table

	log.Println("Create state table...")
	statement, err := db.Prepare(createStateTableSQL) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec() // Execute SQL Statements
	log.Println("state table created")
}

// We are passing db reference connection from main to our method with other parameters
func insertState(db *sql.DB, state State) error {
	log.Println("Inserting state record ...")
	insertStateSQL := `INSERT INTO state(name, checksum, previous) VALUES (?, ?, ?)`
	statement, err := db.Prepare(insertStateSQL) // Prepare statement. This is good to avoid SQL injections
	if err != nil {
		log.Println("error :", err.Error())
		return err
	}
	_, err = statement.Exec(state.Name, state.Checksum, state.Previous)
	if err != nil {
		log.Println("error :", err.Error())
		return err
	}
	return nil
}

func getState(db *sql.DB, name string) (string, error) {
	log.Println("Searching for", name, "state")
	row := db.QueryRow("SELECT * FROM state WHERE name=? ORDER BY id DESC LIMIT 1", name)

	var id int
	var getname string
	var checksum string
	var previous string
	if err := row.Scan(&id, &getname, &checksum, &previous); err != nil {
		log.Println("error", err.Error())
		return "", err
	}
	log.Println("State: ", getname, " ", checksum)
	return checksum, nil
}

func getAllStates(db *sql.DB) []string {
	log.Println("Get all states inside DB")
	rows, _ := db.Query("SELECT DISTINCT name FROM state")
	defer rows.Close()
	names := make([]string, 0)
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			log.Fatal(err)
		}
		names = append(names, name)
	}
	return names
}
