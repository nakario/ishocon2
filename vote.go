package main

// Vote Model
type Vote struct {
	ID          int
	UserID      int
	CandidateID int
	Keyword     string
}

func getVoteCountByCandidateID(candidateID int) (count int) {
	row := db.QueryRow("SELECT SUM(cnt) AS count FROM votes WHERE candidate_id = ?", candidateID)
	row.Scan(&count)
	return
}

func createVote(userID int, candidateID int, keyword string, cnt int) {
	db.Exec("INSERT INTO votes (user_id, candidate_id, keyword, cnt) VALUES (?, ?, ?, ?)",
		userID, candidateID, keyword, cnt)
	voteCount, _ := VoteCountByCandidateIDMap.Load(candidateID)
	VoteCountByCandidateIDMap.Store(candidateID, voteCount.(int)+cnt)
}

func getVoiceOfSupporter(candidateIDs []int) (voices []string) {
	rows, err := db.Query(`
    SELECT keyword
    FROM votes
    WHERE candidate_id IN (?)
    GROUP BY keyword
    ORDER BY SUM(cnt) DESC
    LIMIT 10`)
	if err != nil {
		return nil
	}

	defer rows.Close()
	for rows.Next() {
		var keyword string
		err = rows.Scan(&keyword)
		if err != nil {
			panic(err.Error())
		}
		voices = append(voices, keyword)
	}
	return
}
