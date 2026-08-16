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

func TestNegativeScore(t *testing.T) {
	lb := leaderboard.New()
	mux := NewMux(lb)

	// A negative score must be rejected.
	req := httptest.NewRequest(http.MethodPost, "/score", strings.NewReader(`{"user_id":"alice","score":-50}`))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for negative score, got %d", w.Code)
	}
	// And must not have been written.
	if entry := lb.GetUserEntry("alice"); entry != nil {
		t.Errorf("expected alice to be absent after rejected update, got %+v", entry)
	}

	// Positive scores are unaffected.
	req2 := httptest.NewRequest(http.MethodPost, "/score", strings.NewReader(`{"user_id":"bob","score":50}`))
	w2 := httptest.NewRecorder()
	mux.ServeHTTP(w2, req2)
	if w2.Code != http.StatusNoContent {
		t.Fatalf("expected 204 for positive score, got %d", w2.Code)
	}
	if r := lb.GetRank("bob"); r != 0 {
		t.Errorf("expected bob at rank 0, got %d", r)
	}
}
