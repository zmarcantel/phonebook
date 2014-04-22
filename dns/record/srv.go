package record

import (
    "fmt"
    "time"
    "bytes"
    "errors"
    "strings"
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

func (self *SRVRecord) Header() *RecordHeader {
    return &RecordHeader {
        Name:        self.Name,
        Type:        self.Type,
        Class:       self.Class,
        TTL:         self.TTL,
        RDataLength: self.RDataLength,
    }
}

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

func (self *SRVRecord) Print(indent string) {
    fmt.Printf("%sSRV:\n", indent)
    fmt.Printf("%s\t   Label: %s\n", indent, self.Name)
    fmt.Printf("%s\tPriority: %d\n", indent, self.Priority)
    fmt.Printf("%s\t  Weight: %d\n", indent, self.Weight)
    fmt.Printf("%s\t    Port: %d\n", indent, self.Port)
    fmt.Printf("%s\t  Target: %+v\n", indent, self.Target)
}

func (self SRVRecord) Serialize() ([]byte, error) {
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

func SRV(name, target string, ttl time.Duration, priority, weight, port uint16) (*SRVRecord, error) {
    if len(name) <= 0 {
        return nil, errors.New("A base hostname is required.")
    } else if len(target) <= 0 {
        return nil, errors.New("The record must contain a target hostname.")
    } else if ttl.Seconds() < 5 { // TODO: get actual max class int
        return nil, errors.New(fmt.Sprintf("TTL of <5s is not supported. Received: %d", ttl.Seconds))
    }

    if !strings.HasSuffix(name, ".") { name += "." }
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
    _, err := result.Data()
    return result, err
}
