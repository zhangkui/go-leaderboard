package leaderboard

import "testing"

func TestUpdateAndTopN(t *testing.T) {
	lb := New()
	lb.UpdateScore("alice", 100)
	lb.UpdateScore("bob", 200)
	lb.UpdateScore("carol", 150)

	top := lb.TopN(2)
	if len(top) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(top))
	}
	if top[0].UserID != "bob" || top[0].Score != 200 {
		t.Errorf("expected bob/200 first, got %s/%d", top[0].UserID, top[0].Score)
	}
	if top[1].UserID != "carol" || top[1].Score != 150 {
		t.Errorf("expected carol/150 second, got %s/%d", top[1].UserID, top[1].Score)
	}
}

func TestGetRank(t *testing.T) {
	lb := New()
	lb.UpdateScore("alice", 100)
	r := lb.GetRank("alice")
	if r < 0 {
		t.Errorf("expected non-negative rank, got %d", r)
	}
}

func TestGetRankMissing(t *testing.T) {
	lb := New()
	if r := lb.GetRank("nobody"); r != -1 {
		t.Errorf("expected -1 for missing user, got %d", r)
	}
}

func TestTopNLargeN(t *testing.T) {
	lb := New()
	lb.UpdateScore("alice", 100)
	top := lb.TopN(10)
	if len(top) != 1 {
		t.Errorf("expected 1 entry, got %d", len(top))
	}
}

func TestUpdateScoreChanges(t *testing.T) {
	lb := New()
	lb.UpdateScore("alice", 100)
	lb.UpdateScore("alice", 50)
	entry := lb.GetUserEntry("alice")
	if entry.Score != 50 {
		t.Errorf("expected score 50 after update, got %d", entry.Score)
	}
}

func TestUpdateScoreZeroAllowed(t *testing.T) {
	lb := New()
	if err := lb.UpdateScore("alice", 0); err != nil {
		t.Fatalf("zero score should be allowed, got %v", err)
	}
	entry := lb.GetUserEntry("alice")
	if entry.Score != 0 {
		t.Errorf("expected score 0, got %d", entry.Score)
	}
}

func TestUpdateScoreRejectsNegative(t *testing.T) {
	lb := New()
	if err := lb.UpdateScore("alice", -50); err == nil {
		t.Fatal("expected error for negative score, got nil")
	}
	// A rejected update must not pollute the leaderboard.
	if entry := lb.GetUserEntry("alice"); entry != nil {
		t.Errorf("expected alice to be absent after rejected update, got %+v", entry)
	}
}
