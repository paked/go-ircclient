package main

import (
	"github.com/nsf/termbox-go"
	// "log"
)

const cdef = termbox.ColorDefault

func drawString(x int, y int, str string) {
	length := len(str)
	for i := 0; i < length; i++ {
		termbox.SetCell(x+i, y, rune(str[i]), cdef, cdef)
	}
}

func draw_all() {
	termbox.Clear(cdef, cdef)

	w, h := termbox.Size()

	midy := h / 2
	midx := w / 2

	termbox.SetCell(midx, midy, rune('o'), cdef, cdef)

	drawString(0, 0, "lolol")

	termbox.Flush()
}

func main() {
	err := termbox.Init()

	if err != nil {
		panic(err)
	}

	defer termbox.Close()

	draw_all()
	running := true
	for running {
		ev := termbox.PollEvent()
		switch ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				running = false
			}
		case termbox.EventResize:
			draw_all()
			//RESIZE
		}
	}
}
