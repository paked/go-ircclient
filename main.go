package main

import (
	"flag"
	"github.com/nsf/termbox-go"
	ircevent "github.com/thoj/go-ircevent"
	"log"
)

const cdef = termbox.ColorDefault

var tb TextBox = TextBox{}

var clb ChatLogBox = ChatLogBox{}

var room = flag.String("room", "#pakedtheking", "Room you want to join")
var irc = ircevent.IRC("adwdwandba", "adwdwandba")

// var running = make(chan bool)
var done = make(chan bool)

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
	length := len(t.Content)
	if length != 0 {
		t.Content = t.Content[0:(length - 1)]
	}
}

func (t *TextBox) Clear() {
	t.Content = ""
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

func (cl *ChatLogBox) Format(nick string, msg string) string {
	return "[" + nick + "] " + msg
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

	err = irc.Connect("irc.freenode.net:6667")

	if err != nil {
		panic(err)
	}

	// When we've connected to the IRC server, go join the room!
	irc.AddCallback("001", func(e *ircevent.Event) {
		clb.AddMessage("joining " + *room)
		irc.Join(*room)
	})

	irc.AddCallback("JOIN", func(e *ircevent.Event) {
		clb.AddMessage("Joined finally")
	})

	irc.AddCallback("PRIVMSG", func(e *ircevent.Event) {
		clb.AddMessage(clb.Format(e.Nick, e.Message()))
	})

	_, h := termbox.Size()

	tb.SetPos(0, h-1)

	tb.InsertRune('a')
	tb.InsertRune(' ')
	tb.InsertRune('b')
	tb.InsertRune('c')

	clb.SetPos(0, 0)
	go drawLoop()
	eventLoop()

	go irc.Loop()
}

func drawLoop() {
	running := true
	for running {
		draw_all()
		select {
		case <-done:
			log.Println("DONE YET?")
			running = false
		default:
		}
	}

	log.Println("DCX")
}

func eventLoop() {
	running := true
	for running {
		ev := termbox.PollEvent()
		switch ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				running = false
				done <- true

			case termbox.KeySpace:
				tb.InsertRune(' ')

			case termbox.KeyBackspace2:
				tb.DeleteRune()

			case termbox.KeyEnter:
				clb.AddMessage(clb.Format(irc.GetNick(), tb.Content))
				irc.Privmsg(*room, tb.Content)
				tb.Clear()

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

	termbox.Close()

	log.Println("DCS")
}
