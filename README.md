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
   sender:=tcpsender.New("localhost",6000)
   sender2:=tcpsender.New("localhost",6001)
   go sender.HeartBeat("hb--", 1000)
   time.Sleep(4000 * time.Millisecond)
   sender2.SendMessage("Hallo")
   for {}
}
```

This is only a prrof of concept. Library has to be cleaned up and hardened
