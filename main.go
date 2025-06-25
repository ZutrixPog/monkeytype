package main

import (
	"log"

	"github.com/gdamore/tcell/v2"
)

var (
	BackgroundColor    = tcell.Color239
	AppTextStyle       = tcell.StyleDefault.Background(BackgroundColor).Foreground(tcell.Color252)
	AppYellowTextStyle = tcell.StyleDefault.Background(BackgroundColor).Foreground(tcell.Color214)
	TargetTextStyle    = tcell.StyleDefault.Background(BackgroundColor).Foreground(tcell.Color246)
	CorrectTextStyle   = tcell.StyleDefault.Background(BackgroundColor).Foreground(tcell.Color252)
	WrongTextStyle     = tcell.StyleDefault.Background(BackgroundColor).Foreground(tcell.Color197)
	ExtraTextStyle     = tcell.StyleDefault.Background(BackgroundColor).Foreground(tcell.Color196)
)

func FillBackground(s tcell.Screen, color tcell.Color) {
	width, height := s.Size()

	for row := 0; row < height; row++ {
		for col := 0; col <= width; col++ {
			s.SetContent(col, row, ' ', nil, tcell.StyleDefault.Background(BackgroundColor))
		}
	}
}

func main() {
	// Initialize screen
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}
	s.SetCursorStyle(tcell.CursorStyleBlinkingUnderline)
	s.Clear()

	quit := func() {
		maybePanic := recover()
		s.Fini()
		if maybePanic != nil {
			panic(maybePanic)
		}
	}
	defer quit()

	currElement := NewMenu(s)
	currElement.Init()

	// Event loop
	for {
		s.Clear()
		FillBackground(s, BackgroundColor)
		currElement.Draw()
		s.Show()

		// Poll event
		ev := s.PollEvent()

		switch ev := ev.(type) {
		case *tcell.EventResize:
			s.Clear()
			currElement.Draw()
			s.Sync()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || (ev.Key() == tcell.KeyRune && ev.Rune() == 'q') {
				return
			} else if ev.Key() == tcell.KeyCtrlL {
				s.Sync()
			}

			nextElement := currElement.Update(ev)
			if nextElement != nil {
				currElement = nextElement
				currElement.Init()
			}
		}
	}
}
