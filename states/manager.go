package states

import (
	"database/sql"
	"errors"
	"github.com/thomasdomingos/terraform-state-server/config"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
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
			log.Println("registry directory does not not exist, creating it")
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
	log.Println("closing State Manager")
	if nil == m.db {
		return nil
	}

	return m.db.Close()
}

func (m *Mgr) GetState(name string) ([]byte, error) {
	exists, id, err := getState(m.db, name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("state does not exists")
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

func (m *Mgr) PutState(name string, content []byte) error {
	// Try to recover current state
	exists, id, err := getState(m.db, name)
	if err != nil {
		return err
	}
	var state *State
	// Create state as next or new depending of predecessor existence
	if !exists {
		// Verify the State directory exists
		if err := assertDirExists(filepath.Join(m.cfg.Registry.Path, name)); err != nil {
			return err
		}
		state = NewState(name, content)
	} else {
		oldState := State{Name: name, Previous: id}
		state = NextState(oldState, content)
	}

	log.Println("writing state on disk...")
	// Write state content to file
	err = ioutil.WriteFile(filepath.Join(m.cfg.Registry.Path, name, state.Checksum), content, 0644)
	if err != nil {
		return err
	}
	log.Println("done writing state on disk")
	// Finally insert state into the DB
	if err := insertState(m.db, *state); err != nil {
		return err
	}
	return nil
}

func assertDirExists(path string) error {
	// Test directory containing state
	if _, err := os.Stat(path); err != nil {
		// Create directory that will contain state file(s)
		if os.IsNotExist(err) {
			log.Println("directory does not not exist, creating it")
			if err := os.MkdirAll(path, os.ModePerm); err != nil {
				return err
			}
		}
	}
	return nil
}
