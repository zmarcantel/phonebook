package record

import (
    "bytes"
    "strings"
    "encoding/binary"
)

func CreateMessageLabel(source string) ([]byte, error) {
    var result = make([]byte, 0)
    var buffer = bytes.NewBuffer(result)
    var parts = strings.Split(source, ".")

    for _, part := range parts {
        if part == "" {
            buffer.Write([]byte{0})
            continue
        }

        binary.Write( buffer, binary.BigEndian, uint8(len(part)) )

        _, err := buffer.Write([]byte(part))
        if err != nil { return nil, err }
    }

    if !bytes.HasSuffix(buffer.Bytes(), []byte{0}) {
        buffer.Write([]byte{0})
    }

    return bytes.TrimSpace(buffer.Bytes()), nil
}


func RawToRecord(raw *RawRecord) (Record, error) {
    return nil, nil
}

func Uint16ToBytes(source uint16) []byte {
    var result = make([]byte, 2)
    result[0] = byte(uint8(source >> 8))
    result[1] = byte(uint8(source & 0xFFFF))
    return result
}

func Uint32ToBytes(source uint32) []byte {
    var result = make([]byte, 4)
    var high = Uint16ToBytes(uint16(source >> 16))
    var low = Uint16ToBytes(uint16(source & 0xFFFF))
    result[0] = high[0]
    result[1] = high[1]
    result[2] = low[0]
    result[3] = low[1]
    return result
}