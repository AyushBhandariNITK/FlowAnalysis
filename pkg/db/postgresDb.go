package db

import (
	"database/sql"
	"fmt"
	"strconv"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

const (
	MAX_IDEAL_TIME_PERIOD = 30 * time.Second
)

type PostgreDb struct {
	db *sql.DB
}

var (
	once     sync.Once
	instance *PostgreDb
)

func GetPostgreDbInstance() *PostgreDb {
	once.Do(func() {
		instance = &PostgreDb{}
		err := instance.Connect()
		if err != nil {
			panic(fmt.Sprintf("failed to connect to PostgreSQL: %v", err))
		}
	})
	return instance
}

func (p *PostgreDb) Connect() error {
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "password")
	dbName := getEnv("DB_NAME", "mydb")
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")

	connStr := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		dbUser,
		dbPassword,
		dbName,
		dbHost,
		dbPort,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("error connecting to PostgreSQL: %v", err)
	}
	maxOpenConns := getEnv("DB_MAX_OPEN_CONNS", "100")
	maxOpenConnsInt, _ := strconv.Atoi(maxOpenConns)
	db.SetMaxOpenConns(maxOpenConnsInt)
	db.SetConnMaxLifetime(MAX_IDEAL_TIME_PERIOD)

	p.db = db
	return nil
}

func (p *PostgreDb) Disconnect() error {
	if p.db != nil {
		err := p.db.Close()
		if err != nil {
			return fmt.Errorf("error disconnecting from PostgreSQL: %v", err)
		}
	}
	return nil
}

func (p *PostgreDb) Insert(query string, args ...interface{}) error {
	return p.Execute(query, args...)
}

func (p *PostgreDb) Update(query string, args ...interface{}) error {
	return p.Execute(query, args...)
}

func (p *PostgreDb) Delete(query string, args ...interface{}) error {
	return p.Execute(query, args...)
}

func (p *PostgreDb) Query(query string, args ...interface{}) (interface{}, error) {
	rows, err := p.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %v", err)
	}
	defer rows.Close()

	var result []map[string]interface{}
	columns, _ := rows.Columns()
	for rows.Next() {
		columnsData := make([]interface{}, len(columns))
		columnPtrs := make([]interface{}, len(columns))
		for i := range columnsData {
			columnPtrs[i] = &columnsData[i]
		}
		err = rows.Scan(columnPtrs...)
		if err != nil {
			return nil, err
		}

		rowData := make(map[string]interface{})
		for i, colName := range columns {
			rowData[colName] = columnsData[i]
		}
		result = append(result, rowData)
	}

	return result, nil
}

func (p *PostgreDb) Execute(query string, args ...interface{}) error {
	tx, err := p.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}

	_, err = tx.Exec(query, args...)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to execute query: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}
