package controllers

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/romana/rlog"

	"gitlab.com/voip-services/go-kamailio-api/internal/respond"
)

// App is application structure
type App struct {
	Router          *mux.Router
	DB              *pgxpool.Pool
	Ctx             context.Context
	jsonrpcHTTPAddr string
	httpClient      *http.Client
}

const (
	apiVersion = "v1"
)

// Initialize connect to the database and wire up routes
func (a *App) Initialize(DbURL string) {
	var err error

	a.Ctx = context.Background()

	a.DB, err = pgxpool.Connect(a.Ctx, DbURL)
	if err != nil {
		log.Debugf("Cannot connect to the pool, error: %s", err)
		os.Exit(1)
	}

	a.Router = mux.NewRouter().StrictSlash(true)
	a.initializeRoutes()
}

// NewClient returns exported ProviderRoutes
func (a *App) NewClient(httpAddr string) {

	a.jsonrpcHTTPAddr = fmt.Sprintf("http://%s/RPC/", httpAddr)
	a.httpClient = &http.Client{}

}

func (a *App) initializeRoutes() {

	a.Router.HandleFunc("/", home).Methods("GET")

	s := a.Router.PathPrefix("/api/" + apiVersion).Subrouter()

	s.HandleFunc("/subscribers", a.createSubscriber).Methods("POST")
	s.HandleFunc("/subscribers", a.getSubscribers).Methods("GET")
	s.HandleFunc("/subscribers/{id:[0-9]+}", a.getSubscriber).Methods("GET")
	s.HandleFunc("/subscribers/{id:[0-9]+}", a.updateSubscriber).Methods("PUT")
	s.HandleFunc("/subscribers/{id:[0-9]+}", a.deleteSubscriber).Methods("DELETE")
	s.HandleFunc("/subscribers/online", a.getSubscribersOnline).Methods("GET") // Get online subs. JSONRPC

}

// RunServer start our app
func (a *App) RunServer(listenAddr string) {
	log.Infof("Server starting on [%s]", listenAddr)
	err := http.ListenAndServe(listenAddr, a.Router)
	if err != nil {
		log.Error("Failed to start web service.")
		return
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	respond.JSON(w, http.StatusOK, "go-kamailio-api "+apiVersion)
}
