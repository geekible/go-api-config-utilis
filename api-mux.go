package goapiconfigutilis

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httplog"
)

type ApiMux struct {
	apiName string
	mux     *chi.Mux
}

func InitApiMux(apiName string) *ApiMux {
	chiMux := chi.NewMux()
	chiMux.Use(middleware.Compress(5, "application/json"))
	chiMux.Use(middleware.AllowContentType("application/json", "text/xml"))
	chiMux.Use(middleware.NoCache)
	chiMux.Use(middleware.StripSlashes)
	chiMux.Use(middleware.Logger)
	chiMux.Use(middleware.Recoverer)

	return &ApiMux{
		apiName: apiName,
		mux:     chiMux,
	}
}

func (m *ApiMux) AddHttpLogging() {
	requestLogger := httplog.NewLogger(m.apiName, httplog.Options{
		JSON:     true,
		Concise:  true,
		LogLevel: "debug",
	})

	m.mux.Use(httplog.RequestLogger(requestLogger))
}

func (m *ApiMux) AddCorsPolicy() {
	m.mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
}

func (m *ApiMux) RegisterRoute(pattern string, httpVerb HttpMethod, handlerFunc http.HandlerFunc) {
	switch httpVerb {
	case POST:
		m.mux.Post(pattern, handlerFunc)
	case PUT:
		m.mux.Put(pattern, handlerFunc)
	case DELETE:
		m.mux.Delete(pattern, handlerFunc)
	case GET:
		m.mux.Get(pattern, handlerFunc)
	default:
		log.Fatalf("%s is not a supported HttpMethod", httpVerb)
	}
}

func (m *ApiMux) ListenAndServe(apiPort int) {
	if err := http.ListenAndServe(fmt.Sprintf(":%d", apiPort), m.mux); err != nil {
		log.Fatalf("error starting http server: %v", err)
	}
}
