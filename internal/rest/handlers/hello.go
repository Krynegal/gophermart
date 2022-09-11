package handlers

import "net/http"

func (r *Router) Hello(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("HELLO"))
}
