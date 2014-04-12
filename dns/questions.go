package dns

import (
    "fmt"
    "bytes"
    "encoding/binary"

    "./record"
)


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

func (self Question) Serialize() ([]byte, error) {
    var label, err = record.CreateMessageLabel(self.Name)
    if err != nil { return nil, err }

    var result = make([]byte, 0)
    var buffer = bytes.NewBuffer(result)

    err = binary.Write(buffer, binary.BigEndian, QuestionRaw{
        label,
        self.Type,
        self.Class,
    })
    if err != nil { return nil, err }

    return result, nil
}

type QuestionCollection []Question

func (self QuestionCollection) Serialize() ([]byte, error) {
    var result = make([]byte, 0)
    var buffer = bytes.NewBuffer(result)

    var err error
    var serialized []byte
    for _, q := range self {
        serialized, err = q.Serialize()
        if err != nil { return nil, err }

        _, err = buffer.Write(serialized)
        if err != nil { return nil, err }
    }

    return bytes.TrimSpace(result), nil
}


func (self QuestionCollection) Print(indent string) {
    fmt.Printf("%sQuestions:\n", indent)
    for _, q := range self {
        fmt.Printf("%sName: %s\n", indent, q.Name)
        fmt.Printf("%s\tClass: %d\n", indent, q.Class)
        fmt.Printf("%s\t Type: %d\n", indent, q.Type)
    }
}


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
