package main

import "net/http"

func run_handler(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	_, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		// mute : websocket: the client is not using the websocket protocol: 'upgrade' token not found in 'Connection' header
		return
	}

}
