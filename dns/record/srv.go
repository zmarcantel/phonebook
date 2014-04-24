package record

import (
    "fmt"
    "time"
    "bytes"
    "errors"
    "encoding/binary"
)

//----------------------------------------------
//  A Record
//      Hostname -> IPV4
//----------------------------------------------

type SRVRecord struct {
    RecordHeader
    Priority                uint16
    Weight                  uint16
    Port                    uint16
    Target                  string
}

//
// Print the record to stdout (convenience function)
//
func (self *SRVRecord) Print(indent int) {
    var indentString string
    for i := 0 ; i < indent; i++ { indentString += "\t" }

    fmt.Printf("%sSRV:\n", indentString)
    fmt.Printf("%s\t   Label: %s\n", indentString, self.Name)
    fmt.Printf("%s\t     TTL: %+v\n", indentString, self.TTL)
    fmt.Printf("%s\tPriority: %d\n", indentString, self.Priority)
    fmt.Printf("%s\t  Weight: %d\n", indentString, self.Weight)
    fmt.Printf("%s\t    Port: %d\n", indentString, self.Port)
    fmt.Printf("%s\t  Target: %+v\n", indentString, self.Target)
}

//
// Return the record type
//
func (self *SRVRecord) GetType() uint16 {
    return self.Type
}

//
// Return the record label
//
func (self *SRVRecord) GetLabel() string {
    return self.Name
}

//
// Return (serialized) any data that affect the record's "Data Length" property
//
func (self *SRVRecord) Data() ([]byte, error) {
    var result = make([]byte, 0)
    var buffer = bytes.NewBuffer(result)

    binary.Write(buffer, binary.BigEndian, self.Priority)
    binary.Write(buffer, binary.BigEndian, self.Weight)
    binary.Write(buffer, binary.BigEndian, self.Port)

    label, err := CreateMessageLabel(self.Target)
    if err != nil { return nil, err }
    buffer.Write(label)

    return buffer.Bytes(), nil
}

//
// Translate the record into a byte array to be placed in a DNS packet
//
func (self *SRVRecord) Serialize() ([]byte, error) {
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
// Create a SRV record given the name, target, TTL, priority, weight, and port
//
func SRV(name, target string, ttl time.Duration, priority, weight, port uint16) (*SRVRecord, error) {
    if len(name) <= 0 {
        return nil, errors.New("A base hostname is required.")
    } else if len(target) <= 0 {
        return nil, errors.New("The record must contain a target hostname.")
    } else if ttl.Seconds() < 5 { // TODO: get actual max class int
        return nil, errors.New(fmt.Sprintf("TTL of <5s is not supported. Received: %d", ttl.Seconds))
    }

    var result = &SRVRecord{
        RecordHeader{
            Name:        name,
            Type:        SRV_RECORD,
            Class:       uint16( 1 ),                 // 'IN' class
            TTL:         ttl,
        },
        priority,
        weight,
        port,
        target,
    }

    // serialize to catch errors
    _, err := result.Serialize()
    return result, err
}
