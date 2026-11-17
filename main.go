package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"yamaha-rx-v1500-control/yamaha"
)

func main() {
	serialPort := os.Getenv("SERIAL_PORT")
	if serialPort == "" {
		serialPort = "/dev/ttyUSB0"
	}

	transport := yamaha.NewSerialTransport(serialPort)
	ctrl := yamaha.NewController(transport)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/power", func(w http.ResponseWriter, r *http.Request) {
		on := r.URL.Query().Get("on") == "true"
		if err := ctrl.SetPower(on); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("POST /api/speaker/a", func(w http.ResponseWriter, r *http.Request) {
		on := r.URL.Query().Get("on") == "true"
		if err := ctrl.SetSpeakerA(on); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("POST /api/speaker/b", func(w http.ResponseWriter, r *http.Request) {
		on := r.URL.Query().Get("on") == "true"
		if err := ctrl.SetSpeakerB(on); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("POST /api/volume", func(w http.ResponseWriter, r *http.Request) {
		db, err := strconv.ParseFloat(r.URL.Query().Get("db"), 64)
		if err != nil {
			http.Error(w, "wrong db parameter", http.StatusBadRequest)
			return
		}
		if err := ctrl.SetVolume(db); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("POST /hook/spotify/play", func(w http.ResponseWriter, r *http.Request) {
		if err := ctrl.OnSpotifyPlay(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("POST /hook/spotify/pause", func(w http.ResponseWriter, r *http.Request) {
		ctrl.OnSpotifyPause(10 * time.Minute)
		w.WriteHeader(http.StatusOK)
	})

	log.Printf("Yamaha Controller is starting at :8080 (port: %s)", serialPort)
	log.Fatal(http.ListenAndServe(":8080", mux))
}
