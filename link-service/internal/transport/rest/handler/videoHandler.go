package handler

import "net/http"

func StreamVideoHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "C:\\Users\\leono\\Desktop\\prog\\go\\shortener\\ad\\video.mp4")
}
