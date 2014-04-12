package dns

import (
    "fmt"
    "bytes"
    "errors"
    "encoding/binary"

    "github.com/zmarcantel/phonebook/dns/record"
)

const (
    RAW_RESPONSE   uint8 = 0x80
    RAW_OPCODE     uint8 = 0x78
    RAW_AUTHORITY  uint8 = 0x04
    RAW_TRUNCATED  uint8 = 0x02
    RAW_REC_DES    uint8 = 0x01
)

type MessageHeader struct {
    ID                 uint16
    Response           bool
    Opcode             uint32
    Authoritative      bool
    Truncated          bool
    RecursionDesired   bool
    RecursionAvailable bool
    Zero               bool
    Rcode              int
    QDCount            uint16
    ANCount            uint16
    NSCount            uint16
    ARCount            uint16
}

type MessageHeaderRaw struct {
    ID                 uint16
    LowOpts            uint8
    HighOpts           uint8
    QDCount            uint16
    ANCount            uint16
    NSCount            uint16
    ARCount            uint16
}

func (self *MessageHeader) Serialize() []byte {
    var raw MessageHeaderRaw
    raw.ID                  = self.ID

    raw.LowOpts             = Btoi(self.Response) << 7
    raw.LowOpts            |= uint8(self.Opcode << 6)
    raw.LowOpts            |= Btoi(self.Authoritative) << 2
    raw.LowOpts            |= Btoi(self.Truncated) << 1
    raw.LowOpts            |= Btoi(self.RecursionDesired)

    raw.HighOpts            = Btoi(self.RecursionAvailable) << 7
    raw.HighOpts           |= uint8(self.Rcode)

    var result              = make([]byte, 0)
    var buffer              = bytes.NewBuffer(result)
    binary.Write(buffer, binary.BigEndian, raw.ID)
    binary.Write(buffer, binary.BigEndian, raw.LowOpts)
    binary.Write(buffer, binary.BigEndian, raw.HighOpts)
    binary.Write(buffer, binary.BigEndian, self.QDCount)
    binary.Write(buffer, binary.BigEndian, self.ANCount)
    binary.Write(buffer, binary.BigEndian, self.NSCount)
    binary.Write(buffer, binary.BigEndian, self.ARCount)

    return bytes.TrimSpace(buffer.Bytes())
}

type Message struct {
    Header            MessageHeader
    Questions         QuestionCollection            // Holds the RR(s) of the question section.
    Answers           record.RecordCollection       // Holds the RR(s) of the answer section.
    Ns                record.RecordCollection       // Holds the RR(s) of the authority section.
    Extra             record.RecordCollection       // Holds the RR(s) of the additional section.
}

type MessageRaw struct {
    MessageHeaderRaw
    Questions         QuestionCollection            // Holds the RR(s) of the question section.
    Answers           record.RecordCollection       // Holds the RR(s) of the answer section.
    Ns                record.RecordCollection       // Holds the RR(s) of the authority section.
    Extra             record.RecordCollection       // Holds the RR(s) of the additional section.
}

func (self *Message) Print() {
    fmt.Printf("Message: %d\n", self.Header.ID)

    fmt.Println("\tHeader:")
    fmt.Printf("\t\t     Response: %v\n", self.Header.Response)
    fmt.Printf("\t\t       Opcode: %d\n", self.Header.Opcode)
    fmt.Printf("\t\tAuthoratative: %v\n", self.Header.Authoritative)
    fmt.Printf("\t\t    Truncated: %v\n", self.Header.Truncated)
    fmt.Printf("\t\t    Recursion: %v\n", self.Header.RecursionDesired)
    fmt.Printf("\t\tRec Available: %v\n", self.Header.RecursionAvailable)
    fmt.Printf("\t\t        RCode: %d\n", self.Header.Rcode)
    fmt.Printf("\t\t      Queries: %d\n", self.Header.QDCount)
    fmt.Printf("\t\t      Answers: %d\n", self.Header.ANCount)
    fmt.Printf("\t\t        Names: %d\n", self.Header.NSCount)
    fmt.Printf("\t\t   Additional: %d\n", self.Header.ARCount)

    self.Questions.Print("\t")
    fmt.Println("\tAnswers:")
    self.Answers.Print("\t\t")
    fmt.Println("\tNS:")
    self.Ns.Print("\t\t")
    fmt.Println("\tExtra:")
    self.Extra.Print("\t\t")
}

func (self *Message) Serialize() ([]byte, error) {
    var result = make([]byte, 0)
    var buffer = bytes.NewBuffer(result)
    var header = self.Header.Serialize()

    _, err := buffer.Write(header)
    if err != nil { return nil, err }

    serialized, err := self.Questions.Serialize()
    if err != nil { return nil, errors.New("ERROR: Could not serialize Questions:\n" + err.Error()) }
    _, err = buffer.Write(serialized)
    if err != nil { return nil, err }

    serialized, err = self.Answers.Serialize()
    if err != nil { return nil, errors.New("ERROR: Could not serialize Answers:\n" + err.Error()) }
    _, err = buffer.Write(serialized)
    if err != nil { return nil, err }

    serialized, err = self.Ns.Serialize()
    if err != nil { return nil, errors.New("ERROR: Could not serialize NS:\n" + err.Error()) }
    _, err = buffer.Write(serialized)
    if err != nil { return nil, err }

    serialized, err = self.Extra.Serialize()
    if err != nil { return nil, errors.New("ERROR: Could not serialize Extra:\n" + err.Error()) }
    _, err = buffer.Write(serialized)
    if err != nil { return nil, err }

    return buffer.Bytes(), nil
}

func UnpackMessage(source []byte) *Message {
    header, offset := UnpackHeader(source)
    questions, offset := UnpackQuestions(source[offset:], int(header.QDCount))
    _, offset = UnpackResources(source[offset:], int(header.ANCount))

    return &Message{
        header,
        questions,
        []record.Record{},
        []record.Record{},
        []record.Record{},
    }
}

func UnpackHeader(source []byte) (MessageHeader, int) {
    var raw MessageHeaderRaw
    binary.Read(bytes.NewReader(source), binary.BigEndian, &raw)

    var header MessageHeader
    header.ID = raw.ID
    header.Response            = BitSet(raw.LowOpts, 8)
    header.Opcode              = uint32( (raw.LowOpts & RAW_OPCODE) >> 3 )
    header.Authoritative       = BitSet( (raw.LowOpts & RAW_AUTHORITY), 3 )
    header.Truncated           = BitSet( (raw.LowOpts & RAW_AUTHORITY), 2 )
    header.RecursionDesired    = BitSet( (raw.LowOpts & RAW_AUTHORITY), 1 )
    header.RecursionAvailable  = BitSet( (raw.HighOpts & RAW_AUTHORITY), 8 )
    header.Rcode               = int(raw.HighOpts & 0x0F)
    header.QDCount             = raw.QDCount
    header.ANCount             = raw.ANCount
    header.NSCount             = raw.NSCount
    header.ARCount             = raw.ARCount

    return header, 12 // TODO: non-hardcoded length
}
