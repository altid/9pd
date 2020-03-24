package ninep

import (
	"os"

	"github.com/altid/server/client"
	"github.com/altid/server/command"
	"github.com/altid/server/tail"
)

func (s *service) listenCommands(fp *os.File) {
	defer fp.Close()

	for cmd := range s.command {
		c := s.client.Client(cmd.UUID)
		old := c.Current()

		switch cmd.CmdType {
		case command.OtherCmd:
			cmd.WriteOut(fp)
		case command.OpenCmd:
			c.SetBuffer(cmd.Args[0])
			cmd.WriteOut(fp)
			s.update(c, old)
		case command.BufferCmd:
			c.SetBuffer(cmd.Args[0])
			s.update(c, old)
		case command.CloseCmd:
			// Pop back to the last buffer
			history := c.History()
			if len(history) < 1 {
				c.SetBuffer("none")
			} else {
				c.SetBuffer(history[len(history)-1])
			}
			s.update(c, old)
			cmd.WriteOut(fp)
		case command.LinkCmd:
			c.SetBuffer(cmd.Args[1])
			cmd.WriteOut(fp)
			s.update(c, old)
		case command.ReloadCmd:
			// TODO (halfwit): We want to recreate everything but save our client connections
			// possibly we'll be loading more services, etc
		case command.QuitCmd:
			cmd.WriteOut(fp)
		}
	}
}

// We need to send feed commands at very least
func (s *service) listenEvents() {
	for e := range s.events {
		t := s.tabs.Tab(e.Name)
		s.debug("event name=\"%s\" service=\"%s\"", e.Name, e.Service)

		switch e.EventType {
		case tail.NotifyEvent:
			t.Alert = true
		case tail.DocEvent:
		case tail.FeedEvent:
		default:
		}

		if !t.Active {
			t.Unread++
		}

		go s.sendFeeds(e)
	}
}

// sendFeeds loops through all connected clients
// If they're listening to the current buffer try to send an event
// Shortcut to a bool comparison if inactive
func (s *service) sendFeeds(e *tail.Event) {
	for _, c := range s.client.List() {
		if c.Active && c.Current() == e.Name {
			s.debug("feed name=\"%s\" id=\"%d\"", e.Name, c.UUID)
			s.feed.Send(c.UUID)
		}
	}
}

func (s *service) update(c *client.Client, old string) {
	s.tabs.Tab(old).Active = false

	tab := s.tabs.Tab(c.Current())
	tab.Active = true
	tab.Unread = 0

	s.feed.Done(c.UUID)
}
