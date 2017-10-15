package Cassandra

import (
	"github.com/gocql/gocql"
	"log"
)

var Session *gocql.Session

func init() {
	var err error

	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "horaeapi"
	Session, err = cluster.CreateSession()

	if err != nil {
		panic(err)
	}
	log.Println("Cassandra connected")
}
