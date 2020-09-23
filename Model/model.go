package Model

import (
	"Crewl/Config"
	"database/sql"
	"fmt"
	"github.com/go-redis/redis"
	"log"
)

func Authorized(data *Config.LoginForm) bool {
	var tmp string
	err := db.QueryRow("select id from users where username=? and password=? limit 1", data.Username, data.Password).Scan(&tmp)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}
		log.Fatal(err, "Mysql Auth Error")
	}
	return true
}

func GetComments(page int) ([]Config.Comment, int, int) {
	var data = make([]Config.Comment, 10)
	rows, err := db.Query("select * from daily limit ?,10", (page-1) * 10)
	defer func() {
		if rows != nil {
			rows.Close()   //关闭掉未scan的sql连接
		}
	}()
	if err != nil {
		fmt.Printf("Query failed,err:%v\n", err)
		return nil, 0, 0
	}
	var i int = 0
	for rows.Next() {
		err = rows.Scan(&data[i].Id, &data[i].Bid, &data[i].Uid, &data[i].Content, &data[i].Time)  //不scan会导致连接不释放
		if err != nil {
			fmt.Printf("Scan failed,err:%v\n", err)
			return nil, 0, 0
		}
		//fmt.Println("scan successd:", data[i])
		i++
	}

	var allCount int
	err = db.QueryRow("select Count(*) from daily").Scan(&allCount)
	return data, i, allCount
}

func GetValidComments(page int) ([]Config.InvalidComment, int, int) {
	var data = make([]Config.InvalidComment, 10)
	rows, err := db.Query("select * from exception limit ?,10", (page-1) * 10)
	defer func() {
		if rows != nil {
			rows.Close()   //关闭掉未scan的sql连接
		}
	}()
	if err != nil {
		fmt.Printf("Query failed,err:%v\n", err)
		return nil, 0, 0
	}
	var i int = 0
	for rows.Next() {
		err = rows.Scan(&data[i].Id, &data[i].Bid, &data[i].Uid, &data[i].Content, &data[i].Time, &data[i].Sensitive)  //不scan会导致连接不释放
		if err != nil {
			fmt.Printf("Scan failed,err:%v\n", err)
			return nil, 0, 0
		}
		fmt.Println("scan successd:", data[i])
		i++
	}

	var allCount int
	err = db.QueryRow("select Count(*) from exception").Scan(&allCount)
	return data, i, allCount
}


// ------------- Redis Function --------------- //
func RedisAuth(cookie *Config.UserCookie) bool {
	UserName, err := Client.Get(cookie.BI_UUID).Result()
	if err == redis.Nil {
		return false
	} else if err != nil {
		log.Fatal(err)
	}
	if UserName == cookie.UserName {
		return true
	}
	return false
}
