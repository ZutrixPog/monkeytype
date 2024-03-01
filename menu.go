package main

import (
	"fmt"
	"strconv"

	"github.com/gdamore/tcell/v2"
)

type TestPrompt interface {
	Name() string
	Draw(screen tcell.Screen, startRow, startCol, boxWidth int)
	Update(e *tcell.EventKey) bool
	Config() Config
}

var TestTypes = []TestPrompt{
	&WordPrompt{wordCount: 2, inCheckList: true},
	&TimePrompt{duration: 2, inCheckList: true},
	&QuotePrompt{},
}

const (
	PUNCTUATION int = iota
	NUMBER
)

var (
	BoxStyle      = tcell.StyleDefault.Background(tcell.Color237).Foreground(tcell.Color252)
	SelectedStyle = tcell.StyleDefault.Background(tcell.Color237).Foreground(tcell.Color214)

	CheckListItems   = []string{"punctuation", "numbers"}
	DurationChoices  = []string{"15", "30", "60", "120"}
	WordCountChoices = []string{"10", "25", "50", "100"}
)

var _ Drawable = (*Menu)(nil)

type Menu struct {
	screen   tcell.Screen
	testType int
	config   Config

	inPrompt bool
}

func NewMenu(screen tcell.Screen) Drawable {
	return &Menu{
		screen,
		TEST_WORD,
		Config{},
		false,
	}
}

func (m *Menu) Init() {
}

func (m *Menu) Draw() {
	startingRow := drawTitle(m.screen, MAIN_TITLE)
	text := "choose a style, press s to start test..."

	startingRow = drawTextCentered(m.screen, len(text), startingRow, text, AppTextStyle)
	m.drawChoiceBox(m.screen, startingRow)
}

func (m *Menu) Update(event tcell.Event) Drawable {
	k := event.(*tcell.EventKey)
	if k.Key() == tcell.KeyRune && k.Rune() == 's' {
		conf := TestTypes[m.testType].Config()
		return NewTest(m.screen, m.testType, conf)
	}

	if m.inPrompt {
		if TestTypes[m.testType].Update(k) {
			m.inPrompt = false
		}
	} else {
		switch k.Key() {
		case tcell.KeyLeft:
			if m.testType > 0 {
				m.testType -= 1
			} else {
				m.testType = len(TestTypes) - 1
			}
		case tcell.KeyRight:
			if m.testType < len(TestTypes)-1 {
				m.testType += 1
			} else {
				m.testType = 0
			}
		case tcell.KeyDown:
			m.inPrompt = true
			TestTypes[m.testType].Update(k)
		}

	}
	return nil
}

func (m *Menu) drawChoiceBox(screen tcell.Screen, startingRow int) {
	if len(TestTypes) == 0 {
		return
	}

	sWidth, _ := screen.Size()
	boxWidth := sWidth / 2
	startWidth, _ := drawCenteredBox(screen, startingRow, boxWidth, 2, AppTextStyle, tcell.Style{})

	totalSpace := boxWidth - 2
	space := totalSpace / (len(TestTypes) + 1)
	for i, ch := range TestTypes {
		choice := ch.Name()
		style := AppTextStyle
		if i == m.testType {
			style = AppYellowTextStyle
			if !m.inPrompt {
				choice = fmt.Sprintf("%c %s", tcell.RuneDiamond, ch.Name())
			}
		}

		startW := ((i + 1) * space) + startWidth - (len(choice) / 2)
		drawText(screen, len(choice), startW, startingRow+1, choice, style)
	}

	for i := startWidth; i <= startWidth+boxWidth; i++ {
		screen.SetContent(i, startingRow, tcell.RuneHLine, nil, AppYellowTextStyle)
	}

	for i := startingRow; i <= startingRow+3; i++ {
		screen.SetContent(startWidth, i, tcell.RuneVLine, nil, AppYellowTextStyle)
	}
	screen.SetContent(startWidth, startingRow, '╭', nil, AppYellowTextStyle)

	for i := startWidth; i <= startWidth+boxWidth; i++ {
		screen.SetContent(i, startingRow+3, tcell.RuneHLine, nil, AppYellowTextStyle)
	}
	screen.SetContent(startWidth, startingRow+3, '╰', nil, AppYellowTextStyle)
	screen.SetContent(startWidth+boxWidth, startingRow+3, '╮', nil, AppYellowTextStyle)

	for i := startingRow + 4; i <= startingRow+13; i++ {
		screen.SetContent(startWidth+boxWidth, i, tcell.RuneVLine, nil, AppYellowTextStyle)
	}

	TestTypes[m.testType].Draw(screen, startingRow+4, startWidth+1, boxWidth)
}

// TODO: dedup code for word and time prompts
type WordPrompt struct {
	includedWords []int
	currWord      int

	wordCount     int
	currWordCount int

	inCheckList bool
	inPrompt    bool
}

func (w *WordPrompt) Name() string {
	return "word"
}

func (w *WordPrompt) Draw(screen tcell.Screen, startRow, startCol, boxWidth int) {
	lineWidth := (boxWidth / 2) + startCol - 1

	// draw first column
	chCol := startCol - 6 + (lineWidth-startCol)/2
	for i, t := range CheckListItems {
		item := fmt.Sprintf("[] %s", t)
		style := AppTextStyle
		if Contains(w.includedWords, i) {
			item = fmt.Sprintf("[%c] %s", 'X', t)
		}
		if i == w.currWord && w.inCheckList && w.inPrompt {
			style = AppYellowTextStyle
			item = string(tcell.RuneDiamond) + item
		}
		drawText(screen, len(item), chCol, startRow+((i+1)*3), item, style)
	}

	// draw vertical line
	for i := startRow + 1; i < startRow+9; i++ {
		screen.SetContent(lineWidth, i, tcell.RuneVLine, nil, AppYellowTextStyle)
	}

	// draw second column
	dcCol := lineWidth + (startCol+boxWidth-lineWidth)/2
	for i, t := range WordCountChoices {
		item := t
		style := AppTextStyle
		if !w.inCheckList {
			if i == w.wordCount {
				style = AppYellowTextStyle
			}
			if i == w.currWordCount && w.inPrompt {
				item = fmt.Sprintf("%c %s", tcell.RuneDiamond, t)
			}
		}
		if i == w.wordCount {
			style = AppYellowTextStyle
		}

		drawText(screen, len(item), dcCol, startRow+((i+1)*2), item, style)
	}

}

func (w *WordPrompt) Update(event *tcell.EventKey) bool {
	switch event.Key() {
	case tcell.KeyDown:
		if !w.inPrompt {
			w.currWord = 0
			w.inCheckList = true
			w.inPrompt = true
			return false
		}

		if w.inCheckList {
			if w.currWord < len(CheckListItems)-1 {
				w.currWord += 1
			} else {
				w.currWord = 0
			}
		} else {
			if w.wordCount < len(WordCountChoices)-1 {
				w.currWordCount += 1
			} else {
				w.currWordCount = 0
			}
		}
	case tcell.KeyUp:
		if w.inCheckList {
			if w.currWord > 0 {
				w.currWord -= 1
			} else {
				w.inPrompt = false
				return true
			}
		} else {
			if w.currWordCount > 0 {
				w.currWordCount -= 1
			} else {
				w.inPrompt = false
				return true
			}
		}
	case tcell.KeyLeft, tcell.KeyRight:
		w.inCheckList = !w.inCheckList
	case tcell.KeyEnter:
		if w.inCheckList {
			if !Contains(w.includedWords, w.currWord) {
				w.includedWords = append(w.includedWords, w.currWord)
			} else {
				newList := make([]int, 0)
				for _, item := range w.includedWords {
					if item != w.currWord {
						newList = append(newList, item)
					}
				}
				w.includedWords = newList
			}
		} else {
			w.wordCount = w.currWordCount
		}
	}

	return false
}

func (w *WordPrompt) Config() Config {
	wc, _ := strconv.Atoi(WordCountChoices[w.wordCount])
	return Config{
		Words:       wc,
		Punctuation: Contains(w.includedWords, PUNCTUATION),
		Number:      Contains(w.includedWords, NUMBER),
	}
}

type TimePrompt struct {
	includedWords []int
	currWord      int

	duration     int
	currDuration int

	inCheckList bool
	inPrompt    bool
}

func (w *TimePrompt) Name() string {
	return "time"
}

func (w *TimePrompt) Draw(screen tcell.Screen, startRow, startCol, boxWidth int) {
	lineWidth := (boxWidth / 2) + startCol - 1

	// draw first column
	chCol := startCol - 6 + (lineWidth-startCol)/2
	for i, t := range CheckListItems {
		item := fmt.Sprintf("[] %s", t)
		style := AppTextStyle
		if Contains(w.includedWords, i) {
			item = fmt.Sprintf("[%c] %s", 'X', t)
		}
		if i == w.currWord && w.inCheckList && w.inPrompt {
			style = AppYellowTextStyle
			item = string(tcell.RuneDiamond) + item
		}
		drawText(screen, len(item), chCol, startRow+((i+1)*3), item, style)
	}

	// draw vertical line
	for i := startRow + 1; i < startRow+9; i++ {
		screen.SetContent(lineWidth, i, tcell.RuneVLine, nil, AppYellowTextStyle)
	}

	// draw second column
	dcCol := lineWidth + (startCol+boxWidth-lineWidth)/2
	for i, t := range DurationChoices {
		item := t
		style := AppTextStyle
		if !w.inCheckList {
			if i == w.duration {
				style = AppYellowTextStyle
			}
			if i == w.currDuration && w.inPrompt {
				item = fmt.Sprintf("%c %s", tcell.RuneDiamond, t)
			}
		}
		if i == w.duration {
			style = AppYellowTextStyle
		}

		drawText(screen, len(item), dcCol, startRow+((i+1)*2), item, style)
	}

}

func (w *TimePrompt) Update(event *tcell.EventKey) bool {
	switch event.Key() {
	case tcell.KeyDown:
		if !w.inPrompt {
			w.currWord = 0
			w.inCheckList = true
			w.inPrompt = true
			return false
		}

		if w.inCheckList {
			if w.currWord < len(CheckListItems)-1 {
				w.currWord += 1
			} else {
				w.currWord = 0
			}
		} else {
			if w.duration < len(DurationChoices)-1 {
				w.currDuration += 1
			} else {
				w.currDuration = 0
			}
		}
	case tcell.KeyUp:
		if w.inCheckList {
			if w.currWord > 0 {
				w.currWord -= 1
			} else {
				w.inPrompt = false
				return true
			}
		} else {
			if w.currDuration > 0 {
				w.currDuration -= 1
			} else {
				w.inPrompt = false
				return true
			}
		}
	case tcell.KeyLeft, tcell.KeyRight:
		w.inCheckList = !w.inCheckList
	case tcell.KeyEnter:
		if w.inCheckList {
			if !Contains(w.includedWords, w.currWord) {
				w.includedWords = append(w.includedWords, w.currWord)
			} else {
				newList := make([]int, 0)
				for _, item := range w.includedWords {
					if item != w.currWord {
						newList = append(newList, item)
					}
				}
				w.includedWords = newList
			}
		} else {
			w.duration = w.currDuration
		}
	}

	return false
}

func (w *TimePrompt) Config() Config {
	d, _ := strconv.Atoi(DurationChoices[w.duration])
	return Config{
		Duration:    d,
		Punctuation: Contains(w.includedWords, PUNCTUATION),
		Number:      Contains(w.includedWords, NUMBER),
	}
}

type QuotePrompt struct {
	qType int

	InPrompt bool
}

func (w *QuotePrompt) Name() string {
	return "quote"
}

func (w *QuotePrompt) Update(event *tcell.EventKey) bool {
	switch event.Key() {
	case tcell.KeyUp:
		if w.qType == 0 {
			w.InPrompt = false
			return true
		}
		w.qType -= 1
	case tcell.KeyDown:
		if w.qType < len(QuoteTypes)-1 && w.InPrompt {
			w.qType += 1
		} else {
			w.InPrompt = true
		}
	}

	return false
}

func (w *QuotePrompt) Draw(screen tcell.Screen, startRow, startCol, boxWidth int) {
	startRow += 1

	for i, t := range QuoteTypes {
		text := t
		style := AppTextStyle
		if i == w.qType && w.InPrompt {
			style = AppYellowTextStyle
			text = fmt.Sprintf("%c %s", tcell.RuneDiamond, t)
		}
		drawTextCentered(screen, len(text), startRow+i*3, text, style)
	}
}

func (w *QuotePrompt) Config() Config {
	return Config{
		QuoteLen: w.qType,
	}
}
