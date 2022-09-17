package models

import (
	"YYS_helper/consts"
	"github.com/go-vgo/robotgo"
	"image"
	"log"
	"time"
)

type HeroPieceModel struct {
	Base            BaseModel
	FightSignalChan chan int
	SumClickTimes   int
}

var heroPiece = HeroPieceModel{
	Base: BaseModel{
		Status: consts.COMMON_STATE_IDLE,
		Tick:   time.NewTicker(time.Second),
	},
	FightSignalChan: make(chan int),
	SumClickTimes:   0,
}

func init() {
	modelPool[consts.MODEL_HERO_PIECE] = &heroPiece
}

func (m *HeroPieceModel) Start() {
	heroPieceSymbol, _, _ = robotgo.DecodeImg(consts.HERO_PIECE_SYMBOL)
	matchSymbol, _, _ = robotgo.DecodeImg(consts.HERO_PIECE_MATCH_SYMBOL)
	teamSymbol, _, _ = robotgo.DecodeImg(consts.CHALLENGE_UNABLE_SYMBOL)
	quitSymbol, _, _ = robotgo.DecodeImg(consts.QUIT_SYMBOL)
	fightSymbol, _, _ = robotgo.DecodeImg(consts.FIGHT_SYMBOL)
	sumSymbol, _, _ = robotgo.DecodeImg(consts.SUM_SYMBOL)
	checkSymbol, _, _ = robotgo.DecodeImg(consts.CHECK_SYMBOL)
	snakeSymbol, _, _ = robotgo.DecodeImg(consts.SNAKE_SYMBOL)
	m.Base.Status = consts.COMMON_STATE_IDLE
	go m.checkTeamStatusChange()
	m.Loop()
}

func (m *HeroPieceModel) Loop() {
	var target image.Image
	for {
		m.BeforeLoop()

		switch m.Base.Status {
		case consts.COMMON_STATE_IDLE:
			target = heroPieceSymbol
			m.Base.NextStatus = consts.HP_STATE_MATCH
		case consts.HP_STATE_MATCH:
			target = matchSymbol
			m.Base.NextStatus = consts.HP_STATE_TEAM_SEARCH
		case consts.HP_STATE_TEAM_SEARCH:
			target = teamSymbol
		case consts.HP_STATE_CAPTAIN:
			target = quitSymbol
			m.Base.NextStatus = consts.HP_STATE_WAITING_CHECK
		case consts.HP_STATE_FIGHT:
			target = fightSymbol
			m.Base.NextStatus = consts.HP_STATE_FIGHT
		case consts.HP_STATE_SUM:
			target = sumSymbol
			m.Base.NextStatus = consts.COMMON_STATE_SUM
		case consts.HP_STATE_WAITING_CHECK:
			target = checkSymbol
			m.Base.NextStatus = consts.HP_STATE_MATCH
		default:
			continue
		}
		log.Println("当前处于", m.Base.Status, ", 检测中...")
		screen := robotgo.CaptureImg()
		result := find(target, screen)
		if valid(result) {
			if m.Base.Status == consts.HP_STATE_TEAM_SEARCH {
				m.Base.Status = consts.HP_STATE_CAPTAIN
				break
			}
			onFindSymbol(result, &m.Base.Status, m.Base.NextStatus)
		}
		continue
	}
}

func (m *HeroPieceModel) BeforeLoop() {
	select {
	case signal := <-m.FightSignalChan:
		if signal == 1 {
			log.Println("收到消息，转为战斗状态")
			m.Base.Status = consts.HP_STATE_FIGHT
			m.Base.NextStatus = consts.HP_STATE_SUM
		} else if signal == 2 {
			log.Println("收到消息，转为战斗状态")
			m.Base.Status = consts.COMMON_STATE_SUM
			m.Base.NextStatus = consts.COMMON_STATE_SUM
		} else if signal == 3 {
			log.Println("收到消息，转为空闲状态")
			m.Base.Status = consts.COMMON_STATE_IDLE
			m.Base.NextStatus = consts.HP_STATE_MATCH
		}
	case <-m.Base.Tick.C:
	}
}

// 检查队伍状态变更
func (m *HeroPieceModel) checkTeamStatusChange() {
	for {
		if m.Base.Status == consts.HP_STATE_TEAM_SEARCH { // 判断是否已进入战斗
			screen := robotgo.CaptureImg()
			result := find(fightSymbol, screen)
			if valid(result) {
				m.FightSignalChan <- 1
				log.Println("检测到已进入战场")
			}
		} else if m.Base.Status == consts.HP_STATE_FIGHT { // 判断是否已经点击了准备
			screen := robotgo.CaptureImg()
			result := find(readySymbol, screen)
			if !valid(result) {
				m.FightSignalChan <- 2
				log.Println("检测到已完成准备")
			}
		} else if m.Base.Status == consts.HP_STATE_SUM { // 判断是否已经结算完成
			screen := robotgo.CaptureImg()
			result := find(heroPieceSymbol, screen)
			if valid(result) {
				m.FightSignalChan <- 3
				log.Println("检测到已回到主页")
			}
		}
		time.Sleep(time.Second)
	}
}
