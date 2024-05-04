package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/leehaowei/tolling-micro-service/types"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sirupsen/logrus"
)

type HTTPFunc func(http.ResponseWriter, *http.Request) error

type APIError struct {
	Code int
	Err  error
}

// Error implements the Error interface
func (e APIError) Error() string {
	return e.Err.Error()
}

type HTTPMerticHandler struct {
	reqCounter prometheus.Counter
	errCounter prometheus.Counter
	reqLatency prometheus.Histogram
}

func makeHTTPHandlerFunc(fn HTTPFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			if apiErr, ok := err.(APIError); ok {
				writeJSON(w, apiErr.Code, map[string]string{"error": apiErr.Error()})
			}
		}
	}
}

func newHttpMetricsHandler(reqName string) *HTTPMerticHandler {
	reqCounter := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: fmt.Sprintf("http_%s_%s", reqName, "request_counter"),
		Name:      "aggregator",
	})
	errCounter := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: fmt.Sprintf("http_%s_%s", reqName, "err_counter"),
		Name:      "aggregator",
	})
	reqLatency := promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace: fmt.Sprintf("http_%s_%s", reqName, "request_latency"),
		Name:      "aggregator",
		Buckets:   []float64{0.1, 0.5, 1},
	})
	return &HTTPMerticHandler{
		reqCounter: reqCounter,
		reqLatency: reqLatency,
		errCounter: errCounter,
	}
}

func (h *HTTPMerticHandler) instrument(next HTTPFunc) HTTPFunc {

	return func(w http.ResponseWriter, r *http.Request) error {
		var err error
		defer func(start time.Time) {
			latency := time.Since(start).Seconds()
			logrus.WithFields(logrus.Fields{
				"latency": latency,
				"request": r.RequestURI,
				"err":     err,
			}).Info()
			h.reqLatency.Observe(latency)
			h.reqCounter.Inc()
			if err != nil {
				h.errCounter.Inc()
			}
		}(time.Now())
		err = next(w, r)
		return err
	}
}

func handleGetInovice(svc Aggregator) HTTPFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		if r.Method != "GET" {
			return APIError{
				Code: http.StatusBadRequest,
				Err:  fmt.Errorf("invalid HTTP methid %s", r.Method),
			}
		}
		values := r.URL.Query().Get("obu")
		if len(values) == 0 {
			return APIError{
				Code: http.StatusBadRequest,
				Err:  fmt.Errorf("missing OBU id"),
			}
		}
		obuID, err := strconv.Atoi(values)
		if err != nil {
			return APIError{
				Code: http.StatusBadRequest,
				Err:  fmt.Errorf("invalid OBU id %d", obuID),
			}
		}
		invoice, err := svc.CalculateInvoice(obuID)
		if err != nil {
			return APIError{
				Code: http.StatusInternalServerError,
				Err:  err,
			}
		}
		return writeJSON(w, http.StatusOK, invoice)
	}
}

func handleAggregate(svc Aggregator) HTTPFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		if r.Method != "POST" {
			return APIError{
				Code: http.StatusBadRequest,
				Err:  fmt.Errorf("invalid HTTP methid %s", r.Method),
			}
		}
		var distance types.Distance
		if err := json.NewDecoder(r.Body).Decode(&distance); err != nil {
			return APIError{
				Code: http.StatusBadRequest,
				Err:  fmt.Errorf("failed to decode the response body: %s", err),
			}
		}
		if err := svc.AggregateDistance(distance); err != nil {
			return APIError{
				Code: http.StatusInternalServerError,
				Err:  err,
			}
		}
		return writeJSON(w, http.StatusOK, map[string]string{"msg": "ok"})
	}
}
