package handlers

import "net/http"

func (h *Handler) Hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("HELLO"))
}
