package server

import (
	"context"
	"time"

	"log"

	"net/http"

	"github.com/aredoff/reagate/internal/app"
	"github.com/aredoff/reagate/pkg/httptracer"
	"github.com/gorilla/mux"
)

func New() *Server {
	r := mux.NewRouter()

	srv := &http.Server{
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	return &Server{
		server: srv,
		router: r,
		tracer: httptracer.New(),
		sites:  app.New(),
	}
}

type Server struct {
	server *http.Server
	router *mux.Router
	tracer httptracer.HttpTracer
	sites  *app.URLMonitor
}

func (a *Server) Serve(addr string) {
	a.initializeRoutes()
	a.server.Addr = addr
	if err := a.server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func (a *Server) Shutdown(ctx context.Context) error {
	return a.server.Shutdown(ctx)
}

func (a *Server) initializeRoutes() {
	a.router.HandleFunc("/", a.homeHandler).Methods("GET")
	a.router.HandleFunc("/trace", a.traceHandler).Methods("GET")
}

// func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
// 	response, _ := json.Marshal(payload)

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(code)
// 	w.Write(response)
// }
