package game

import (
	"fmt"
	"math"
	"reflect"
	"regexp"
	"strconv"
	"time"

	"PZ_GameServer/common/gift"
	"PZ_GameServer/common/random_name"
	"PZ_GameServer/common/testCard"
	"PZ_GameServer/common/util"
	"PZ_GameServer/model"
	"PZ_GameServer/protocol/def"
	"PZ_GameServer/protocol/pb"
	"PZ_GameServer/sdk"
	//"PZ_GameServer/server/game/common"
	err "PZ_GameServer/server/game/error"
	room "PZ_GameServer/server/game/room"
	"PZ_GameServer/server/game/room/ningbo"
	"PZ_GameServer/server/game/room/pinshi"
	"PZ_GameServer/server/game/room/srddz"
	"PZ_GameServer/server/game/room/xiangshan"
	"PZ_GameServer/server/game/room/xizhou"
	rb "PZ_GameServer/server/game/roombase"
	"PZ_GameServer/server/user"

	"github.com/golang/protobuf/proto"
)

// 初始化基础路由
func InitBaseRouter() {

	var normal int32 = 0

	//---- server
	Bind(mjgame.MsgID_MSG_Login, normal, &mjgame.Login{}, Login)                                   //登录
	Bind(mjgame.MsgID_MSG_Create_Room, normal, &mjgame.Create_Room{}, CreateRoom)                  //创建房间
	Bind(mjgame.MsgID_MSG_Match_Room, normal, &mjgame.Match_Room{}, MatchRoom)                     //匹配玩家
	Bind(mjgame.MsgID_MSG_Cancel_Match_Room, normal, &mjgame.Cancel_Match_Room{}, CancelMatchRoom) //匹配玩家
	Bind(mjgame.MsgID_MSG_Room_List, normal, &mjgame.Room_List{}, RoomList)                        //房间列表(未使用)
	Bind(mjgame.MsgID_MSG_Into_Room, normal, &mjgame.Into_Room{}, IntoRoom)                        //进入房间
	Bind(mjgame.MsgID_MSG_Into_MatchRoom, normal, &mjgame.Into_MatchRoom{}, Into_MatchRoom)        //进入自动匹配场房间(未使用)
	Bind(mjgame.MsgID_MSG_Room_Info, normal, &mjgame.Room_Info{}, RoomInfo)                        //房间信息(未使用)
	Bind(mjgame.MsgID_MSG_User_Info, normal, &mjgame.User_Info{}, UserInfo)                        //用户信息
	Bind(mjgame.MsgID_MSG_Find_Room, normal, &mjgame.Find_Room{}, FindRoom)                        //查询房间(未使用)
	Bind(mjgame.MsgID_MSG_Standup, normal, &mjgame.StandUp{}, StandUp)                             //起立
	Bind(mjgame.MsgID_MSG_Sitdown, normal, &mjgame.SitDown{}, SitDown)                             //坐下
	Bind(mjgame.MsgID_MSG_Exit_Room, normal, &mjgame.ExitRoom{}, ExitRoom)                         //离开房间
	Bind(mjgame.MsgID_MSG_Start_Game, normal, &mjgame.Start_Game{}, StartGame)                     //开始游戏
	Bind(mjgame.MsgID_MSG_User_Ready, normal, &mjgame.UserReady{}, Ready)                          //准备
	Bind(mjgame.MsgID_MSG_Put_Card, normal, &mjgame.Put_Card{}, PutCard)                           //出牌
	Bind(mjgame.MsgID_MSG_Chow, normal, &mjgame.Chow{}, Chow)                                      // 吃
	Bind(mjgame.MsgID_MSG_Peng, normal, &mjgame.Peng{}, Peng)                                      // 碰
	Bind(mjgame.MsgID_MSG_Pass, normal, &mjgame.Pass{}, Pass)                                      // 过
	Bind(mjgame.MsgID_MSG_Kong, normal, &mjgame.Kong{}, Kong)                                      // 杠
	Bind(mjgame.MsgID_MSG_Win, normal, &mjgame.Win{}, Win)                                         // 胡
	Bind(mjgame.MsgID_MSG_Chat, normal, &mjgame.Chat{}, Chat)                                      // 聊天
	Bind(mjgame.MsgID_MSG_Notice, normal, &mjgame.Notice{}, Notice)                                // 公告跑马灯
	Bind(mjgame.MsgID_MSG_Disband, normal, &mjgame.Disband{}, DisbandRoom)                         // 解散房间
	Bind(mjgame.MsgID_MSG_Roomowner_Disband_Room, normal, &mjgame.Roomowner_Disband_Room{}, RoomownerDisbandRoom)
	Bind(mjgame.MsgID_MSG_Vote, normal, &mjgame.Vote{}, Vote)                                     // 投票
	Bind(mjgame.MsgID_MSG_Battle_Record, normal, &mjgame.BattleRecordRequest{}, BattleRecordList) // 战绩记录
	Bind(mjgame.MsgID_MSG_Home_Owner, normal, &mjgame.HomeOwnerRequest{}, HomeOwnerList)          // 房主记录列表
	Bind(mjgame.MsgID_MSG_Room_Summary, normal, &mjgame.RoomSummaryRequest{}, RoomSummary)        // 房间结算统计
	Bind(mjgame.MsgID_MSG_BattleDetail, normal, &mjgame.BattleDetail{}, BattleDetail)
	Bind(mjgame.MsgID_MSG_Room_Kick, normal, &mjgame.KickRequest{}, Kick)                    // 踢人
	Bind(mjgame.MsgID_MSG_MessageJson, normal, &mjgame.MessageJson{}, MessageJson)           // messagejson
	Bind(mjgame.MsgID_MSG_NOTIFY_RECHARGE, normal, &mjgame.NotifyRecharge{}, NotifyRecharge) // 通知充值
	Bind(mjgame.MsgID_MSG_Gift, normal, &mjgame.Gift{}, Gift)                                // 发送礼物

	//三人斗地主相关
	Bind(mjgame.MsgID_MSG_Sddz_Jiaofen, normal, &mjgame.Sddz_Jiaofen{}, Sddz_Jiaofen)
	Bind(mjgame.MsgID_MSG_Sddz_Jiabei, normal, &mjgame.Sddz_Jiabei{}, Sddz_Jiabei)
	Bind(mjgame.MsgID_MSG_Sddz_Mingpai, normal, &mjgame.Sddz_Mingpai{}, Sddz_Mingpai)
	Bind(mjgame.MsgID_MSG_Sddz_Chupai, normal, &mjgame.Sddz_Chupai{}, Sddz_Chupai)
	Bind(mjgame.MsgID_MSG_Sddz_Pass, normal, &mjgame.Sddz_Chupai{}, Sddz_Pass)

	//上海四人斗地主相关
	Bind(mjgame.MsgID_MSG_Srddz_Baodao, normal, &mjgame.Srddz_Baodao{}, Srddz_Baodao)
	Bind(mjgame.MsgID_MSG_Srddz_StrictWin, normal, &mjgame.Srddz_StrictWin{}, Srddz_StrictWin)
	//拼十相关
	Bind(mjgame.MsgID_MSG_Nn_Xiazhu, normal, &mjgame.Nn_Xiazhu{}, Nn_Xiazhu)

	//---- 特殊规则
	//	Bind(mjgame.MsgID_MSG_Benefits, &mjgame.Benefits{}, Benefits)          //救济金
	//	Bind(mjgame.MsgID_MSG_Change3Card, &mjgame.Change3Card{}, Change3Card) //换3张
	//	Bind(mjgame.MsgID_MSG_FixMiss, &mjgame.FixMiss{}, FixMiss)             //定缺

}

//登录请求
func Login(args ...interface{}) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Login Error...")
			fmt.Println(err)
		}
	}()
	m := &mjgame.Login{}
	er := proto.Unmarshal(*(args[0].(*[]byte)), m)
	a := args[1].(*user.User)
	if er != nil || a == nil {
		a.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrMessageError)
		return
	}

	if len(m.Openid) == 0 {
		a.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrInvalidOpenId)
		return
	}
	var userInfo *sdk.UserInfo

	if len(m.Openid) != 0 && len(m.Token) != 0 {
		r, err1 := sdk.CheckTokenIsValid(m.Openid, m.Token)
		if err1 != nil {
			a.SendMessage(mjgame.MsgID_MSG_ACK_Error, &mjgame.ErrorItem{1, r.Errmsg})
		}
		if r.Errcode != 0 {
			a.SendMessage(mjgame.MsgID_MSG_ACK_Error, &mjgame.ErrorItem{int32(r.Errcode), r.Errmsg})
		}
		userInfo, err1 = sdk.GetUserInfo(m.Openid, m.Token)
		if err1 != nil {
			a.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrOpenIdNotEmpty)
		}
	}

	var unionId string
	if userInfo == nil {
		unionId = m.Openid
	}

	user, _ := model.GetUserModel().GetUserByOpenId(m.Openid)
	if user == nil {
		//a.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrInvalidOpenId)
		//return
		// 数据库中无用户, 有可能是零时删表造成的, 但客户端还保留openid.
		unionId = util.GetSID()
		if m.Openid == "" {
			m.Openid = util.GetSID()
		}
		//fmt.Println("chuangjianxinyonghu", m.Openid)
		user = CreateUser(m.GPS_LAT, m.GPS_LNG, a.Conn.RemoteAddr().String(), userInfo, unionId, m.Openid) //创建新用户
		fmt.Println("chuangjianxinyonghu", user.ID)
	} else {
		//fmt.Println("yonghucunzai:: ", m.Openid)
	}

	if user.Sid == "" {
		user.Sid = user.OpenID
	}

	//	fmt.Println("user=", user == nil)
	//	fmt.Println("CheckUserList=", GServer.CheckUserList == nil)
	//重复登录
	//if v, ok := GServer.CheckUserList[user.ID]; ok {
	if v, ok := GServer.GetLockCheckUser(user.ID); ok {
		//已经存在用户
		v.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrHasLoginOtherDevice)
		v.Conn.Close()
		time.Sleep(100 * time.Millisecond)
	}

	a.User = user
	//if v, ok := GServer.UserList[user.Sid]; ok {
	if v, ok := GServer.GetLockUser(user.Sid); ok {
		a.GameType = v.GameType
		a.User.RoomId = v.RoomId
	}

	GServer.mux.Lock()
	GServer.CheckUserList[user.ID] = a
	GServer.UserList[user.Sid] = a
	GServer.mux.Unlock()

	ackLogin := mjgame.ACK_Login{
		Sid: user.Sid,
	}
	user.State = def.Online
	a.SendMessage(mjgame.MsgID_MSG_ACK_Login, &ackLogin)
	UserInfo(nil, a)

	//断线重连
	if user.RoomId > 0 {
		//roomHandle, ok := room.RoomList[user.RoomId]
		roomHandle, ok := room.GetLockRoomHandle(user.RoomId)
		if ok {
			BroadcastUserState(a, user.RoomId)
			//roomHandle.Room.(*xiangshan.RoomXiangshan).IntoUser(a)
			FunCall(roomHandle.Room, "IntoUser", []reflect.Value{reflect.ValueOf(a)})
		}
	}

	match := room.GetMatchInfo(a)

	fmt.Println("wanjiapipeizhong00", match)
	if match != nil {
		fmt.Println("wanjiapipeizhong")
		//提示玩家正在匹配中
		ack := mjgame.ACK_Match_Room{UID: string(a.ID), Type: match.Type, City: match.City, PWD: match.PWD, Rule: match.Rule}
		a.SendMessage(mjgame.MsgID_MSG_ACK_Match_Room, &ack)
		return
	}
}

//创建新的用户																																																																																																																				```
func CreateUser(GPS_LNG, GPS_LAT float32, ip string, userInfo *sdk.UserInfo, unionID string, openid string) *model.User {
	var user model.User

	user.Sid = util.GetSID()
	user.LastIp = ip
	user.Coin = 20
	user.Diamond = 6
	user.GPS_LAT = GPS_LNG
	user.GPS_LNG = GPS_LAT
	if userInfo != nil {
		user.NickName = userInfo.Nickname
		user.Sex = userInfo.Sex
		user.Province = userInfo.Province
		user.City = userInfo.City
		user.Country = userInfo.Country
		user.Icon = userInfo.Headimgurl
	} else {
		user.NickName = random_name.GetRandomName()
		//user.Icon = "icon_" + strconv.Itoa(util.RandInt(0, 5))
		user.Icon = ""
		user.Sex = 1
	}
	user.OpenID = openid
	user.UnionID = unionID

	err := model.GetUserModel().Create(&user)
	if err != nil {
		fmt.Println("charushibai::", err.Error())
	}

	return &user
}

//创建新的AI RoBot
func CreateNewAIRoBot(gType int32) *user.User {
	//	//m.GPS_LAT, m.GPS_LNG, a.Conn.RemoteAddr().String(), userInfo, unionId
	//	ai := CreateUser(0, 0, "AI", nil, "unionId", "testopenid") //创建新用户

	//	a := &user.User{}
	//	a.User = ai
	//	a.GameType = gType

	//	a.ID = util.GetUID()
	//	fmt.Println("ai id = ", a.ID)
	//	a.Sid = util.GetSID()
	//	a.NickName = "机器人" + string(a.ID)
	//	a.LastIp = ""

	//	a.Coin = 50000
	//	a.Diamond = 10000
	//	a.GPS_LAT = util.RandFloat(0, 10)
	//	a.GPS_LNG = util.RandFloat(0, 10)
	//	a.Icon = "icon_" + strconv.Itoa(util.RandInt(0, 5))
	//	return a
	return nil
}

//用户信息
func UserInfo(args ...interface{}) {
	a := args[1].(*user.User)
	if a == nil {
		//fmt.Printf("denglucuowu!!")
		a.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrMessageError)
		return
	}

	ack := mjgame.ACK_User_Info{
		Name:      a.NickName,
		Uid:       strconv.Itoa(a.ID),
		RoomId:    int32(a.RoomId),
		Ip:        "",
		Index:     0,
		Icon:      a.Icon,
		Coin:      int32(a.Coin),
		Type:      0,
		Diamond:   int32(a.Diamond),
		Email:     "",
		Sex:       int32(a.Sex),
		State:     int32(a.State),
		Level:     int32(a.Level),
		ParentUid: strconv.Itoa(a.ParentUid),
	}
	//fmt.Printf("dengluchenggong::", a.ID)

	a.SendMessage(mjgame.MsgID_MSG_ACK_User_Info, &ack)
}

//房间列表
func RoomList(args ...interface{}) {
	m := &mjgame.Room_List{}
	er := proto.Unmarshal(*(args[0].(*[]byte)), m)
	a := args[1].(*user.User)
	if er != nil || a == nil {
		a.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrMessageError)
		return
	}

	ack := mjgame.ACK_Room_List{
		Type: m.Type,
		City: 0,
	}
	a.SendMessage(mjgame.MsgID_MSG_ACK_Room_List, &ack)
}

//创建房间
func CreateRoom(args ...interface{}) {
	//	m := args[0].(*mjgame.Create_Room)
	//	a := args[1].(*user.User)
	m := &mjgame.Create_Room{}

	er := proto.Unmarshal(*(args[0].(*[]byte)), m)
	fmt.Println("roomType::", m.Type)
	//fmt.Println("m.Type ", m.Type)
	a := args[1].(*user.User)
	//停服不可以创建房间并发送公告
	if StopGameNotice != nil {
		a.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrMessageServerRepairs)
		SendStopGameNotice(a)
		return
	}

	//Rule:25 Rule:19 Rule:23
	if er != nil || a == nil {
		a.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrMessageError)
		return
	}
	roomId := GServer.CreateRoom(m, a, -1)

	fmt.Println("roomid::" + strconv.Itoa(roomId))

	if roomId > 0 {
		ack := mjgame.ACK_Create_Room{RID: int32(roomId), PWD: "", Rule: m.Rule, Type: m.Type}
		a.SendMessage(mjgame.MsgID_MSG_ACK_Create_Room, &ack)
	} else {
		a.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrCreateRoom)
	}
}

//匹配房间
func MatchRoom(args ...interface{}) {

	fmt.Println("v....MatchRoom")

	m := &mjgame.Match_Room{}
	er := proto.Unmarshal(*(args[0].(*[]byte)), m)
	fmt.Println("roomType::", m.Type)
	fmt.Println("m.Type ", m.Type)
	a := args[1].(*user.User)
	//停服不可以创建房间并发送公告
	if StopGameNotice != nil {
		a.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrMessageServerRepairs)
		SendStopGameNotice(a)
		return
	}

	//Rule:25 Rule:19 Rule:23
	if er != nil || a == nil {
		a.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrMessageError)
		return
	}

	//如果玩家正在匹配中，则不让匹配了，必须先取消匹配
	fmt.Println("v....1qian000")
	match := room.GetMatchInfo(a)

	fmt.Println("v....1qian")

	if match != nil {
		//提示玩家正在匹配中
		ack := mjgame.ACK_Match_Room{UID: string(a.ID), Type: match.Type, City: match.City, PWD: match.PWD, Rule: match.Rule}
		a.SendMessage(mjgame.MsgID_MSG_ACK_Match_Room, &ack)
		return
	}
	fmt.Println("v....1hou")

	roomHandle, users := room.MatchRoom(m, a)
	//给users做个保护
	for _, u := range users {
		user, ok := GServer.GetLockUser(u.Sid)
		if ok {
			u = user
		}
	}

	if roomHandle != nil && users != nil {
		for index, v := range users {
			FunCall(roomHandle.Room, "IntoUser", []reflect.Value{reflect.ValueOf(v)})
			fmt.Println("v....0", v)
			sitdown := &mjgame.SitDown{
				Sid:   v.Sid,
				Index: int32(index),
			}

			fmt.Println("v....1", v.ID)

			FunCall(roomHandle.Room, "SitDown", []reflect.Value{reflect.ValueOf(v), reflect.ValueOf(sitdown)})
		}
		FunCall(roomHandle.Room, "Start", []reflect.Value{reflect.ValueOf(users[0])})
	} else {
		ack := mjgame.ACK_Match_Room{UID: string(a.ID), Type: m.Type, City: m.City, PWD: m.PWD, Rule: m.Rule}
		a.SendMessage(mjgame.MsgID_MSG_ACK_Match_Room, &ack)
	}

}

//取消匹配
func CancelMatchRoom(args ...interface{}) {
	m := &mjgame.Cancel_Match_Room{}
	er := proto.Unmarshal(*(args[0].(*[]byte)), m)
	a := args[1].(*user.User)
	if er != nil || a == nil {
		a.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrMessageError)
		return
	}

	//将玩家从匹配列表中删除
	result := room.DeleMatchFromList(a)
	if result {
		ack := mjgame.ACK_Cancel_Match_Room{UID: string(a.ID)}
		a.SendMessage(mjgame.MsgID_MSG_ACK_Cancel_Match_Room, &ack)
	}
}

// 查找房间
func FindRoom(args ...interface{}) {
	//	m := args[0].(*mjgame.Find_Room)
	//	a := args[1].(*user.User)

	//	if !CheckUser(a.Sid) {
	//		a.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrNeedLogin)
	//		return
	//	}

	//	rid := int(m.GetRID())
	//	//ms := ""
	//	usercount := -1

	//	r, ok := GServer.RoomList[rid]
	//	if ok && r != nil {
	//		usercount = r.GetUserCount()

	//	} else {
	//		//ms = "没有找到房间"
	//	}

	//	//a.Conn.WriteMsg(j.Marshal(msg.ACK_Find_Room{RID: rid, UserCount: usercount, MSG: ms}))
	//	a.SendMessage(mjgame.MsgID_MSG_ACK_Find_Room, &mjgame.ACK_Find_Room{
	//		RID:       int32(rid),
	//		UserCount: int32(usercount),
	//	})
}

//直接开始游戏, (空缺的位置机器自动打)
func StartGame(args ...interface{}) {
	a := args[1].(*user.User)

	if !CheckUser(a.Sid) {
		a.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrNeedLogin)
		return
	}

	//roomHandle, ok := room.RoomList[a.RoomId]
	roomHandle, ok := room.GetLockRoomHandle(a.RoomId)
	if ok {
		FunCall(roomHandle.Room, "Start", []reflect.Value{reflect.ValueOf(a)})
	} else {
		a.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrRoomNotExist) // 房间不存在
	}

}

// 重新开始游戏
func RestartGame(args ...interface{}) {
	//m := args[0].(*mjgame.Restart_Game)
	//a := args[1].(*user.User)

	//	if !CheckUser(a.Sid) {
	//		a.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrNeedLogin)
	//		return
	//	}

	//	r, ok := GServer.RoomList[a.RoomId]
	//	if ok {
	//		if r.GetIndex(string(a.ID)) < 0 {
	//			a.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrSelfNotInRoom)
	//			return
	//		}

	//		ok := r.CheckCanStart()
	//		if ok {
	//			r.Restart()
	//		} else {
	//			a.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrPlayHasNotEnoughGold)
	//		}
	//	} else {
	//		//房间不存在
	//		a.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrRoomNotExist)
	//	}
}

//检查用户是否存在
// true = 存在
func CheckUser(sid string) bool {
	if sid == "" {
		return false
	} else {
		//v, ok := GServer.UserList[sid]
		v, ok := GServer.GetLockUser(sid)
		if ok == true && v != nil && v.User != nil {
			return true
		}
	}
	return false
}

//------------------------------------------
// 进入房间
func IntoRoom(args ...interface{}) {
	m := &mjgame.Into_Room{}
	er := proto.Unmarshal(*(args[0].(*[]byte)), m)
	//m := args[0].(*mjgame.Into_Room)
	a := args[1].(*user.User)

	fmt.Println("IntoRoom ", int(m.GetRID()))

	if er != nil || a == nil || m.RID <= 0 {
		a.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrMessageError)
		return
	}

	//停服不可以进入房间并发送公告
	if StopGameNotice != nil {
		a.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrMessageServerRepairs)
		SendStopGameNotice(a)
		return
	}

	if a.User.RoomId > 0 {
		//_, oldroom_ok := room.RoomList[a.User.RoomId] // 如果以前的room已经被销毁
		_, oldroom_ok := room.GetLockRoomHandle(a.User.RoomId) // 如果以前的room已经被销毁
		if oldroom_ok {
			m.RID = int32(a.User.RoomId)
		}
	}

	//roomHandle, hasRoom := room.RoomList[int(m.GetRID())]           // 内存中是否存在
	roomHandle, hasRoom := room.GetLockRoomHandle(int(m.GetRID()))    // 内存中是否存在
	redisRoomid, roomStruct := rb.Redis_CheckRoom(int(m.RID), m.Code) // Redis中是否存在

	// 如果内存中没有, Redis有,  则是新创建房间
	// 如果内存中有, Redis有, union_code 也对的上, 则是进入房间
	// 如果内存中有, Redis没有, 则直接进入房间
	// 如果内存中没有, redis也没有, 则返回失败(房间已经解散, 或者结束)
	if !hasRoom && redisRoomid > 0 {

		// 创建新房间
		roomreq := mjgame.Create_Room{Type: int32(roomStruct.Game_type), Rule: roomStruct.Rules}
		rid := GServer.CreateRoom(&roomreq, a, redisRoomid)
		fmt.Println("创建新房间 ", redisRoomid)

		if rid > 0 {
			ack := mjgame.ACK_Create_Room{RID: int32(rid), Type: int32(roomStruct.Game_type), Rule: roomStruct.Rules}
			a.SendMessage(mjgame.MsgID_MSG_ACK_Create_Room, &ack)
			return
		} else {
			a.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrCreateRoom)
			return
		}
		roomHandle, hasRoom = room.GetLockRoomHandle(redisRoomid)

	}
	if hasRoom && roomHandle != nil {

		// 进入房间
		//		if len(m.Code) > 0 {
		//			var uniqueCode string
		//			if roomHandle.GameType == int(mjgame.MsgID_GTYPE_ZheJiang_XiangShan) {
		//				uniqueCode = roomHandle.Room.(*xiangshan.RoomXiangshan).UniqueCode
		//			} else if roomHandle.GameType == int(mjgame.MsgID_GTYPE_ZheJiang_XiZhou) {
		//				uniqueCode = roomHandle.Room.(*xizhou.RoomXiZhou).UniqueCode
		//			}
		//			if uniqueCode != m.Code {
		//				a.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrRoomNotExist)
		//				return
		//			} else {
		//				a.User.RoomId = redisRoomid
		//			}
		//		} else {
		//			// 进入房间
		//			a.User.RoomId = redisRoomid
		//		}

		a.User.RoomId = redisRoomid
		FunCall(roomHandle.Room, "IntoUser", []reflect.Value{reflect.ValueOf(a)})

	} else {

		a.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrRoomNotExist) // 房间不存在

	}

}

// 查看战绩
func ShowRoomScore(u *user.User, code string) {

}

// 进入自动匹配场房间
func Into_MatchRoom(args ...interface{}) {
	m := &mjgame.Into_MatchRoom{}
	er := proto.Unmarshal(*(args[0].(*[]byte)), m)
	//m := args[0].(*mjgame.Into_MatchRoom)
	a := args[1].(*user.User)
	if er != nil || a == nil {
		a.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrMessageError)
		return
	}

	if !CheckUser(a.Sid) {
		a.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrNeedLogin)
		return
	}

	gameType := m.GetType()
	roomreq := mjgame.Create_Room{Type: gameType, Rule: []int32{6}}

	rid := GServer.CreateRoom(&roomreq, a, -1)
	//rHandle, ok := room.RoomList[rid]
	rHandle, ok := room.GetLockRoomHandle(rid)

	if ok {
		FunCall(rHandle.Room, "IntoUser", []reflect.Value{reflect.ValueOf(a)})
		FunCall(rHandle.Room, "IntoUser", []reflect.Value{reflect.ValueOf(CreateNewAIRoBot(gameType))})

		//FunCall(rHandle.Room, "IntoUser", []reflect.Value{reflect.ValueOf(CreateNewAIRoBot(gameType))})
		//FunCall(rHandle.Room, "IntoUser", []reflect.Value{reflect.ValueOf(CreateNewAIRoBot(gameType))})

	} else {
		a.SendMessage(mjgame.MsgID_MSG_ACKBC_Into_Room, err.ErrRoomNotExist) // 房间不存在
	}

	//	fmt.Println("新开匹配房间 ", rid)

	//	if ok && room != nil { // 进入房间

	//		a.RoomId = rid

	//		room.IntoUser(a)
	//		room.IntoUser(CreateNewAIRoBot())
	//		room.IntoUser(CreateNewAIRoBot())
	//		room.IntoUser(CreateNewAIRoBot())

	//	} else {
	//		a.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrAutoIntoRoomFailed)
	//	}
}

//退出房间
func ExitRoom(args ...interface{}) {
	a := args[1].(*user.User)

	if !CheckUser(a.Sid) {
		a.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrNeedLogin)
		return
	}

	//roomHandle, ok := room.RoomList[a.RoomId]
	roomHandle, ok := room.GetLockRoomHandle(a.RoomId)
	if ok {
		FunCall(roomHandle.Room, "ExitUser", []reflect.Value{reflect.ValueOf(a)})
	} else {
		a.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrRoomNotExist) // 房间不存在
	}
}

//房间信息
func RoomInfo(args ...interface{}) {
	//	m := args[0].(*mjgame.Room_Info)
	//	a := args[1].(*user.User)

	//	room := GServer.GetRoom(a.RoomId, m.GetPWD())

	//	if room != nil {
	//		roomInfo := mjgame.ACK_Room_Info{
	//			RoomId:    int32(room.RID),
	//			Type:      int32(room.Type),
	//			City:      int32(0),
	//			Level:     int32(0),
	//			SeatCount: int32(len(room.Seat)),
	//		}
	//		a.SendMessage(mjgame.MsgID_MSG_ACK_RoomInfo, &roomInfo)
	//	} else {
	//		a.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrRoomNotExist)
	//	}
}

//救济金
func Benefits(args ...interface{}) {

	a := args[1].(*user.User)
	if a.Coin <= 0 {
		reward := 4000
		a.Coin += reward
		ack := mjgame.ACK_Benefits{UID: string(a.ID), Reward: int32(reward), Coin: int32(a.Coin)}
		a.SendMessage(mjgame.MsgID_MSG_ACK_Benefits, &ack)
	} else {
		ack := mjgame.ACK_Benefits{UID: string(a.ID), Reward: int32(0), Coin: int32(a.Coin)}
		a.SendMessage(mjgame.MsgID_MSG_ACK_Benefits, &ack)
	}
}

//进入房间坐下
func SitDown(args ...interface{}) {
	m := &mjgame.SitDown{}
	er := proto.Unmarshal(*(args[0].(*[]byte)), m)
	//m := args[0].(*mjgame.SitDown)
	a := args[1].(*user.User)
	if er != nil || a == nil {
		fmt.Println("Sitdown m ", m, a)
		a.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrMessageError)
		return
	}

	if m.Index < 0 {
		fmt.Println("Sitdown m.Index ", m.Index)
		a.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrInvalidParam)
		return
	}

	//roomHandle, ok := room.RoomList[a.RoomId]
	roomHandle, ok := room.GetLockRoomHandle(a.RoomId)
	if ok {
		//uesrp := user.RunParam{Room: roomHandle.Room, FunName: "SitDown", Params: []reflect.Value{reflect.ValueOf(a), reflect.ValueOf(m)}}
		//FunCall(roomHandle.Room, "AddMsgList", []reflect.Value{reflect.ValueOf(a), reflect.ValueOf(uesrp)})
		FunCall(roomHandle.Room, "SitDown", []reflect.Value{reflect.ValueOf(a), reflect.ValueOf(m)})
	} else {
		a.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrRoomNotExist) // 房间不存在
	}
}

//进入房间起立
func StandUp(args ...interface{}) {
	m := &mjgame.StandUp{}
	er := proto.Unmarshal(*(args[0].(*[]byte)), m)
	//m := args[0].(*mjgame.StandUp)
	a := args[1].(*user.User)
	if er != nil || a == nil {
		a.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrMessageError)
		return
	}

	if m.Index < 0 {
		a.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrInvalidParam)
		return
	}

	//roomHandle, ok := room.RoomList[a.RoomId]
	roomHandle, ok := room.GetLockRoomHandle(a.RoomId)
	if ok {
		FunCall(roomHandle.Room, "StandUp", []reflect.Value{reflect.ValueOf(m), reflect.ValueOf(a)})
	} else {
		a.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrRoomNotExist) // 房间不存在
	}
}

//玩家准备
func Ready(args ...interface{}) {
	//_ = args[0].(*mjgame.UserReady)
	a := args[1].(*user.User)
	if a == nil || a.RoomId == 0 {
		a.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrMessageError)
		return
	}

	//roomHandle, ok := room.RoomList[a.RoomId]
	roomHandle, ok := room.GetLockRoomHandle(a.RoomId)
	if ok {
		FunCall(roomHandle.Room, "Ready", []reflect.Value{reflect.ValueOf(a)})
	} else {
		a.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrRoomNotExist) // 房间不存在
	}
}

// 出牌
func PutCard(args ...interface{}) {
	//m := args[0].(*mjgame.Put_Card)
	//	data := args[0].(*[]byte)
	//	a := args[1].(*user.User)
	//	var m = &mjgame.Put_Card{}
	//	proto.Unmarshal(*data, m)

	m := &mjgame.Put_Card{}
	er := proto.Unmarshal(*(args[0].(*[]byte)), m)
	a := args[1].(*user.User)
	if er != nil || a == nil {
		a.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrMessageError)
		return
	}

	//roomHandle, ok := room.RoomList[a.RoomId]
	roomHandle, ok := room.GetLockRoomHandle(a.RoomId)
	if ok {
		//uesrp := user.RunParam{Room: roomHandle.Room, FunName: "PutCard", Params: []reflect.Value{reflect.ValueOf(a), reflect.ValueOf(m)}}
		//FunCall(roomHandle.Room, "AddMsgList", []reflect.Value{reflect.ValueOf(a), reflect.ValueOf(uesrp)})
		FunCall(roomHandle.Room, "PutCard", []reflect.Value{reflect.ValueOf(a), reflect.ValueOf(m)})
	} else {
		a.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrRoomNotExist) // 房间不存在
	}
}

// 吃
func Chow(args ...interface{}) {
	m := &mjgame.Chow{}
	er := proto.Unmarshal(*(args[0].(*[]byte)), m)
	//m := args[0].(*mjgame.Chow)
	a := args[1].(*user.User)
	if er != nil || a == nil {
		a.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrMessageError)
		return
	}

	//roomHandle, ok := room.RoomList[a.RoomId]
	roomHandle, ok := room.GetLockRoomHandle(a.RoomId)
	if ok {
		FunCall(roomHandle.Room, "ChowCard", []reflect.Value{reflect.ValueOf(a), reflect.ValueOf(m)})
	} else {
		a.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrRoomNotExist) // 房间不存在
	}
}

// 碰
func Peng(args ...interface{}) {
	m := &mjgame.Peng{}
	er := proto.Unmarshal(*(args[0].(*[]byte)), m)
	//m := args[0].(*mjgame.Peng)
	a := args[1].(*user.User)
	if er != nil || a == nil {
		a.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrMessageError)
		return
	}

	if m.Cid < 0 {
		a.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrInvalidParam)
		return
	}
	//roomHandle, ok := room.RoomList[a.RoomId]
	roomHandle, ok := room.GetLockRoomHandle(a.RoomId)
	if ok {
		//uesrp := user.RunParam{Room: roomHandle.Room, FunName: "PengCard", Params: []reflect.Value{reflect.ValueOf(a), reflect.ValueOf(m)}}
		//FunCall(roomHandle.Room, "AddMsgList", []reflect.Value{reflect.ValueOf(a), reflect.ValueOf(uesrp)})
		FunCall(roomHandle.Room, "PengCard", []reflect.Value{reflect.ValueOf(a), reflect.ValueOf(int(m.Cid))})
	} else {
		a.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrRoomNotExist) // 房间不存在
	}
}

// 过、取消
func Pass(args ...interface{}) {
	a := args[1].(*user.User)
	if a == nil {
		a.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrMessageError)
		return
	}

	//roomHandle, ok := room.RoomList[a.RoomId]
	roomHandle, ok := room.GetLockRoomHandle(a.RoomId)
	if ok {
		//uesrp := user.RunParam{Room: roomHandle.Room, FunName: "Pass", Params: []reflect.Value{reflect.ValueOf(a)}}
		//FunCall(roomHandle.Room, "AddMsgList", []reflect.Value{reflect.ValueOf(a), reflect.ValueOf(uesrp)})
		FunCall(roomHandle.Room, "Pass", []reflect.Value{reflect.ValueOf(a)})
	} else {
		a.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrRoomNotExist) // 房间不存在
	}

}

// 杠
func Kong(args ...interface{}) {
	m := &mjgame.Kong{}
	er := proto.Unmarshal(*(args[0].(*[]byte)), m)
	//m := args[0].(*mjgame.Kong)
	a := args[1].(*user.User)
	if er != nil || a == nil {
		a.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrMessageError)
		return
	}

	//fmt.Println("kong ", m.Cid, m.Sid, a)
	//roomHandle, ok := room.RoomList[a.RoomId]
	roomHandle, ok := room.GetLockRoomHandle(a.RoomId)
	if ok {
		FunCall(roomHandle.Room, "KongCard", []reflect.Value{reflect.ValueOf(a), reflect.ValueOf(int(m.Cid))})
	} else {
		a.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrRoomNotExist) // 房间不存在
	}
}

// 胡牌
func Win(args ...interface{}) {
	m := &mjgame.Win{}
	er := proto.Unmarshal(*(args[0].(*[]byte)), m)
	//m := args[0].(*mjgame.Win)
	a := args[1].(*user.User)
	if er != nil || a == nil {
		a.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrMessageError)
		return
	}

	//roomHandle, ok := room.RoomList[a.RoomId]
	roomHandle, ok := room.GetLockRoomHandle(a.RoomId)
	if ok {
		FunCall(roomHandle.Room, "WinCard", []reflect.Value{reflect.ValueOf([]*user.User{a}), reflect.ValueOf(int(m.Cid))})
	} else {
		a.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrRoomNotExist) // 房间不存在
	}
}

//三人斗地主相关
func Sddz_Jiaofen(args ...interface{}) {
	m := &mjgame.Sddz_Jiaofen{}
	er := proto.Unmarshal(*(args[0].(*[]byte)), m)
	a := args[1].(*user.User)
	if er != nil || a == nil {
		a.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrMessageError)
		return
	}

	roomHandle, ok := room.GetLockRoomHandle(a.RoomId)
	if ok {
		FunCall(roomHandle.Room, "Jiaofen", []reflect.Value{reflect.ValueOf(a), reflect.ValueOf(int(m.Fen))})
	} else {
		a.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrRoomNotExist) // 房间不存在
	}
}

//拼十 下注
func Nn_Xiazhu(args ...interface{}) {
	m := &mjgame.Nn_Xiazhu{}
	er := proto.Unmarshal(*(args[0].(*[]byte)), m)
	a := args[1].(*user.User)
	if er != nil || a == nil {
		a.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrMessageError)
		return
	}

	roomHandle, ok := room.GetLockRoomHandle(a.RoomId)
	if ok {
		FunCall(roomHandle.Room, "Xiazhu", []reflect.Value{reflect.ValueOf(a), reflect.ValueOf(int(m.Fen))})
	} else {
		a.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrRoomNotExist) // 房间不存在
	}
}

//加倍
func Sddz_Jiabei(args ...interface{}) {
	m := &mjgame.Sddz_Jiabei{}
	er := proto.Unmarshal(*(args[0].(*[]byte)), m)
	a := args[1].(*user.User)
	if er != nil || a == nil {
		a.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrMessageError)
		return
	}

	//roomHandle, ok := room.RoomList[a.RoomId]
	roomHandle, ok := room.GetLockRoomHandle(a.RoomId)
	if ok {
		FunCall(roomHandle.Room, "Jiabei", []reflect.Value{reflect.ValueOf(a), reflect.ValueOf(bool(m.Jiabei))})
	} else {
		a.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrRoomNotExist) // 房间不存在
	}
}

//明牌
func Sddz_Mingpai(args ...interface{}) {
	//todo
	m := &mjgame.Sddz_Mingpai{}
	er := proto.Unmarshal(*(args[0].(*[]byte)), m)
	a := args[1].(*user.User)
	if er != nil || a == nil {
		a.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrMessageError)
		return
	}

	//roomHandle, ok := room.RoomList[a.RoomId]
	roomHandle, ok := room.GetLockRoomHandle(a.RoomId)
	if ok {
		FunCall(roomHandle.Room, "Mingpai", []reflect.Value{reflect.ValueOf(a)})
	} else {
		a.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrRoomNotExist) // 房间不存在
	}
}

//出牌
func Sddz_Chupai(args ...interface{}) {
	//todo

	m := &mjgame.Sddz_Chupai{}
	er := proto.Unmarshal(*(args[0].(*[]byte)), m)
	a := args[1].(*user.User)
	if er != nil || a == nil {
		a.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrMessageError)
		return
	}
	roomHandle, ok := room.GetLockRoomHandle(a.RoomId)
	if ok {
		FunCall(roomHandle.Room, "Chupai", []reflect.Value{reflect.ValueOf(a), reflect.ValueOf(m)})
	} else {
		a.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrRoomNotExist) // 房间不存在
	}
}

//不出
func Sddz_Pass(args ...interface{}) {
	a := args[1].(*user.User)
	if a == nil {
		a.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrMessageError)
		return
	}

	//roomHandle, ok := room.RoomList[a.RoomId]
	roomHandle, ok := room.GetLockRoomHandle(a.RoomId)
	if ok {
		//uesrp := user.RunParam{Room: roomHandle.Room, FunName: "Pass", Params: []reflect.Value{reflect.ValueOf(a)}}
		//FunCall(roomHandle.Room, "AddMsgList", []reflect.Value{reflect.ValueOf(a), reflect.ValueOf(uesrp)})
		FunCall(roomHandle.Room, "Sddz_Pass", []reflect.Value{reflect.ValueOf(a)})
	} else {
		a.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrRoomNotExist) // 房间不存在
	}
}

//报到
func Srddz_Baodao(args ...interface{}) {
	m := &mjgame.Srddz_Baodao{}
	er := proto.Unmarshal(*(args[0].(*[]byte)), m)
	a := args[1].(*user.User)
	if er != nil || a == nil {
		a.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrMessageError)
		return
	}
	roomHandle, ok := room.GetLockRoomHandle(a.RoomId)
	if ok {
		FunCall(roomHandle.Room, "Baodao", []reflect.Value{reflect.ValueOf(a), reflect.ValueOf(m)})
	} else {
		a.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrRoomNotExist) // 房间不存在
	}
}

//直接赢
func Srddz_StrictWin(args ...interface{}) {
	m := &mjgame.Srddz_StrictWin{}
	er := proto.Unmarshal(*(args[0].(*[]byte)), m)
	a := args[1].(*user.User)
	if er != nil || a == nil {
		a.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrMessageError)
		return
	}
	roomHandle, ok := room.GetLockRoomHandle(a.RoomId)
	if ok {
		FunCall(roomHandle.Room, "StrictWin", []reflect.Value{reflect.ValueOf(a), reflect.ValueOf(m)})
	} else {
		a.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrRoomNotExist) // 房间不存在
	}
}

//聊天消息
func Chat(args ...interface{}) {
	m := &mjgame.Chat{}
	er := proto.Unmarshal(*(args[0].(*[]byte)), m)
	//m := args[0].(*mjgame.Chat)
	a := args[1].(*user.User)
	if er != nil || a == nil {
		a.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrMessageError)
		return
	}

	//roomHandle, ok := room.RoomList[a.RoomId]
	roomHandle, ok := room.GetLockRoomHandle(a.RoomId)
	if ok {
		var index int

		if a.GameType == xiangshan.IFCXiangShanType {
			handle := roomHandle.Room.(*xiangshan.RoomXiangshan)
			index = handle.GetSeatIndexById(a.ID)
		} else if a.GameType == xizhou.IFCXiZhouType {
			handle := roomHandle.Room.(*xizhou.RoomXiZhou)
			index = handle.GetSeatIndexById(a.ID)
		} else if a.GameType == ningbo.IFCNingBoType {
			handle := roomHandle.Room.(*ningbo.RoomNingBo)
			index = handle.GetSeatIndexById(a.ID)
		} else if a.GameType == srddz.IFCSrddzType {
			handle := roomHandle.Room.(*srddz.RoomSrddz)
			index = handle.GetSeatIndexById(a.ID)
		} else if a.GameType == pinshi.IFCPinshiType {
			handle := roomHandle.Room.(*pinshi.RoomPinshi)
			index = handle.GetSeatIndexById(a.ID)
		}

		if index < 0 && m.Type == def.Text { //围观用户禁止发言
			a.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrStandUserCanNotChat)
			return
		}
		message := &mjgame.Chat{
			Message: m.Message,
		}
		if a.GameType == xiangshan.IFCXiangShanType {
			handle := roomHandle.Room.(*xiangshan.RoomXiangshan)
			handle.RoomBase.BCMessage(mjgame.MsgID_MSG_ACKBC_Chat, message)
		} else if a.GameType == xizhou.IFCXiZhouType {
			handle := roomHandle.Room.(*xizhou.RoomXiZhou)
			handle.RoomBase.BCMessage(mjgame.MsgID_MSG_ACKBC_Chat, message)
		} else if a.GameType == ningbo.IFCNingBoType {
			handle := roomHandle.Room.(*ningbo.RoomNingBo)
			handle.RoomBase.BCMessage(mjgame.MsgID_MSG_ACKBC_Chat, message)
		} else if a.GameType == srddz.IFCSrddzType {
			handle := roomHandle.Room.(*srddz.RoomSrddz)
			handle.RoomBase.BCMessage(mjgame.MsgID_MSG_ACKBC_Chat, message)
		} else if a.GameType == pinshi.IFCPinshiType {
			handle := roomHandle.Room.(*pinshi.RoomPinshi)
			handle.RoomBase.BCMessage(mjgame.MsgID_MSG_ACKBC_Chat, message)
		}

	} else {
		a.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrRoomNotExist) // 房间不存在
	}
}

//公告
func Notice(args ...interface{}) {
	//	m := args[0].(*mjgame.Notice)
	//	a := args[1].(*user.User)
	m := &mjgame.Notice{}
	er := proto.Unmarshal(*(args[0].(*[]byte)), m)
	a := args[1].(*user.User)
	if er != nil || a == nil {
		a.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrMessageError)
		return
	}

	noticesInfo := GetNoticesById(uint(m.Id))

	notice := &mjgame.AckNotice{
		Notices: noticesInfo,
	}
	a.SendMessage(mjgame.MsgID_MSG_ACK_Notice, notice)
}

//房主解散房间
func RoomownerDisbandRoom(args ...interface{}) {
	m := &mjgame.Roomowner_Disband_Room{}
	//m := args[0].(*mjgame.Disband)
	er := proto.Unmarshal(*(args[0].(*[]byte)), m)
	a := args[1].(*user.User)
	if er != nil || a == nil {
		a.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrMessageError)
		return
	}

	if !CheckUser(a.Sid) {
		a.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrNeedLogin)
		return
	}

	//roomHandle, ok := room.RoomList[a.RoomId]
	roomHandle, ok := room.GetLockRoomHandle(a.RoomId)
	if ok {
		FunCall(roomHandle.Room, "OwnerDisbandRoom", []reflect.Value{reflect.ValueOf(a), reflect.ValueOf(m)})
	} else {
		a.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrRoomNotExist) // 房间不存在
	}
}

// 解散房间
func DisbandRoom(args ...interface{}) {
	m := &mjgame.Disband{}
	er := proto.Unmarshal(*(args[0].(*[]byte)), m)
	//m := args[0].(*mjgame.Disband)
	a := args[1].(*user.User)
	if er != nil || a == nil {
		a.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrMessageError)
		return
	}

	if !CheckUser(a.Sid) {
		a.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrNeedLogin)
		return
	}

	//roomHandle, ok := room.RoomList[a.RoomId]
	roomHandle, ok := room.GetLockRoomHandle(a.RoomId)
	if ok {
		FunCall(roomHandle.Room, "DisbandRoom", []reflect.Value{reflect.ValueOf(a), reflect.ValueOf(m)})
	} else {
		a.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrRoomNotExist) // 房间不存在
	}
}

// 投票
func Vote(args ...interface{}) {
	m := &mjgame.Vote{}
	er := proto.Unmarshal(*(args[0].(*[]byte)), m)
	//m := args[0].(*mjgame.Vote)
	a := args[1].(*user.User)
	if er != nil || a == nil {
		a.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrMessageError)
		return
	}

	if !CheckUser(a.Sid) {
		a.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrNeedLogin)
		return
	}

	//roomHandle, ok := room.RoomList[a.RoomId]
	roomHandle, ok := room.GetLockRoomHandle(a.RoomId)
	if ok {
		FunCall(roomHandle.Room, "Vote", []reflect.Value{reflect.ValueOf(a), reflect.ValueOf(m)})
	} else {
		a.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrRoomNotExist) // 房间不存在
	}
}

//战绩记录
func BattleRecordList(args ...interface{}) {
	a := args[1].(*user.User)
	if a == nil {
		a.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrMessageError)
		return
	}

	resultAll, err := model.GetRoomRecordModel().QueryAll(a.ID)
	if err != nil {
		a.SendMessage(mjgame.MsgID_MSG_ACK_Error, nil)
		return
	}
	result, err := model.GetRoomRecordModel().Query(a.ID, model.Win)
	if err != nil {
		a.SendMessage(mjgame.MsgID_MSG_ACK_Error, nil)
		return
	}

	var data mjgame.BattleRecordResponse

	for k, v := range resultAll {
		o := &mjgame.BattleRecord{
			Type:            int32(k),
			TotalRoundCount: int32(v),
			Ratio:           int32(math.Ceil(float64(result[k]*100) / float64(v))),
		}
		data.List = append(data.List, o)
	}

	a.SendMessage(mjgame.MsgID_MSG_ACK_Battle_Record, &data)
}

//房主记录
func HomeOwnerList(args ...interface{}) {
	m := &mjgame.HomeOwnerRequest{}
	er := proto.Unmarshal(*(args[0].(*[]byte)), m)
	//m := args[0].(*mjgame.HomeOwnerRequest)
	a := args[1].(*user.User)
	if er != nil || a == nil {
		a.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrMessageError)
		return
	}

	roomIds, _ := model.GetRoomRecordModel().GetRoomIdsById(a.ID)

	//rooms, err := model.GetRoomModel().Query(roomIds, m.Type)
	rooms, err := model.GetRoomModel().QueryAll(roomIds)
	if err != nil {
		a.SendMessage(mjgame.MsgID_MSG_ACK_Error, nil)
		return
	}

	var data mjgame.HomeOwnerResponse

	for _, v := range rooms {
		o := &mjgame.HomeOwner{
			Id:       strconv.Itoa(v.User.ID),
			RoomId:   int64(v.ID),
			NickName: v.User.NickName,
			//Icon:       v.User.Icon,
			Rule:       v.Rules,
			Timestamp:  v.CreatedAt.Unix(),
			UniqueCode: v.UniqueCode,
			Rid:        int32(v.ServerRoomID),
			SelfRecord: int32(getUserInRoomRecord(int32(v.ID), a.ID)),
			Roomtype:   int32(v.Type),
		}
		reg := regexp.MustCompile("http://thirdwx.qlogo.cn/mmopen/vi_32/")
		v.User.Icon = reg.ReplaceAllString(v.User.Icon, "")
		o.Icon = v.User.Icon
		data.List = append(data.List, o)

		if len(data.List) > 50 {
			break
		}

	}

	a.SendMessage(mjgame.MsgID_MSG_ACK_Home_Owner, &data)
}

//根据房号和玩家id获取玩家的分数
func getUserInRoomRecord(roomid int32, userid int) int {

	battles, err := model.GetBattleRecordModel().GetBattleRecordByRoomId(roomid)
	if err != nil {
		return 0
	}

	if len(battles) == 0 {
		return 0
	}

	var userScore = make(map[int]int)
	for _, v := range battles {
		for k, score := range v.Result {
			if _, ok := userScore[k]; ok {
				userScore[k] += score
			} else {
				userScore[k] = score
			}
		}
	}

	for k, v := range userScore {
		if k == userid {
			return v
		}
	}

	return 0
}

//房间结算记录
func RoomSummary(args ...interface{}) {
	m := &mjgame.RoomSummaryRequest{}
	er := proto.Unmarshal(*(args[0].(*[]byte)), m)
	//m := args[0].(*mjgame.RoomSummaryRequest)
	a := args[1].(*user.User)
	if er != nil || a == nil {
		a.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrMessageError)
		return
	}

	room, err := model.GetRoomModel().GetRoomByUniqueCode(m.UniqueCode)
	if err != nil {
		a.SendMessage(mjgame.MsgID_MSG_ACK_Room_Summary, nil)
		return
	}

	battles, err := model.GetBattleRecordModel().GetBattleRecordByRoomId(int32(room.ID))
	if err != nil {
		a.SendMessage(mjgame.MsgID_MSG_ACK_Room_Summary, nil)
		return
	}

	if len(battles) == 0 {
		var data mjgame.RoomSummaryResponse
		a.SendMessage(mjgame.MsgID_MSG_ACK_Room_Summary, &data)
		return
	}

	var userScore = make(map[int]int)
	for _, v := range battles {
		for k, score := range v.Result {
			if _, ok := userScore[k]; ok {
				userScore[k] += score
			} else {
				userScore[k] = score
			}
		}
	}

	var ids []int
	for k, _ := range userScore {
		ids = append(ids, k)
	}

	userMap, err := model.GetUserModel().GetUsersByIds(ids)
	if err != nil {
		a.SendMessage(mjgame.MsgID_MSG_ACK_Room_Summary, nil)
		return
	}

	var data mjgame.RoomSummaryResponse
	for _, v := range userMap {
		o := &mjgame.Summary{
			Id:    strconv.Itoa(v.ID),
			Name:  v.NickName,
			Icon:  v.Icon,
			Score: int32(userScore[v.ID]),
		}

		reg := regexp.MustCompile("http://thirdwx.qlogo.cn/mmopen/vi_32/")
		v.Icon = reg.ReplaceAllString(v.Icon, "")
		o.Icon = v.Icon

		if v.ID == room.UserID {
			o.Houseowner = true
		}
		data.List = append(data.List, o)
	}
	data.HomeOwner = &mjgame.HomeOwner{
		Id:         strconv.Itoa(room.User.ID),
		RoomId:     int64(room.ID),
		NickName:   room.User.NickName,
		Icon:       room.User.Icon,
		Rule:       room.Rules,
		Timestamp:  room.CreatedAt.Unix(),
		UniqueCode: room.UniqueCode,
	}
	data.RoomId = int32(room.ServerRoomID)
	data.RoomType = int32(room.Type)

	a.SendMessage(mjgame.MsgID_MSG_ACK_Room_Summary, &data)
}

//获取单个房间的战绩列表
func BattleDetail(args ...interface{}) {

	m := &mjgame.BattleDetail{}
	er := proto.Unmarshal(*(args[0].(*[]byte)), m)
	a := args[1].(*user.User)

	if er != nil || a == nil {
		a.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrMessageError)
		return
	}

	battles, err := model.GetBattleRecordModel().GetBattleRecordByRoomId(m.RoomId)
	if err != nil {
		a.SendMessage(mjgame.MsgID_MSG_ACK_Room_Summary, nil)
		return
	}

	room, errm := model.GetRoomModel().GetRoomById(m.RoomId)
	if errm != nil {
		a.SendMessage(mjgame.MsgID_MSG_ACK_Room_Summary, nil)
		return
	}

	var data mjgame.ACK_BattleDetail
	var ids []int
	var userMap map[int]model.User

	bLen := len(battles)
	idx := 0

	for k, _ := range battles[len(battles)-1].Result {
		ids = append(ids, k)
	}

	userMap, err = model.GetUserModel().GetUsersByIds(ids)
	if err != nil {
		a.SendMessage(mjgame.MsgID_MSG_ACK_BattleDetail, nil)
		return
	}

	for i := 0; i < bLen; i++ {
		battle := battles[i]
		if i == 0 {
			data.RoomId = int64(room.ID)
			data.UniqueCode = room.UniqueCode
			data.RoomType = int32(room.Type)
			data.Rid = int32(room.ServerRoomID)
			data.RoundCount = int32(len(battles))

		}

		o := &mjgame.PlayerBattleDetail{
			PlayBack: "", //battle.PlayBack
			//ReviewCode: battle.ReviewCode,
		}

		for k, score := range battle.Result {
			p := &mjgame.PlayerBattleInfo{
				UserId:   strconv.Itoa(userMap[k].ID),
				NickName: userMap[k].NickName,
				Score:    int32(score),
			}
			reg := regexp.MustCompile("http://thirdwx.qlogo.cn/mmopen/vi_32/")
			p.Icon = userMap[k].Icon
			p.Icon = reg.ReplaceAllString(p.Icon, "")

			if i == 0 {
				//				reg := regexp.MustCompile("http://thirdwx.qlogo.cn/mmopen/vi_32/")
				//				p.Icon = userMap[k].Icon

				//				p.Icon = reg.ReplaceAllString(p.Icon, "")
				idx = k
			} else if k > idx {
				//				reg := regexp.MustCompile("http://thirdwx.qlogo.cn/mmopen/vi_32/")
				//				p.Icon = userMap[k].Icon

				//				p.Icon = reg.ReplaceAllString(p.Icon, "")
				idx = k
			}

			o.List = append(o.List, p)
		}
		data.List = append(data.List, o)
	}

	//fmt.Println(&data)
	a.SendMessage(mjgame.MsgID_MSG_ACK_BattleDetail, &data)
}

//踢人
func Kick(args ...interface{}) {
	m := &mjgame.KickRequest{}
	er := proto.Unmarshal(*(args[0].(*[]byte)), m)
	//m := args[0].(*mjgame.KickRequest)
	a := args[1].(*user.User)
	if er != nil || a == nil {
		a.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrMessageError)
		return
	}

	if !CheckUser(a.Sid) {
		a.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrNeedLogin)
		return
	}

	//roomHandle, ok := room.RoomList[a.RoomId]
	roomHandle, ok := room.GetLockRoomHandle(a.RoomId)
	if ok {
		FunCall(roomHandle.Room, "Kick", []reflect.Value{reflect.ValueOf(a), reflect.ValueOf(m)})
	} else {
		a.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrRoomNotExist) // 房间不存在
	}
}

//messageJson
//handletype=10 用于测试设置下一张牌
func MessageJson(args ...interface{}) {
	return
	m := &mjgame.MessageJson{}
	er := proto.Unmarshal(*(args[0].(*[]byte)), m)
	//m := args[0].(*mjgame.MessageJson)
	a := args[1].(*user.User)
	if er != nil || a == nil {
		a.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrMessageError)
		return
	}

	if m.GetSID() == "1000" {
		getTotalUserOnline(m, a)
	}
	if m.GetSID() == "1010" {
		GetNewRoomId(m, a)
	}
	if m.GetSID() == "1020" {
		GetRoomRecord(m, a)
	}

	if !rb.Debug {
		print("this is not debug")
		return
	}

	if m.GetSID() == "10" {
		testCard.SetNextCard(m, a)
	} else if m.GetSID() == "20" {
		testCard.SetInitCards(m, a)
	}

}

//通知充值
func NotifyRecharge(args ...interface{}) {
	m := &mjgame.NotifyRecharge{}
	er := proto.Unmarshal(*(args[0].(*[]byte)), m)
	//m := args[0].(*mjgame.NotifyRecharge)
	if er != nil {
		return
	}

	GServer.mux.Lock()
	defer GServer.mux.Unlock()

	for _, v := range GServer.UserList {
		if v.ID == int(m.Id) && v.State == def.Online {
			user, err := model.GetUserModel().GetUserById(v.ID)
			if err == nil {
				v.Diamond = user.Diamond
				v.SendMessage(mjgame.MsgID_MSG_ACK_NOTIFY_RECHARGE,
					&mjgame.NotifyRechargeResponse{Id: m.Id, Diamond: int32(v.Diamond)})
				break
			}
		}
	}
}

//送礼物
func Gift(args ...interface{}) {

	m := &mjgame.Gift{}
	er := proto.Unmarshal(*(args[0].(*[]byte)), m)
	a := args[1].(*user.User)
	if er != nil || a == nil {
		a.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrMessageError)
		return
	}

	//costDiamondType := 10

	//roomHandle, ok := room.RoomList[a.RoomId]
	roomHandle, ok := room.GetLockRoomHandle(a.RoomId)
	if !ok {
		a.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrRoomNotExist) // 房间不存在
		return
	}

	//礼物不存在
	if m.Id <= 0 {
		a.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrMessageIsNotGift)
		return
	}
	item, ok := gift.Gifts[m.Id]
	//礼物不存在
	if !ok {
		a.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrMessageIsNotGift)
		return

	}
	//余额不足
	if int32(a.Diamond) < item.Price {
		//		a.SendMessage(mjgame.MsgID_MSG_ACK_Error, err.ErrDiamondNotEnough)
		//		return
	}
	ackMessage := &mjgame.ACK_Gift{
		Id:      m.Id,
		Diamond: int32(a.Diamond),
		Uid:     strconv.Itoa(a.ID),
		TUid:    m.TUid,
	}

	if a.GameType == xiangshan.IFCXiangShanType {
		handle := roomHandle.Room.(*xiangshan.RoomXiangshan)
		if !handle.RoomBase.IsSitDownUser(strconv.Itoa(a.ID)) {
			return
		}
		handle.RoomBase.BCMessage(mjgame.MsgID_MSG_ACK_Gift, ackMessage)
	} else if a.GameType == xizhou.IFCXiZhouType {
		handle := roomHandle.Room.(*xizhou.RoomXiZhou)
		if !handle.RoomBase.IsSitDownUser(strconv.Itoa(a.ID)) {
			return
		}
		handle.RoomBase.BCMessage(mjgame.MsgID_MSG_ACK_Gift, ackMessage)
	} else if a.GameType == ningbo.IFCNingBoType {
		handle := roomHandle.Room.(*ningbo.RoomNingBo)
		if !handle.RoomBase.IsSitDownUser(strconv.Itoa(a.ID)) {
			return
		}
		handle.RoomBase.BCMessage(mjgame.MsgID_MSG_ACK_Gift, ackMessage)
	} else if a.GameType == srddz.IFCSrddzType {
		handle := roomHandle.Room.(*srddz.RoomSrddz)
		if !handle.RoomBase.IsSitDownUser(strconv.Itoa(a.ID)) {
			return
		}
		handle.RoomBase.BCMessage(mjgame.MsgID_MSG_ACK_Gift, ackMessage)
	} else if a.GameType == pinshi.IFCPinshiType {
		handle := roomHandle.Room.(*pinshi.RoomPinshi)
		if !handle.RoomBase.IsSitDownUser(strconv.Itoa(a.ID)) {
			return
		}
		handle.RoomBase.BCMessage(mjgame.MsgID_MSG_ACK_Gift, ackMessage)
	}

	//costDiamondType = int(item.Id) + 20
	//common.AddDiamondLog(a, costDiamondType, -int(item.Price)) //成功发送礼物 扣款 costDiamondType

}
