package main

import (
	"github.com/emicklei/go-restful"
	"io"
	"net/http"
	"github.com/gocql/gocql"
	_ "github.com/dimiro1/banner/autoload"
	//"fmt"
	"log"
)

//var cluster *gocql.ClusterConfig
//var session *gocql.Session

func main() {
	cluster := gocql.NewCluster("localhost")
	cluster.Keyspace = "examples"
	cluster.Consistency = gocql.One

	session, _ := cluster.CreateSession()
	defer session.Close()

	ws := new(restful.WebService)
	ws.Route(ws.GET("/hello").To(hello))
	ws.Route(ws.GET("/scylla").To(
		func(req *restful.Request, resp *restful.Response) {
			scylla(req, resp, session)
		}))

	restful.Add(ws)

	print("Web engine started at 8080")
	http.ListenAndServe(":8080", nil)
}

func hello(req *restful.Request, resp *restful.Response) {
	io.WriteString(resp, "world")
}

func scylla(req *restful.Request, resp *restful.Response, session *gocql.Session) {

	if err := session.Query(`INSERT INTO basic (txt, id, val) VALUES (?, ?, ?)`,
		"me", gocql.TimeUUID(), 10).Exec(); err != nil {
		log.Fatal(err)
	}

	/*
	var id gocql.UUID
	var text string

	if err := session.Query(`SELECT id, txt FROM basic LIMIT 1`).Scan(&id, &text); err != nil {
		log.Fatal(err)
	}
	io.WriteString(resp, fmt.Sprintf("Basic:", id, text))
	*/

	/*
	iter := session.Query(`SELECT id, txt FROM basic `).Iter()
	for iter.Scan(&id, &text) {
		fmt.Println("Basic:", id, text)
	}
	if err := iter.Close(); err != nil {
		log.Fatal(err)
	}
	*/
}