package storage

import (
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"os"
	"testing"
)

// MockDB interface for database operations
type MockDB interface {
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
}

// Mock implementation of the KeyValueStorage interface for testing
type MockKeyValueStorage struct {
	data   map[string]string
	mockDB MockDB
}

func (s *MockKeyValueStorage) Set(key, value string) error {
	if _, ok := s.data[key]; ok {
		return fmt.Errorf("Key already exists")
	}
	s.data[key] = value
	return nil
}

func (s *MockKeyValueStorage) GetKeyByValue(value string) (string, error) {
	for key, val := range s.data {
		if val == value {
			return key, nil
		}
	}
	return "", fmt.Errorf("Key not found")
}

func TestInMemoryKeyValueStorage(t *testing.T) {
	storage := NewInMemoryKeyValueStorage()

	// Test the Set method
	err := storage.Set("key1", "value1")
	if err != nil {
		t.Errorf("Unexpected error while setting key-value pair: %v", err)
	}

	// Test the GetKeyByValue method
	key, err := storage.GetKeyByValue("value1")
	if err != nil {
		t.Errorf("Unexpected error while getting key by value: %v", err)
	}
	if key != "key1" {
		t.Errorf("Unexpected key. Expected: key1, Got: %s", key)
	}
}

func TestPostgreSQLKeyValueStorage(t *testing.T) {
	// Create a mock *sql.DB object using sqlmock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	mockStorage := &PostgreSQLKeyValueStorage{
		db: db,
	}

	// Mock the database queries and expectations using sqlmock
	mock.ExpectExec("INSERT INTO data").
		WithArgs("key1", "value1").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectQuery("SELECT key FROM data").
		WithArgs("value1").
		WillReturnRows(sqlmock.NewRows([]string{"key"}).AddRow("key1"))

	// Test the Set method
	err = mockStorage.Set("key1", "value1")
	if err != nil {
		t.Errorf("Unexpected error while setting key-value pair: %v", err)
	}

	// Test the GetKeyByValue method
	key, err := mockStorage.GetKeyByValue("value1")
	if err != nil {
		t.Errorf("Unexpected error while getting key by value: %v", err)
	}
	if key != "key1" {
		t.Errorf("Unexpected key. Expected: key1, Got: %s", key)
	}

	// Verify that all expectations were met
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("Unfulfilled database expectations: %v", err)
	}
}

func TestMain(m *testing.M) {
	// Run the tests
	result := m.Run()

	// Perform any teardown or cleanup after the tests

	// Exit with the test result
	os.Exit(result)
}
