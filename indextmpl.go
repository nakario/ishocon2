package main
import (
	"net/http"
	"github.com/gin-gonic/gin"
	"strconv"
)
var code1 = []byte(`<!DOCTYPE html><html><head><meta http-equiv="Content-Type" content="text/html" charset="utf-8"><link rel="stylesheet" href="/css/bootstrap.min.css"><title>ISUCON選挙結果</title></head><body><nav class="navbar navbar-inverse navbar-fixed-top"><div class="container"><div class="navbar-header"><a class="navbar-brand" href="/">ISUCON選挙結果</a></div><div class="header clearfix"><nav><ul class="nav nav-pills pull-right"><li role="presentation"><a href="/vote">投票する</a></li></ul></nav></div></div></nav><div class="jumbotron"><div class="container"><h1>選挙の結果を大発表！</h1></div></div><div class="container"><h2>個人の部</h2><div id="people" class="row">`)
var code2 = []byte(`</div><h2>政党の部</h2><div id="parties" class="row">`)
var code3 = []byte(`</div><h2>男女比率</h2><div id="sex_ratio" class="row"><div class="col-md-6"><div class="panel panel-default"><div class="panel-heading"><p>男性</p></div><div class="panel-body"><p>得票数: `)
var code4 = []byte(`</p></div></div></div><div class="col-md-6"><div class="panel panel-default"><div class="panel-heading"><p>女性</p></div><div class="panel-body"><p>得票数: `)
var code5 = []byte(`</p></div></div></div></div></div></body></html>`)

func WriteIndexHTML(c *gin.Context, candidates []CandidateElectionResult, ratioMen int, ratioWomen int, partyResults []PartyElectionResult) {

	c.Status(http.StatusOK)
	c.Writer.Header()["Content-Type"] = []string{"text/html; charset=utf-8"}
	c.Writer.Write(code1)
	dom := ""
	for i,candidate := range candidates {
		dom += `<div class="col-md-3"><div class="panel panel-default"><div class="panel-heading">`
		if i < 10 {
			dom += `<p>`
			dom += strconv.Itoa(i + 1)
			dom += `. <a href="/candidates/`
			dom += strconv.Itoa(candidate.ID)
			dom += `">`
			dom += candidate.Name
			dom += `</a></p>`
		} else {
			dom += `<p>最下位. <a href="/candidates/`
			dom += strconv.Itoa(candidate.ID)
			dom += `">`
			dom += candidate.Name
			dom += `</a></p>`
		}
		dom += `</div><div class="panel-body"><p>得票数: `
		dom += strconv.Itoa(candidate.VoteCount)
		dom += `</p><p>政党: `
		dom += candidate.PoliticalParty
		dom += `</p></div></div></div>`
	}
	c.Writer.Write([]byte(dom))
	c.Writer.Write(code2)
  dom = ""
	for i,party := range partyResults {
		dom += `<div class="col-md-3"><div class="panel panel-default"><div class="panel-heading"><p>`
		dom += strconv.Itoa(i + 1)
		dom += `. <a href="/political_parties/`
		dom += party.PoliticalParty
		dom += `">`
		dom += party.PoliticalParty
		dom += `</a></p></div><div class="panel-body"><p>得票数: `
		dom += strconv.Itoa(party.VoteCount)
		dom += `</p></div></div></div>`
	}
	c.Writer.Write([]byte(dom))
	c.Writer.Write(code3)
	c.Writer.Write([]byte(strconv.Itoa(ratioMen)))
	c.Writer.Write(code4)
	c.Writer.Write([]byte(strconv.Itoa(ratioWomen)))
	c.Writer.Write(code5)
	c.Writer.WriteHeaderNow()
}


