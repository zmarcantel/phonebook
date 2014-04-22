package dns

import (
    "fmt"
    "bytes"
    "errors"
    "strings"
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

//----------------------------------------------
// Message Header Structures
//----------------------------------------------

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
    raw.LowOpts             = Btoi(self.Response) << 7
    raw.LowOpts            |= uint8(self.Opcode << 6)
    raw.LowOpts            |= Btoi(self.Authoritative) << 2
    raw.LowOpts            |= Btoi(self.Truncated) << 1
    raw.LowOpts            |= Btoi(self.RecursionDesired)

    raw.HighOpts            = Btoi(self.RecursionAvailable) << 7
    raw.HighOpts           |= uint8(self.Rcode)

    var result              = make([]byte, 0)
    var buffer              = bytes.NewBuffer(result)

    buffer.Write(Uint16ToBytes(self.ID))
    buffer.Write([]byte{byte(raw.LowOpts)})
    buffer.Write([]byte{byte(raw.HighOpts)})
    buffer.Write(Uint16ToBytes(self.QDCount))
    buffer.Write(Uint16ToBytes(self.ANCount))
    buffer.Write(Uint16ToBytes(self.NSCount))
    buffer.Write(Uint16ToBytes(self.ARCount))

    return bytes.TrimSpace(buffer.Bytes())
}

//----------------------------------------------
// Message Structures
//----------------------------------------------

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

func (self *Message) Print(indent string) {
    fmt.Printf("Message: %d\n", self.Header.ID)

    var shorter = strings.TrimSuffix(indent, "\n")
    fmt.Printf("%sHeader:\n", shorter)
    fmt.Printf("%s     Response: %v\n", shorter, self.Header.Response)
    fmt.Printf("%s       Opcode: %d\n", shorter, self.Header.Opcode)
    fmt.Printf("%sAuthoratative: %v\n", shorter, self.Header.Authoritative)
    fmt.Printf("%s    Truncated: %v\n", shorter, self.Header.Truncated)
    fmt.Printf("%s    Recursion: %v\n", shorter, self.Header.RecursionDesired)
    fmt.Printf("%sRec Available: %v\n", shorter, self.Header.RecursionAvailable)
    fmt.Printf("%s        RCode: %d\n", shorter, self.Header.Rcode)
    fmt.Printf("%s      Queries: %d\n", shorter, self.Header.QDCount)
    fmt.Printf("%s      Answers: %d\n", shorter, self.Header.ANCount)
    fmt.Printf("%s        Names: %d\n", shorter, self.Header.NSCount)
    fmt.Printf("%s   Additional: %d\n", shorter, self.Header.ARCount)

    fmt.Println("\tQuestions:")
    self.Questions.Print(indent)
    fmt.Println("\tAnswers:")
    self.Answers.Print(indent)
    fmt.Println("\tNS:")
    self.Ns.Print(indent)
    fmt.Println("\tExtra:")
    self.Extra.Print(indent)
}

func (self *Message) Serialize() ([]byte, error) {
    var result = make([]byte, 0)
    var buffer = bytes.NewBuffer(result)
    var header = self.Header.Serialize()
    buffer.Write(header)

    que, err := self.Questions.Serialize()
    if err != nil { return nil, errors.New("ERROR: Could not serialize Questions:\n" + err.Error()) }
    buffer.Write(que)

    ans, err := self.Answers.Serialize()
    if err != nil { return nil, errors.New("ERROR: Could not serialize Answers:\n" + err.Error()) }
    buffer.Write(ans)

    ns, err := self.Ns.Serialize()
    if err != nil { return nil, errors.New("ERROR: Could not serialize NS:\n" + err.Error()) }
    buffer.Write(ns)

    extra, err := self.Extra.Serialize()
    if err != nil { return nil, errors.New("ERROR: Could not serialize Extra:\n" + err.Error()) }
    buffer.Write(extra)

    return buffer.Bytes(), nil
}

//----------------------------------------------
// Functions
//----------------------------------------------

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
    header.Truncated           = BitSet( (raw.LowOpts & RAW_TRUNCATED), 2 )
    header.RecursionDesired    = BitSet( (raw.LowOpts & RAW_REC_DES), 1 )
    header.RecursionAvailable  = BitSet( (raw.HighOpts), 8 )
    header.Rcode               = int(raw.HighOpts & 0x0F)
    header.QDCount             = raw.QDCount
    header.ANCount             = raw.ANCount
    header.NSCount             = raw.NSCount
    header.ARCount             = raw.ARCount

    return header, 12 // TODO: non-hardcoded length
}

func reverse(input []byte) []byte {
    for i, j := 0, len(input)-1; i < j; i, j = i+1, j-1 {
        input[i], input[j] = input[j], input[i]
    }
    return input
}
