package gmHttp

import (
	"PZ_GameServer/server/game"
	room "PZ_GameServer/server/game/room"
	//	"flag"
	"PZ_GameServer/model"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

const (
	getRoomInfo = 123
)

//开启服务
func StartHttp(ip string, port int) {
	//	host := flag.String("host", ip, "listen host")
	//	port := flag.String("port", strconv.Itoa(p), "listen port")
	http.HandleFunc("/getServerInfo", getServerInfo)
	http.HandleFunc("/checkInThisServer", checkInThisServer)
	http.HandleFunc("/getReviewData", getReviewData)
	err := http.ListenAndServe(ip+":"+strconv.Itoa(port), nil)

	if err != nil {
		panic(err)
	}
}

func getServerInfo(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() //解析参数，默认是不会解析的
	//	fmt.Println(r.Form) //这些信息是输出到服务器端的打印信息
	//	fmt.Println("path", r.URL.Path)
	//	fmt.Println("scheme", r.URL.Scheme)
	//	fmt.Println(r.Form["id"])
	var id string
	var pwd string
	var cmd int
	for k, v := range r.Form {
		if k == "id" {
			id = strings.Join(v, "")
		} else if k == "pwd" {
			pwd = strings.Join(v, "")
		} else if k == "cmd" {
			str := strings.Join(v, "")
			var err error
			cmd, err = strconv.Atoi(str)
			if err != nil {
				fmt.Fprintf(w, "caonima")
				return
			}
		}
	}
	if id != "admin" && pwd != "bestPZ123" {
		fmt.Fprintf(w, "caonimabi")
		return
	}
	switch cmd {
	case getRoomInfo:
		fmt.Fprintf(w, strconv.Itoa(len(game.GServer.CheckUserList))+","+strconv.Itoa(len(room.RoomList)))
		return
	}

	fmt.Fprintf(w, "请输入正确的指令！") //这个写入到w的是输出到客户端的
}

func checkInThisServer(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() //解析参数，默认是不会解析的

}

func getReviewData(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	var reviewCode string

	for k, v := range r.Form {
		if k == "reviewCode" {
			reviewCode = strings.Join(v, "")
		}
	}

	battle, err := model.GetBattleRecordModel().GetBattleRecordByReviewCode(reviewCode)
	if err != nil {
		return
	}

	reviewData := battle.PlayBack
	fmt.Fprintf(w, reviewData)

}
