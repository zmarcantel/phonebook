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

//
// Print the record to stdout (convenience function)
//
func (self *ARecord) Print(indent int) {
    var indentString string
    for i := 0 ; i < indent; i++ { indentString += "\t" }

    fmt.Printf("%sA:\n", indentString)
    fmt.Printf("%s\tLabel: %s\n", indentString, self.Name)
    fmt.Printf("%s\t  TTL: %+v\n", indentString, self.TTL)
    fmt.Printf("%s\t   IP: %+v\n", indentString, self.IP)
}

//
// Return the record type
//
func (self *ARecord) GetType() uint16 {
    return self.Type
}

//
// Return the record label
//
func (self *ARecord) GetLabel() string {
    return self.Name
}

//
// Return (serialized) any data that affect the record's "Data Length" property
//
func (self *ARecord) Data() ([]byte, error) {
    return []byte(self.IP)[len(self.IP) - 4: ], nil
}

//
// Translate the record into a byte array to be placed in a DNS packet
//
func (self *ARecord) Serialize() ([]byte, error) {
    var result = make([]byte, 0)
    var buffer = bytes.NewBuffer(result)

    label, err := CreateMessageLabel(self.Name)
    if err != nil { return nil, err }
    buffer.Write(label)

    buffer.Write(Uint16ToBytes(self.Type))
    buffer.Write(Uint16ToBytes(self.Class))
    buffer.Write(Uint32ToBytes(uint32(self.TTL.Seconds())))

    data, err := self.Data()
    if err != nil { return nil, err }

    self.RDataLength = uint16(len(data))
    buffer.Write(Uint16ToBytes(self.RDataLength))
    buffer.Write(data)

    return buffer.Bytes(), nil
}

//
// Create an A record given the hostname, TTL, and target IP
//
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

//
// Print the record to stdout (convenience function)
//
func (self *AAAARecord) Print(indent int) {
    var indentString string
    for i := 0 ; i < indent; i++ { indentString += "\t" }

    fmt.Printf("%sAAAA:\n", indentString)
    fmt.Printf("%s\tLabel: %s\n", indentString, self.Name)
    fmt.Printf("%s\t  TTL: %+v\n", indentString, self.TTL)
    fmt.Printf("%s\t   IP: %+v\n", indentString, self.IP)
}

//
// Return the record type
//
func (self *AAAARecord) GetType() uint16 {
    return self.Type
}

//
// Return the record label
//
func (self *AAAARecord) GetLabel() string {
    return self.Name
}

//
// Return (serialized) any data that affect the record's "Data Length" property
//
func (self *AAAARecord) Data() ([]byte, error) {
    return []byte(self.IP)[len(self.IP) - 16: ], nil
}

//
// Translate the record into a byte array to be placed in a DNS packet
//
func (self *AAAARecord) Serialize() ([]byte, error) {
    var result = make([]byte, 0)
    var buffer = bytes.NewBuffer(result)

    label, err := CreateMessageLabel(self.Name)
    if err != nil { return nil, err }
    buffer.Write(label)

    buffer.Write(Uint16ToBytes(self.Type))
    buffer.Write(Uint16ToBytes(self.Class))
    buffer.Write(Uint32ToBytes(uint32(self.TTL.Seconds())))

    data, err := self.Data()
    if err != nil { return nil, err }

    self.RDataLength = uint16(len(data))
    buffer.Write(Uint16ToBytes(self.RDataLength))
    buffer.Write(data)

    return buffer.Bytes(), nil
}

//
// Create an AAAA record given the hostname, TTL, and target IP
//
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
