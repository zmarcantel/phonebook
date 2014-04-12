package record

import (
    "fmt"
    "time"
    "bytes"
    "encoding/binary"
)

//----------------------------------------------
// Basic Record Structures
//----------------------------------------------

type RecordHeader struct {
    Name            string
    Class           uint16
    Type            uint16
    TTL             time.Duration       // will be converted into 32 bit unsigned int representing seconds
    RDataLength     uint16
}

func (self *RecordHeader) String() string {
    return fmt.Sprintf("A: %s\n\tClass: %d\n\tTTL: %+v\n", self.Name, self.Class, self.TTL)
}

type Record interface {
    Header()        *RecordHeader
    Print()
    Data()          []byte
}

type RawRecord struct {
    RecordHeader
    Data            []byte
}

type RecordCollection []Record

func (self RecordCollection) Serialize() ([]byte, error) {
    var result = make([]byte, 0)
    var buffer = bytes.NewBuffer(result)

    for _, rec := range self {
        var header = rec.Header()
        var data = rec.Data()
        var label, err = CreateMessageLabel(header.Name)
        if err != nil { return nil, err }

        buffer.Write(label)
        buffer.Write([]byte{byte(0)})
        binary.Write(buffer, binary.BigEndian, header.Type)
        binary.Write(buffer, binary.BigEndian, header.Class)
        binary.Write(buffer, binary.BigEndian, uint32(header.TTL.Seconds()))

        binary.Write(buffer, binary.BigEndian, uint16(len(data)))
        buffer.Write(data)
    }

    return bytes.TrimSpace(buffer.Bytes()), nil
}

func (self RecordCollection) Print(indent string) {
    for _, a := range self {
        fmt.Printf("%sName: %s\n", indent, a.Header().Name)
        fmt.Printf("%s\tClass: %d\n", indent, a.Header().Class)
        fmt.Printf("%s\t Type: %d\n", indent, a.Header().Type)
    }
}
