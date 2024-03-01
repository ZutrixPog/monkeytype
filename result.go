package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
)

const RES_TITLE = ` _____                            ___    __      
/\  ___\                         /\_ \  /\ \__   
\ \ \L\ \     __    ____  __  __\//\ \ \ \ ,_\  
 \ \ ,  /   /'__'\ /',__\/\ \/\ \ \ \ \ \ \ \/  
  \ \ \\ \ /\  __//\__, '\ \ \_\ \ \_\ \_\ \ \_ 
   \ \_\ \_\ \____\/\____/\ \____/ /\____\\ \__\
    \/_/\/ /\/____/\/___/  \/___/  \/____/ \/__/
`

type Metric struct {
	duration     time.Duration
	allChars     int
	correctChars int
}

var _ Drawable = (*Result)(nil)

type Result struct {
	screen  tcell.Screen
	metrics Metric

	rawWpm   int
	wpm      int
	accuracy int
}

func NewResult(screen tcell.Screen, metrics Metric) Drawable {
	return &Result{
		screen:  screen,
		metrics: metrics,
	}
}

func (r *Result) Init() {
	r.rawWpm, r.wpm = r.calcWpm()
	r.accuracy = r.calcAccuracy()
}

func (r *Result) Draw() {
	drawTitle(r.screen, RES_TITLE)
	drawDashedBox(r.screen, r.wpm, r.accuracy, int(r.metrics.duration.Seconds()), r.rawWpm)

	txt := "press enter to continue or q to exit..."
	drawTextCentered(r.screen, len(txt), 15, txt, AppTextStyle)
}

func (r *Result) Update(e tcell.Event) (next Drawable) {
	key := e.(*tcell.EventKey)
	switch key.Key() {
	case tcell.KeyEnter:
		return NewMenu(r.screen)
	default:
		if key.Rune() == 'q' {
			os.Exit(0)
		}
	}
	return nil
}

func (r *Result) calcWpm() (int, int) {
	raw := (float64(r.metrics.allChars) / 5) / r.metrics.duration.Minutes()
	adjusted := raw * (float64(r.metrics.correctChars) / float64(r.metrics.allChars))
	return int(raw), int(adjusted)
}

func (r *Result) calcAccuracy() int {
	return int(float64(r.metrics.correctChars) / float64(r.metrics.allChars) * 100)
}

func drawDashedBox(screen tcell.Screen, wpm, accuracy, duration, raw int) {
	swidth, _ := screen.Size()
	boxLen := swidth / 2
	startWidth := (swidth - boxLen) / 2
	space := 20

	for i := startWidth; i < startWidth+boxLen; i += 2 {
		screen.SetContent(i, 9, tcell.RuneHLine, nil, AppYellowTextStyle)
		screen.SetContent(i, 14, tcell.RuneHLine, nil, AppYellowTextStyle)
	}
	for i := 10; i <= 13; i += 2 {
		screen.SetContent(startWidth, i, tcell.RuneVLine, nil, AppYellowTextStyle)
		screen.SetContent(startWidth+boxLen-1, i, tcell.RuneVLine, nil, AppYellowTextStyle)
	}

	af := "Accuracy: "
	av := fmt.Sprintf("%d", accuracy) + " %"
	wf := "WPM: "
	wv := fmt.Sprintf("%d", wpm)
	tf := "Time: "
	tv := fmt.Sprintf("%d", duration) + "s"
	rf := "Raw: "
	rv := fmt.Sprintf("%d", raw)
	innerStartWidth := (startWidth + (boxLen-space-len(wf)-len(wv)-len(rf)-len(rv))/2)

	drawText(screen, len(af), innerStartWidth+3, 11, af, AppTextStyle)
	drawText(screen, len(av), innerStartWidth+3+len(af), 11, av, AppYellowTextStyle)

	drawText(screen, len(wf), innerStartWidth+3, 12, wf, AppTextStyle)
	drawText(screen, len(wv), innerStartWidth+8+len(wf), 12, wv, AppYellowTextStyle)

	drawText(screen, len(tf), innerStartWidth+space+3, 11, tf, AppTextStyle)
	drawText(screen, len(tv), innerStartWidth+4+space+len(wf), 11, tv, AppYellowTextStyle)

	drawText(screen, len(rf), innerStartWidth+space+3, 12, rf, AppTextStyle)
	drawText(screen, len(rv), innerStartWidth+space+4+len(wf), 12, rv, AppYellowTextStyle)
}
