package record

import (
    "fmt"
    "time"
    "bytes"
    "errors"
)

//----------------------------------------------
//  TXT Record
//      Hostname -> Text Data
//----------------------------------------------

type TXTRecord struct {
    RecordHeader
    Text            string
}

//
// Print the record to stdout (convenience function)
//
func (self *TXTRecord) Print(indent int) {
    var indentString string
    for i := 0 ; i < indent; i++ { indentString += "\t" }

    fmt.Printf("%sTXT:\n", indentString)
    fmt.Printf("%s\tLabel: %s\n", indentString, self.Name)
    fmt.Printf("%s\t  TTL: %+v\n", indentString, self.TTL)
    fmt.Printf("%s\t Text: %+v\n", indentString, self.Text)
}

//
// Return the record type
//
func (self *TXTRecord) GetType() uint16 {
    return self.Type
}

//
// Return the record label
//
func (self *TXTRecord) GetLabel() string {
    return self.Name
}

//
// Return (serialized) any data that affect the record's "Data Length" property
//
func (self *TXTRecord) Data() ([]byte, error) {
    return []byte(self.Text), nil
}

//
// Translate the record into a byte array to be placed in a DNS packet
//
func (self *TXTRecord) Serialize() ([]byte, error) {
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
// Create an TXT record given the hostname, TTL, and target IP
//
func TXT(hostname string, ttl time.Duration, text string) (*TXTRecord, error) {
    if len(hostname) <= 0 {
        return nil, errors.New(fmt.Sprintf("The record must contain a hostname. Received: '%s'.", hostname))
    } else if ttl.Seconds() < 5 { // TODO: get actual max class int
        return nil, errors.New(fmt.Sprintf("TTL of <5s is not supported. Received: %d", ttl.Seconds))
    } else if len(text) == 0 {
        return nil, errors.New("Cannot store emtpy TXT record")
    }

    // TODO: add checks on target -- are we remapping the current IP and some other security stuff

    // TODO: verify text data length fits into 16bits
    return &TXTRecord{
        RecordHeader{
            hostname,
            TXT_RECORD,
            uint16( 1 ),                 // 'IN' class
            ttl,
            uint16(len(text)),
        },
        text,
    }, nil
}
