package record

import (
    "fmt"
    "time"
    "bytes"
    "errors"
)

//----------------------------------------------
//  MX Record
//      IP -> Hostname
//----------------------------------------------

type MXRecord struct {
    RecordHeader
    Priority                uint16
    Target                  string
}

//
// Print the record to stdout (convenience function)
//
func (self *MXRecord) Print(indent int) {
    var indentString string
    for i := 0 ; i < indent; i++ { indentString += "\t" }

    fmt.Printf("%MX:\n", indentString)
    fmt.Printf("%s\t   Label: %s\n", indentString, self.Name)
    fmt.Printf("%s\t     TTL: %+v\n", indentString, self.TTL)
    fmt.Printf("%s\tPriority: %+v\n", indentString, self.Priority)
    fmt.Printf("%s\t  Target: %+v\n", indentString, self.Target)
}

//
// Return the record type
//
func (self *MXRecord) GetType() uint16 {
    return self.Type
}

//
// Return the record label
//
func (self *MXRecord) GetLabel() string {
    return self.Name
}

//
// Return (serialized) any data that affect the record's "Data Length" property
//
func (self *MXRecord) Data() ([]byte, error) {
    var result = make([]byte, 0)
    var buffer = bytes.NewBuffer(result)

    buffer.Write(Uint16ToBytes(self.Priority))

    label, err := CreateMessageLabel(self.Target)
    if err != nil { return nil, err }
    buffer.Write(label)

    return buffer.Bytes(), nil
}

//
// Translate the record into a byte array to be placed in a DNS packet
//
func (self *MXRecord) Serialize() ([]byte, error) {
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
// Create an MX record given the name, target, priority, and TTL
//
func MX(name, target string, priority uint16, ttl time.Duration) (*MXRecord, error) {
    if len(name) <= 0 {
        return nil, errors.New("A top-level MX name is required.")
    } else if len(target) <= 0 {
        return nil, errors.New("The MX record must contain a target mail server.")
    } else if ttl.Seconds() < 5 { // TODO: get actual max class int
        return nil, errors.New(fmt.Sprintf("TTL of <5s is not supported. Received: %d", ttl.Seconds))
    }

    var result = &MXRecord{
        RecordHeader{
            Name:        name,
            Type:        MX_RECORD,
            Class:       uint16( 1 ),                 // 'IN' class
            TTL:         ttl,
        },
        priority,
        target,
    }

    // serialize to catch errors
    _, err := result.Serialize()
    return result, err
}
