package main

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"sync/atomic"
	"time"
)

const VotesFile = "/home/ishocon/votes.csv"

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

func createVote(myNumber string, candidateID int, cnt int) {
	insertVoteCh <- vote{myNumber, candidateID, cnt}

	voteCount := voteCounts[candidateID-1]
	atomic.AddInt64(voteCount, int64(cnt))

	car := candidateElectionResults[candidateID-1]
	car.Lock()
	car.VoteCount += cnt
	car.Unlock()
}

func getVoiceOfSupporter(candidateIDs []int) (voices []string) {
	return []string{}
}

type vote struct{
	myNumber string
	candidateID int
	cnt int
}
var insertVoteCh = make(chan vote, 1000000)

func voteManager() {
	votes := make([][]string, 0, 1000000)
	ticker := time.Tick(2 * time.Minute)
	for {
		select {
		case req := <- insertVoteCh:
			votes = append(votes, []string{req.myNumber, strconv.Itoa(req.candidateID), strconv.Itoa(req.cnt)})
		case <- ticker:
			f, err := os.OpenFile(VotesFile, os.O_APPEND, 0777)
			if err != nil {
				log.Println("Failed to open votes.csv:", err)
			}
			w := csv.NewWriter(f)
			err = w.WriteAll(votes)
			if err != nil {
				log.Println("Failed to write votes.csv:", err)
			}
			f.Close()
			votes = make([][]string, 0, 1000000)
		}
	}
}
