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
        binary.Write( buffer, binary.BigEndian, uint8(len(part)) )

        _, err := buffer.Write([]byte(part))
        if err != nil { return nil, err }
    }

    return bytes.TrimSpace(buffer.Bytes()), nil
}
