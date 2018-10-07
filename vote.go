package main

import "sync/atomic"

// Vote Model
type Vote struct {
	ID          int
	UserID      int
	CandidateID int
	Keyword     string
}

var voteCounts [30]*int64

func getVoteCountByCandidateID(candidateID int) (count int) {
	return int(atomic.LoadInt64(voteCounts[candidateID-1]))
}

func createVote(userID int, candidateID int, keyword string, cnt int) {
	go func() {
		db.Exec("INSERT INTO votes (user_id, candidate_id, keyword, cnt) VALUES (?, ?, ?, ?)",
			userID, candidateID, keyword, cnt)
	}()

	voteCount := voteCounts[candidateID-1]
	atomic.AddInt64(voteCount, int64(cnt))

	car := candidateElectionResults[candidateID]
	car.Lock()
	car.VoteCount += cnt
	candidateElectionResults[candidateID] = car
	car.Unlock()
}

func getVoiceOfSupporter(candidateIDs []int) (voices []string) {
	return []string{}
}
