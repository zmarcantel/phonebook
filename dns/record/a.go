package record

import (
    "fmt"
    "net"
    "time"
    "bytes"
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

func (self *ARecord) Data() ([]byte, error) {
    return []byte(self.IP)[len(self.IP) - 4: ], nil
}

func (self *ARecord) Print(indent string) {
    fmt.Printf("%sHeader:\n%s\t%s\n", indent, indent, self.Header().String())
    fmt.Printf("%s\t     IP: %+v", self.IP)
}

func (self *ARecord) Serialize() ([]byte, error) {
    var result = make([]byte, 0)
    var buffer = bytes.NewBuffer(result)

    var header, err = self.Header().Serialize()
    if err != nil { return nil, err }
    buffer.Write(header)

    data, err := self.Data()
    if err != nil { return nil, err }
    buffer.Write(data)

    return buffer.Bytes(), nil
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
            A_RECORD,
            uint16( 1 ),                 // 'IN' class
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
        Name:        self.Name,
        Class:       self.Class,
        Type:        self.Type,
        TTL:         self.TTL,
        RDataLength: self.RDataLength,
    }
}

func (self *AAAARecord) Data() ([]byte, error) {
    return []byte(self.IP)[len(self.IP) - 16: ], nil
}

func (self *AAAARecord) Print(indent string) {
    fmt.Printf("%sHeader:\n%s\t%s\n", indent, indent, self.Header().String())
    fmt.Printf("%s\t     IP: %+v", self.IP)
}

func (self *AAAARecord) Serialize() ([]byte, error) {
    var result = make([]byte, 0)
    var buffer = bytes.NewBuffer(result)

    var header, err = self.Header().Serialize()
    if err != nil { return nil, err }
    buffer.Write(header)

    data, err := self.Data()
    if err != nil { return nil, err }
    buffer.Write(data)

    return buffer.Bytes(), nil
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
            Name:        hostname,
            Type:        AAAA_RECORD,
            Class:       1,                           // 'IN' class
            TTL:         ttl,
            RDataLength: 16,                          // AAAA Records send the target IP in a 16-octet data section
        },
        target,
    }, nil
}
