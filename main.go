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

func main() {

    //
    // Add testing records, one of each type, until test suite built
    //
    var a, err = record.A("zed.io", 10 * time.Second, net.ParseIP("127.0.0.1"))
    if err != nil { panic(err) }
    server.AddRecord(a)

    aaaa, err := record.AAAA("zed.io", 10 * time.Second, net.ParseIP("::1"))
    if err != nil { panic(err) }
    server.AddRecord(aaaa)

    srv, err := record.SRV("_test._tcp.zed.io", "zed.io", 10 * time.Second, 5, 5, 8053)
    if err != nil { panic(err) }
    server.AddRecord(srv)

    cname, err := record.CNAME("app.production", "zed.io", 10 * time.Second)
    if err != nil { panic(err) }
    server.AddRecord(cname)

    fmt.Printf("\n%+v\n\n", server.GetCache())


    //
    // Setup the signal handlers and start the server
    // The channel serves as an unhandled exception
    //

    var lock = make(chan error, 10)
    watchSignals(lock)
    server.Start("127.0.0.1", 53, lock)
    // shorthand for the above would be "server.Local(lock)"

    // wait for either unhandled exception or nil (signal)
    err = <-lock
    die(err)
}



//
// Responds to an error or signal being put on the server-lock
//
func die(err error) {
    // signals put nil on the channel, so ignore those
    if err != nil {
        fmt.Printf("ERROR: %s\n", err)
    }

    // exit
    os.Exit(0)
}


//
// Defines handlers for OS signals
//
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
