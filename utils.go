package main

import (
	"strings"
	"unicode"

	"github.com/gdamore/tcell/v2"
)

const MAIN_TITLE = `  _   _                  __                       __                              
 /'\_/.\                /\ \                     /\ \__                           
/\      \    ___     ___\ \ \/'\      __   __  __\ \ ,_\  __  __  _____      __   
\ \ \__\ \  / __.\ /' _ .\ \ , <    /'__.\/\ \/\ \\ \ \/ /\ \/\ \/\ '__.\  /'__.\ 
 \ \ \_/\ \/\ \L\ \/\ \/\ \ \ \\.\ /\  __/\ \ \_\ \\ \ \_\ \ \_\ \ \ \L\ \/\  __/ 
  \ \_\\ \_\ \____/\ \_\ \_\ \_\ \_\ \____\\/.____ \\ \__\\/.____ \ \ ,__/\ \____\
   \/_/ \/_/\/___/  \/_/\/_/\/_/\/_/\/____/ ./___/> \\/__/ ./___/> \ \ \/  \/____/
											                                    /\___/         /\___/\ \_\         
											                                    \/__/          \/__/  \/_/         `

func drawTitle(screen tcell.Screen, title string) int {
	lines := strings.Split(title, "\n")

	x := centerWidth(screen, len(lines[0]))
	y := 1

	for i, line := range lines {
		for j, char := range line {
			screen.SetContent(x+j, y+i, char, nil, AppYellowTextStyle)
		}
	}

	return len(lines) + 2
}

func drawText(screen tcell.Screen, lineLen, startW, startH int, text string, style tcell.Style) int {
	lines := splitTextIntoLines(text, lineLen)

	for i, line := range lines {
		for j, ch := range line {
			screen.SetContent(startW+j, startH+i, rune(ch), nil, style)
		}
	}

	return startH + 2
}

func drawTextCentered(screen tcell.Screen, lineLen, startH int, text string, style tcell.Style) int {
	sWidth, _ := screen.Size()
	startW := (sWidth - lineLen) / 2

	return drawText(screen, lineLen, startW, startH, text, style)
}

func splitTextIntoLines(text string, maxLength int) []string {
	var lines []string
	words := strings.FieldsFunc(text, unicode.IsSpace)

	currentLine := ""
	currentLength := 0

	for _, word := range words {
		wordLength := len(word)

		if currentLength+wordLength > maxLength && currentLength > 0 {
			lines = append(lines, currentLine)
			currentLine = ""
			currentLength = 0
		}

		if currentLength > 0 {
			currentLine += " "
			currentLength++
		}
		currentLine += word
		currentLength += wordLength
	}

	if currentLength > 0 {
		lines = append(lines, currentLine)
	}

	return lines
}

func drawBox(screen tcell.Screen, startWidth, startHeight, width, height int, bgStyle, borderStyle tcell.Style) int {
	// TODO: safety checks
	// sWidth, _ := screen.Size()
	endWidth := startWidth + width
	endHeight := startHeight + height

	for i := startWidth; i <= endWidth; i++ {
		for j := startHeight; j <= endHeight; j++ {
			screen.SetContent(i, j, ' ', nil, bgStyle)
		}
	}

	defStyle := tcell.Style{}
	if borderStyle == defStyle {
		return endHeight + 1
	}
	for i := startWidth; i <= endWidth; i++ {
		screen.SetContent(i, startHeight, tcell.RuneHLine, nil, borderStyle)
		screen.SetContent(i, endHeight, tcell.RuneHLine, nil, borderStyle)
	}

	for j := startHeight; j <= endHeight; j++ {
		screen.SetContent(startWidth, j, tcell.RuneVLine, nil, borderStyle)
		screen.SetContent(endWidth, j, tcell.RuneVLine, nil, borderStyle)
	}

	screen.SetContent(startWidth, startHeight, '╭', nil, borderStyle)
	screen.SetContent(endWidth, startHeight, '╮', nil, borderStyle)
	screen.SetContent(startWidth, endHeight, '╰', nil, borderStyle)
	screen.SetContent(endWidth, endHeight, '╯', nil, borderStyle)
	return endHeight + 1
}

func drawCenteredBox(screen tcell.Screen, startHeight, width, height int, bgStyle, borderStyle tcell.Style) (int, int) {
	sWidth, _ := screen.Size()
	startWidth := (sWidth - width) / 2

	return startWidth, drawBox(screen, startWidth, startHeight, width, height, bgStyle, borderStyle)
}

func centerWidth(screen tcell.Screen, textLen int) int {
	width, _ := screen.Size()
	return (width - textLen) / 2
}

func Contains(arr []int, target int) bool {
	for i, _ := range arr {
		if arr[i] == target {
			return true
		}
	}

	return false
}
