// 牌类操作

package roombase

import (
	"fmt"

	al "PZ_GameServer/common/util/arrayList"
	"PZ_GameServer/protocol/pb"
	st "PZ_GameServer/server/game/statement"
	"strconv"
)

const (
	W = 0 // 万字牌 0-8   36
	B = 1 // 饼字牌 0-8   36
	T = 2 // 条字牌 0-8   36
	F = 3 // 风字牌 0-3 东 南 西 北 (逆时针)   16
	J = 4 // 箭字牌 0-3 中 发 白   12
	H = 5 // 花牌 0=春 1=夏 2=秋 3=冬 4=梅 5=兰 6=竹 7=菊  8
)

const (
	Hu   = 0
	Kong = 1
	Peng = 2
	Chow = 3
	Put  = 4
	Pass = 5
)

var (
	numbers    = []string{"一", "二", "三", "四", "五", "六", "七", "八", "九"}
	kinds      = []string{"万", "饼", "条"}
	directions = []string{"东", "南", "西", "北"}
	jians      = []string{"中", "发", "白"}
	huaPais    = []string{"春", "夏", "秋", "冬", "梅", "兰", "竹", "菊"}
)

const (
	InitUserCardsNumber     = 13 //玩家初始化手牌数（麻将）
	InitUserPocksNumber     = 17 //玩家初始化手牌数（三人斗地主）
	InitUserlongPocksNumber = 25 // 四人斗地主用的
	InitPinshiPocksNumber   = 5
)

// 牌
type Card struct {
	ID     int    // id
	Type   int    // 类型 w=0 b=1 t=2
	Num    int    // 字数
	TIndex int    // 碰吃家的座位Index
	Status int    // 杠(0=明杠 1=暗杠 2=碰杠) 吃(0:自己的牌 1:吃掉的牌) 状态类型
	MSG    string // 说明
}

type ChowData struct {
	Card1 *Card //另外的牌
	Card2 *Card //另外的牌
	Card3 *Card //另外的牌
}

// 牌的列表
type UserCard struct {
	List    *al.ArrayList // 牌的列表
	Kong    *al.ArrayList // 杠的牌
	Peng    *al.ArrayList // 碰的牌
	Chow    *al.ArrayList // 吃的牌
	Out     *al.ArrayList // 打出的牌
	Hua     *al.ArrayList // 花牌
	Hu      *al.ArrayList // 胡的牌
	OutStep [][]int       //出牌记录
}

// 牌的控制
type CardCtl struct {
	AllCard    []Card // 全部牌
	CurCard    *Card  // 当前牌(用于吃碰杠胡)
	CurMJIndex int    // 当前拿牌的位置
	StartIndex int    // 开始的牌的位置 (可以算出总共拿了多少张)
	EndBlank   int    // 结尾拿掉的牌(杠后从结尾拿掉的牌)
}

// 初始化随机牌
// 牌的列表指针
// 最大张数
func (r *RoomBase) InitRandCard() {

}

//初始化用户的牌
// uid 顺序
// sid 顺序
func (r *RoomBase) InitUserCard() {
	var index int = 0
	for _, seat := range r.Seats {
		for {
			if seat.Cards.List.Count == InitUserCardsNumber {
				break
			}
			if index > len(r.AllCards)-1 {
				fmt.Println("RoomBase InitUserCard index > 143,index = ", index)
				return
			}
			if &r.AllCards[index] == nil {
				fmt.Println("RoomBase InitUserCard r.AllCards[index] == nil,index = ",
					index, "allcardlen = ", len(r.AllCards), "r.AllCards =", r.AllCards)
			}
			//			if &r.AllCards[index] != nil {
			if r.AllCards[index].Type == H {
				seat.Cards.Hua.Add(&r.AllCards[index])
			} else {
				seat.Cards.List.Add(&r.AllCards[index])
			}

			(r.StlCtrl).(st.IStatement).AddTool(st.T_Deal, seat.Index, -1, []int{r.AllCards[index].ID})
			//			} else {
			//				break
			//			}
			index++
		}
	}

	r.StartIndex = index
	r.CurMJIndex = index
}

//初始化用户的牌
// uid 顺序
// sid 顺序
func (r *RoomBase) InitUserPock() {
	var index int = 0

	var cardsLen int = 0

	if r.Type == int32(mjgame.MsgID_GTYPE_SirenDizhu) {
		cardsLen = InitUserlongPocksNumber
	} else if r.Type == int32(mjgame.MsgID_GTYPE_SanDizhu) {
		cardsLen = InitUserPocksNumber
	} else if r.Type == int32(mjgame.MsgID_GTYPE_Pinshi) {
		cardsLen = InitPinshiPocksNumber
	}
	for _, seat := range r.Seats {
		for {

			if seat.Cards.List.Count == cardsLen {
				break
			}
			if index > len(r.AllCards)-1 {
				fmt.Println("RoomBase InitUserCard index > 143,index = ", index)
				return
			}
			if &r.AllCards[index] == nil {
				fmt.Println("RoomBase InitUserCard r.AllCards[index] == nil,index = ",
					index, "allcardlen = ", len(r.AllCards), "r.AllCards =", r.AllCards)
			}
			//			if &r.AllCards[index] != nil {
			seat.Cards.List.Add(&r.AllCards[index])

			(r.StlCtrl).(st.IStatement).AddTool(st.T_Deal, seat.Index, -1, []int{r.AllCards[index].ID})
			//			} else {
			//				break
			//			}
			index++
		}
	}

	r.StartIndex = index //斗地主这两个没有用
	r.CurMJIndex = index //斗地主这两个没有用
}

//从所有牌里面拿出一张不是花的牌，并且放到第一位去
func (r *RoomBase) GetFanCard() {
	for i := 0; i < len(r.AllCards); i++ {
		if r.AllCards[i].Type == H {
			continue
		} else {
			r.FanCard = &r.AllCards[i]
			r.AllCards = append(r.AllCards[:i], r.AllCards[i+1:]...) // append(r.AllCards[:i,(i+1):])
			//将其从数组中删除
			fmt.Println("翻完牌剩余牌组的数量：：" + strconv.Itoa(len(r.AllCards)))
			break
		}
	}
}

//得到牌的int32数组
func (r *RoomBase) GetListArray(list *al.ArrayList) []*mjgame.Card {
	l := list.Count
	arr := make([]*mjgame.Card, l)
	for i := 0; i < l; i++ {
		if *list.Index(i) != nil {
			pcard := (*list.Index(i)).(*Card)
			arr[i] = &mjgame.Card{
				Cid:         int32(pcard.ID),
				TargetIndex: int32(pcard.TIndex),
				Type:        int32(pcard.Status),
			}
		}

	}
	return arr
}

//得到牌的int数组
func (r *RoomBase) GetIntArray(list *al.ArrayList) []int {
	l := list.Count
	arr := make([]int, l)
	for i := 0; i < l; i++ {
		if *list.Index(i) != nil {
			pcard := (*list.Index(i)).(*Card)
			arr[i] = pcard.ID
		}

	}
	return arr
}

//将int数组转化成int32数组
func (r *RoomBase) GetInt32Arr(from []int) []int32 {
	arr := []int32{}
	for _, v := range from {
		arr = append(arr, int32(v))
	}
	return arr

}

//将int32数组转化成int数组
func (r *RoomBase) GetIntArr(from []int32) []int {
	arr := []int{}
	for _, v := range from {
		arr = append(arr, int(v))
	}
	return arr

}

//得到下一张牌
func (r *RoomBase) GetNewCard(bForward bool, uIndex int) *Card {
	var card Card

	if r.CurMJIndex > len(r.AllCards)-1 {
		return nil
	} else {
		if bForward {
			card = r.AllCards[r.CurMJIndex]
			r.CurMJIndex++
		} else {
			card = r.AllCards[(len(r.AllCards) - 1 - r.EndBlank)]
			r.EndBlank++
		}
		return &card
	}
	return nil
}

//设置下一张牌
func (r *RoomBase) SetNextCard(cid int, isBack int) {
	if cid <= 41 && cid >= 0 && len(r.AllCards) > 0 {
		curCard := GetCardById(cid)
		if isBack == 1 {
			r.AllCards[(len(r.AllCards) - 1 - r.EndBlank)] = *curCard
		} else {
			r.AllCards[r.CurMJIndex] = *curCard
		}
	}
}

//得到牌
func GetCardById(id int) *Card {
	c := Card{ID: id}
	c.Type, c.Num = st.GetMjTypeNum(id)
	c.Status = -1
	c.TIndex = -1
	c.MSG = st.GetMjNameForIndex(id)
	return &c
}

//
func (r *RoomBase) GetCardsByType(cardType, outer, inner int) []Card {
	cards := make([]Card, 0) //所有牌
	for i := 0; i < outer; i++ {
		for j := 0; j < inner; j++ {
			c := Card{Type: cardType, Num: j}
			c.ID = st.GetMjIndex(c.Type, c.Num)
			c.MSG = st.GetMjNameForIndex(c.ID)
			cards = append(cards, c)
		}
	}
	return cards
}

// 得到牌2
func (r *RoomBase) GetCard(index int, cid int) *Card {
	t, n := st.GetMjTypeNum(cid)
	pCard := r.GetUserCard(r.Seats[index].Cards.List, t, n)
	return pCard
}

// 得到用户牌的数量
func (r *RoomBase) GetUserCardCount(uIndex int, cid int) int {
	var count int
	length := r.Seats[uIndex].Cards.List.Count
	for i := 0; i < length; i++ {
		if *r.Seats[uIndex].Cards.List.Index(i) != nil {
			card := (*r.Seats[uIndex].Cards.List.Index(i)).(*Card)
			if card.ID == cid {
				count++
			}
		}

	}
	return count
}

//得到牌的数量
func (r *RoomBase) GetCardCount(list *al.ArrayList, t int, n int) int {
	var count int
	for i := 0; i < list.Count; i++ {
		if *list.Index(i) != nil {
			card := (*list.Index(i)).(*Card)
			if card.Type == t && card.Num == n {
				count++
			}
		}

	}
	return count
}

//得到牌
func (r *RoomBase) GetUserCard(list *al.ArrayList, t int, n int) *Card {
	for i := 0; i < list.Count; i++ {
		if *list.Index(i) != nil {
			card := (*list.Index(i)).(*Card)
			if card.Type == t && card.Num == n {
				return card
			}
		}

	}
	return nil
}

//查找用户牌
func (r *RoomBase) FindUserCard(uIndex int, cid int) *Card {
	l := r.Seats[uIndex].Cards.List.Count
	for i := 0; i < l; i++ {
		if *r.Seats[uIndex].Cards.List.Index(i) != nil {
			card := (*r.Seats[uIndex].Cards.List.Index(i)).(*Card)
			if card.ID == cid {
				return card
			}
		}

	}
	return nil
}

//得到牌指针
func (r *RoomBase) GetCardPoint(cid int) *Card {
	for i := 0; i < len(r.AllCards); i++ {
		if r.AllCards[i].ID == cid {
			return &r.AllCards[i]
		}
	}
	return nil
}

// 移动r牌
func (r *RoomBase) MoveToList(listForm *al.ArrayList, cards []*Card, listTo *al.ArrayList) {

	for _, v := range cards {
		for i := 0; i < listForm.Count; i++ {
			if *listForm.Index(i) != nil {
				card := (*listForm.Index(i)).(*Card)
				if card.ID == v.ID {
					//r.StateMutex.Lock()
					listTo.Add(v)
					listForm.RemoveAt(i)
					//r.StateMutex.Unlock()
					break
				}
			}

		}
	}

}

// 添加可以操作的用户
func (r *RoomBase) AddToolUser(uIndex int, iWin int, iKong int, iPeng int, iChow int, iPut int, iPass int) {
	r.WaitOptTool.AddCanTool(uIndex, iWin, iKong, iPeng, iChow, iPut, iPass)
}

// 得到卡片的数组
func GetCardArray(list *al.ArrayList) []int {
	mjs := make([]int, 27)
	for i := 0; i < list.Count; i++ {
		if *list.Index(i) != nil {
			c := (*list.Index(i)).(*Card)
			mjs[c.ID]++ // 类型 w=0 b=1 t=2
		}

	}
	return mjs
}

// 胡牌判断
func IsWin(pM []int) int {
	//递归退出条件
	iNum := 0
	//获得牌的张数
	for i := 0; i < 27; i++ {
		iNum += pM[i]
	}

	if iNum == 0 {
		return 1
	}

	//--- 判断7对
	duiCount := 0
	for i := 0; i < 27; i++ {

		if pM[i] == 0 {
			continue
		}
		if pM[i]%2 != 0 {
			duiCount = 0
			break
		} else {
			duiCount++
		}
	}

	if duiCount == 7 {
		return 1
	}

	//找到有牌的位置
	iIndex := 0
	for iIndex = 0; iIndex < 27; iIndex++ {
		if pM[iIndex] > 0 {
			break
		}
	}

	//刻子判断
	if pM[iIndex] >= 3 {

		pM[iIndex] -= 3 //减去刻子

		//若余下的牌能够胡牌,返回1
		if 1 == IsWin(pM) {
			return 1
		} else { //若余下的牌不能胡牌,还原牌
			pM[iIndex] += 3
		}
	}

	//将对判断
	if pM[iIndex] >= 2 {
		//减去将对,并设置标记
		pM[iIndex] -= 2

		//若余下的牌能够胡牌,返回1
		if 1 == IsWin(pM) {
			return 1
		} else { //若余下的牌不能胡牌,还原牌
			pM[iIndex] += 2
		}
	}

	//顺子判断
	if iIndex != 25 && iIndex != 26 {
		if pM[iIndex] >= 1 && pM[iIndex+1] >= 1 && pM[iIndex+2] >= 1 && iIndex != 7 && iIndex != 8 && iIndex != 16 && iIndex != 17 {
			//减去顺子
			pM[iIndex] -= 1
			pM[iIndex+1] -= 1
			pM[iIndex+2] -= 1

			//若余下的牌能够胡牌,返回1
			if 1 == IsWin(pM) {
				return 1
			} else { //若余下的牌不能胡牌,还原牌
				pM[iIndex] += 1
				pM[iIndex+1] += 1
				pM[iIndex+2] += 1
			}
		}
	}

	return 0
}

//
func (u *UserCard) GetCardByCardId(id int) *Card {
	for _, v := range u.List.Array() {
		card := (*v).(*Card)
		if card != nil && card.ID == id {
			return card
		}
	}
	return nil
}
