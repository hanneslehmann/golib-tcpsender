package tcpsender

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"
	"strings"
	"strconv"
)

type Sender struct {
   Port int
	 Host string
	 Wait int
	 msgchan chan string
	 addchan chan Client
	 rmchan chan Client
	 ln net.Listener
}

type Client struct {
	conn     net.Conn
	ch       chan string
}

func New(host string,port int) (sender *Sender){
	ln, err := net.Listen("tcp", host+":"+strconv.Itoa(port))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	msgchan := make(chan string)
	addchan := make(chan Client)
	rmchan := make(chan Client)
	go handleMessages(msgchan, addchan, rmchan)
  return &Sender {
     Port: port,
		 Host: host,
		 msgchan : msgchan,
		 addchan : addchan,
		 rmchan : rmchan,
		 ln: ln,
   }
}

func (s Sender) SendMessage(message string) {
	go sendMessage(s.msgchan, message)
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go handleConnection(conn, s.msgchan, s.addchan, s.rmchan)
	}
}


func (s Sender) HeartBeat(message string, wait int) {
	  s.Wait = wait
		go  sendHeartBeat(s.msgchan, message, s.Wait)
		for {
			conn, err := s.ln.Accept()
			if err != nil {
				fmt.Println(err)
				continue
			}
			go handleConnection(conn, s.msgchan, s.addchan, s.rmchan)
		}
}

func (c Client) WriteLinesFrom(ch <-chan string) {
	for msg := range ch {
		_, err := io.WriteString(c.conn, msg)
		if err != nil {
			return
		}
	}
}

func sendHeartBeat(msgchan chan<- string, message string, wait int) {
	for {
		t := time.Now()
    sendMessage(msgchan, message + strings.ToUpper(t.Format("20060102150405")))
    time.Sleep(time.Duration(wait) * time.Millisecond)
	}
}

func sendMessage(msgchan chan<- string, message string) {
    msgchan <- message
}

func handleConnection(c net.Conn, msgchan chan<- string, addchan chan<- Client, rmchan chan<- Client) {
	//bufc := bufio.NewReader(c)
	defer c.Close()
	client := Client{
		conn:     c,
		ch:       make(chan string),
	}
	// Register client
	addchan <- client
	defer func() {
		rmchan <- client
	}()
	// I/O
	client.WriteLinesFrom(client.ch)
}

func handleMessages(msgchan <-chan string, addchan <-chan Client, rmchan <-chan Client) {
	clients := make(map[net.Conn]chan<- string)

	for {
		select {
		case msg := <-msgchan:
			//log.Printf("New message: %s", msg)
			for _, ch := range clients {
				go func(mch chan<- string) { mch <- msg}(ch)
			}
		case client := <-addchan:
			log.Printf("New connection: %v\n", client.conn.RemoteAddr())
			clients[client.conn] = client.ch
		case client := <-rmchan:
			log.Printf("Connection closed: %v\n", client.conn.RemoteAddr())
			delete(clients, client.conn)
		}
	}
}
