package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq" // Анонимный импорт пакета драйвера
	"github.com/sirupsen/logrus"
)

type PostgresDBParams struct {
	dbName   string
	host     string
	user     string
	password string
}

type PostgresTransactionLogger struct {
	events chan<- Event // Канал только для записи; для передачи событий
	errors <-chan error // Канал только для чтения; для приема ошибок
	db     *sql.DB      // Интерфейс доступа к базе данных
}

func (l *PostgresTransactionLogger) WritePut(key, value string) {
	logrus.Info("exec put ...")
	l.events <- Event{EventType: EventPut, Key: key, Value: value}
}

func (l *PostgresTransactionLogger) WriteDelete(key string) {
	l.events <- Event{EventType: EventDelete, Key: key}
}

func (l *PostgresTransactionLogger) Err() <-chan error {
	return l.errors
}

func (l *PostgresTransactionLogger) Run() {
	events := make(chan Event, 16) // Создать канал событий
	l.events = events
	errors := make(chan error, 1) // Создать канал ошибок
	l.errors = errors
	go func() { // Запрос INSERT
		query := `INSERT INTO transactions (event_type, key, value) VALUES ($1, $2, $3)`
		for e := range events { // Извлечь следующее событие Event
			logrus.Infof("executing query: INSERT INTO transactions (event_type, key, value) VALUES (%d, %s, %s)", e.EventType, e.Key, e.Value)
			_, err := l.db.Exec( // Выполнить запрос INSERT
				query,
				e.EventType, e.Key, e.Value)
			logrus.Info("query error:", err)
			if err != nil {
				errors <- err
			}
		}
	}()
}

func (l *PostgresTransactionLogger) ReadEvents() (<-chan Event, <-chan error) {
	outEvent := make(chan Event)    // Небуферизованный канал событий
	outError := make(chan error, 1) // Буферизованный канал ошибок
	go func() {
		defer close(outEvent) // Закрыть каналы
		defer close(outError) // по завершении сопрограммы
		query := `SELECT sequence, event_type, key, value FROM transactions ORDER BY sequence`
		rows, err := l.db.Query(query) // Выполнить запрос; получить набор результатов
		if err != nil {
			outError <- fmt.Errorf("sql query error: %w", err)
			return
		}
		defer rows.Close() // Это важно!
		e := Event{}       // Создать пустой экземпляр Event
		for rows.Next() {  // Цикл по записям
			err = rows.Scan(&e.Sequence, &e.EventType, &e.Key, &e.Value)
			if err != nil {
				outError <- fmt.Errorf("error reading row: %w", err)
				return
			}
			outEvent <- e // Послать e в канал
		}
		err = rows.Err()
		if err != nil {
			outError <- fmt.Errorf("transaction log read failure: %w", err)
		}
	}()
	return outEvent, outError
}

func (l *PostgresTransactionLogger) createTableIfNotExists() error {
	query := `CREATE TABLE IF NOT EXISTS transactions (
	    sequence SERIAL PRIMARY KEY,
	    event_type integer NOT NULL,
		key varchar(45) NOT NULL,
		value varchar(450) NOT NULL
	)`
	logrus.Info("create db if not exists ...")
	_, err := l.db.Exec(query) // Выполнить запрос; получить набор результатов
	if err != nil {
		return err
	}
	return nil
}

func NewPostgresTransactionLogger(config PostgresDBParams) (TransactionLogger, error) {
	connStr := fmt.Sprintf("host=%s dbname=%s user=%s password=%s sslmode=disable", config.host, config.dbName, config.user, config.password)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}

	err = db.Ping() // Проверка соединения с базой данных
	if err != nil {
		return nil, fmt.Errorf("failed to open db connection: %w", err)
	}

	logger := &PostgresTransactionLogger{db: db}

	if err = logger.createTableIfNotExists(); err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}
	return logger, nil
}
