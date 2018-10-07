package main

import (
	"errors"
	"html/template"
	"sort"
)

// Candidate Model
type Candidate struct {
	ID             int
	Name           string
	PoliticalParty string
	Sex            string
}

// CandidateElectionResult type
type CandidateElectionResult struct {
	ID             int
	Name           string
	PoliticalParty string
	Sex            string
	VoteCount      int
}

// PartyElectionResult type
type PartyElectionResult struct {
	PoliticalParty string
	VoteCount      int
}

type CandidateWithVote struct {
	Candidate
	VoteCount int
}

var candidates = []Candidate{
	Candidate{1, "佐藤 一郎", "夢実現党", "男"},
	Candidate{2, "佐藤 次郎", "国民10人大活躍党", "女"},
	Candidate{3, "佐藤 三郎", "国民10人大活躍党", "女"},
	Candidate{4, "佐藤 四郎", "国民10人大活躍党", "男"},
	Candidate{5, "佐藤 五郎", "国民元気党", "女"},
	Candidate{6, "鈴木 一郎", "国民平和党", "男"},
	Candidate{7, "鈴木 次郎", "国民元気党", "女"},
	Candidate{8, "鈴木 三郎", "国民10人大活躍党", "女"},
	Candidate{9, "鈴木 四郎", "国民元気党", "女"},
	Candidate{10, "鈴木 五郎", "国民元気党", "女"},
	Candidate{11, "高橋 一郎", "国民平和党", "男"},
	Candidate{12, "高橋 次郎", "夢実現党", "男"},
	Candidate{13, "高橋 三郎", "夢実現党", "男"},
	Candidate{14, "高橋 四郎", "国民平和党", "女"},
	Candidate{15, "高橋 五郎", "国民10人大活躍党", "女"},
	Candidate{16, "田中 一郎", "夢実現党", "男"},
	Candidate{17, "田中 次郎", "国民平和党", "女"},
	Candidate{18, "田中 三郎", "夢実現党", "女"},
	Candidate{19, "田中 四郎", "国民元気党", "男"},
	Candidate{20, "田中 五郎", "夢実現党", "女"},
	Candidate{21, "渡辺 一郎", "夢実現党", "女"},
	Candidate{22, "渡辺 次郎", "国民平和党", "女"},
	Candidate{23, "渡辺 三郎", "夢実現党", "男"},
	Candidate{24, "渡辺 四郎", "国民平和党", "女"},
	Candidate{25, "渡辺 五郎", "国民10人大活躍党", "男"},
	Candidate{26, "伊藤 一郎", "夢実現党", "女"},
	Candidate{27, "伊藤 次郎", "国民10人大活躍党", "女"},
	Candidate{28, "伊藤 三郎", "国民平和党", "女"},
	Candidate{29, "伊藤 四郎", "国民10人大活躍党", "男"},
	Candidate{30, "伊藤 五郎", "国民元気党", "男"},
}

var VoteCountByCandidateIDMap = map[int]int{
	1:  0,
	2:  0,
	3:  0,
	4:  0,
	5:  0,
	6:  0,
	7:  0,
	8:  0,
	9:  0,
	10: 0,
	11: 0,
	12: 0,
	13: 0,
	14: 0,
	15: 0,
	16: 0,
	17: 0,
	18: 0,
	19: 0,
	20: 0,
	21: 0,
	22: 0,
	23: 0,
	24: 0,
	25: 0,
	26: 0,
	27: 0,
	28: 0,
	29: 0,
	30: 0,
}

func getInitAllCandidatesDOM() template.HTML {
	result := ""
	for _, candidate := range candidates {
		result += `<option value="` + candidate.Name + `">` + candidate.Name + `</option>`
	}
	return template.HTML(result)
}

var getAllCandidatesDOM = getInitAllCandidatesDOM()

func getCandidate(candidateID int) (c Candidate, err error) {
	if candidateID <= 0 || candidateID >= 30 {
		err = errors.New("yee")
	}
	c = candidates[candidateID-1]
	err = nil
	return
}

func initCandidateByNameMap() map[string]Candidate {
	result := make(map[string]Candidate)
	for _, candidate := range candidates {
		result[candidate.Name] = candidate
	}
	return result
}

var candidadeByNameMap = initCandidateByNameMap()

func getCandidateByName(name string) (c Candidate, err error) {
	c, ok := candidadeByNameMap[name]
	if !ok {
		err = errors.New("no candidate")
	} else {
		err = nil
	}
	return
}

func getAllPartyName() (partyNames []string) {
	return []string{
		"国民10人大活躍党",
		"国民元気党",
		"国民平和党",
		"夢実現党",
	}
}

func getCandidatesByPoliticalParty(party string) (candidates []Candidate) {
	rows, err := db.Query("SELECT * FROM candidates WHERE political_party = ?", party)
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		c := Candidate{}
		err = rows.Scan(&c.ID, &c.Name, &c.PoliticalParty, &c.Sex)
		if err != nil {
			panic(err.Error())
		}
		candidates = append(candidates, c)
	}
	return
}

func (c Results) Len() int {
	return len(c)
}

func (c Results) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c Results) Less(i, j int) bool {
	return c[i].VoteCount < c[j].VoteCount
}

type Results []CandidateElectionResult

func getElectionResult() (result []CandidateElectionResult) {
	for _, candidate := range candidates {
		r := CandidateElectionResult{}
		r.ID = candidate.ID
		r.Name = candidate.Name
		r.PoliticalParty = candidate.PoliticalParty
		r.Sex = candidate.Sex
		r.VoteCount = VoteCountByCandidateIDMap[candidate.ID]
		result = append(result, r)
	}
	result = Results(result)
	sort.Sort(sort.Reverse(result))
	result = []CandidateElectionResult(result)
	return
}
