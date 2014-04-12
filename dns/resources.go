package dns

import (
    "bytes"
    "encoding/binary"

    "github.com/zmarcantel/phonebook/dns/record"
)

func UnpackResources(source []byte, count int) ([]*record.RawRecord, int) {
    var result = make([]*record.RawRecord, count)

    var offset int
    for i := 0 ; i < count ; i++ {
        name, offset := GetMessageLabel(source)

        var rType, rClass, rLength uint16
        var rTTL uint32

        binary.Read(bytes.NewReader(source[offset:offset+2]), binary.BigEndian, &rType)
        offset += 2
        binary.Read(bytes.NewReader(source[offset:offset+2]), binary.BigEndian, &rClass)
        offset += 2
        binary.Read(bytes.NewReader(source[offset:offset+4]), binary.BigEndian, &rTTL)
        offset += 4
        binary.Read(bytes.NewReader(source[offset:offset+2]), binary.BigEndian, &rLength)
        offset += 2

        var rData = source[offset:offset + int(rLength)]

        result[i] = &record.RawRecord{
            Name:            []byte(name),
            Type:            rType,
            Class:           rClass,
            TTL:             rTTL,
            Length:          rLength,
            Data:            rData,
        }
    }

    return result, offset
}
