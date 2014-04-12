package record

import (
    "fmt"
    "net"
    "time"
    "errors"
)

//----------------------------------------------
//  A Record
//      Hostname -> IPV4
//----------------------------------------------

type ARecord struct {
    RecordHeader
    IP              net.IP
}
func (self *ARecord) Header() *RecordHeader {
    return &RecordHeader {
        self.Name,
        self.Class,
        self.Type,
        self.TTL,
        self.RDataLength,
    }
}
func (self *ARecord) Data() []byte { return []byte(self.IP)[len(self.IP) - 4: ] }
func (self *ARecord) Print() {
    fmt.Printf("%s\tIP: %+v\n", self.Header().String(), self.IP)
}

func A(hostname string, ttl time.Duration, target net.IP) (*ARecord, error) {
    if len(hostname) <= 0 {
        return nil, errors.New(fmt.Sprintf("The record must contain a hostname. Received: '%s'.", hostname))
    } else if ttl.Seconds() < 5 { // TODO: get actual max class int
        return nil, errors.New(fmt.Sprintf("TTL of <5s is not supported. Received: %d", ttl.Seconds))
    }
    // TODO: add checks on target -- are we remapping the current IP and some other security stuff

    return &ARecord{
        RecordHeader{
            hostname,
            uint16( 1 ),                 // 'IN' class
            1,                           // A Records are class 1
            ttl,
            4,                           // A Records send the target IP in a 4-octet data section
        },
        target,
    }, nil
}



//----------------------------------------------
//  AAAA Record
//      Hostname -> IPV6
//----------------------------------------------

type AAAARecord struct {
    RecordHeader
    IP              net.IP
}
func (self *AAAARecord) Header() *RecordHeader {
    return &RecordHeader {
        self.Name,
        self.Class,
        self.Type,
        self.TTL,
        self.RDataLength,
    }
}
func (self *AAAARecord) Data() []byte { return []byte(self.IP)[len(self.IP) - 16: ] }
func (self *AAAARecord) Print() {
    var header = self.Header()
    fmt.Printf("AAAA: %s\n\tIP: %+v\n\tClass: %d\n\tTTL: %+v\n", header.Name, self.IP, header.Class, header.TTL)
}

func AAAA(hostname string, ttl time.Duration, target net.IP) (*AAAARecord, error) {
    if len(hostname) <= 0 {
        return nil, errors.New(fmt.Sprintf("The record must contain a hostname. Received: '%s'.", hostname))
    } else if ttl.Seconds() < 5 { // TODO: get actual max class int
        return nil, errors.New(fmt.Sprintf("TTL of <5s is not supported. Received: %d", ttl.Seconds))
    }
    // TODO: add checks on target -- are we remapping the current IP and some other security stuff

    return &AAAARecord{
        RecordHeader{
            hostname,
            1,                           // 'IN' class
            28,                          // AAA Records are type 28
            ttl,
            16,                          // AAAA Records send the target IP in a 16-octet data section
        },
        target,
    }, nil
}
