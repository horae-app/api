package Cassandra

import (
	"github.com/gocql/gocql"
	"log"

	settings "github.com/horae-app/api/Settings"
)

var Session *gocql.Session

func init() {
	var err error

	cluster := gocql.NewCluster(settings.DB_ENDPOINT)
	cluster.Keyspace = settings.DB_NAME
	Session, err = cluster.CreateSession()

	if err != nil {
		panic(err)
	}
	log.Println("Cassandra connected")
}
