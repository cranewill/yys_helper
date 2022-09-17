package tests

import (
	"fmt"
	"github.com/go-vgo/robotgo"
	"github.com/vcaesar/gcv"
	"log"
	"testing"
	"time"
)

func TestRobotgo(t *testing.T) {
	img := robotgo.CaptureImg(100, 100, 200, 200)
	robotgo.SavePng(img, "./img.png")

	part := robotgo.CaptureImg(100, 100, 100, 100)
	robotgo.SavePng(part, "./part.png")
	log.Println(fmt.Sprintf("%+v", gcv.Find(part, img)))
}

func TestGetPos(t *testing.T) {
	for {
		time.Sleep(time.Second)
		x, y := robotgo.GetMousePos()
		log.Println(x, " : ", y)
	}
}
