package main

import (
	"flag"
	"github.com/nsf/termbox-go"
	ircevent "github.com/thoj/go-ircevent"
	"log"
	"os"
)

const (
	cdef      = termbox.ColorDefault
	MSG_NOTIF = "1"
	MSG_NORM  = "0"
)

var (
	tb     TextBox
	clb    ChatLogBox
	locked bool

	room = flag.String("room", "#pakedtheking", "Room you want to join")
	irc  = ircevent.IRC("adwdwandba", "adwdwandba")
	done = make(chan bool)
)

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
	Content []ChatLogMessage
	X       int
	Y       int
}

func (cl *ChatLogBox) SetPos(x int, y int) {
	cl.X = x
	cl.Y = y
}

func (cl *ChatLogBox) AddMessage(msg ChatLogMessage) {
	cl.Content = append(cl.Content, msg)
}

func (cl *ChatLogBox) Draw() {
	// w, h := termbox.Size()
	for i := 0; i < len(cl.Content); i++ {
		cl.DrawMessage(cl.X, cl.Y+i, cl.Content[i])
	}
}

type ChatLogMessage struct {
	Message string
	Nick    string
	Type    string
}

func (cl *ChatLogBox) format(msg ChatLogMessage) string {
	str := ""
	if msg.Type == MSG_NORM {
		str = "[" + msg.Nick + "] " + msg.Message
	} else if msg.Type == MSG_NOTIF {
		str = msg.Message
	}
	return str
}

func (cl *ChatLogBox) DrawMessage(x int, y int, msg ChatLogMessage) {
	str := cl.format(msg)
	length := len(str)
	colourOne := termbox.ColorDefault
	colourTwo := colourOne
	if msg.Type == MSG_NOTIF {
		colourOne = termbox.ColorMagenta
	}

	for i := 0; i < length; i++ {
		termbox.SetCell(x+i, y, rune(str[i]), colourOne, colourTwo)
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

	_, h := termbox.Size()

	drawString(0, h-1, "->")

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
		clb.AddMessage(ChatLogMessage{"Joining " + *room, "", MSG_NOTIF})
		irc.Join(*room)
	})

	irc.AddCallback("JOIN", func(e *ircevent.Event) {

		clb.AddMessage(ChatLogMessage{"Successfully joined " + *room, "", MSG_NOTIF})
	})

	irc.AddCallback("PRIVMSG", func(e *ircevent.Event) {
		clb.AddMessage(ChatLogMessage{e.Message(), e.Nick, MSG_NORM})
	})

	_, h := termbox.Size()

	tb.SetPos(3, h-1)

	clb.SetPos(0, 0)
	go drawLoop()
	eventLoop()

	// go irc.Loop()
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
				clb.AddMessage(ChatLogMessage{tb.Content, irc.GetNick(), MSG_NORM})
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
	os.Exit(0)
}
