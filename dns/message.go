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

//
// Transform the MessageHeader struct into a transmittable DNS packet header
//
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

//
// Print the message to the stdout (convenience function)
//
func (self *Message) Print(indent int) {
    var indentString string
    for i := 0 ; i < indent; i++ { indentString += "\t" }

    fmt.Printf("%sMessage: %d\n", indentString, self.Header.ID)
    fmt.Printf("%sHeader:\n", indentString)
    fmt.Printf("%s\t     Response: %v\n", indentString, self.Header.Response)
    fmt.Printf("%s\t       Opcode: %d\n", indentString, self.Header.Opcode)
    fmt.Printf("%s\tAuthoratative: %v\n", indentString, self.Header.Authoritative)
    fmt.Printf("%s\t    Truncated: %v\n", indentString, self.Header.Truncated)
    fmt.Printf("%s\t    Recursion: %v\n", indentString, self.Header.RecursionDesired)
    fmt.Printf("%s\tRec Available: %v\n", indentString, self.Header.RecursionAvailable)
    fmt.Printf("%s\t        RCode: %d\n", indentString, self.Header.Rcode)
    fmt.Printf("%s\t      Queries: %d\n", indentString, self.Header.QDCount)
    fmt.Printf("%s\t      Answers: %d\n", indentString, self.Header.ANCount)
    fmt.Printf("%s\t        Names: %d\n", indentString, self.Header.NSCount)
    fmt.Printf("%s\t   Additional: %d\n", indentString, self.Header.ARCount)

    fmt.Printf("%sQuestions:\n", indentString)
    self.Questions.Print(indent + 1)
    fmt.Printf("%sAnswers:\n", indentString)
    self.Answers.Print(indent + 1)
    fmt.Printf("%sNS:\n", indentString)
    self.Ns.Print(indent + 1)
    fmt.Printf("%sExtra:\n", indentString)
    self.Extra.Print(indent + 1)
}


//
// Transform the entire message into a transmittable DNS packet
// This is a crucial function to the response to queries
//
func (self *Message) Serialize() ([]byte, error) {
    // make an empty buffer
    var result = make([]byte, 0)
    var buffer = bytes.NewBuffer(result)

    // first, grab the header
    var header = self.Header.Serialize()
    buffer.Write(header)

    // then copy over any questions
    que, err := self.Questions.Serialize()
    if err != nil { return nil, errors.New("ERROR: Could not serialize Questions:\n" + err.Error()) }
    buffer.Write(que)

    // fill in our answers
    ans, err := self.Answers.Serialize()
    if err != nil { return nil, errors.New("ERROR: Could not serialize Answers:\n" + err.Error()) }
    buffer.Write(ans)

    // we won't have any nameserver records, but do that just in case
    ns, err := self.Ns.Serialize()
    if err != nil { return nil, errors.New("ERROR: Could not serialize NS:\n" + err.Error()) }
    buffer.Write(ns)

    // and any additional features, records, etc that need to be communicated
    extra, err := self.Extra.Serialize()
    if err != nil { return nil, errors.New("ERROR: Could not serialize Extra:\n" + err.Error()) }
    buffer.Write(extra)

    return buffer.Bytes(), nil
}

//----------------------------------------------
// Functions
//----------------------------------------------

//
// Translate a DNS packet into a readable message
//
func UnpackMessage(source []byte) *Message {
    header, offset := UnpackHeader(source)
    questions, offset := UnpackQuestions(source[offset:], int(header.QDCount))

    return &Message{
        header,
        questions,
        []record.Record{},
        []record.Record{},
        []record.Record{},
    }
}

//
// Extract the header from the DNS packet
//
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

