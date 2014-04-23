package dns

import (
    "fmt"
    "bytes"
    "encoding/binary"

    "github.com/zmarcantel/phonebook/dns/record"
)

//----------------------------------------------
// Question Structures
//----------------------------------------------

type Question struct {
    Name   string
    Type   uint16
    Class  uint16
}

type QuestionRaw struct {
    Name   []byte
    Type   uint16
    Class  uint16
}

//
// Transform the struct into a transmittable DNS record byte array
//
func (self Question) Serialize() ([]byte, error) {
    var result = make([]byte, 0)
    var buffer = bytes.NewBuffer(result)

    var label, err = record.CreateMessageLabel(self.Name)
    if err != nil { return nil, err }
    buffer.Write(label)

    buffer.Write(Uint16ToBytes(self.Type))
    buffer.Write(Uint16ToBytes(self.Class))

    return buffer.Bytes(), nil
}


//----------------------------------------------
// Question Collection Structures
//----------------------------------------------

type QuestionCollection []Question

//
// Serializes a set of questions using each query's .Serialize() function
//
func (self QuestionCollection) Serialize() ([]byte, error) {
    var result = make([]byte, 0)
    var buffer = bytes.NewBuffer(result)

    for _, q := range self {
        serialized, err := q.Serialize()
        if err != nil { return nil, err }

        buffer.Write(serialized)
    }

    return bytes.TrimSpace(buffer.Bytes()), nil
}

//
// Print out the list of queries (convenience function)
//
func (self QuestionCollection) Print(indent int) {
    var indentString string
    for i := 0 ; i < indent; i++ { indentString += "\t" }
    for _, q := range self {
        fmt.Printf("%sName: %s\n", indentString, q.Name)
        fmt.Printf("%s\tClass: %d\n", indentString, q.Class)
        fmt.Printf("%s\t Type: %d\n", indentString, q.Type)
    }
}

//
// Transform the queries in a DNS packet into a question structure
//
func UnpackQuestions(source []byte, count int) ([]Question, int) {
    var result = make([]Question, count)

    var offset int
    for i := 0 ; i < count ; i++ {
        name, offset := GetMessageLabel(source)

        var qType, qClass uint16
        binary.Read(bytes.NewReader(source[offset : offset + 2]), binary.BigEndian, &qType)
        binary.Read(bytes.NewReader(source[offset + 2 : offset + 5]), binary.BigEndian, &qClass)

        result[i] = Question {
            Name:        name,
            Type:        uint16(qType),
            Class:       uint16(qClass),
        }
        source = source[offset:]
    }

    return result, offset
}
