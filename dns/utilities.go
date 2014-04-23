package dns

import (
)

//
// Transform a series of bytes into a qualified DNS label
//
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

//
// Test if a bit (i) is set in the byte (num)
//
func BitSet(num uint8, i uint) bool {
    return Itob( (num >> (i - 1)) & 0x01 )
}

//
// Translate 0/1 into true/false
//
func Itob(i uint8) bool {
    if i == 1 {
        return true
    }
    return false
 }

//
// Translate true/false to 0/1
//
func Btoi(b bool) uint8 {
    if b { return 1 }
    return 0
}

//
// Correctly serialize a 16bit unsigned integer into a byte array of length 2
//
func Uint16ToBytes(source uint16) []byte {
    var result = make([]byte, 2)
    result[0] = byte(uint8(source >> 8))
    result[1] = byte(uint8(source & 0xFFFF))
    return result
}
