package main
import (
	"net/http"
	"github.com/gin-gonic/gin"
)
var codeA = []byte(
`<!DOCTYPE html>
<html>
	<head>
		<meta http-equiv="Content-Type" content="text/html" charset="utf-8">
		<link rel="stylesheet" href="/css/bootstrap.min.css">
		<title>ISUCON選挙結果</title>
	</head>

	<body>
		<nav class="navbar navbar-inverse navbar-fixed-top">
			<div class="container">
				<div class="navbar-header">
					<a class="navbar-brand" href="/">ISUCON選挙結果</a>
				</div>
				<div class="header clearfix">
					<nav>
						<ul class="nav nav-pills pull-right">
							<li role="presentation"><a href="/vote">投票する</a></li>
						</ul>
					</nav>
				</div>
			</div>
		</nav>
<div class="jumbotron">
	<div class="container">
		<h1>選挙の結果を大発表！</h1>
	</div>
</div>
<div class="container">
	<h2>個人の部</h2>
	<div id="people" class="row">
		`)
var codeB = []byte(`</div>
<h2>政党の部</h2>
<div id="parties" class="row">`)
var codeC =[]byte(  `</div>
<h2>男女比率</h2>
<div id="sex_ratio" class="row">
	<div class="col-md-6">
		<div class="panel panel-default">
			<div class="panel-heading">
				<p>男性</p>
			</div>
			<div class="panel-body">
				<p>得票数: `)


func WriteIndexHTML(c *gin.Context,candidates Candidate[],ratioMen int,ratioWomen int,partyResults []PartyElectionResult{}) {
	c.Status(http.StatusOK)
	c.Writer.Header()["Content-Type"] = []string{"text/html; charset=utf-8"}
	c.Writer.Write(codeA)
	c.Writer.Write([]byte(getAllCandidatesDOM))
	c.Writer.Write(codeB)
	c.Writer.Write([]byte(message))
	c.Writer.Write(codeC)
	c.Writer.WriteHeaderNow()
}


/*
{{ range $index, $candidate := .candidates }}
      <div class="col-md-3">
        <div class="panel panel-default">
          <div class="panel-heading">
            {{ if lt $index 10 }}
              <p>{{ $index | indexPlus1 }}. <a href="/candidates/{{ $candidate.ID }}">{{ $candidate.Name }}</a></p>
            {{ else }}
              <p>最下位. <a href="/candidates/{{ $candidate.ID }}">{{ $candidate.Name }}</a></p>
            {{ end }}
          </div>
          <div class="panel-body">
            <p>得票数: {{ $candidate.VoteCount }}</p>
            <p>政党: {{ $candidate.PoliticalParty }}</p>
          </div>
        </div>
      </div>
		{{ end }}
*/
/*
  `{{ range $index, $party := .parties }}
      <div class="col-md-3">
        <div class="panel panel-default">
          <div class="panel-heading">
            <p>{{ $index | indexPlus1 }}. <a href="/political_parties/{{ $party.PoliticalParty }}">{{ $party.PoliticalParty }}</a></p>
          </div>
          <div class="panel-body">
            <p>得票数: {{ $party.VoteCount }}</p>
          </div>
        </div>
      </div>
		{{ end }}`
*/
{{ .sexRatio.men }}`</p>
        </div>
      </div>
    </div>
    <div class="col-md-6">
      <div class="panel panel-default">
        <div class="panel-heading">
          <p>女性</p>
        </div>
        <div class="panel-body">
          <p>得票数: {{ .sexRatio.women }}</p>
        </div>
      </div>
    </div>
  </div>
</div>
</body>
</html>`
