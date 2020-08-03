package main

import (
	"database/sql"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

var (
	pgHost     = os.Getenv("POSTGRES_HOST")
	pgPort     = 5432
	pgUser     = os.Getenv("POSTGRES_USER")
	pgDbname   = os.Getenv("POSTGRES_DB")
	pgPassword = os.Getenv("POSTGRES_PASSWORD")
)
/* Create Message(id,sender_id,receiver_id, text, timestamp)
	and Groups(grp_id,user_id,timestamp) relations */
func createRelations(db *sql.DB) error{
	sql := `create table if not exists Message( id serial primary key,
					sender_id int not null, receiver_id int not null,text varchar , epoch_timestamp int)`
	_,err:=db.Exec(sql)
	if err!= nil{
		return err
	}
	sql =`create table if not exists Groups(group_id int, user_id int,joined_at int )`
	_,err = db.Exec(sql)
	if err!= nil{
		return err
	}
	sql = `create index if not exists grp_idx on Groups(group_id)`
	_,err = db.Exec(sql)
	return err
}
// Log messages info into Messages relation
func insertMsgDb(db *sql.DB,msg Message) error{
	sql:= `insert into Message(sender_id,receiver_id,text,epoch_timestamp)
			values ($1,$2,$3,$4)`

	var err error = nil
	if msg.Sender.Id != msg.Receiver.Id{
		_,err = db.Exec(sql,msg.Sender.Id,msg.Receiver.Id,msg.Text,msg.Timestamp.Unix())
	}
	return err
}

// Log Group join event as an entry in Groups table.
func joinGroup(db *sql.DB, userId, grpId int) error{
	//sql:=`insert into Groups values($1,$2,$3)`
	sql := `insert into Groups select $1,$2,$3 where not exists 
			(select user_id from Groups where group_id=$1 and user_id=$2)`
	_,err:=db.Exec(sql,grpId,userId,time.Now().Unix())
	return err
}

/* Log Group leave event for a user by deleting the mapping between the
	user and the group in Groups relation */
func leaveGroup(db *sql.DB, userId, grpId int) error{
	sql:=`delete from Groups where group_id=$1 and user_id=$2`
	_,err:=db.Exec(sql,grpId,userId,time.Now().Unix())
	return err
}

// Gets the user Id of all the members of the grpId for group messaging
func getGroupMembers(db *sql.DB,grpId int) ([]int,error){
	sql:= `select user_id from Groups where group_id=$1`
	rows,err:=db.Query(sql,grpId)
	if err!= nil{
		return []int{},err
	}
	members := make([]int,0,5)
	for rows.Next() {
		var userId int
		err := rows.Scan(&userId)
		if err!= nil {
			return []int{}, err
		}
		logrus.Info("user_id",userId)
		members= append(members, userId)
	}
	return members,nil
}
