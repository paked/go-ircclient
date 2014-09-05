package main

import (
	"flag"
	"github.com/nsf/termbox-go"
	ircevent "github.com/thoj/go-ircevent"
)

const cdef = termbox.ColorDefault

var tb TextBox = TextBox{}

var clb ChatLogBox = ChatLogBox{}

var room = flag.String("room", "paked", "Room you want to join")

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
	t.Content = t.Content[0:newLength]
}

func (t *TextBox) Draw() {
	drawString(t.X, t.Y, t.Content)
}

type ChatLogBox struct {
	Content []string
	X       int
	Y       int
}

func (cl *ChatLogBox) SetPos(x int, y int) {
	cl.X = x
	cl.Y = y
}

func (cl *ChatLogBox) AddMessage(msg string) {
	cl.Content = append(cl.Content, msg)
}

func (cl *ChatLogBox) Draw() {
	// w, h := termbox.Size()
	for i := 0; i < len(cl.Content); i++ {
		drawString(cl.X, cl.Y+i, cl.Content[i])
	}
}

func drawString(x int, y int, str string) {
	length := len(str)
	for i := 0; i < length; i++ {
		termbox.SetCell(x+i, y, rune(str[i]), cdef, cdef)
	}
}

func draw_all() {
	termbox.Clear(cdef, cdef)

	tb.Draw()
	clb.Draw()

	termbox.Flush()
}

func main() {
	err := termbox.Init()

	if err != nil {
		panic(err)
	}

	defer termbox.Close()

	irc := ircevent.IRC("heyoitsmeabot", "paked")
	err = irc.Connect("irc.freenode.net:6667")

	if err != nil {
		panic(err)
	}

	// When we've connected to the IRC server, go join the room!
	irc.AddCallback("001", func(e *ircevent.Event) {
		clb.AddMessage("joining devsofa")
		irc.Join("devsofa")
	})

	irc.AddCallback("JOIN", func(e *ircevent.Event) {
		clb.AddMessage("Joined finally")
		tb.Content = "JOINED"
	})

	// Check each message to see if it contains a URL, and return the title
	irc.AddCallback("PRIVMSG", func(e *ircevent.Event) {
		// log.Printf("[%v] %v", e.Nick, e.Message())
	})

	go irc.Loop()

	_, h := termbox.Size()

	tb.SetPos(0, h-1)

	tb.InsertRune('a')
	tb.InsertRune(' ')
	tb.InsertRune('b')
	tb.InsertRune('c')

	clb.SetPos(0, 0)

	// clb.AddMessage("Hey boys")
	// clb.AddMessage("Hey man")

	// draw_all()
	running := true
	for running {
		// clb.AddMessage("yolo")
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
		default:
		}

	}
}
