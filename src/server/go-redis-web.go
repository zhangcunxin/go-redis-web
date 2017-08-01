package main

import (
	"net/http"
	"log"
	"html/template"
	"strings"
	"encoding/json"
	"fmt"
)

type RedisObj struct {
	Key   string
	Value string
	Msg   string
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
	if request.Method != "GET" {
		http.Error(writer, "The method not allowed!", 405)
	}
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	key := strings.TrimSpace(request.FormValue("key"))
	value, _ := Get(key)
	if value == nil {
		value = []byte("nil")
	}
	strObj := RedisObj{Key: key, Value: string(value)}
	json.NewEncoder(writer).Encode(strObj)
}

func saveValue(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		http.Error(writer, "The method not allowed!", 405)
	}

	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	var postParams RedisObj
	decoder := json.NewDecoder(request.Body)
	decoder.Decode(&postParams)

	key := strings.TrimSpace(postParams.Key)
	value := strings.TrimSpace(postParams.Value)
	fmt.Println(key, ":", value)
	err := Set(key, value)

	if err != nil {
		http.Error(writer, "Set value to redis failed!", 500)
		return
	}
	postParams.Msg = "ok"
	json.NewEncoder(writer).Encode(postParams)
}

func deleteValue(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		http.Error(writer, "The method not allowed!", 405)
	}

	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	var postParams RedisObj
	decoder := json.NewDecoder(request.Body)
	decoder.Decode(&postParams)
	key := strings.TrimSpace(postParams.Key)
	originalValue, _ := Get(key)
	if originalValue == nil {
		json.NewEncoder(writer).Encode(RedisObj{Key: key, Value: "this key not existed from redis!", Msg: "ok"})
		return
	}
	v, _ := Del(key)
	if v == 1 {
		json.NewEncoder(writer).Encode(RedisObj{Key: key, Value: "delete success!", Msg: "ok"})
		return
	}
	json.NewEncoder(writer).Encode(RedisObj{Key: key, Value: "未知异常！"})
}

func main() {
	http.HandleFunc("/", home)
	http.HandleFunc("/queryValue", queryValue)
	http.HandleFunc("/saveValue", saveValue)
	http.HandleFunc("/deleteValue", deleteValue)
	// 静态资源处理
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("../template"))))
	err := http.ListenAndServe(":9001", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
