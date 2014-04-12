package server

import (
    "net"
    "fmt"
    "time"
    "errors"

    "../dns"
    "../dns/record"
)

const dnsTimeout time.Duration = 2 * 1e9
var dnsHost = "127.0.0.1:domain"


var ErrShortRead error = errors.New("short read")


func Start(die chan error) {
    a, e := net.ResolveUDPAddr("udp", dnsHost)
    if e != nil { panic (e) }
    conn, err := net.ListenUDP("udp", a)
    if err != nil {
        // handle error
        panic(err)
    }

    go listenOnConnection(conn, die)
    fmt.Printf("DNS Server listening on: %s\n", dnsHost)
}

func listenOnConnection(conn *net.UDPConn, die chan error) {
    defer conn.Close()
    for {
        var content = make([]byte, 512)
        var readLength, addr, err = conn.ReadFromUDP(content)
        if err != nil {
            panic(err)
        }
        if readLength == 0 { panic(ErrShortRead) } // TODO: don't panic

        content = content[:readLength]
        go serve(conn, addr, content)
    }
    panic("DNS passed on connection.")
}


func serve(conn *net.UDPConn, addr net.Addr, content []byte) {
    // TODO: catch and respond to packet errors
    var message = dns.UnpackMessage(content)

    fmt.Printf("REQUEST: %d\n", message.Header.ID)
    message.Questions.Print("\t")

    if !message.Header.Response {
        var answers, err = AnswerQuestions(message.Questions)
        if err != nil {
            panic(err) // TODO: keep calm
        }

        var response = generateAnswerMessage(message, answers)
        serialized, err := response.Serialize()
        if err != nil {
            panic(err) // TODO: keep calm
        }

        _, err = conn.WriteTo(serialized, addr)
        if err != nil {
            fmt.Printf("ERROR: There was an error responding to request:\n%s\n\n", err)
            return
        }
    }
}


func generateAnswerMessage(request *dns.Message, answers []record.Record) dns.Message {
    var header = dns.MessageHeader {
        ID: request.Header.ID,
        Response: true,
        Opcode: request.Header.Opcode,
        Authoritative: true,         // TODO: set this truthfully
        Truncated: false,            // TODO: set this truthfully
        RecursionDesired: request.Header.RecursionDesired,
        RecursionAvailable: request.Header.RecursionAvailable,
        Rcode: 0, // no error
        QDCount: 0,
        ANCount: uint16(len(answers)),
        NSCount: 0,
        ARCount: 0,
    }

    return dns.Message{
        Header: header,
        Questions: nil,
        Answers: answers,
        Ns: nil,
        Extra:nil,
    }
}
