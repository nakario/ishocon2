package main
import (
	"net/http"
	"github.com/gin-gonic/gin"
)

var codeA = []byte(`<!DOCTYPE html><html><head><meta http-equiv="Content-Type" content="text/html" charset="utf-8"><link rel="stylesheet" href="/css/bootstrap.min.css"><title>ISUCON選挙結果</title></head><body><nav class="navbar navbar-inverse navbar-fixed-top"><div class="container"><div class="navbar-header"><a class="navbar-brand" href="/">ISUCON選挙結果</a></div><div class="header clearfix"><nav><ul class="nav nav-pills pull-right"><li role="presentation"><a href="/vote">投票する</a></li></ul></nav></div></div></nav><div class="jumbotron"><div class="container"><h1>清き一票をお願いします！！！</h1></div></div><div class="container"><div class="row"><div class="col-md-6 col-md-offset-3"><div class="login-panel panel panel-default"><div class="panel-heading"><h3 class="panel-title">投票フォーム</h3></div><div class="panel-body"><form method="POST" action="/vote"><fieldset><label>氏名</label><div class="form-group"><input class="form-control" name="name" autofocus></div><label>住所</label><div class="form-group"><input class="form-control" name="address" value=""></div><label>私の番号</label><div class="form-group"><input class="form-control" name="mynumber" value=""></div><label>候補者</label><div class="form-group"><select name="candidate">									`)
var codeB =[]byte(`</select></div><label>投票理由</label><div class="form-group"><input class="form-control" name="keyword" value=""></div><label>投票数</label><div class="form-group"><input class="form-control" name="vote_count" value=""></div><div class="text-danger">`)
var codeC = []byte(`</div><input class="btn btn-lg btn-success btn-block" type="submit" name="vote" value="投票" /></fieldset></form></div></div></div></div></div></body></html>`)

func WriteVoteHTML(c *gin.Context,message string) {
	c.Status(http.StatusOK)
	c.Writer.Header()["Content-Type"] = []string{"text/html; charset=utf-8"}
	c.Writer.Write(codeA)
	c.Writer.WriteString(getAllCandidatesDOM)
	c.Writer.Write(codeB)
	c.Writer.WriteString(message)
	c.Writer.Write(codeC)
	c.Writer.WriteHeaderNow()
}