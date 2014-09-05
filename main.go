package main

import (
	"github.com/nsf/termbox-go"
	// "strings"
	// "log"
)

const cdef = termbox.ColorDefault

var tb TextBox = TextBox{"", 10, 10, 0}

var bottom string = ""

type TextBox struct {
	Content string
	X       int `default: 10`
	Y       int `default: 20`
	Cursor  int
}

func (t *TextBox) SetPos(x int, y int) {
	t.X = x
	t.Y = y
}

func (t *TextBox) InsertRune(o rune) {
	t.Content += string(o)
}

func (t *TextBox) DeleteRune() {
	newLength := len(t.Content) - 1
	newContent := ""
	for i := 0; i < newLength; i++ {
		newContent += string(t.Content[i])
	}

	t.Content = newContent
}

func (t *TextBox) Draw() {
	drawString(t.X, t.Y, t.Content)
}

func drawString(x int, y int, str string) {
	length := len(str)
	for i := 0; i < length; i++ {
		termbox.SetCell(x+i, y, rune(str[i]), cdef, cdef)
	}
}

func draw_all() {
	termbox.Clear(cdef, cdef)

	// termbox.SetCell(midx, midy, rune('o'), cdef, cdef)

	tb.Draw()

	termbox.Flush()
}

func main() {
	err := termbox.Init()

	if err != nil {
		panic(err)
	}

	defer termbox.Close()

	_, h := termbox.Size()

	tb.SetPos(0, h-1)

	tb.InsertRune('a')
	tb.InsertRune(' ')
	tb.InsertRune('b')
	tb.InsertRune('c')

	// draw_all()
	running := true
	for running {
		draw_all()
		// log.Println("hey")
		ev := termbox.PollEvent()
		switch ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				running = false
			case termbox.KeySpace:
				tb.InsertRune(' ')
			case termbox.KeyBackspace2:
				tb.DeleteRune()
			default:
				if ev.Ch != 0 {
					tb.InsertRune(ev.Ch)
				}
			}
		case termbox.EventResize:
			draw_all()
		}
	}
}
