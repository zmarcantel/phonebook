package main

import (
    "os"
    "fmt"
    "net"
    "time"
    "os/signal"

    "github.com/zmarcantel/phonebook/server"
    "github.com/zmarcantel/phonebook/dns/record"

)

var lock = make(chan error)

func main() {
    var a, err = record.A("zed.io", 10 * time.Second, net.ParseIP("127.0.0.1"))
    if err != nil { panic(err) }
    server.AddRecord(a)

    aaaa, err := record.AAAA("zed.io", 10 * time.Second, net.ParseIP("::1"))
    if err != nil { panic(err) }
    server.AddRecord(aaaa)

    srv, err := record.SRV("_test._tcp.zed.io", "zed.io", 10 * time.Second, 5, 5, 8053)
    if err != nil { panic(err) }
    server.AddRecord(srv)

    watchSignals(lock)
    server.Start(lock)

    err = <-lock
    die(err)
}

func die(err error) {
    if err == nil {
        os.Exit(0)
    }

    fmt.Printf("ERROR: %s\n", err)
}

func watchSignals(done chan error) {
    var sigint = make(chan os.Signal, 1)
    signal.Notify(sigint, os.Interrupt)

    var sigkill = make(chan os.Signal, 1)
    signal.Notify(sigkill, os.Interrupt)

    go func(){
        var handled = false
        for {
            select {
                case sig := <-sigint:
                    if handled { break }
                    fmt.Printf("\nReceived interrupt signal: %s\n", sig)
                    handled = true
                    done <- nil
                    break

                case sig := <-sigkill:
                    if handled { break }
                    fmt.Printf("\nReceived termination signal: %s\n", sig)
                    handled = true
                    done <- nil
                    break
            }
        }
    }()
}
