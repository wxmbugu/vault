package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	//"strconv"
	"time"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	//"github.com/lib/pq"
	//	"time"
)

type Server struct {
	Services Service
	Router   *mux.Router
}

func SetupDb(conn string) *sql.DB {
	db, err := sql.Open("postgres", conn)
	if err != nil {
		log.Fatal(err)
	}
	db.Ping()
	db.SetMaxOpenConns(35)
	db.SetMaxIdleConns(35)
	db.SetConnMaxLifetime(time.Hour)
	return db
}

func NewServer() *Server {
	var wait time.Duration
	mux := mux.NewRouter()

	//postgres://{username}:{password}@{hostname}:{port}/{database}?options
	conn := SetupDb("postgres://mpxaqnjcqewmve:849658d51852ea38573b12b5d2cb5973760507f4beb29638707a29071418771f@ec2-44-205-112-253.compute-1.amazonaws.com:5432/d58gq8lh7l5ru2")
	services := NewService(conn)
	server := Server{
		Router:   mux,
		Services: services,
	}

	server.Routes()
	srve := http.Server{
		Addr:         "localhost:9000",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	fmt.Println("serving at port :9000")
	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srve.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()
	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c
	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srve.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
	return &server
}

func (server *Server) Routes() {
	http.Handle("/", server.Router)
	server.Router.Use(corsmiddleware)
	server.Router.Use(jsonmiddleware)
	server.Router.HandleFunc("/v1/secret", server.createsecret).Methods("POST", "OPTIONS")
	server.Router.HandleFunc("/v1/secret/{id}", server.getsecret).Methods("GET", "OPTIONS")

}

//TODO: SWAGGER

type Vaultreq struct {
	Secret   string `json:"secret" validate:"required"`
	Duration string `json:"duration" validate:"required"`
}

func (server *Server) createsecret(w http.ResponseWriter, r *http.Request) {
	//	var c chan int
	var req Vaultreq
	jsondec := json.NewDecoder(r.Body)
	jsondec.DisallowUnknownFields()
	err := jsondec.Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	data := vault{
		secret:   req.Secret,
		duration: req.Duration,
		uuid:     String(100),
	}
	vault, err := server.Services.CreateSecret(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(vault.uuid)
	w.Write(reqBodyBytes.Bytes())
	log.Println(http.StatusOK, "Secret created with id: ", vault.id)

}

func (server *Server) getsecret(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	vault, err := server.Services.FindSecret(id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, errors.New("no such secret").Error(), http.StatusNotFound)
			log.Print("No such secret", r.URL.Path, http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Print(err.Error(), r.URL.Path, http.StatusBadRequest)
		return
	}

	//TODO: Make this into a function to avoid repition
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(vault.secret)
	w.Write(reqBodyBytes.Bytes())
	log.Println(http.StatusOK, "Secret  received with uuid: ", id)
}

func corsmiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS,PUT")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")
		w.Header().Set("Content-Type", "application/json")
		if r.Method == "OPTIONS" {
			return
		}

		next.ServeHTTP(w, r)
	})
}

// xml set header middleware
func xmlmiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml")
		next.ServeHTTP(w, r)
	})
}

// yaml set header middleware
func yamlmiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/x-yaml")
		next.ServeHTTP(w, r)
	})
}

// json set header middleware
func jsonmiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
