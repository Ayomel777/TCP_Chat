package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
)

type server struct {
	rooms    map[string]*room
	commands chan command
}

func newServer() *server {
	return &server{
		rooms:    make(map[string]*room),
		commands: make(chan command),
	}
}

func (s *server) run() {
	for cmd := range s.commands {
		switch cmd.id {
		case CMD_NICK:
			s.nick(cmd.client, cmd.args)
		case CMD_JOIN:
			s.join(cmd.client, cmd.args)
		case CMD_ROOMS:
			s.listRooms(cmd.client)
		case CMD_MSG:
			s.msg(cmd.client, cmd.args)
		case CMD_QUIT:
			s.quit(cmd.client, cmd.args)
		}
	}
}
func (s *server) newClient(conn net.Conn) {
	log.Printf("new client has connected: %s", conn.RemoteAddr().String())

	c := &client{
		conn:     conn,
		nick:     "anonymous",
		commands: s.commands,
	}
	c.readInput()
}

func (s *server) nick(c *client, args []string) {
	if len(args) == 2 {
		if strings.TrimSpace(args[1]) == "" {
			c.err(fmt.Errorf("nick cant be empty\n"))
			return
		}
		c.nick = args[1]
		c.msg(fmt.Sprintf("your nick is %s", c.nick))
	} else {
		c.err(fmt.Errorf("wrong command format! /nick <nickname>"))
		return
	}
}

func (s *server) join(c *client, args []string) {
	if len(args) == 2 {
		if strings.TrimSpace(args[1]) == "" {
			c.err(fmt.Errorf("room name cant be empty or consist only of spaces!\n"))
			return
		}
		roomName := args[1]
		r, ok := s.rooms[roomName]
		if !ok {
			r = &room{
				name:    roomName,
				members: make(map[net.Addr]*client),
			}
			s.rooms[roomName] = r
		}

		r.members[c.conn.RemoteAddr()] = c

		s.quitCurrentRoom(c)
		c.room = r

		r.broadcast(c, fmt.Sprintf("%s has joined the room", c.nick))
		c.msg(fmt.Sprintf("welcome to %s", r.name))
	} else {
		c.err(fmt.Errorf("wrong command format! /join <room_name>"))
		return
	}
}

func (s *server) listRooms(c *client) {
	var rooms []string
	for name := range s.rooms {
		rooms = append(rooms, name)
	}
	if len(rooms) > 0 {
		c.msg(fmt.Sprintf("available rooms are: %s", strings.Join(rooms, ", ")))
	} else {
		c.msg(fmt.Sprintf("no available rooms"))
	}
}

func (s *server) msg(c *client, args []string) {
	if len(args) >= 2 {
		if c.room == nil {
			c.err(errors.New("you must join the room first"))
			return
		}
		message := strings.Trim(strings.Join(args[1:], " "), " ")
		if message == "" {
			c.err(fmt.Errorf("message cannot be empty"))
			return
		}

		c.room.broadcast(c, c.nick+": "+message)
	} else {
		c.err(fmt.Errorf("wrong command format! /msg <your message>"))
		return
	}
}

func (s *server) quit(c *client, args []string) {
	if len(args) == 1 {
		log.Printf("client has disconnected: %s", c.conn.RemoteAddr().String())

		s.quitCurrentRoom(c)

		c.msg(":(")
		c.conn.Close()
	} else {
		c.err(fmt.Errorf("wrong command format! /quit"))
	}
}

func (s *server) quitCurrentRoom(c *client) {
	if c.room != nil {
		delete(c.room.members, c.conn.RemoteAddr())
		c.room.broadcast(c, fmt.Sprintf("%s has lefted the room", c.nick))
	}
}
