package main

import (
	"database/sql"
	"encoding/csv"
	"github.com/gin-contrib/pprof"
	"html/template"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"sync"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var usersMap map[string]*User

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func loadUsers() {
	log.Println("Start loading users")
	usersMap = make(map[string]*User, 4000000)
	f, err := os.Open("users.csv")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	r := csv.NewReader(f)

	records, err := r.ReadAll()
	if err != nil {
		panic(err)
	}

	for _, record := range records {
		id, _ := strconv.Atoi(record[0])
		votes, _ := strconv.Atoi(record[4])
		usersMap[record[3]] = &User{id, record[1], record[2], record[3], votes, 0, sync.Mutex{}}
	}

	rows, err := db.Query("SELECT u.mynumber, v.cnt FROM votes AS v INNER JOIN users as u WHERE v.user_id = u.id")
	if err != nil && err != sql.ErrNoRows {
		panic(err)
	}

	for rows.Next() {
		var myNumber string
		var cnt int
		err := rows.Scan(&myNumber, &cnt)
		if err != nil {
			panic(err)
		}

		user, ok := usersMap[myNumber]
		if !ok {
			panic("hoge")
		}

		// user.Lock()
		user.Voted += cnt
		// user.Unlock()
	}
	log.Println("Finished loading users")
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	// database setting
	user := getEnv("ISHOCON2_DB_USER", "ishocon")
	pass := getEnv("ISHOCON2_DB_PASSWORD", "ishocon")
	dbname := getEnv("ISHOCON2_DB_NAME", "ishocon2")
	db, _ = sql.Open("mysql", user+":"+pass+"@/"+dbname)
	db.SetMaxIdleConns(5)

	loadUsers()

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(static.Serve("/css", static.LocalFile("public/css", true)))

	// session store
	store := sessions.NewCookieStore([]byte("mysession"))
	store.Options(sessions.Options{HttpOnly: true})
	r.Use(sessions.Sessions("showwin_happy", store))

	funcs := template.FuncMap{"indexPlus1": func(i int) int { return i + 1 }}
	r.SetFuncMap(funcs)
	r.LoadHTMLGlob("templates/*.tmpl")

	pprof.Register(r)

	// GET /
	r.GET("/", GetIndex)

	// GET /candidates/:candidateID(int)
	r.GET("/candidates/:candidateID", GetCandidateByID)

	// GET /political_parties/:name(string)
	r.GET("/political_parties/:name", GetPoliticalPartyByName)

	// GET /vote
	r.GET("/vote", GetVote)

	// POST /vote
	r.POST("/vote", PostVote)

	r.GET("/initialize", GetInitialize)

	r.RunTLS(":443", "/etc/nginx/ssl/server.crt", "/etc/nginx/ssl/server.key")
}

func GetIndex(c *gin.Context) {
	electionResults := getElectionResult()

	// 上位10人と最下位のみ表示
	tmp := make([]CandidateElectionResult, len(electionResults))
	copy(tmp, electionResults)
	candidates := tmp[:10]
	candidates = append(candidates, tmp[len(tmp)-1])

	partyNames := getAllPartyName()
	partyResultMap := map[string]int{}
	for _, name := range partyNames {
		partyResultMap[name] = 0
	}
	for _, r := range electionResults {
		partyResultMap[r.PoliticalParty] += r.VoteCount
	}
	partyResults := []PartyElectionResult{}
	for name, count := range partyResultMap {
		r := PartyElectionResult{}
		r.PoliticalParty = name
		r.VoteCount = count
		partyResults = append(partyResults, r)
	}
	// 投票数でソート
	sort.Slice(partyResults, func(i, j int) bool { return partyResults[i].VoteCount > partyResults[j].VoteCount })

	sexRatio := map[string]int{
		"men":   0,
		"women": 0,
	}
	for _, r := range electionResults {
		if r.Sex == "男" {
			sexRatio["men"] += r.VoteCount
		} else if r.Sex == "女" {
			sexRatio["women"] += r.VoteCount
		}
	}

	c.HTML(http.StatusOK, "templates/index.tmpl", gin.H{
		"candidates": candidates,
		"parties":    partyResults,
		"sexRatio":   sexRatio,
	})
}

func GetCandidateByID(c *gin.Context) {
	candidateID, _ := strconv.Atoi(c.Param("candidateID"))
	candidate, err := getCandidate(candidateID)
	if err != nil {
		c.Redirect(http.StatusFound, "/")
	}
	votes := getVoteCountByCandidateID(candidateID)
	candidateIDs := []int{candidateID}
	keywords := getVoiceOfSupporter(candidateIDs)

	c.HTML(http.StatusOK, "templates/candidate.tmpl", gin.H{
		"candidate": candidate,
		"votes":     votes,
		"keywords":  keywords,
	})
}

func GetPoliticalPartyByName(c *gin.Context) {
	partyName := c.Param("name")
	var votes int
	electionResults := getElectionResult()
	for _, r := range electionResults {
		if r.PoliticalParty == partyName {
			votes += r.VoteCount
		}
	}

	candidates := getCandidatesByPoliticalParty(partyName)
	candidateIDs := []int{}
	for _, c := range candidates {
		candidateIDs = append(candidateIDs, c.ID)
	}
	keywords := getVoiceOfSupporter(candidateIDs)

	c.HTML(http.StatusOK, "templates/political_party.tmpl", gin.H{
		"politicalParty": partyName,
		"votes":          votes,
		"candidates":     candidates,
		"keywords":       keywords,
	})
}

func GetVote(c *gin.Context) {
	WriteVoteHTML(c,"")
}


func postForm(c *gin.Context,key string) string {
	if values := c.Request.PostForm[key]; len(values) > 0 {
		return values[0]
	}
	return ""
}


func PostVote(c *gin.Context) {
	c.Request.ParseForm()
	user, userErr := getUser(postForm(c,"name"), postForm(c,"address"), postForm(c,"mynumber"))
	candidate, cndErr := getCandidateByName(postForm(c,"candidate"))
	votedCount := 0
	if user != nil {
		user.Lock()
		defer user.Unlock()
		votedCount = user.Voted
	}
	voteCount, _ := strconv.Atoi(postForm(c,"vote_count"))

	var message string
	if userErr != nil {
		message = "個人情報に誤りがあります"
	} else if user.Votes < voteCount+votedCount {
		message = "投票数が上限を超えています"
	} else if postForm(c,"candidate") == "" {
		message = "候補者を記入してください"
	} else if cndErr != nil {
		message = "候補者を正しく記入してください"
	} else if postForm(c,"keyword") == "" {
		message = "投票理由を記入してください"
	} else {
		createVote(user.ID, candidate.ID, postForm(c,"keyword"), voteCount)
		user.Voted += voteCount
		message = "投票に成功しました"
	}
	WriteVoteHTML(c,message)
}

func GetInitialize(c *gin.Context) {
	db.Exec("DELETE FROM votes")
	for _, u := range usersMap {
		// u.L.Lock()
		u.Voted = 0
		// u.L.Unlock()
	}

	c.String(http.StatusOK, "Finish")
}
