package states

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
	"log"
	"os"
	"strconv"
)

func initDB(path string) (*sql.DB, error) {
	// Verify the DB file exists
	if info, err := os.Stat(path); os.IsNotExist(err) {
		// Create if the file does not exists
		log.Println("creating db", path)
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
    "previous" TEXT
  );` // SQL Statement for Create Table

	log.Println("create state table...")
	statement, err := db.Prepare(createStateTableSQL) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec() // Execute SQL Statements
	log.Println("done creating state table")

	createLockTableSQL := `CREATE TABLE IF NOT EXISTS lock (
    "name" TEXT NOT NULL PRIMARY KEY,
    "lock" INTEGER
  );` // SQL Statement for Create Table

	log.Println("create lock table...")
	statement, err = db.Prepare(createLockTableSQL) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec() // Execute SQL Statements
	log.Println("done creating lock table")
}

func bool2int(b bool) int {
	if b {
		return 1
	} else {
		return 0
	}
}

func setLockState(db *sql.DB, name string, lock bool) error {
	log.Println("setting lock state to", strconv.FormatBool(lock))
	setLockStateSQL := `INSERT INTO lock(name, lock) VALUES (?, ?) ON CONFLICT(name) DO UPDATE SET lock=?`
	statement, err := db.Prepare(setLockStateSQL) // Prepare statement. This is good to avoid SQL injections
	if err != nil {
		log.Println("error :", err.Error())
		return err
	}
	_, err = statement.Exec(name, bool2int(lock), bool2int(lock))
	if err != nil {
		log.Println("error :", err.Error())
		return err
	}
	log.Println("done setting lock state")
	return nil
}

func getLockState(db *sql.DB, name string) (bool, error) {
	log.Println("searching lock state for", name)
	row := db.QueryRow("SELECT lock FROM lock WHERE name=?", name)
	// Parse the resulting row
	var lock int
	err := row.Scan(&lock)
	switch {
	case err == sql.ErrNoRows:
		// the state does not exists so it is not locked
		return false, nil
	case err != nil:
		return false, err
	}
	return lock != 0, nil
}

// We are passing db reference connection from main to our method with other parameters
func insertState(db *sql.DB, state State) error {
	log.Println("inserting state record ...")
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
	log.Println("done inserting state")
	return nil
}

func getState(db *sql.DB, name string) (bool, string, error) {
	log.Println("searching for", name, "state")
	row := db.QueryRow("SELECT * FROM state WHERE name=? ORDER BY id DESC LIMIT 1", name)
	// Parse the resulting row
	var id int
	var getname string
	var checksum string
	var previous string
	exists := true
	err := row.Scan(&id, &getname, &checksum, &previous)
	switch {
	case err == sql.ErrNoRows:
		// the state does not exists
		exists = false
	case err != nil:
		return false, "", err
	}
	if exists {
		log.Println("GetState:", getname, "(", checksum, ")")
	} else {
		log.Println("cannot find state", getname)
	}
	return exists, checksum, nil
}

func getAllStates(db *sql.DB) []string {
	log.Println("get all states inside DB")
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

func getHistory(db *sql.DB, name string) []string {
	log.Println("get history for state", name)
	rows, _ := db.Query("SELECT checksum FROM state WHERE name=? ORDER BY id DESC", name)
	defer rows.Close()
	layers := make([]string, 0)
	for rows.Next() {
		var checksum string
		if err := rows.Scan(&checksum); err != nil {
			log.Fatal(err)
		}
		layers = append(layers, checksum)
	}
	return layers
}
