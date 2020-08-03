package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

const (
	ChatMessage = iota
	GroupJoinMessage
	GroupLeaveMessage
)

type Person struct {
	Id int		`json:"id"`
	Name string	`json:"name"`
}
type Message struct {
	Text string	`json:"text"`
	Sender Person	`json:"sender"`
	Receiver Person	`json:"receiver"`
	Group int		`json:"group"`
	Timestamp time.Time	`json:"timestamp"`
}

type Payload struct{
	Type int	`json:"type"`
	Message
}
type SystemConfig struct{
	host string
	port string
	broadcast chan Message
	router *mux.Router
	server http.Server
	clients map[int]*websocket.Conn
	pgsql *sql.DB
}

// initializes the server specific specific values
func newSystem() (*SystemConfig,error){
	sys := new(SystemConfig)
	flag.StringVar(&sys.host,"h","0.0.0.0","Host IP")
	flag.StringVar(&sys.port,"p","8020","Host port")
	pgsqlInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",pgHost, pgPort, pgUser,
			pgDbname, pgPassword)
	db,_:= sql.Open("postgres",pgsqlInfo)
	err := db.Ping()
	if err!= nil {
		logrus.Error(err)
		return sys,err
	}
	sys.pgsql = db
	sys.broadcast = make(chan Message)
	sys.router = mux.NewRouter()
	sys.server = http.Server{
		Addr: sys.host + ":" + sys.port,
		Handler: sys.router,
	}
	sys.clients = make(map[int]*websocket.Conn)
	return sys,nil
}

