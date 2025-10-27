package main

import (
	"net/http"
)

func handleErr(w http.ResponseWriter, r *http.Request) {
	responseWithErr(w, 500, "Somthing is wrong")
}
