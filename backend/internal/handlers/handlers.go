package handlers

import "net/http"

func (h *Handler) Main(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world"))
}
