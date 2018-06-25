package main

import (
	"bufio"
	"log"
	"path"
	"os"
	"time"
)

// Act on it!
// TODO: Make sure we have enough context to act on correct server
func (st *State) Control(b []byte) {
	switch b[0] {
	// Join
	// Part
	// Quit
	// Reconnect
	// Connect
	// TODO: ???
	}
}

func (st *State) CtlLoop() {
	filePath := path.Join(*inPath, "ctl")
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, 0666)
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}
	buffer := bufio.NewReader(f)
	// Cheapo epoll
	for {
		bytes, _, _ := buffer.ReadLine()
		if len(bytes) != 0 {
			st.Control(bytes)
		}
		time.Sleep(50 * time.Millisecond)
	}
}
