package server

import (
    "fmt"
    "net"
    "time"
    "errors"
    "strconv"

    "github.com/zmarcantel/phonebook/dns"
    "github.com/zmarcantel/phonebook/dns/record"

    "github.com/zmarcantel/phonebook/server/store"
)

const (
    DNSTimeout      time.Duration   = 2 * 1e8

    ERR_FORMAT      int             = 1
    ERR_INTERNAL    int             = 2
    ERR_NOEXIST     int             = 3
    ERR_NOIMPL      int             = 4
    ERR_REFUSED     int             = 5
)

var ErrShortRead    error           = errors.New("ERROR: short read")


type Server struct {
    Fatal           chan error
    Error           chan error
    Address         net.Addr
    Store           store.DNSStore
    Connection      *net.UDPConn
}



//
// Start up the server and use the given channel as a killswitch
//    Bind: string representation of interface to bind to ("127.0.0.1", "localhost", ":::1", "0.0.0.0", etc)
//    Port: port to listen for queries on
//
func Start(bind string, port int, backing store.DNSStore, die chan error) *Server {
    // get the DNS address for the DNS host
    if len(bind) == 0 || bind == "localhost" {
        fmt.Printf("Using DEFULT_HOST:127.0.0.1 based on binding input: %s\n", bind)
        bind = "127.0.0.1"
    }

    if port <= 0 {
        fmt.Printf("Using DEFULT_PORT:53 based on port input: %d\n", port)
        port = 53
    }

    // check if the store option is nil
    // if so, that means default storage agent -- MapStore
    if backing == nil {
        fmt.Println("No backing given -- using default MapStore")
        backing = store.Map()
    }

    addr, err := net.ResolveUDPAddr("udp", bind + ":" + strconv.Itoa(port))
    if err != nil {
        // there was an error parsing binding address
        // pass the error and abort
        die <- err
        return nil
    }

    // create a UDP "listener"
    // simply waits on a channel for packets
    conn, err := net.ListenUDP("udp", addr)
    if err != nil {
        // if there was an error, abort
        die <- err
        return nil
    }

    // make the server we will return
    var result = &Server{
        die,
        make(chan error),
        addr,
        backing,
        conn,
    }

    // start watching for errors
    go result.WatchErrors()

    // get the server listening before we return (convenience)
    // do the listening in a goroutine
    go result.Listen()

    return result
}

//
// Shorthand for listening on localhost:53
//
func Local(backing store.DNSStore, die chan error) *Server {
    return Start("127.0.0.1", 53, backing, die)
}


//
// Listen on the net.UDPConn for incoming packets
// Responsible for intake only
//
func (self *Server) Listen() {
    // announce the listener
    fmt.Printf("DNS Server listening on: %s\n", self.Connection.LocalAddr())

    // defer closing the connection until the below for loop exits
    // only happens on kill/int signal, error, etc
    defer self.Connection.Close()

    // round and round it goes, when it stops, only the program knows!!
    for {
        // make a 512 byte buffer (max DNS packet size as per RFC 1035)
        var content = make([]byte, 512)

        // read our packet into the buffer
        var readLength, addr, err = self.Connection.ReadFromUDP(content)
        if err != nil {
            // report the issue if it exists
            self.Fatal <- err
        }
        if readLength == 0 {
            // got a short read... abort
            // TODO: send error response
            self.Error <- err
        }

        // trim of any buffer fat and respond in an isolated goroutine
        content = content[:readLength]
        go self.Serve(addr, content)
    }

    // report the escaping for the foor loop without an error on the channel
    fmt.Println("DNS passed on connection.")
}

//
// Take in a DNS query and respond with a Record, [], or error code
// Runs in isolated/concurrent thread
//
func (self *Server) Serve(addr net.Addr, query []byte) {
    // TODO: catch and respond to packet errors
    var message = dns.UnpackMessage(query)

    // TODO: logging verbosity
    // print the request to logs
    fmt.Printf("\n\nREQUEST: %d\n", message.Header.ID)
    message.Questions.Print(1)

    // verify it's a query...
    if !message.Header.Response {
        // get the answers to the questions posed
        var answers, err = self.Answer(message.Questions)
        if err != nil {
            if err == store.ErrNotFound {
                message.Header.Rcode = ERR_NOEXIST
            } else if err == store.ErrInvalidType {
                message.Header.Rcode = ERR_NOIMPL
            } else {
                message.Header.Rcode = ERR_INTERNAL
            }
            self.Error <- err
        }

        // format the response(s) we found into a DNS packet to
        // be served to the client
        var response = generateAnswerMessage(message, answers)

        // serialize the message for wire transfer
        serialized, err := response.Serialize()
        if err != nil {
            self.Fatal <- err
        }

        // print the response to logs
        fmt.Println("\n\nRESPONSE:")
        response.Print(1)

        // write the serialized DNS packet to the address given in the request
        // this ends the cycle of the DNS request
        _, err = self.Connection.WriteTo(serialized, addr)
        if err != nil {
            fmt.Printf("ERROR: There was an error responding to request:\n%s\n\n", err)
        }
    }
}

func (self *Server) WatchErrors() {
    for {
        fmt.Println(<- self.Error)
    }
}


func generateAnswerMessage(message *dns.Message, answers []record.Record) dns.Message {
    var header = dns.MessageHeader {
        ID: message.Header.ID,
        Response: true,
        Opcode: message.Header.Opcode,
        Authoritative: true,            // TODO: set this truthfully
        Truncated: false,               // TODO: set this truthfully
        RecursionDesired: message.Header.RecursionDesired,
        RecursionAvailable: false,      // we will NEVER go to other DNS servers. RAFT baby....
        Rcode: message.Header.Rcode,    // no error
        QDCount: uint16(len(message.Questions)),
        ANCount: uint16(len(answers)),
        NSCount: 0,
        ARCount: 0,
    }

    return dns.Message{
        Header: header,
        Questions: message.Questions,
        Answers: answers,
        Ns: nil,
        Extra:nil,
    }
}
