package main

import (
	//"github.com/emicklei/go-restful"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"github.com/gocql/gocql"
	_ "github.com/dimiro1/banner/autoload"
	"log"
	"fmt"
)

func main() {
	cluster := gocql.NewCluster("localhost")
	cluster.Keyspace = "examples"
	cluster.Consistency = gocql.One

	session, _ := cluster.CreateSession()
	defer session.Close()

	setupWebServer(session);
}

func setupWebServer(session *gocql.Session) {
	router := httprouter.New()
	router.GET("/hello", hello)
	router.GET("/scylla", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		scylla(w, r, ps, session)
	})

	log.Fatal(http.ListenAndServe(":8080", router))
}

func hello(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "World!\n")
}

func scylla(w http.ResponseWriter, r *http.Request, _ httprouter.Params, session *gocql.Session) {
	if err := session.Query(`INSERT INTO basic (txt, id, val) VALUES (?, ?, ?)`,
		"me", gocql.TimeUUID(), 10).Exec(); err != nil {
		log.Fatal(err)
	}

	fmt.Fprint(w, "OK")
}