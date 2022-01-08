package main

import (
	"fmt"
)

var logger TransactionLogger

func initializeTransactionLog() error {
	var err error
	logger, err = NewPostgresTransactionLogger(PostgresDBParams{
		host:     "localhost",
		dbName:   "postgres",
		user:     "postgres",
		password: "pgpwd4habr",
	})

	//logger, err = NewFileTransactionLogger("transaction.log")
	if err != nil {
		return fmt.Errorf("failed to create event logger: %w", err)
	}
	events, errors := logger.ReadEvents()
	e, ok := Event{}, true
	for ok && err == nil {
		select {
		case err, ok = <-errors: // Получает ошибки
		case e, ok = <-events:
			switch e.EventType {
			case EventDelete: // Получено событие DELETE!
				err = Delete(e.Key)
			case EventPut: // Получено событие PUT!
				err = Put(e.Key, e.Value)
			}
		}
	}
	logger.Run()
	return err
}
