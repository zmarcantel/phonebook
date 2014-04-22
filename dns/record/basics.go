package record

import (
    "fmt"
    "time"
    "bytes"
)

const (
    A_RECORD uint16        = 1
    AAAA_RECORD uint16     = 28
    SRV_RECORD uint16      = 33
)

var TypeIntToString = map[uint16]string {
    A_RECORD:           "A",
    AAAA_RECORD:        "AAAA",
    SRV_RECORD:         "SRV",
}

//----------------------------------------------
// Record Header Structures
//----------------------------------------------

type RecordHeader struct {
    Name            string
    Type            uint16
    Class           uint16
    TTL             time.Duration       // will be converted into 32 bit unsigned int representing seconds
    RDataLength     uint16
}

func (self *RecordHeader) String() string {
    return fmt.Sprintf("%s: %s\n\tClass: %d\n\tType: %d\n\tTTL: %+v\n\tLength: %d\n", TypeIntToString[self.Type], self.Name, self.Class, self.Type, self.TTL, self.RDataLength)
}

func (self *RecordHeader) Serialize() ([]byte, error) {
    var result = make([]byte, 0)
    var buffer = bytes.NewBuffer(result)

    label, err := CreateMessageLabel(self.Name)
    if err != nil { return nil, err }
    buffer.Write(label)

    buffer.Write(Uint16ToBytes(self.Type))
    buffer.Write(Uint16ToBytes(self.Class))
    buffer.Write(Uint32ToBytes(uint32(self.TTL.Seconds())))
    buffer.Write(Uint16ToBytes(self.RDataLength))

    return buffer.Bytes(), nil
}

//----------------------------------------------
// Basic Record Structures
//----------------------------------------------

type Record interface {
    Header()        *RecordHeader
    Print(indent string)
    Data()          ([]byte, error)
    Serialize()     ([]byte, error)
}

type RawRecord struct {
    RecordHeader
    Data            []byte
}

//----------------------------------------------
// Record Array
//----------------------------------------------

type RecordCollection []Record

func (self RecordCollection) Serialize() ([]byte, error) {
    if len(self) < 1 {
        return nil, nil
    }

    var result = make([]byte, 0)
    var buffer = bytes.NewBuffer(result)

    for _, rec := range self {
        data, err := rec.Serialize()
        if err != nil { return nil, err }
        buffer.Write(data)
    }

    return buffer.Bytes(), nil
}

func (self RecordCollection) Print(indent string) {
    for _, a := range self {
        a.Print(indent)
    }
}
