package analytics_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/nbalakrcloud/visitor-analytics/internal/analytics"
)

// Benchmark the RecordVisitor method
func BenchmarkVisitorStore_RecordVisitor(b *testing.B) {
	store := analytics.NewVisitorStore()

	for i := 0; i < b.N; i++ {
		store.RecordVisitor("http://foo.com", "Alice")
	}
}

// Benchmark the GetUniqueVisitors method
func BenchmarkVisitorStore_GetUniqueVisitors(b *testing.B) {
	store := analytics.NewVisitorStore()
	store.RecordVisitor("http://foo.com", "Alice")
	store.RecordVisitor("http://foo.com", "Bob")
	store.RecordVisitor("http://bar.com", "Alice")

	for i := 0; i < b.N; i++ {
		store.GetUniqueVisitors()
	}
}

// Benchmark the getAnalytics method
func BenchmarkVisitorStore_getAnalytics(b *testing.B) {
	service := analytics.NewAnalyticsService()
	router := mux.NewRouter()
	service.WireUpAnalytics(router)

	service.VisitorStore.RecordVisitor("http://foo.com", "Alice")
	service.VisitorStore.RecordVisitor("http://foo.com", "Bob")
	service.VisitorStore.RecordVisitor("http://bar.com", "Alice")

	req, _ := http.NewRequest("GET", "/analytics", nil)

	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
	}
}
