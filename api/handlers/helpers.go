package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/natefinch/lumberjack"
)

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
