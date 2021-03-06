package common

import (
	"PZ_GameServer/model"
	"PZ_GameServer/protocol/pb"
	"PZ_GameServer/server/user"
	"fmt"
	"strconv"
)

//记录钻石操作
func AddDiamondLog(user *user.User, typ int, diamond int) {
	if user.IsRobot == 1 {
		return
	}

	fmt.Println("AddDiamondLog :::", diamond)

	//临时使用，先查询数据库
	if user == nil {
		fmt.Println("错误: user为nil ")
		return
	}
	dbUser, err := model.GetUserModel().GetUserById(user.ID)
	if err == nil {
		user.Diamond = dbUser.Diamond
	} else {
		fmt.Println("错误: AddDiamondLog 用户数据为空 ", user.ID)
		return
	}
	user.Diamond += diamond

	logDiamond := &model.LogDiamond{
		UserID:  user.ID,
		Type:    typ,
		Diamond: diamond,
	}

	model.GetUserModel().Save(user.User)
	model.GetLogDiamondModel().Create(logDiamond)

	userInfo := &mjgame.ACK_User_Info{
		Uid:     strconv.Itoa(user.ID),
		Diamond: int32(user.Diamond),
		Icon:    user.Icon,
		Coin:    int32(user.Coin),
	}
	user.SendMessage(mjgame.MsgID_MSG_UPDATE_USERINFO, userInfo)
	fmt.Println("fasongle ... ")
}
