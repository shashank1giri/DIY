package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"time"
)

func init(){
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)
}

func readPayload(msg []byte,user int, sys *SystemConfig){
	var payload Payload
	err:= json.Unmarshal(msg,&payload)
	payload.Timestamp = time.Now()
	payload.Sender.Id = user

	if err== nil{
		switch payload.Type {
		case ChatMessage:
			sys.broadcast <- payload.Message
			logrus.Info("Text message received",payload.Message)

		case GroupJoinMessage:
			err:=joinGroup(sys.pgsql,user,payload.Group)
			if err!= nil{
				logrus.WithFields(logrus.Fields{
					"error":err,
					"user":user,
					"payload":payload,
				}).Error("Group Join Request failed")
			}

		case GroupLeaveMessage:
			err:=leaveGroup(sys.pgsql,user,payload.Group)
			if err!= nil{
				logrus.WithFields(logrus.Fields{
					"error":err,
					"user":user,
					"payload":payload,
				}).Error("Group Leave Request failed")
			}

		}
	}
}

func webSocketHandler(sys *SystemConfig) func(http.ResponseWriter, *http.Request){
	return func( w http.ResponseWriter, req *http.Request){
		vars := mux.Vars(req)
		logrus.Info("%v",vars["token"])
		jwtToken := vars["token"]
		isValid, user := verifyToken(jwtToken)
		if !isValid {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Request unauthorized"))
			return
		}

		upg :=websocket.Upgrader{
			CheckOrigin: func(req *http.Request) bool{
				return  true
			},
		}

		ws ,err := upg.Upgrade(w,req,nil)
		defer ws.Close()

		if err != nil{
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		sys.clients[user] = ws
		logrus.WithFields(logrus.Fields{
			"user" :user,
		}).Info("Socket connection successful")

		for{
			messagetype , msg, err:=ws.ReadMessage()
			logrus.WithFields(logrus.Fields{
				"type":messagetype,
				"msg":string(msg),

			}).Info("Server received this message from user:",user)

			if err!=nil{
				logrus.Error(err)
				return
			}
			readPayload(msg,user,sys)
		}
	}
}

func handleBroadcast(sys *SystemConfig){
	for{
		msg := <- sys.broadcast
		if msg.Group!=0 {
			members, _ := getGroupMembers(sys.pgsql, msg.Group)
			logrus.Info(members)
			jsonMsg,_ := json.Marshal(msg)
			for _,member:= range members {
				if sys.clients[member] == nil {
					continue
				}
				msg.Receiver.Id = member
				logrus.Info("Message sent to receiver:",member, string(jsonMsg))
				_=(sys.clients[member]).WriteMessage(1,jsonMsg)
				_=insertMsgDb(sys.pgsql,msg)
			}

		}else{

			for id,client := range sys.clients{
				msg.Receiver.Id = id
				jsonStr,_ := json.Marshal(msg)
				logrus.Info("Message sent to receiver:",id, string(jsonStr))
				_=client.WriteMessage(1,jsonStr)
				_=insertMsgDb(sys.pgsql,msg)
			}
		}
	}
}

func main(){
	sys,err := newSystem()
	if err!= nil {
		logrus.WithFields(logrus.Fields{
			"system config":sys,
			"err":err,
		}).Error("Unable to initialize the server")
		return
	}

	sys.router.HandleFunc("/ws/{token}",webSocketHandler(sys))
	logrus.Info("Http Server starting at addr: ", sys.host + ":" + sys.port)

	err = createRelations(sys.pgsql)
	if err != nil{
		logrus.Error(err)
		return
	}

	go handleBroadcast(sys)

	logrus.Fatal(sys.server.ListenAndServe())
}

