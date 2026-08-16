package leaderboard

import "sort"

// Entry represents a single user's leaderboard entry.
type Entry struct {
	UserID string
	Score  int64
}

// Leaderboard maintains an in-memory real-time leaderboard.
// Scores are kept in a map for O(1) lookup, and a sorted slice
// of user IDs (descending by score) provides ranked access.
type Leaderboard struct {
	scores map[string]int64
	order  []string // userIDs sorted descending by score
}

// New creates an empty Leaderboard.
func New() *Leaderboard {
	return &Leaderboard{scores: make(map[string]int64)}
}

// UpdateScore sets the user's score and re-sorts the leaderboard.
func (lb *Leaderboard) UpdateScore(userID string, score int64) error {
	lb.scores[userID] = score
	lb.resort()
	return nil
}

// GetRank returns the 1-based rank of the user, or -1 if not found.
func (lb *Leaderboard) GetRank(userID string) int {
	for i, uid := range lb.order {
		if uid == userID {
			return i
		}
	}
	return -1
}

// GetUserEntry returns the entry for a user, or nil if not found.
func (lb *Leaderboard) GetUserEntry(userID string) *Entry {
	score, ok := lb.scores[userID]
	if !ok {
		return nil
	}
	return &Entry{UserID: userID, Score: score}
}

// TopN returns the top n entries in descending score order.
func (lb *Leaderboard) TopN(n int) []Entry {
	if n <= 0 {
		return nil
	}
	if n > len(lb.order) {
		n = len(lb.order)
	}
	entries := make([]Entry, 0, n)
	for i := 0; i < n; i++ {
		uid := lb.order[i]
		entries = append(entries, Entry{UserID: uid, Score: lb.scores[uid]})
	}
	return entries
}

// All returns all entries in ranked order.
func (lb *Leaderboard) All() []Entry {
	entries := make([]Entry, 0, len(lb.order))
	for _, uid := range lb.order {
		entries = append(entries, Entry{UserID: uid, Score: lb.scores[uid]})
	}
	return entries
}

// RemoveUser removes a user from the leaderboard.
func (lb *Leaderboard) RemoveUser(userID string) {
	delete(lb.scores, userID)
}

func (lb *Leaderboard) resort() {
	lb.order = lb.order[:0]
	for uid := range lb.scores {
		lb.order = append(lb.order, uid)
	}
	sort.Slice(lb.order, func(i, j int) bool {
		return lb.scores[lb.order[i]] > lb.scores[lb.order[j]]
	})
}
