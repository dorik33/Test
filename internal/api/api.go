package api

import (
	"net/http"

	_ "github.com/dorik33/Test/docs"
	"github.com/dorik33/Test/internal/config"
	"github.com/dorik33/Test/internal/middleware"
	"github.com/dorik33/Test/internal/store"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
)

type API struct {
	config *config.Config
	router *mux.Router
	logger *logrus.Logger
	store  *store.Store
}

func New(cfg *config.Config) *API {
	api := &API{
		config: cfg,
		router: mux.NewRouter(),
		logger: logrus.New(),
	}
	return api
}

func (api *API) Start() error {
	if err := api.configureLogger(); err != nil {
		return err
	}

	dbStore, err := store.NewConnection(api.config, api.logger)
	if err != nil {
		return err
	}
	api.store = dbStore

	api.logger.Debug("Successful connection to database")

	defer api.store.Close()

	api.configureRouter()
	server := &http.Server{
		Handler:      api.router,
		Addr:         api.config.Addr,
		WriteTimeout: api.config.WriteTimeout,
	}

	api.logger.Debug("Server is running with addr: ", api.config.Addr)
	return server.ListenAndServe()
}

func (api *API) configureLogger() error {
	api.logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	level, err := logrus.ParseLevel(api.config.LogLevel)
	if err != nil {
		return err
	}
	api.logger.SetLevel(level)

	return nil
}

func (api *API) configureRouter() {
	api.router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	api.router.Use(middleware.JSONContentTypeMiddleware)
	api.router.Use(middleware.LoggingMiddleware(api.logger))

	api.router.HandleFunc("/songs", api.GetSongsHandler).Methods("GET")
	api.router.HandleFunc("/songText/{id}", api.GetTextSongByIDHandler).Methods("GET")
	api.router.HandleFunc("/song/{id}", api.DeleteSonghandler).Methods("DELETE")
	api.router.HandleFunc("/song/{id}", api.UpdateSongHandler).Methods("PUT")
	api.router.HandleFunc("/song", api.AddSongHandler).Methods("POST")
}
