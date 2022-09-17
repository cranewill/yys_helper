package models

import (
	"YYS_helper/consts"
	"github.com/go-vgo/robotgo"
	"image"
	"log"
	"time"
)

type SoulModel struct {
	Base            BaseModel
	FightSignalChan chan int
}

var soul = SoulModel{
	Base: BaseModel{
		Status: consts.SOUL_STATE_TEAM,
		Tick:   time.NewTicker(time.Second),
	},
	FightSignalChan: make(chan int),
}

func init() {
	modelPool[consts.MODEL_SOUL] = &soul
}

func (m *SoulModel) Start() {
	teamSymbol, _, _ = robotgo.DecodeImg(consts.CHALLENGE_ABLE_SYMBOL)
	fightSymbol, _, _ = robotgo.DecodeImg(consts.FIGHT_SYMBOL)
	readySymbol, _, _ = robotgo.DecodeImg(consts.READY_SYMBOL)
	sumSymbol, _, _ = robotgo.DecodeImg(consts.SUM_SYMBOL)
	checkSymbol, _, _ = robotgo.DecodeImg(consts.CHECK_SYMBOL)
	m.Base.Status = consts.SOUL_STATE_TEAM
	go m.checkTeamStatusChange()
	m.Loop()
}

func (m *SoulModel) Loop() {
	var target image.Image
	for {
		m.BeforeLoop()

		switch m.Base.Status {
		case consts.SOUL_STATE_TEAM:
			target = teamSymbol
			m.Base.NextStatus = consts.SOUL_STATE_TEAM
		case consts.SOUL_STATE_FIGHT:
			m.Base.Status = consts.SOUL_STATE_SUM
			continue
		case consts.SOUL_STATE_SUM:
			target = sumSymbol
			m.Base.NextStatus = consts.SOUL_STATE_SUM
		case consts.SOUL_STATE_INVITE_MEMBER:
			target = checkSymbol
			m.Base.NextStatus = consts.SOUL_STATE_TEAM
		default:
			continue
		}
		log.Println("当前处于", m.Base.Status, ", 检测中...")
		screen := robotgo.CaptureImg()
		result := find(target, screen)
		if valid(result) {
			onFindSymbol(result, &m.Base.Status, m.Base.NextStatus)
		}
		continue
	}
}

func (m *SoulModel) BeforeLoop() {
	select {
	case signal := <-m.FightSignalChan:
		if signal == 1 {
			log.Println("收到消息，转为战斗状态")
			m.Base.Status = consts.SOUL_STATE_FIGHT
			m.Base.NextStatus = consts.SOUL_STATE_SUM
		} else {
			log.Println("收到消息，转为邀请队友标志")
			m.Base.Status = consts.SOUL_STATE_INVITE_MEMBER
			m.Base.NextStatus = consts.SOUL_STATE_TEAM
		}
	case <-m.Base.Tick.C:
	}
}

// 检查队伍状态变更
func (m *SoulModel) checkTeamStatusChange() {
	for {
		if m.Base.Status == consts.SOUL_STATE_TEAM { // 判断是否已进入战斗
			screen := robotgo.CaptureImg()
			result := find(readySymbol, screen)
			if valid(result) {
				m.FightSignalChan <- 1
				log.Println("检测到已进入战场")
			}
		} else if m.Base.Status == consts.SOUL_STATE_SUM { // 判断是否有邀请队员弹窗
			screen := robotgo.CaptureImg()
			result := find(checkSymbol, screen)
			if valid(result) {
				m.FightSignalChan <- 2
				log.Println("检测到默认邀请队友弹窗")
			}
		}
		time.Sleep(time.Second)
	}
}
