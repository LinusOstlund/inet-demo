package main

import (
	"fmt"
	"net"
	"os"
	"encoding/json"
	"time"
)

/* Define types and structs used between server and client */
var players map[string]Position
var clients map[*Client]*Data
var state *Data

type commandID int

/* Possible commands */
const (
  CMD_LEFT commandID = iota
  CMD_RIGHT
  CMD_UP
  CMD_DOWN
  CMD_QUIT
)

type Command struct {
	Id           commandID
  }

type Client struct {
	conn     net.Conn
	addr     string
	outgoing_ch chan Data
	Pos 	 Position
}

type Position struct {
	X, Y int
}

type Data struct {
	Players	     map[string]Position
}

/* The initial state*/

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide a port number!")
		return
	}

	PORT := ":" + arguments[1]

	clients = make(map[*Client]*Data)
	players = make(map[string]Position) 

	init := &Data{
		Players: players,
	}

	state = init
  
	/* Two parallell routines which 1) sends graphic to client and 2) set up the connection */
  	go graph2Client()
  	go ConnectionController(PORT)
  	select{}
}

func ConnectionController(PORT string){
  listener, err := net.Listen("tcp4", PORT)

  if err != nil {
    fmt.Println(err)
    return
  }

  defer listener.Close()

  fmt.Println("Welcome!")
  fmt.Println(listener.Addr())

  /* Add new clients and make sure each are recieving and sending data */
  for {
		conn, err := listener.Accept()

		if err != nil {
			fmt.Println(err.Error())
		}
		
		client := &Client{
			conn:     conn,
			addr:     conn.RemoteAddr().String(),
			outgoing_ch: make(chan Data),
			Pos: Position{X: 10, Y: 10},
		}

		if _, ok := clients[client] ; !ok {
		state.Players[client.addr] = client.Pos /* Add this player's init pos to current state */
		clients[client] = state
		}

		fmt.Println("A client joined!")

		/* Starting routines for send & recieve */
		go SendData(client)
    	go ListenClient(client)
	}
}

/* Recieve input from client and update server side */
func ListenClient(client *Client){
loop:
  for client != nil{
	decoder := json.NewDecoder(client.conn)
	pos := players[client.addr]
    var cmd Command
    err := decoder.Decode(&cmd)
    if err == nil {
		switch cmd.Id {
		case CMD_UP:
				fmt.Println("Client moved up: ", client.conn.RemoteAddr())
				pos.Y--
		case CMD_DOWN:
				fmt.Println("Client moved down: ", client.conn.RemoteAddr())
				pos.Y++
		case CMD_LEFT:
				fmt.Println("Client moved up: ", client.conn.RemoteAddr())
				pos.X--
		case CMD_RIGHT:
				fmt.Println("Client moved down: ", client.conn.RemoteAddr())
				pos.X++
		case CMD_QUIT:
				fmt.Println("Client has disconnected")
				delete(clients, client)
				client.conn.Close()
				close(client.outgoing_ch)
				client = nil
		}
			state.Players[client.addr] = pos
    } else {
          fmt.Println("ListenClient: ", err)
					break loop
    	}
	}

}

func SendData(client *Client) {
	encoder := json.NewEncoder(client.conn)
	for data := range client.outgoing_ch { // waits for data to send to clients
		encoder.Encode(data)
	}
}

/* Pushes the state to all clients every 100 ms */
func graph2Client(){
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	for _ = range ticker.C {
		for c, _ := range clients {
		c.outgoing_ch <- *state
		}
	}
}