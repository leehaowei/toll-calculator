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

type HTTPMerticHandler struct {
	reqCounter prometheus.Counter
	reqLatency prometheus.Histogram
}

func newHttpMetricsHandler(reqName string) *HTTPMerticHandler {
	reqCounter := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: fmt.Sprintf("http_%s_%s", reqName, "request_counter"),
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
	}
}

func (h *HTTPMerticHandler) instrument(next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		defer func(start time.Time) {
			latency := time.Since(start).Seconds()
			logrus.WithFields(logrus.Fields{
				"latency": latency,
				"request": r.RequestURI,
			}).Info()
			h.reqLatency.Observe(latency)
		}(time.Now())
		h.reqCounter.Inc()
		next(w, r)
	}
}

func handleGetInovice(svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "method not supported"})
			return
		}
		obuID, err := strconv.Atoi(r.URL.Query().Get("obu"))
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid OBU ID"})
			return
		}
		invoice, err := svc.CalculateInvoice(obuID)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, invoice)
	}
}

func handleAggregate(svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "method not supported"})
			return
		}
		var distance types.Distance
		if err := json.NewDecoder(r.Body).Decode(&distance); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		if err := svc.AggregateDistance(distance); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
	}
}
