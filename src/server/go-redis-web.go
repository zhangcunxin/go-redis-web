package main

import (
	"net/http"
	"log"
	"html/template"
)

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


func main() {
	http.HandleFunc("/", home)
	// 静态资源处理
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("../template"))))
	err := http.ListenAndServe(":9001", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
