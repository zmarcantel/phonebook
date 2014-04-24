package record

import (
    "fmt"
    "time"
    "bytes"
    "errors"
)

//----------------------------------------------
//  PTR Record
//      IP -> Hostname
//----------------------------------------------

type PTRRecord struct {
    RecordHeader
    Target                  string
}

//
// Print the record to stdout (convenience function)
//
func (self *PTRRecord) Print(indent int) {
    var indentString string
    for i := 0 ; i < indent; i++ { indentString += "\t" }

    fmt.Printf("%sPTR:\n", indentString)
    fmt.Printf("%s\t   Label: %s\n", indentString, self.Name)
    fmt.Printf("%s\t     TTL: %+v\n", indentString, self.TTL)
    fmt.Printf("%s\t  Target: %+v\n", indentString, self.Target)
}

//
// Return the record type
//
func (self *PTRRecord) GetType() uint16 {
    return self.Type
}

//
// Return the record label
//
func (self *PTRRecord) GetLabel() string {
    return self.Name
}

//
// Return (serialized) any data that affect the record's "Data Length" property
//
func (self *PTRRecord) Data() ([]byte, error) {
    var result = make([]byte, 0)
    var buffer = bytes.NewBuffer(result)

    label, err := CreateMessageLabel(self.Target)
    if err != nil { return nil, err }
    buffer.Write(label)

    return buffer.Bytes(), nil
}

//
// Translate the record into a byte array to be placed in a DNS packet
//
func (self *PTRRecord) Serialize() ([]byte, error) {
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
// Create a CNAME record given the name, target, and TTL
//
func PTR(name, target string, ttl time.Duration) (*PTRRecord, error) {
    if len(name) <= 0 {
        return nil, errors.New("An IP, zone, or combo is required.")
    } else if len(target) <= 0 {
        return nil, errors.New("The record must contain a target hostname.")
    } else if ttl.Seconds() < 5 { // TODO: get actual max class int
        return nil, errors.New(fmt.Sprintf("TTL of <5s is not supported. Received: %d", ttl.Seconds))
    }

    var result = &PTRRecord{
        RecordHeader{
            Name:        name,
            Type:        PTR_RECORD,
            Class:       uint16( 1 ),                 // 'IN' class
            TTL:         ttl,
        },
        target,
    }

    // serialize to catch errors
    _, err := result.Data()
    return result, err
}
