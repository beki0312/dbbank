package handlers

import (
	"encoding/json"
	"github.com/natefinch/lumberjack"
	"log"
	"net/http"
)

func RespondBadRequest(w http.ResponseWriter, message interface{}) {
	data, err := json.Marshal(message)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	_, err = w.Write(data)
	if err != nil {
		log.Println(err)
	}
}

func RespondServerError(w http.ResponseWriter, message interface{}) {
	data, err := json.Marshal(message)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	_, err = w.Write(data)
	if err != nil {
		log.Println(err)
	}
}

func RespondUnauthorized(w http.ResponseWriter, message interface{}) {
	data, err := json.Marshal(message)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	_, err = w.Write(data)
	if err != nil {
		log.Println(err)
	}
}
func RespondNotImplemented(w http.ResponseWriter, message interface{}) {
	data, err := json.Marshal(message)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotImplemented)
	_, err = w.Write(data)
	if err != nil {
		log.Println(err)
	}
}
func RespondNotFound(w http.ResponseWriter, message interface{}) {
	data, err := json.Marshal(message)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	_, err = w.Write(data)
	if err != nil {
		log.Println(err)
	}
}

//respondJSON - ответ от JSON.
func RespondJSON(w http.ResponseWriter, item interface{}) {
	data, err := json.Marshal(item)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		log.Println(err)
	}
}

func LogInit() {
	log.SetOutput(&lumberjack.Logger{
		Filename:   "Logger/logs.log",
		MaxSize:    25,
		MaxBackups: 5,
		MaxAge:     60,
		Compress:   true,
	})

}
