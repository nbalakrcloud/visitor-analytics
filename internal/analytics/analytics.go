package analytics

import (
	"bytes"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
)

type AnalyticsService struct {
	VisitorStore VisitorStore
}
type VisitorStore struct {
	Store map[string]map[string]bool
	mu    sync.Mutex
}

func NewVisitorStore() *VisitorStore {
	return &VisitorStore{
		Store: make(map[string]map[string]bool),
	}
}

func NewAnalyticsService() *AnalyticsService {
	return &AnalyticsService{
		VisitorStore: *NewVisitorStore(),
	}
}

func (a *AnalyticsService) WireUpAnalytics(router *mux.Router) {
	router.HandleFunc("/track", a.VisitorStore.trackVisitor).Methods("GET")
	router.HandleFunc("/analytics", a.VisitorStore.getAnalytics).Methods("GET")
}

func (v *VisitorStore) RecordVisitor(url, visitorID string) {
	v.mu.Lock()
	defer v.mu.Unlock()

	if _, ok := v.Store[url]; !ok {
		v.Store[url] = make(map[string]bool)
	}
	v.Store[url][visitorID] = true
}

func (v *VisitorStore) trackVisitor(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	visitorID := r.URL.Query().Get("visitorID")

	if url == "" || visitorID == "" {
		http.Error(w, "url and visitorID parameters are required", http.StatusBadRequest)
		return
	}
	v.RecordVisitor(url, visitorID)
	w.WriteHeader(http.StatusOK)
}

func (v *VisitorStore) getAnalytics(w http.ResponseWriter, r *http.Request) {

	visitorMap := v.GetUniqueVisitors()
	var buf bytes.Buffer

	for key, value := range visitorMap {
		buf.WriteString(key + ":" + strconv.Itoa(value) + "\n")
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	buf.WriteTo(w)
}

func (v *VisitorStore) GetUniqueVisitors() map[string]int {
	v.mu.Lock()
	defer v.mu.Unlock()
	returnMap := make(map[string]int)
	for k := range v.Store {
		returnMap[k] = len(v.Store[k])
	}

	return returnMap
}
