package states

import (
  "log"
  "database/sql"

	"github.com/thomasdomingos/terraform-state-server/config"
)

type Mgr struct {
  db *sql.DB
}

func (m *Mgr) Init(cfg config.Config) error {
  db, err := initDB(cfg.Database.Path)
  if err != nil {
    return err
  }
  m.db = db
  return nil
}

func (m *Mgr) Close() error {
  log.Println("Closing State Manager")
  if nil == m.db {
    return nil
  }
  return m.db.Close()
}


