package dbdao

import (
  "database/sql"
  "fmt"
  "log"
  "sync"
)

var _ = fmt.Println

type DBConnector interface {
  Ping() error
  Close() error
  GetRawDB() *sql.DB

  Exec(query string, args ...any) (sql.Result, error)
  QueryRow(query string, args ...any) *sql.Row
}

type SQLConnector struct {
  db *sql.DB
}

func (c *SQLConnector) Ping() error {
  return c.db.Ping()
}

func (c *SQLConnector) Close() error {
  return c.db.Close()
}

func (c *SQLConnector) GetRawDB() *sql.DB {
  return c.db
}

func (c *SQLConnector) Exec(query string, args ...any) (sql.Result, error){
  return c.db.Exec(query, args...)
}

func (c *SQLConnector) QueryRow(query string, args ...any) *sql.Row {
  return c.db.QueryRow(query, args...)
}

type DBManager struct {
  Driver string
  DSN    string
  mu     sync.Mutex
}

var (
  dbInstance DBConnector
  once       sync.Once
)

func (m *DBManager) GetDBInstance() DBConnector {
  m.mu.Lock()
  defer m.mu.Unlock()

  if dbInstance != nil {
    if err := dbInstance.Ping(); err == nil {
      return dbInstance
    }

    // Ping failure -> DB is closed
    _ = dbInstance.Close()
    dbInstance = nil
  }

  conn, err := sql.Open(m.Driver, m.DSN)
  if err != nil {
    log.Fatalf("Failed to re-connect: %v", err)
  }
  dbInstance = &SQLConnector{db: conn}

  return dbInstance
}

// func NewDBManagerWithConnector(driver, dsn string, conn DBConnector) (*DBManager, error) {
//   if err := conn.Ping(); err != nil {
//     return nil, err
//   }
// 
//   return &DBManager{
//     DB:     conn,
//     Driver: driver,
//     DSN:    dsn,
//   }, nil
// }
// 
// func NewDBManager(driver, dsn string) (*DBManager, error) {
//   db, err := sql.Open(driver, dsn)
//   if err != nil { return nil, err }
// 
//   conn := &SQLConnector{db: db}
//   return NewDBManagerWithConnector(driver, dsn, conn)
// }
