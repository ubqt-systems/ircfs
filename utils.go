package main

import (
	"strings"

	"github.com/go-irc/irc"
	"github.com/ubqt-systems/cleanmark"
	"github.com/ubqt-systems/fslib"
)

type fname int
const (
	faction fname = iota
	fbuffer
	fhighlight
	fself
	fserver
	fsidebar
	fstatus
	ftitle
)

type msg struct {
	buff string
	data string
	from string
	fn   fname
}

// Private message
func pm(s *server, msg string) error {
	token := strings.Fields(msg)
	m := &irc.Message{
		Command: "PRIVMSG",
		Prefix: &irc.Prefix{
			Name: s.conf.Name,
		},
		Params: token[:1],
	}
	// Param[1] is the body of the msg
	m.Params = append(m.Params, strings.Join(token[1:], " "))
	return sendmsg(s, m)
}

func sendmsg(s *server, m *irc.Message) error {
	w := irc.NewWriter(s.conn)
	return w.WriteMessage(m)
}

func title(name string, s *server, m *irc.Message) {
	s.m <- &msg{
		buff: name,
		data: m.Trailing(),
		fn:   ftitle,
	}
}

func feed(fn fname, name string, s *server, m *irc.Message) {
	s.m <- &msg{
		buff: name,
		data: m.Trailing(),
		from: m.Prefix.Name,
		fn:   fn,
	}
}

func status(s *server, m *irc.Message) {
	// Just use m.Params[0] for the fname
}

func fileWriter(c *fslib.Control, m *msg) {
	c.CreateBuffer(m.buff, "feed")
	var w *fslib.WriteCloser
	switch m.fn {
	case fbuffer, faction, fhighlight:
		w = c.MainWriter(m.buff, "feed")
		feed := cleanmark.NewCleaner(w)
		defer feed.Close()
		// Color - use m.from and friends
		switch m.fn {
		case fself:
			feed.WritefEscaped("[%s](grey): ", m.from)
		case fbuffer:
			feed.WritefEscaped("[%s](blue): ", m.from)
		case faction:
			feed.WritefEscaped(" * [%s](blue) ", m.from)
		case fhighlight:
			feed.WritefEscaped("[%s](red): ", m.from)
		}
		feed.WriteStringEscaped(m.data + "\n")
		return
	case fserver:
		w = c.MainWriter("server", "feed")
	case fstatus:
		w = c.StatusWriter(m.buff)
	case fsidebar:
		w = c.SideWriter(m.buff)
	case ftitle:
		w = c.TitleWriter(m.buff)
	}
	if w == nil {
		return
	}
	cleaner := cleanmark.NewCleaner(w)
	defer cleaner.Close()
	// if m.from write it
	cleaner.WriteStringEscaped(m.data + "\n")
}
