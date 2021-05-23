package main

import (
	//"bufio"
	"fmt"
	"net"
	"os"
	//"strconv"
	"encoding/json"
	"time"
	//"strings"
)

// state behöver en mutex

//var clients map[*Client]Data
var players map[string]Position
var clients map[*Client]*Data // ska man lägga position här?
var state *Data

type commandID int



const (
  CMD_LEFT commandID = iota
  CMD_RIGHT
  CMD_UP
  CMD_DOWN
  CMD_QUIT
)

type Client struct {
	conn     net.Conn
	addr     string
	outgoing_ch chan Data
	Pos 	 Position
}

type Position struct {
	X, Y int
}

type Command struct {
	Id           commandID
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

	/* These could be done as pointers but I am lazy */
	clients = make(map[*Client]*Data)
	players = make(map[string]Position) 

	init := &Data{
		Players: players, // blir det här en pointer?
	}
	state = init
  
  	go graph2Client() // denna kan heta tick
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

  fmt.Println("Hello :)")
  fmt.Println(listener.Addr())

  /* Add new clients and make sure each are 
  recieving and sending data */
  for {
		conn, err := listener.Accept()

		if err != nil {
			fmt.Println(err.Error())
		}
		
		client := &Client{
			conn:     conn,
			addr:     conn.RemoteAddr().String(),
			outgoing_ch: make(chan Data), // borde inte denna ligga utanför for loopen? varför göra en ny varje gång
			Pos: Position{X: 10, Y: 10},
		}

		// Detta kommer loopa i all oändlighet så appendar
		// https://stackoverflow.com/questions/10485743/contains-method-for-a-slice
		// vill egenltigen ha map, mnen då måste nyckeln vara string för att kunna skickas över json
		// kan lägga in in map innan som ser till att det endas kör detta ifall clienten redan ligger i 
		if _, ok := clients[client] ; !ok {
		state.Players[client.addr] = client.Pos /* Add this players init pos to current state */
		clients[client] = state
		}
		
// rad 97
// Den här behövs egentligen inte, det sköts väl av graph2Client?
/*
		for k, v := range players {
			c.outgoing_ch <- Data{
				Players: players,
			}
			fmt.Println("A client joined!")
		}
*/
		fmt.Println("A client joined!")
		/* Starting routines for send & recieve */
		go SendData(client)
    	go ListenClient(client)
	}
}

func ListenClient(client *Client){
loop:
  for client != nil{
	decoder := json.NewDecoder(client.conn)
	pos := players[client.addr]
	//var newState Data // tänker att man egenligen vill update(state) typ...
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

			//client.Mutex.Lock()
			state.Players[client.addr] = pos
			//client.Mutex.Unlock()
			//client.outgoing_ch <- newState
		
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

/* TODO varför måste jag skicka in *s? */
/* Pushes the state to all clients every 100 ms */
func graph2Client(){
	ticker := time.NewTicker(100* time.Millisecond)
	defer ticker.Stop()
	for _ = range ticker.C {
		for c, _ := range clients {
		//fmt.Println("Sent graphic!")
		c.outgoing_ch <- *state // hmm detta lär krångla my man
		}
	}
}