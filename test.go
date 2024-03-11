package main

import (
	_ "embed"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
)

//go:embed res/words.txt
var words string

//go:embed res/long_quotes.txt
var longQuotes string

//go:embed res/medium_quotes.txt
var mediumQuotes string

//go:embed res/short_quotes.txt
var shortQuotes string

const (
	TEST_WORD int = iota
	TEST_TIME
	TEST_QUOTE
)

const (
	MAX_NUM             = 100000
	NUM_FACTOR          = 8
	PUNCTUATIONS_FACTOR = 4
	WRAPPER_FACTOR      = 7
	WRONG_CHAR          = "|"
)

var QuoteTypes = []string{"short", "medium", "long"}

type Config struct {
	Punctuation bool
	Number      bool
	Words       int
	Duration    int
	QuoteLen    int
}

var _ Drawable = (*Test)(nil)

// TODO: we need more abstraction here
type Test struct {
	screen tcell.Screen
	kind   int
	config Config

	startTime  time.Time
	txt        string
	typedTxt   string
	words      int
	typedWords int
}

func NewTest(screen tcell.Screen, kind int, config Config) Drawable {
	return &Test{
		screen: screen,
		kind:   kind,
		config: config,
	}
}

func (t *Test) Init() {
	t.generateText()

	if t.kind == TEST_TIME {
		ti := time.NewTicker(time.Second)
		go func() {
			for {
				<-ti.C
				_ = t.screen.PostEvent(tcell.NewEventMouse(0, 0, tcell.RuneS1, tcell.ModAlt))
			}
		}()
	}
}

func (t *Test) Draw() {
	swidth, _ := t.screen.Size()
	lineLen := swidth / 2

	startW := centerWidth(t.screen, lineLen)
	startH := 4
	t.drawCounter(startW-2, startH-1)

	var x int
	drawFn := func(ch byte, style tcell.Style) {
		t.screen.SetContent(startW+(x%lineLen), startH+(x/lineLen), rune(ch), nil, style)
		x++
	}

	targetWords := strings.Fields(t.txt)
	typedWords := strings.Fields(t.typedTxt)
	for i := 0; i < len(targetWords); i++ {
		targetWord := targetWords[i]
		typedWord := ""
		if i < len(typedWords) {
			typedWord = typedWords[i]
		}

		for j := 0; j < len(targetWord); j++ {
			if j < len(typedWord) {
				if targetWord[j] == typedWord[j] {
					drawFn(targetWord[j], CorrectTextStyle)
				} else {
					drawFn(targetWord[j], WrongTextStyle)
				}
			} else {
				drawFn(targetWord[j], TargetTextStyle)
			}
		}

		if len(typedWord) > len(targetWord) {
			extraChars := typedWord[len(targetWord):]
			for _, ch := range extraChars {
				if ch == ' ' || string(ch) == WRONG_CHAR {
					break
				}
				drawFn(byte(ch), WrongTextStyle)
			}
		}

		if i != len(targetWords)-1 {
			drawFn(' ', AppTextStyle)
		}
	}

	for i := len(targetWords); i < len(typedWords); i++ {
		for _, ch := range typedWords[i] {
			drawFn(byte(ch), WrongTextStyle)
		}

		if i != len(typedWords)-1 {
			drawFn(' ', AppTextStyle)
		}
	}
}

func (t *Test) drawCounter(w, h int) {
	counter := fmt.Sprintf("%d/%d", t.typedWords+1, t.words)
	if t.kind == TEST_TIME {
		dt := time.Time{}
		td := t.config.Duration - int(time.Since(t.startTime).Seconds())
		if t.startTime == dt {
			td = t.config.Duration
		}
		counter = fmt.Sprintf("%d", td)
	}

	drawText(t.screen, len(counter), w, h, counter, AppYellowTextStyle)
}

func (t *Test) Update(event tcell.Event) Drawable {
	key := event.(*tcell.EventKey)
	if key.Key() == tcell.KeyRune {
		if key.Rune() == ' ' {
			for len(t.typedTxt) < len(t.txt) && t.txt[len(t.typedTxt)] != ' ' {
				t.typedTxt += WRONG_CHAR
			}
			t.typedWords += 1
		}
		dt := time.Time{}
		if t.startTime == dt {
			t.startTime = time.Now()
		}
		t.typedTxt += string(key.Rune())
	}

	switch key.Key() {
	case tcell.KeyBackspace, tcell.KeyBackspace2:
		if len(t.typedTxt) != 0 {
			if t.typedTxt[len(t.typedTxt)-1] == ' ' {
				t.typedWords -= 1
			}
			t.typedTxt = t.typedTxt[:len(t.typedTxt)-1]
		}
	}

	if next := t.finish(); next != nil {
		return next
	}

	return nil
}

func (t *Test) finish() Drawable {
	if t.txt == t.typedTxt || t.words == t.typedWords ||
		(t.kind == TEST_TIME && time.Now().After(t.startTime.Add(time.Second*time.Duration(t.config.Duration)))) {
		duration := time.Since(t.startTime)
		if t.kind == TEST_TIME {
			duration = time.Duration(t.config.Duration) * time.Second
		}
		correct := 0
		for i := 0; i < len(t.typedTxt); i++ {
			if i >= len(t.txt) {
				break
			}
			if t.typedTxt[i] == t.txt[i] {
				correct += 1
			}
		}

		return NewResult(t.screen, Metric{
			duration:     duration,
			allChars:     len(t.typedTxt),
			correctChars: correct,
		})
	}

	return nil
}

func (t *Test) generateText() {
	if t.kind == TEST_QUOTE {
		t.txt = generateQuote(t.config)
	} else {
		if t.config.Words == 0 {
			t.config.Words = t.config.Duration + t.config.Duration/2
		}
		t.txt = generateWords(t.config)
	}

	t.words = len(strings.Split(t.txt, " "))
}

func generateQuote(conf Config) string {
	categoryToFile := map[string]string{
		"short":  shortQuotes,
		"medium": mediumQuotes,
		"long":   longQuotes,
	}

	content := categoryToFile[QuoteTypes[conf.QuoteLen]]

	quotes := strings.Split(string(content), "\n")
	var nonEmptyQuotes []string
	for _, quote := range quotes {
		if quote != "" {
			nonEmptyQuotes = append(nonEmptyQuotes, quote)
		}
	}

	selectedQuote := nonEmptyQuotes[rand.Intn(len(nonEmptyQuotes))]
	selectedQuote = strings.TrimSpace(selectedQuote)
	return selectedQuote
}

func generateWords(conf Config) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	words := strings.Fields(words)
	r.Shuffle(len(words), func(i, j int) {
		words[i], words[j] = words[j], words[i]
	})

	selectedWords := words[:conf.Words]

	if conf.Punctuation {
		punctuation := []string{".", ",", "!", "?", ";", ":"}
		indexes := generateRandomNumbers(conf.Words/PUNCTUATIONS_FACTOR, 0, conf.Words-1)
		for _, i := range indexes {
			selectedWords[i] += punctuation[r.Intn(len(punctuation))]
		}

		wrappers := []string{"[]", "()", "{}", `""`, `''`}
		indexes = generateRandomNumbers(conf.Words/WRAPPER_FACTOR, 0, conf.Words-1)
		for _, i := range indexes {
			wrap := wrappers[r.Intn(len(wrappers))]
			selectedWords[i] = fmt.Sprintf("%c%s%c", wrap[0], selectedWords[i], wrap[1])
		}
	}

	if conf.Number {
		indexes := generateRandomNumbers(conf.Words/NUM_FACTOR, 0, conf.Words-1)
		for _, i := range indexes {
			selectedWords[i] = fmt.Sprintf("%d", rand.Intn(MAX_NUM))
		}
	}

	result := strings.Join(selectedWords, " ")

	return result
}

func generateRandomNumbers(n, min, max int) []int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	randomNumbers := make([]int, n)
	for i := 0; i < n; i++ {
		randomNumbers[i] = r.Intn(max-min+1) + min
	}

	return randomNumbers
}
