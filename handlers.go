package main

import (
	"errors"
	"github.com/gorilla/mux"
	"io"
	"net/http"
)

// keyValuePutHandler ожидает получить PUT-запрос с
// ресурсом "/v1/key/{key}".
func keyValuePutHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // Получить ключ из запроса
	key := vars["key"]
	value, err := io.ReadAll(r.Body) // Тело запроса хранит значение
	defer r.Body.Close()
	if err != nil { // Если возникла ошибка, сообщить о ней
		http.Error(w,
			err.Error(),
			http.StatusInternalServerError)
		return
	}
	err = Put(key, string(value)) // Сохранить значение как строку
	if err != nil {               // Если возникла ошибка, сообщить о ней
		http.Error(w,
			err.Error(),
			http.StatusInternalServerError)
		return
	}
	logger.WritePut(key, string(value))
	w.WriteHeader(http.StatusCreated) // Все хорошо! Вернуть StatusCreated
}

func keyValueGetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // Извлечь ключ из запроса
	key := vars["key"]
	value, err := Get(key) // Получить значение для данного ключа
	if errors.Is(err, ErrorNoSuchKey) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(value)) // Записать значение в ответ
}

func keyValueDeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // Получить ключ из запроса
	key := vars["key"]

	err := Delete(key) // Сохранить значение как строку
	if err != nil {    // Если возникла ошибка, сообщить о ней
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	logger.WriteDelete(key)
	w.WriteHeader(http.StatusNoContent) // Все хорошо! Вернуть StatusCreated
}
