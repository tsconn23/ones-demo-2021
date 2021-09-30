package transitor

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/project-alvarium/provider-logging/pkg/interfaces"
	"github.com/project-alvarium/provider-logging/pkg/logging"
	"net/http"
	"time"
)

func LoadRestRoutes(r *mux.Router, logger interfaces.Logger) {
	r.HandleFunc("/",
		func(w http.ResponseWriter, r *http.Request) {
			getIndexHandler(w, r, logger)
		}).Methods(http.MethodGet)
}

func getIndexHandler(w http.ResponseWriter, r *http.Request, logger interfaces.Logger) {
	defer r.Body.Close()
	start := time.Now()
	w.Header().Add("Content-Type", "text/html")
	w.Write([]byte("<html><head><title>Transitor API</title></head><body><h2>Transitor API</h2></body></html>"))

	elapsed := time.Now().Sub(start)
	logger.Write(logging.TraceLevel, fmt.Sprintf("getIndexHandler OK %v", elapsed))
}
