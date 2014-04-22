package dns

import (
)

func GetMessageLabel(source []byte) (string, int) {
    var name string
    var offset int

    for offset = 0 ; ; {
        var length = int(source[offset])

        if length == 0 {
            break
        }
        var start = offset + 1
        var finish = start + length
        if len(name) > 0 { name += "." }
        name += string(source[start:finish])
        offset += length + 1
    }
    offset += 1

    return name, offset
}

func BitSet(num uint8, i uint) bool {
    return Itob( (num >> (i - 1)) & 0x01 )
}

func Itob(i uint8) bool {
    if i == 1 {
        return true
    }
    return false
 }

func Btoi(b bool) uint8 {
    if b { return 1 }
    return 0
}

func Uint16ToBytes(source uint16) []byte {
    var result = make([]byte, 2)
    result[0] = byte(uint8(source >> 8))
    result[1] = byte(uint8(source & 0xFFFF))
    return result
}
