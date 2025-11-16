package app

import "net/http"

func RunServer() {
	http.ListenAndServe(":8080", nil)
}
