package transitor

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/project-alvarium/alvarium-sdk-go/pkg/contracts"
	"github.com/project-alvarium/alvarium-sdk-go/pkg/interfaces"
	"github.com/project-alvarium/ones-demo-2021/internal/models"
	"io/ioutil"
	"log/slog"
	"net/http"
	"time"
)

func LoadRestRoutes(r *mux.Router, sdk interfaces.Sdk, logger interfaces.Logger) {
	r.HandleFunc("/",
		func(w http.ResponseWriter, r *http.Request) {
			getIndexHandler(w, r, logger)
		}).Methods(http.MethodGet)

	r.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		postReceiveDataHandler(w, r, sdk, logger)
	}).Methods(http.MethodPost)
}

func getIndexHandler(w http.ResponseWriter, r *http.Request, logger interfaces.Logger) {
	defer r.Body.Close()
	start := time.Now()
	w.Header().Add("Content-Type", "text/html")
	w.Write([]byte("<html><head><title>Transitor API</title></head><body><h2>Transitor API</h2></body></html>"))

	elapsed := time.Now().Sub(start)
	logger.Write(slog.LevelDebug, fmt.Sprintf("getIndexHandler OK %v", elapsed))
}

func postReceiveDataHandler(w http.ResponseWriter, r *http.Request, sdk interfaces.Sdk, logger interfaces.Logger) {
	if r.Body != nil {
		defer func() { _ = r.Body.Close() }()
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	//logger.Write(logging.DebugLevel, string(b))
	var sample models.SampleData
	err = json.Unmarshal(b, &sample)
	if err != nil {
		logger.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	//Having validated that the received payload can be unmarshaled to the correct type, we simply use the
	//[]byte from the request body here.
	ctx := context.WithValue(r.Context(), contracts.AnnotationTLS, r.TLS)
	sdk.Transit(ctx, b)

	r.Close = true
	w.WriteHeader(http.StatusAccepted)
	return
}
