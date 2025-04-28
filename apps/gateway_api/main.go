package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

func main() {
	// Создаем прокси для сервиса-монолита
	monolithUrl, _ := url.Parse(getEnv("MONOLITH_URL", "http://localhost:8000"))
	monolithProxy := httputil.NewSingleHostReverseProxy(monolithUrl)

	// Создаем прокси для сервиса temperature-api
	temperatureApiUrl, _ := url.Parse(getEnv("TEMPERATURE_API_URL", "http://localhost:8001"))
	temperatureApiProxy := httputil.NewSingleHostReverseProxy(temperatureApiUrl)

	// Настраиваем маршрутизацию
	http.HandleFunc("/service1/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Перенаправление запроса в Service1")
		r.URL.Path = r.URL.Path[len("/service1/"):] // Удаляем префикс
		monolithProxy.ServeHTTP(w, r)
	})

	http.HandleFunc("/service2/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Перенаправление запроса в Service2")
		r.URL.Path = r.URL.Path[len("/service2/"):] // Удаляем префикс
		temperatureApiProxy.ServeHTTP(w, r)
	})

	// Корневой маршрут
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Gateway is running"))
	})

	// Запускаем сервер
	log.Println("Gateway запущен на порту 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
