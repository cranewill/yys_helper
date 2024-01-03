package models

import (
	"github.com/go-vgo/robotgo"
	"github.com/vcaesar/gcv"
	"image"
	"log"
	"time"
)

type IModel interface {
	Start()
	Loop()
	BeforeLoop()
}

type BaseModel struct {
	Status     string
	NextStatus string
	Tick       *time.Ticker
}

var modelPool = make(map[int]IModel)

func GetModel(id int) (IModel, bool) {
	model, ok := modelPool[id]
	return model, ok
}

var heroPieceSymbol,
	matchSymbol,
	teamSymbol,
	quitSymbol,
	fightSymbol,
	readySymbol,
	sumSymbol,
	checkSymbol,
	snakeSymbol,
	emptyTeamSymbol image.Image

func findMulti(targets []image.Image, source image.Image) []gcv.Result {
	results := make([]gcv.Result, len(targets))
	matSource, _ := gcv.ImgToMat(source)
	for i, target := range targets {
		mat, _ := gcv.ImgToMat(target)
		mResult := gcv.FindAllTemplate(matSource, mat, 0.9)
		if len(mResult) == 0 {
			results[i] = gcv.Result{MaxVal: []float32{0}}
		} else {
			results[i] = mResult[0]
		}
	}
	return results
}

func find(target, source image.Image) (result gcv.Result) {
	defer func() {
		if x := recover(); x != nil {
			log.Println("fail, retry...")
			result = gcv.Result{MaxVal: []float32{0}}
		}
	}()
	matT, _ := gcv.ImgToMat(target)
	matS, _ := gcv.ImgToMat(source)
	results := gcv.FindAllTemplate(matS, matT)
	if len(results) > 0 {
		result = results[0]
	} else {
		result = gcv.Result{
			MaxVal: []float32{0},
		}
	}
	return
}

func valid(result gcv.Result) bool {
	return result.MaxVal[0] > 1.0
}

func onFindSymbol(result gcv.Result, state *string, nextState string) {
	robotgo.Move(result.Middle.X/2, result.Middle.Y/2)
	robotgo.MilliSleep(500)
	robotgo.Click()
	*state = nextState
	log.Println("进入", *state)
}
