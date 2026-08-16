package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"go-leaderboard/internal/leaderboard"
)

// Handler exposes the leaderboard over HTTP.
type Handler struct {
	lb *leaderboard.Leaderboard
}

// NewMux builds the HTTP mux for the leaderboard service.
func NewMux(lb *leaderboard.Leaderboard) http.Handler {
	h := &Handler{lb: lb}
	mux := http.NewServeMux()
	mux.HandleFunc("/score", h.score)
	mux.HandleFunc("/rank/", h.rank)
	mux.HandleFunc("/topn", h.topn)
	mux.HandleFunc("/leaderboard", h.leaderboard)
	return mux
}

func (h *Handler) score(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		UserID string `json:"user_id"`
		Score  int64  `json:"score"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	if err := h.lb.UpdateScore(req.UserID, req.Score); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) rank(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Path[len("/rank/"):]
	if userID == "" {
		http.Error(w, "missing user", http.StatusBadRequest)
		return
	}
	entry := h.lb.GetUserEntry(userID)
	if entry == nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}
	rank := h.lb.GetRank(userID)
	resp := map[string]any{
		"user_id": entry.UserID,
		"score":   entry.Score,
		"rank":    rank,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) topn(w http.ResponseWriter, r *http.Request) {
	n, err := strconv.Atoi(r.URL.Query().Get("n"))
	if err != nil || n <= 0 {
		http.Error(w, "bad n", http.StatusBadRequest)
		return
	}
	entries := h.lb.TopN(n)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entries)
}

func (h *Handler) leaderboard(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	size, _ := strconv.Atoi(r.URL.Query().Get("size"))
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}
	all := h.lb.All()
	offset := page * size
	if offset >= len(all) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]leaderboard.Entry{})
		return
	}
	end := offset + size
	if end > len(all) {
		end = len(all)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(all[offset:end])
}
