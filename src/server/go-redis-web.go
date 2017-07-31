package main

import (
	"net/http"
	"log"
	"html/template"
	"strings"
	"encoding/json"
)

type RedisStr struct {
	Key   string
	Value string
}

func home(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/" {
		http.Error(writer, "Not found!", 404)
		return
	}
	if request.Method != "GET" {
		http.Error(writer, "The method not allowed!", 405)
		return
	}
	template, err := template.ParseFiles("../template/home.html")
	if err != nil {
		log.Fatal("parse template err: ", err)
	}
	template.Execute(writer, nil)
}

func queryValue(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	key := strings.TrimSpace(request.FormValue("key"))
	value, err := Get(key)
	if err != nil {
		http.Error(writer, "Get value from redis failed", 400)
	}
	strObj := RedisStr{Key: key, Value: string(value)}
	json.NewEncoder(writer).Encode(strObj)

}

func main() {
	http.HandleFunc("/", home)
	http.HandleFunc("/queryValue", queryValue)
	// 静态资源处理
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("../template"))))
	err := http.ListenAndServe(":9001", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
