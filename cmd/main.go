package main

import (
	"b1a/server/handler"
	"fmt"
	"net/http"
	"os"

	"github.com/MadAppGang/httplog"
)

var (
	loginHandler    http.Handler = http.HandlerFunc(handler.Login)
	registerHandler http.Handler = http.HandlerFunc(handler.Register)
	linksHandler    http.Handler = http.HandlerFunc(handler.Links)
	shortenHandler  http.Handler = http.HandlerFunc(handler.Shorten)
	getUrlHandler   http.Handler = http.HandlerFunc(handler.GetUrl)
	getInfo         http.Handler = http.HandlerFunc(handler.GetInfo)
)

func main() {
	mux := http.NewServeMux()

	loggerWithFormatter := httplog.LoggerWithFormatter(httplog.DefaultLogFormatterWithRequestHeader)
	mux.Handle("/login", loggerWithFormatter(loginHandler))
	mux.Handle("/register", loggerWithFormatter(registerHandler))
	mux.Handle("/user", loggerWithFormatter(linksHandler))
	mux.Handle("/shorten", loggerWithFormatter(shortenHandler))
	mux.Handle("/", loggerWithFormatter(getUrlHandler))
	mux.Handle("/url", loggerWithFormatter(getInfo))

	corsHandler := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, PATCH, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, r)
		})
	}

	err := http.ListenAndServe(":8080", corsHandler(mux))
	if err != nil {
		if err == http.ErrServerClosed {
			fmt.Println("server closed")
		} else {
			fmt.Printf("error starting server: %s\n", err)
			os.Exit(1)
		}
	}
}
