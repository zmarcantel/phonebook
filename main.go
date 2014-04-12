package main

import (
    "fmt"
    "net"
    "time"

    "./server"
    "./dns/record"

)

func main() {
    var rec, err = record.A("zed.io", 10 * time.Second, net.ParseIP("127.0.0.1"))
    if err != nil { panic(err) }
    server.AddRecord(rec)

    var lock = make(chan error)
    server.Start(lock)

    err = <-lock
    die(err)
}

func die(err error) {
    if err == nil { return }

    fmt.Printf("ERROR: %s\n", err)
}
