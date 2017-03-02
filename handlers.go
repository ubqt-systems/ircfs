package main

import (
	"bytes"
	"fmt"
)

// handleInput - append valid runes to input type, curtail input at [history]input lines.
func (st *State) handleInput(data []byte, client string) (int, error) {
	// Strip out initial forward slash of command, test for literal slash input
	if data[0] == '/' {
		data = data[1:]
		if data[0] != '/' {
			return st.handleCtl(data, client)
		}
	}
	c := st.irc[st.clients[client].server]
	cl := st.clients[client]
	c.Commands.Message(cl.channel, string(data))
	st.input = append(st.input, data...)
	return len(data), nil
}

func (st *State) handleCtl(b []byte, client string) (int, error) {
	arr := bytes.Fields(b)
	switch string(arr[0]) {
	case "set":
		// Set for client specif
		st.handleSet(arr[1:], client)
	// Handle -server, default to current [client]
	case "q":
		st.handleMsg(arr[1:], client)
	case "msg":
		st.handleMsg(arr[1:], client)
	case "join":
		// We only need current irc connection here
		st.handleJoin(string(arr[1]), client)
	case "part":
		// We only need current irc connection here
		st.handlePart(string(arr[1]), client)
	case "buffer":
		// Buffer swapping
		st.handleBuffer(string(arr[1]), client)
	case "ignore":
		// This will be a global blacklist that we just don't log messages with, won't need client. Will just be `st.AddIgnore(b) and such
		// Store to file, such as `irc/freenode/ignore`
		st.handleIgnore(arr[1:], client)
	case "connect":

	}
	return len(b), nil
}

func (st *State) ctl(client string) ([]byte, error) {
	return []byte("part\njoin\nquit\nbuffer\nignore\n"), nil
}

func (st *State) status(client string) ([]byte, error) {
	var buf []byte
	cl := st.irc[st.clients[client].server]
	channel := cl.Lookup(st.clients[client].channel)
	if channel == nil {
		return nil, nil
	}
	//TODO: text/template to design the status bar
	buf = append(buf, '\\')
	buf = append(buf, []byte(channel.Name)...)
	buf = append(buf, []byte(channel.Modes.String())...)
	buf = append(buf, '\n')
	return buf, nil
}

func (st *State) sidebar(client string) ([]byte, error) {
	cl := st.irc[st.clients[client].server]
	channel := cl.Lookup(st.clients[client].channel)
	if channel == nil {
		return nil, nil
	}
	var buf []byte
	list := channel.NickList()
	for _, item := range list {
		buf = append(buf, []byte(item)...)
		buf = append(buf, '\n')
	}
	return buf, nil
}

func (st *State) buff(client string) ([]byte, error) {
	//TODO: Format either here, or have the logs formatted.
	//os.Open() make path based whichever current thing we're on
	return []byte("buffer file\n"), nil
}

func (st *State) title(client string) ([]byte, error) {
	cl := st.irc[st.clients[client].server]
	channel := cl.Lookup(st.clients[client].channel)
	buf := []byte(channel.Topic)
	buf = append(buf, '\n')
	return buf, nil
}

func (st *State) handlePrivmsg(server string, b []byte) {
	fmt.Println(string(b))
}
