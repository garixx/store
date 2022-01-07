package main

import (
	"fmt"
	"github.com/garixx/store/loggers"
)

var logger loggers.TransactionLogger

func initializeTransactionLog() error {
	var err error
	logger, err = loggers.NewFileTransactionLogger("transaction.log")
	if err != nil {
		return fmt.Errorf("failed to create event logger: %w", err)
	}
	events, errors := logger.ReadEvents()
	e, ok := loggers.Event{}, true
	for ok && err == nil {
		select {
		case err, ok = <-errors: // Получает ошибки
		case e, ok = <-events:
			switch e.EventType {
			case loggers.EventDelete: // Получено событие DELETE!
				err = Delete(e.Key)
			case loggers.EventPut: // Получено событие PUT!
				err = Put(e.Key, e.Value)
			}
		}
	}
	logger.Run()
	return err
}
