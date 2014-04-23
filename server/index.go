package server

import (
    "net"
    "fmt"
    "time"
    "errors"

    "github.com/zmarcantel/phonebook/dns"
    "github.com/zmarcantel/phonebook/dns/record"
)

const (
    DNSTimeout      time.Duration   = 2 * 1e8
    DEFAULT_HOST    string          = "127.0.0.1:domain"
)

var ErrShortRead    error           = errors.New("short read")


//
// Start up the server and use the given channel as a killswitch
//
func Start(die chan error) {
    // get the DNS address for the DNS host
    // TODO: parameterize the address, not just localhost
    addr, err := net.ResolveUDPAddr("udp", DEFAULT_HOST)
    if err != nil {
        // if there was an error, abort
        die <- err
        return
    }

    // create a UDP "listener"
    // simply waits on a channel for packets
    conn, err := net.ListenUDP("udp", addr)
    if err != nil {
        // if there was an error, abort
        die <- err
        return
    }

    // do the listening in a goroutine
    go Listen(conn, die)
}

func Listen(conn *net.UDPConn, die chan error) {
    // announce the listener
    fmt.Printf("DNS Server listening on: %s\n", conn.LocalAddr())
    
    // defer closing the connection until the below for loop exits
    // only happens on kill/int signal, error, etc
    defer conn.Close()

    // round and round it goes, when it stops, only the program knows!!
    for {
        // make a 512 byte buffer (max DNS packet size as per RFC 1035)
        var content = make([]byte, 512)

        // read our packet into the buffer
        var readLength, addr, err = conn.ReadFromUDP(content)
        if err != nil {
            // report the issue if it exists
            die <- err
        }
        if readLength == 0 {
            // got a short read... abort
            // TODO: send error response, don't panic
            die <- err
        }

        // trim of any buffer fat and respond in an isolated goroutine
        content = content[:readLength]
        go Serve(conn, addr, content)
    }

    // report the escaping for the foor loop without an error on the channel
    fmt.Println("DNS passed on connection.")
}

//
// Take in a DNS query and respond with a Record, [], or error code
//
func Serve(conn *net.UDPConn, addr net.Addr, query []byte) {
    // TODO: catch and respond to packet errors
    var message = dns.UnpackMessage(query)

    // TODO: logging verbosity
    // print the request to logs
    fmt.Printf("\n\nREQUEST: %d\n", message.Header.ID)
    message.Questions.Print(1)

    // verify it's a query...
    if !message.Header.Response {
        // get the answers to the questions posed
        var answers, err = AnswerQuestions(message.Questions)
        if err != nil {
            panic(err) // TODO: keep calm
        }

        // format the response(s) we found into a DNS packet to
        // be served to the client
        var response = generateAnswerMessage(message, answers)

        // serialize the message for wire transfer
        serialized, err := response.Serialize()
        if err != nil {
            panic(err) // TODO: keep calm
        }

        // print the response to logs
        fmt.Println("\n\nRESPONSE:")
        response.Print(1)

        // write the serialized DNS packet to the address given in the request
        // this ends the cycle of the DNS request
        _, err = conn.WriteTo(serialized, addr)
        if err != nil {
            fmt.Printf("ERROR: There was an error responding to request:\n%s\n\n", err)
        }
    }
}


func generateAnswerMessage(request *dns.Message, answers []record.Record) dns.Message {
    var header = dns.MessageHeader {
        ID: request.Header.ID,
        Response: true,
        Opcode: request.Header.Opcode,
        Authoritative: true,            // TODO: set this truthfully
        Truncated: false,               // TODO: set this truthfully
        RecursionDesired: request.Header.RecursionDesired,
        RecursionAvailable: false,      // we will NEVER go to other DNS servers. RAFT baby....
        Rcode: 0,                       // no error
        QDCount: uint16(len(request.Questions)),
        ANCount: uint16(len(answers)),
        NSCount: 0,
        ARCount: 0,
    }

    return dns.Message{
        Header: header,
        Questions: request.Questions,
        Answers: answers,
        Ns: nil,
        Extra:nil,
    }
}
