package main

import (
        "fmt"
        "net"
        "os"
        "encoding/json"
        "github.com/rthornton128/goncurses"
        "log"
        //"strings"
)

type commandID int

const (
  CMD_LEFT commandID = iota
  CMD_RIGHT
  CMD_UP
  CMD_DOWN
  CMD_QUIT
)

type Position struct {
  X, Y int
}

type Command struct {
  Id           commandID
}

type Data struct {
  Players	     map[string]Position
}


func main() {
    arguments := os.Args
    if len(arguments) == 1 {
            fmt.Println("Please provide host:port.")
            return
  }
    
    stdscr, err := goncurses.Init()
    if err != nil {
      log.Fatal("init", err)
  }
  
  defer goncurses.End()
  
  stdscr.MovePrint(3, 0, "Welcome to Roguelike. You are the '@'")
  

  // Turns off echo of characters and cursor
  goncurses.Echo(false)
  goncurses.Cursor(0)
	// Refresh() flushes output to the screen. Internally, it is the same as
	// calling NoutRefresh() on the window followed by a call to Update()
  stdscr.Refresh()

  CONNECT := arguments[1]
  stdscr.MovePrint(5, 0, CONNECT)
  connection, err := net.Dial("tcp", CONNECT)
  if err != nil {
          fmt.Println(err)
          return
  }

  stdscr.MovePrint(5,0,"Connection to server has been established")
  stdscr.Refresh()

  /* Channels and Go Routines */
  outgoing_ch := make(chan Command)
  incoming_ch := make(chan Data)

  go ReceiveData(connection, incoming_ch)
  go SendData(connection, outgoing_ch)
  
  /* Listen to input */
  go func(){
    encoder := json.NewEncoder(connection)
loop:
      for{
        switch  byte(stdscr.GetChar()){
        case 'a':
            outgoing_ch<- Command{
              Id:   CMD_LEFT,
            }
      	case 'd':
            outgoing_ch <- Command{
              Id:   CMD_RIGHT,
            }
      	case 'w':
            outgoing_ch <- Command{
              Id:   CMD_UP,
            }
      	case 's':
            outgoing_ch <- Command{
              Id:   CMD_DOWN,
            }
        case 'q':
            encoder.Encode(Command{
              Id:   CMD_QUIT,
            })
            break loop
        }
      }
    }()

    /* Act on incoming data */
    /* This is where the lazy drawing happens */
  go func(){
      for data := range incoming_ch {
        //fmt.Println(data)
        for _, v := range data.Players {
          stdscr.MoveAddChar(v.Y, v.X, '@')
          stdscr.Refresh()
        }
      }
    }()
  select{}
}

func SendData(connection net.Conn, outgoing_ch <-chan Command){
  encoder := json.NewEncoder(connection)
  for cmd := range outgoing_ch {
      encoder.Encode(cmd)
  }
}

func ReceiveData(connection net.Conn, incoming_ch chan<- Data){
  for {
      d := json.NewDecoder(connection)
      var data Data
      err := d.Decode(&data)
      if err == nil {
            // fmt.Println(data)
            incoming_ch <- data
      } else {
            fmt.Println(err)
            fmt.Println("Connection to server lost, the only solution is to restart.")
            break
      }
    }
}
