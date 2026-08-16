package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go-leaderboard/internal/leaderboard"
)

func TestScoreAndTopN(t *testing.T) {
	lb := leaderboard.New()
	mux := NewMux(lb)

	req := httptest.NewRequest(http.MethodPost, "/score", strings.NewReader(`{"user_id":"alice","score":100}`))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", w.Code)
	}

	req2 := httptest.NewRequest(http.MethodGet, "/topn?n=1", nil)
	w2 := httptest.NewRecorder()
	mux.ServeHTTP(w2, req2)
	var entries []leaderboard.Entry
	if err := json.NewDecoder(w2.Body).Decode(&entries); err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 || entries[0].UserID != "alice" {
		t.Errorf("expected alice, got %v", entries)
	}
}

func TestRank(t *testing.T) {
	lb := leaderboard.New()
	lb.UpdateScore("alice", 100)
	mux := NewMux(lb)

	req := httptest.NewRequest(http.MethodGet, "/rank/alice", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestBadScore(t *testing.T) {
	lb := leaderboard.New()
	mux := NewMux(lb)
	req := httptest.NewRequest(http.MethodPost, "/score", strings.NewReader(`not json`))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}
