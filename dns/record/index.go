package record

import (
    "fmt"
    "time"
    "bytes"
    "errors"
)

// constants representing record type values
const (
    A_RECORD uint16        = 1
    AAAA_RECORD uint16     = 28
    SRV_RECORD uint16      = 33
    CNAME_RECORD uint16    = 5
    PTR_RECORD uint16      = 12
    MX_RECORD uint16       = 15
    TXT_RECORD uint16      = 16
)


// map from record type value to string
var TypeIntToString = map[uint16]string {
    A_RECORD:           "A",
    AAAA_RECORD:        "AAAA",
    SRV_RECORD:         "SRV",
    CNAME_RECORD:       "CNAME",
    PTR_RECORD:         "PTR",
    MX_RECORD:          "MX",
    TXT_RECORD:         "TXT",
}

var ErrInvalidIP = errors.New("Invalid IP type for record")

//----------------------------------------------
// Record Header Structures
//----------------------------------------------

// All (most) records contain these fields
// They are given to a Record type and are available during initializion and serialization
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

//----------------------------------------------
// Basic Record Structures
//----------------------------------------------

//
// Records must be able to:
//    1. print themself
//    2. retrieve type and/or label
//    3. retreive "non-header data" (data affecting self.RDataLength)
//    4. fully serialize itself
//
type Record interface {
    Print(indent int)

    GetType()       uint16
    GetLabel()      string

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

// Created for convenience. RecordCollection is able to serialize the included records
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

func (self RecordCollection) Print(indent int) {
    for _, a := range self {
        a.Print(indent)
    }
}
