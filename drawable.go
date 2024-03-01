package main

import "github.com/gdamore/tcell/v2"

type Drawable interface {
	Init()
	Draw()
	Update(tcell.Event) (next Drawable)
}
