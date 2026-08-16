package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
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

func TestLeaderboardPaginationUsesOneBasedPages(t *testing.T) {
	lb := leaderboard.New()
	for i := 1; i <= 25; i++ {
		userID := "user" + strconv.Itoa(i)
		score := int64(1000 - i)
		if err := lb.UpdateScore(userID, score); err != nil {
			t.Fatalf("UpdateScore(%s): %v", userID, err)
		}
	}

	mux := NewMux(lb)

	req := httptest.NewRequest(http.MethodGet, "/leaderboard?page=1&size=10", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var page1 []leaderboard.Entry
	if err := json.NewDecoder(w.Body).Decode(&page1); err != nil {
		t.Fatal(err)
	}
	if len(page1) != 10 {
		t.Fatalf("expected 10 entries, got %d", len(page1))
	}
	if page1[0].UserID != "user1" || page1[9].UserID != "user10" {
		t.Fatalf("unexpected page 1 entries: first=%s last=%s", page1[0].UserID, page1[9].UserID)
	}

	req = httptest.NewRequest(http.MethodGet, "/leaderboard?page=2&size=10", nil)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	var page2 []leaderboard.Entry
	if err := json.NewDecoder(w.Body).Decode(&page2); err != nil {
		t.Fatal(err)
	}
	if len(page2) != 10 {
		t.Fatalf("expected 10 entries, got %d", len(page2))
	}
	if page2[0].UserID != "user11" || page2[9].UserID != "user20" {
		t.Fatalf("unexpected page 2 entries: first=%s last=%s", page2[0].UserID, page2[9].UserID)
	}
}

func TestLeaderboardPaginationIsStableForTiedScores(t *testing.T) {
	lb := leaderboard.New()
	for i := 1; i <= 20; i++ {
		userID := "user" + strconv.Itoa(i)
		if err := lb.UpdateScore(userID, 100); err != nil {
			t.Fatalf("UpdateScore(%s): %v", userID, err)
		}
	}

	mux := NewMux(lb)

	req := httptest.NewRequest(http.MethodGet, "/leaderboard?page=1&size=10", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	var page1 []leaderboard.Entry
	if err := json.NewDecoder(w.Body).Decode(&page1); err != nil {
		t.Fatal(err)
	}

	req = httptest.NewRequest(http.MethodGet, "/leaderboard?page=2&size=10", nil)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	var page2 []leaderboard.Entry
	if err := json.NewDecoder(w.Body).Decode(&page2); err != nil {
		t.Fatal(err)
	}

	if len(page1) != 10 || len(page2) != 10 {
		t.Fatalf("expected 10 entries on each page, got %d and %d", len(page1), len(page2))
	}

	seen := make(map[string]bool, 20)
	for _, entry := range page1 {
		seen[entry.UserID] = true
	}
	for _, entry := range page2 {
		if seen[entry.UserID] {
			t.Fatalf("duplicate user across pages: %s", entry.UserID)
		}
		seen[entry.UserID] = true
	}
	if len(seen) != 20 {
		t.Fatalf("expected 20 unique users across two pages, got %d", len(seen))
	}
}
