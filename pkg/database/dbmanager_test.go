package database_test

import (
  "fmt"
  "testing"

  _ "github.com/mattn/go-sqlite3"

  "dunn-finance/pkg/database"
)

var _ = fmt.Println

func TestSQLite3_CreateDBInstance_Success(t *testing.T) {
  manager := &database.DBManager{ Driver: "sqlite3", DSN: ":memory:", }

  db1 := manager.GetDBInstance()
  defer db1.Close()

  if err := db1.Ping(); err != nil {
    t.Errorf("Ping failed after connecting: %v", err)
  }

  db2 := manager.GetDBInstance()

  if db1 != db2 {
    t.Errorf("Expected the same connection. But got the different connection.")
  }
}

func TestSQLite3_Reconnect_Success(t *testing.T) {
  manager := &database.DBManager{ Driver: "sqlite3", DSN: ":memory:", }

  db1 := manager.GetDBInstance()
  err := db1.Close()
  if err != nil { t.Fatalf("Failed to close the first connection: %v", err) }

  db2 := manager.GetDBInstance()
  defer db2.Close()

  if db1 == db2 {
    t.Errorf("Expected a new connection after close. But got the same connection")
  }
  if err := db2.Ping(); err != nil {
    t.Errorf("Ping failed after reconnecting: %v", err)
  }
}

func TestSQLite3_Exec_Success(t *testing.T) {
  manager := &database.DBManager{ Driver: "sqlite3", DSN: ":memory:", }
  db := manager.GetDBInstance()
  defer db.Close()

  query := "CREATE TABLE IF NOT EXISTS dummy (id INTEGER)"
  _, err := db.Exec(query)

  if err != nil {
    t.Errorf("Exec failed: %v", err)
  }
}

func TestSQLite3_Exec_Failure(t *testing.T) {
  manager := &database.DBManager{ Driver: "sqlite3", DSN: ":memory:", }
  db := manager.GetDBInstance()
  defer db.Close()

  query := "THIS IS NOT A VALID SQL"
  _, err := db.Exec(query)

  if err == nil {
    t.Errorf("Expected error for invalid SQL, got nil")
  }
}

func TestSQLite3_QueryRow_Success(t *testing.T) {
  manager := &database.DBManager{ Driver: "sqlite3", DSN: ":memory:", }
  db := manager.GetDBInstance()
  defer db.Close()

  query := "SELECT 1"
  row := db.QueryRow(query)
  var result int
  err := row.Scan(&result)

  if err != nil {
    t.Errorf("QueryRow failed: %v", err)
  }
  if result != 1 {
    t.Errorf("Expected 1, got %d", result)
  }
}

func TestSQLite3_QueryRow_Failure(t *testing.T) {
  manager := &database.DBManager{ Driver: "sqlite3", DSN: ":memory:", }
  db := manager.GetDBInstance()
  defer db.Close()

  query := "THIS IS NOT A VALID SQL"
  row := db.QueryRow(query)
  var result any
  err := row.Scan(&result)

  if err == nil {
    t.Errorf("Expected error for invalid SQL, got nil")
  }
}
