package storage

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

// Interface for defining the methods of the key-value storage
type KeyValueStorage interface {
	Set(key, value string) error
	GetKeyByValue(value string) (string, error)
}

// In-memory implementation of the key-value storage
type InMemoryKeyValueStorage struct {
	data map[string]string
}

func NewInMemoryKeyValueStorage() *InMemoryKeyValueStorage {
	return &InMemoryKeyValueStorage{
		data: make(map[string]string),
	}
}

func (s *InMemoryKeyValueStorage) Set(key, value string) error {
	if _, ok := s.data[key]; ok {
		return fmt.Errorf("Key already exists")
	}
	s.data[key] = value
	return nil
}

func (s *InMemoryKeyValueStorage) GetKeyByValue(value string) (string, error) {
	for key, val := range s.data {
		if val == value {
			return key, nil
		}
	}
	return "", fmt.Errorf("Key not found")
}

// PostgreSQL implementation of the key-value storage
type PostgreSQLKeyValueStorage struct {
	db *sql.DB
}

func NewPostgreSQLKeyValueStorage() (*PostgreSQLKeyValueStorage, error) {
	db, err := sql.Open("postgres", "user=postgres password=your_password dbname=gRPC-api sslmode=disable")
	if err != nil {
		return nil, err
	}

	return &PostgreSQLKeyValueStorage{
		db: db,
	}, nil
}

func (s *PostgreSQLKeyValueStorage) Set(key, value string) error {
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM data WHERE key = $1", key).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("Key already exists")
	}

	_, err = s.db.Exec("INSERT INTO data (key, value) VALUES ($1, $2)", key, value)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostgreSQLKeyValueStorage) GetKeyByValue(value string) (string, error) {
	var key string
	err := s.db.QueryRow("SELECT key FROM data WHERE value = $1", value).Scan(&key)
	if err != nil {
		return "", err
	}
	return key, nil
}

func UseKeyValueStorage(storage KeyValueStorage) {
	/*err := storage.Set("", "")
	if err != nil {
		log.Fatal(err)
	}

	key, err := storage.GetKeyByValue("")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(key)*/
}
