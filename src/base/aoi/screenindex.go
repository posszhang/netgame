package aoi

import (
	"math"
)

//方向
const (
	DIR_T     = 0
	DIR_RT    = 1
	DIR_R     = 2
	DIR_RB    = 3
	DIR_B     = 4
	DIR_LB    = 5
	DIR_L     = 6
	DIR_LT    = 7
	DIR_OWNER = 8
)

//向量
var DIR_ADJUST = [...][2]int{{0, -1}, {1, -1}, {1, 0}, {1, 1}, {0, 1}, {-1, 1}, {-1, 0}, {-1, -1}, {0, 0}}

const (
	SCREEN_WIDTH  = 10
	SCREEN_HEIGHT = 7
)

//场景元素
type ISceneEntry interface {
	GetID() uint32
	GetPos() *Pos
}

type Pos struct {
	X uint32
	Y uint32
}

func NewPos(x uint32, y uint32) *Pos {
	pos := &Pos{
		X: x,
		Y: y,
	}

	return pos
}

func GetDirect(pos1 *Pos, pos2 *Pos) uint32 {

	if pos1.X == pos2.X && pos1.Y > pos2.Y {
		return DIR_T
	} else if pos1.X < pos2.X && pos1.Y > pos2.Y {
		return DIR_RT
	} else if pos1.X < pos2.X && pos1.Y == pos2.Y {
		return DIR_R
	} else if pos1.X < pos2.X && pos1.Y < pos2.Y {
		return DIR_RB
	} else if pos1.X == pos2.X && pos1.Y < pos2.Y {
		return DIR_B
	} else if pos1.X > pos2.X && pos1.Y < pos2.Y {
		return DIR_LB
	} else if pos1.X > pos2.X && pos1.Y == pos2.Y {
		return DIR_L
	} else if pos1.X > pos2.X && pos1.Y > pos2.Y {
		return DIR_LT
	}

	return DIR_T
}

func GetDistance(pos1 *Pos, pos2 *Pos) uint32 {
	maxX := pos1.X
	if pos2.X > pos1.X {
		maxX = pos2.X
	}

	minX := pos1.X
	if pos2.X < pos1.X {
		minX = pos2.X
	}

	maxY := pos1.Y
	if pos2.Y > pos1.Y {
		maxY = pos2.Y
	}

	minY := pos1.Y
	if pos2.Y < pos1.Y {
		minY = pos2.Y
	}

	ret := math.Sqrt((float64(maxX)-float64(minX))*(float64(maxX)-float64(minX)) + (float64(maxY)-float64(minY))*(float64(maxY)-float64(minY)))
	return uint32(ret)
}

type zInterface struct {
	npcMap     map[uint64]*Npc
	userMap    map[uint64]*SceneUser
	xingjunMap map[uint64]*XingJun
}

func (data *zInterface) AddNpc(npc *Npc) bool {

	if data.npcMap == nil {
		data.npcMap = make(map[uint64]*Npc)
	}

	data.npcMap[npc.GetID()] = npc

	return true
}

func (data *zInterface) DelNpc(npc *Npc) bool {
	if data.npcMap == nil {
		return true
	}

	delete(data.npcMap, npc.GetID())

	return true
}

func (data *zInterface) GetNpcList() map[uint64]*Npc {
	return data.npcMap
}

func (data *zInterface) AddUser(user *SceneUser) bool {
	if data.userMap == nil {
		data.userMap = make(map[uint64]*SceneUser)
	}

	data.userMap[user.GetID()] = user

	return true
}

func (data *zInterface) DelUser(user *SceneUser) bool {
	if data.userMap == nil {
		return true
	}

	delete(data.userMap, user.GetID())

	return true
}

func (data *zInterface) GetSceneUserList() map[uint64]*SceneUser {
	return data.userMap
}

func (data *zInterface) AddXingJun(xingjun *XingJun) bool {

	if data.xingjunMap == nil {
		data.xingjunMap = make(map[uint64]*XingJun)
	}

	data.xingjunMap[xingjun.GetID()] = xingjun

	return true
}

func (data *zInterface) DelXingJun(xingjun *XingJun) bool {
	if data.xingjunMap == nil {
		return true
	}

	delete(data.xingjunMap, xingjun.GetID())

	return true
}

func (data *zInterface) GetXingJunList() map[uint64]*XingJun {
	return data.xingjunMap
}

type zCell struct {
	zInterface

	//坐标
	pos *Pos
	//所在屏编号
	screen uint32
}

func NewCell(x uint32, y uint32, index uint32) *zCell {
	cell := &zCell{
		pos:    NewPos(x, y),
		screen: index,
	}

	return cell
}

func (cell *zCell) GetScreen() uint32 {
	return cell.screen
}

type zScreen struct {
	zInterface

	//屏编号
	index uint32
	//屏中的格子
	cellList []*zCell
}

func NewScreen(i uint32) *zScreen {

	screen := &zScreen{
		index:    i,
		cellList: make([]*zCell, 0, SCREEN_WIDTH*SCREEN_HEIGHT),
	}

	return screen
}

func (screen *zScreen) GetPosI() uint32 {
	return screen.index
}

func (screen *zScreen) AddCell(cell *zCell) {
	screen.cellList = append(screen.cellList, cell)
}

type zNineScreen struct {
	nineScreen map[uint32][]*zScreen
}

func NewNineScreen() *zNineScreen {
	screen := &zNineScreen{
		nineScreen: make(map[uint32][]*zScreen),
	}

	return screen
}

func (screen *zNineScreen) Set(index uint32, screenList []*zScreen) {

	screen.nineScreen[index] = screenList
}

func (screen *zNineScreen) Get(index uint32) []*zScreen {
	if _, ok := screen.nineScreen[index]; !ok {
		return nil
	}

	return screen.nineScreen[index]
}

func (screen *zNineScreen) Add(index uint32, s *zScreen) {

	if _, ok := screen.nineScreen[index]; !ok {
		screen.nineScreen[index] = make([]*zScreen, 0, 0)
	}

	screen.nineScreen[index] = append(screen.nineScreen[index], s)
}

type zScreenIndex struct {
	zInterface

	//格子宽
	width uint32
	//格子高
	height uint32
	//屏数量
	screen uint32
	//屏x
	screenx uint32
	//屏y
	screeny uint32
	//格子
	cellList []*zCell
	//屏
	screenList []*zScreen

	//9屏
	nineScreen *zNineScreen
	//正向
	directScreen [8]*zNineScreen
	//反向
	reverseScreen [8]*zNineScreen
}

func (screenindex *zScreenIndex) Pos2PosI(pos *Pos) uint32 {
	//return ((screenindex.width+SCREEN_WIDTH-1)/SCREEN_WIDTH)*(pos.Y/SCREEN_HEIGHT) + (pos.Y / SCREEN_WIDTH)

	return (pos.Y/SCREEN_HEIGHT)*screenindex.screenx + pos.X/SCREEN_WIDTH
}

func (screenindex *zScreenIndex) InitScreen(width uint32, height uint32) {

	//格子宽高
	screenindex.width = width
	screenindex.height = height
	screenindex.cellList = make([]*zCell, width*height)

	//屏宽高
	screenindex.screenx = (width + SCREEN_WIDTH - 1) / SCREEN_WIDTH
	screenindex.screeny = (height + SCREEN_HEIGHT - 1) / SCREEN_HEIGHT
	screenindex.screen = screenindex.screenx * screenindex.screeny
	screenindex.screenList = make([]*zScreen, screenindex.screen, screenindex.screen)

	screenindex.nineScreen = NewNineScreen()
	for i := 0; i != len(screenindex.directScreen); i++ {
		screenindex.directScreen[i] = NewNineScreen()
	}
	for i := 0; i != len(screenindex.reverseScreen); i++ {
		screenindex.reverseScreen[i] = NewNineScreen()
	}
	for i := 0; i != int(screenindex.screen); i++ {
		screenindex.screenList[i] = NewScreen(uint32(i))
	}

	for i := uint32(0); i != screenindex.screen; i++ {

		sx := uint32(i % screenindex.screenx)
		sy := uint32(i / screenindex.screenx)

		ff := 0

		//初始化本屏数据
		for x := sx * SCREEN_WIDTH; x < (sx+1)*SCREEN_WIDTH; x++ {

			for y := sy * SCREEN_HEIGHT; y < (sy+1)*SCREEN_HEIGHT; y++ {

				if x >= width || y >= height {
					continue
				}

				index := y*width + x
				cell := NewCell(x, y, i)

				screenindex.cellList[index] = cell
				screenindex.screenList[i].AddCell(cell)
				ff++
			}
		}

		//计算九屏
		for j := 0; j != 9; j++ {

			x := int(sx) + DIR_ADJUST[j][0]
			y := int(sy) + DIR_ADJUST[j][1]
			index := y*int(screenindex.screenx) + x

			if x < 0 || x >= int(screenindex.screenx) || y < 0 || y >= int(screenindex.screeny) {
				continue
			}
			screenindex.nineScreen.Add(i, screenindex.screenList[index])
		}

		//计算正向变化五屏或者三屏
		for dir := 0; dir < 8; dir++ {

			start := 0
			end := 0

			pv := make([]*zScreen, 0, 0)

			if 1 == dir%2 {
				start = 6
				end = 10
			} else {
				start = 7
				end = 9
			}

			for j := start; j <= end; j++ {

				x := int(sx) + DIR_ADJUST[(j+dir)%8][0]
				y := int(sy) + DIR_ADJUST[(j+dir)%8][1]
				index := y*int(screenindex.screenx) + x

				if x < 0 || x >= int(screenindex.screenx) || y < 0 || y >= int(screenindex.screeny) {
					continue
				}

				pv = append(pv, screenindex.screenList[index])
			}
			screenindex.directScreen[dir].Set(i, pv)
		}

		// 计算反向变化五屏或者三屏
		for dir := 0; dir < 8; dir++ {

			start := 0
			end := 0

			pv := make([]*zScreen, 0, 0)

			if 1 == dir%2 {
				start = 2
				end = 6
			} else {
				start = 3
				end = 5
			}

			for j := start; j <= end; j++ {

				x := int(sx) + DIR_ADJUST[(j+dir)%8][0]
				y := int(sy) + DIR_ADJUST[(j+dir)%8][1]
				index := y*int(screenindex.screenx) + x

				if x < 0 || x >= int(screenindex.screenx) || y < 0 || y >= int(screenindex.screeny) {
					continue
				}

				pv = append(pv, screenindex.screenList[index])
			}
			screenindex.reverseScreen[dir].Set(i, pv)
		}

	}
}

func (screen *zScreenIndex) GetCell(pos *Pos) *zCell {

	if pos.X >= screen.width || pos.Y >= screen.height {
		return nil
	}

	index := pos.Y*screen.width + pos.X
	return screen.cellList[index]
}

func (screen *zScreenIndex) GetScreen(i uint32) *zScreen {
	if i >= screen.screen {
		return nil
	}

	return screen.screenList[i]
}

func (screen *zScreenIndex) GetNineScreen(i uint32) []*zScreen {
	return screen.nineScreen.Get(i)
}

func (screen *zScreenIndex) AddNpc(npc *Npc) bool {

	cell := screen.GetCell(npc.GetPos())
	if cell == nil {
		log.Println("坐标错误", npc.GetPos(), screen.width, screen.height)
		return false
	}

	s := screen.GetScreen(cell.GetScreen())
	if s == nil {
		log.Println("未找到对应的屏", npc.GetPos())
		return false
	}

	if !cell.AddNpc(npc) {
		log.Println("格子添加失败", npc.GetPos())
		return false
	}

	if !s.AddNpc(npc) {
		cell.DelNpc(npc)
		log.Println("屏添加失败", npc.GetPos())
		return false
	}

	if !screen.zInterface.AddNpc(npc) {
		cell.DelNpc(npc)
		s.DelNpc(npc)
		log.Println("场景添加失败", npc.GetPos())
		return false
	}

	return true
}

func (screen *zScreenIndex) DelNpc(npc *Npc) bool {

	cell := screen.GetCell(npc.GetPos())
	if cell == nil {
		log.Println("坐标错误", npc.GetPos())
		return false
	}

	s := screen.GetScreen(cell.GetScreen())
	if s == nil {
		log.Println("未找到对应的屏", npc.GetPos())
		return false
	}

	if !cell.DelNpc(npc) {
		log.Println("格子删除失败", npc.GetPos())
		return false
	}

	if !s.DelNpc(npc) {
		cell.DelNpc(npc)
		log.Println("屏删除失败", npc.GetPos())
		return false
	}

	if !screen.zInterface.DelNpc(npc) {
		cell.DelNpc(npc)
		s.DelNpc(npc)
		log.Println("场景删除失败", npc.GetPos())
		return false
	}

	return true
}

func (screen *zScreenIndex) RefreshNpc(npc *Npc, oldPos *Pos, newPos *Pos) bool {

	oldcell := screen.GetCell(oldPos)
	if oldcell == nil {
		log.Println("找不到格子", oldPos)
		return false
	}

	newcell := screen.GetCell(newPos)
	if newcell == nil {
		log.Println("找不到格子", oldPos)
		return false
	}

	if !oldcell.DelNpc(npc) {
		log.Println("删除npc失败", oldPos)
		return false
	}

	if !newcell.AddNpc(npc) {
		log.Println("添加npc失败", newPos)
		oldcell.AddNpc(npc)
		return false
	}

	newi := newcell.GetScreen()
	oldi := oldcell.GetScreen()

	//跨屏
	if newi != oldi {
		oldscreen := screen.GetScreen(oldi)
		if oldscreen == nil {
			log.Println("找不到对应的屏", oldPos)
			oldcell.AddNpc(npc)
			newcell.DelNpc(npc)
			return false
		}

		newscreen := screen.GetScreen(newi)
		if newscreen == nil {
			log.Println("找不到对应的屏", newPos)
			oldcell.AddNpc(npc)
			newcell.DelNpc(npc)
			return false
		}

		if !oldscreen.DelNpc(npc) {
			log.Println("屏删除NPC失败", oldPos)
			oldcell.AddNpc(npc)
			newcell.DelNpc(npc)
			return false
		}

		if !newscreen.AddNpc(npc) {
			log.Println("屏添加NPC失败", newPos)
			oldcell.AddNpc(npc)
			newcell.DelNpc(npc)
			oldscreen.AddNpc(npc)
			return false
		}
	}

	return true
}

func (screen *zScreenIndex) AddUser(user *SceneUser) bool {

	cell := screen.GetCell(user.GetPos())
	if cell == nil {
		log.Println("坐标错误", user.GetPos(), screen.width, screen.height)
		return false
	}

	s := screen.GetScreen(cell.GetScreen())
	if s == nil {
		log.Println("未找到对应的屏", user.GetPos())
		return false
	}

	if !cell.AddUser(user) {
		log.Println("格子添加失败", user.GetPos())
		return false
	}

	log.Println(s.GetPosI(), "添加用户")

	if !s.AddUser(user) {
		cell.DelUser(user)
		log.Println("屏添加失败", user.GetPos())
		return false
	}

	if !screen.zInterface.AddUser(user) {
		cell.DelUser(user)
		s.DelUser(user)
		log.Println("场景添加失败", user.GetPos())
		return false
	}

	return true
}

func (screen *zScreenIndex) DelUser(user *SceneUser) bool {

	cell := screen.GetCell(user.GetPos())
	if cell == nil {
		log.Println("坐标错误", user.GetPos())
		return false
	}

	s := screen.GetScreen(cell.GetScreen())
	if s == nil {
		log.Println("未找到对应的屏", user.GetPos())
		return false
	}

	if !cell.DelUser(user) {
		log.Println("格子删除失败", user.GetPos())
		return false
	}

	if !s.DelUser(user) {
		cell.DelUser(user)
		log.Println("屏删除失败", user.GetPos())
		return false
	}

	if !screen.zInterface.DelUser(user) {
		cell.DelUser(user)
		s.DelUser(user)
		log.Println("场景删除失败", user.GetPos())
		return false
	}

	return true
}

func (screen *zScreenIndex) RefreshUser(user *SceneUser, oldPos *Pos, newPos *Pos) bool {

	oldcell := screen.GetCell(oldPos)
	if oldcell == nil {
		log.Println("找不到格子", oldPos)
		return false
	}

	newcell := screen.GetCell(newPos)
	if newcell == nil {
		log.Println("找不到格子", oldPos)
		return false
	}

	if !oldcell.DelUser(user) {
		log.Println("删除user失败", oldPos)
		return false
	}

	if !newcell.AddUser(user) {
		log.Println("添加user失败", newPos)
		oldcell.AddUser(user)
		return false
	}

	newi := newcell.GetScreen()
	oldi := oldcell.GetScreen()

	//跨屏
	if newi != oldi {
		oldscreen := screen.GetScreen(oldi)
		if oldscreen == nil {
			log.Println("找不到对应的屏", oldPos)
			oldcell.AddUser(user)
			newcell.DelUser(user)
			return false
		}

		newscreen := screen.GetScreen(newi)
		if newscreen == nil {
			log.Println("找不到对应的屏", newPos)
			oldcell.AddUser(user)
			newcell.DelUser(user)
			return false
		}

		if !oldscreen.DelUser(user) {
			log.Println("屏删除user失败", oldPos)
			oldcell.AddUser(user)
			newcell.DelUser(user)
			return false
		}

		if !newscreen.AddUser(user) {
			log.Println("屏添加User失败", newPos)
			oldcell.AddUser(user)
			newcell.DelUser(user)
			oldscreen.AddUser(user)
			return false
		}
	}

	return true
}

func (screen *zScreenIndex) AddXingJun(xingjun *XingJun) bool {

	cell := screen.GetCell(xingjun.GetPos())
	if cell == nil {
		log.Println("坐标错误", xingjun.GetPos(), screen.width, screen.height)
		return false
	}

	s := screen.GetScreen(cell.GetScreen())
	if s == nil {
		log.Println("未找到对应的屏", xingjun.GetPos())
		return false
	}

	if !cell.AddXingJun(xingjun) {
		log.Println("格子添加失败", xingjun.GetPos())
		return false
	}

	if !s.AddXingJun(xingjun) {
		cell.DelXingJun(xingjun)
		log.Println("屏添加失败", xingjun.GetPos())
		return false
	}

	if !screen.zInterface.AddXingJun(xingjun) {
		cell.DelXingJun(xingjun)
		s.DelXingJun(xingjun)
		log.Println("场景添加失败", xingjun.GetPos())
		return false
	}

	return true
}

func (screen *zScreenIndex) DelXingJun(xingjun *XingJun) bool {

	cell := screen.GetCell(xingjun.GetPos())
	if cell == nil {
		log.Println("坐标错误", xingjun.GetPos())
		return false
	}

	s := screen.GetScreen(cell.GetScreen())
	if s == nil {
		log.Println("未找到对应的屏", xingjun.GetPos())
		return false
	}

	if !cell.DelXingJun(xingjun) {
		log.Println("格子删除失败", xingjun.GetPos())
		return false
	}

	if !s.DelXingJun(xingjun) {
		cell.DelXingJun(xingjun)
		log.Println("屏删除失败", xingjun.GetPos())
		return false
	}

	if !screen.zInterface.DelXingJun(xingjun) {
		cell.DelXingJun(xingjun)
		s.DelXingJun(xingjun)
		log.Println("场景删除失败", xingjun.GetPos())
		return false
	}

	return true
}

func (screen *zScreenIndex) RefreshXingJun(xingjun *XingJun, oldPos *Pos, newPos *Pos) bool {

	oldcell := screen.GetCell(oldPos)
	if oldcell == nil {
		log.Println("找不到格子", oldPos)
		return false
	}

	newcell := screen.GetCell(newPos)
	if newcell == nil {
		log.Println("找不到格子", oldPos)
		return false
	}

	if !oldcell.DelXingJun(xingjun) {
		log.Println("删除xingjun失败", oldPos)
		return false
	}

	if !newcell.AddXingJun(xingjun) {
		log.Println("添加xingjun失败", newPos)
		oldcell.AddXingJun(xingjun)
		return false
	}

	newi := newcell.GetScreen()
	oldi := oldcell.GetScreen()

	//跨屏
	if newi != oldi {
		oldscreen := screen.GetScreen(oldi)
		if oldscreen == nil {
			log.Println("找不到对应的屏", oldPos)
			oldcell.AddXingJun(xingjun)
			newcell.DelXingJun(xingjun)
			return false
		}

		newscreen := screen.GetScreen(newi)
		if newscreen == nil {
			log.Println("找不到对应的屏", newPos)
			oldcell.AddXingJun(xingjun)
			newcell.DelXingJun(xingjun)
			return false
		}

		if !oldscreen.DelXingJun(xingjun) {
			log.Println("屏删除xingjun失败", oldPos)
			oldcell.AddXingJun(xingjun)
			newcell.DelXingJun(xingjun)
			return false
		}

		if !newscreen.AddXingJun(xingjun) {
			log.Println("屏添加XingJun失败", newPos)
			oldcell.AddXingJun(xingjun)
			newcell.DelXingJun(xingjun)
			oldscreen.AddXingJun(xingjun)
			return false
		}
	}

	return true
}
