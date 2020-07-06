package states

import (
	"database/sql"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/thomasdomingos/terraform-state-server/config"
)

type Mgr struct {
	db  *sql.DB
	cfg config.Config
}

func (m *Mgr) Init(cfg config.Config) error {
	// Initialize database
	db, err := initDB(cfg.Database.Path)
	if err != nil {
		return err
	}
	m.db = db

	// Initialize directory to store states
	if err := assertDirExists(cfg.Registry.Path); err != nil {
		return err
	}
	if _, err := os.Stat(cfg.Registry.Path); err != nil {
		// Create directory that will contain state file(s)
		if os.IsNotExist(err) {
			log.Println("Registry directory does not not exist, creating it")
			if err := os.MkdirAll(cfg.Registry.Path, os.ModePerm); err != nil {
				m.db.Close()
				return err
			}
		}
	}
	m.cfg = cfg
	return nil
}

func (m *Mgr) Close() error {
	log.Println("Closing State Manager")
	if nil == m.db {
		return nil
	}
	return m.db.Close()
}

func (m *Mgr) GetState(name string) ([]byte, error) {
	id, err := getState(m.db, name)
	if err != nil {
		return nil, err
	}
	// Verify the State directory exists
	if err := assertDirExists(filepath.Join(m.cfg.Registry.Path, name)); err != nil {
		return nil, err
	}
	// Read state file and return its content
	file, err := os.Open(filepath.Join(m.cfg.Registry.Path, name, id))
	if err != nil {
		return nil, err
	}
	defer file.Close()
	b, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func assertDirExists(path string) error {
	// Test directory containing state
	if _, err := os.Stat(path); err != nil {
		// Create directory that will contain state file(s)
		if os.IsNotExist(err) {
			log.Println("Directory does not not exist, creating it")
			if err := os.MkdirAll(path, os.ModePerm); err != nil {
				return err
			}
		}
	}
	return nil
}
