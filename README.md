# golib-tcpsender
Go Library for a tcp sender to publish messages over TCP, asynchron and without need of listening consumers


Example with 2 independent sender (one non blocking go-routine, sending heartbeats, the other just sending one message. Library is intended to provide heartbeats and a stream of log messages via tcp.

```go
package main

import (
    tcpsender "github.com/hanneslehmann/golib-tcpsender"
  "time"
)

func main() {
   sender_one:=tcpsender.New("localhost",6000)
   sender_two:=tcpsender.New("localhost",6001)

   go sender_one.StartAndListen()
   go sender_two.StartAndListen()

   sender_one.HeartBeat("- - HeartBeat\n", 1000)
   time.Sleep(4000 * time.Millisecond)
   sender_two.SendMessage("Hello")
   for {}
}

```

This is only a proof of concept. Library has to be cleaned up and hardened
