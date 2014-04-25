package dns

import (
    "bytes"
    "testing"
)

//----------------------------------------------
// Message Header Serializing Tests
//----------------------------------------------

func TestMessage_HeaderValid(t *testing.T) {
    var header = MessageHeader{
        ID:                1234,
        Response:          true,
        Opcode:            0,
        Authoritative:     true,
        Truncated:         false,
        RecursionDesired:  true,
        RecursionAvailable:false,
        Zero:              false,
        Rcode:             0,
        QDCount:           1,
        ANCount:           3,
        NSCount:           0,
        ARCount:           1,
    }

    var knownID                 = []byte{ 0x04, 0xD2 }
    var knownLowOpts            = []byte{ 0x85 }// 0b10000101
    var knownHighOpts           = []byte{ 0x00 }// 0b00000000
    var knownQuery              = []byte{ 0x00, 0x01 }
    var knownAnswer             = []byte{ 0x00, 0x03 }
    var knownNS                 = []byte{ 0x00, 0x00 }
    var knownAdditional         = []byte{ 0x00, 0x01 }

    testHeaderSerialize(t, header, knownID, knownLowOpts, knownHighOpts, knownQuery, knownAnswer, knownNS, knownAdditional)
}

func TestMessage_HeaderNil(t *testing.T) {
    var header = MessageHeader{
        ID:                0,
        Response:          false,
        Opcode:            0,
        Authoritative:     false,
        Truncated:         false,
        RecursionDesired:  false,
        RecursionAvailable:false,
        Zero:              false,
        Rcode:             0,
        QDCount:           0,
        ANCount:           0,
        NSCount:           0,
        ARCount:           0,
    }

    var knownID                 = []byte{ 0x00, 0x00 }
    var knownLowOpts            = []byte{ 0x00 }
    var knownHighOpts           = []byte{ 0x00 }
    var knownQuery              = []byte{ 0x00, 0x00 }
    var knownAnswer             = []byte{ 0x00, 0x00 }
    var knownNS                 = []byte{ 0x00, 0x00 }
    var knownAdditional         = []byte{ 0x00, 0x00 }

    testHeaderSerialize(t, header, knownID, knownLowOpts, knownHighOpts, knownQuery, knownAnswer, knownNS, knownAdditional)
}

func TestMessage_HeaderZeroIgnored(t *testing.T) {
    var header = MessageHeader{
        ID:                1234,
        Response:          true,
        Opcode:            0,
        Authoritative:     true,
        Truncated:         false,
        RecursionDesired:  true,
        RecursionAvailable:false,
        Zero:              true,
        Rcode:             0,
        QDCount:           1,
        ANCount:           3,
        NSCount:           0,
        ARCount:           1,
    }

    var knownID                 = []byte{ 0x04, 0xD2 }
    var knownLowOpts            = []byte{ 0x85 }// 0b10000101
    var knownHighOpts           = []byte{ 0x00 }// 0b00000000
    var knownQuery              = []byte{ 0x00, 0x01 }
    var knownAnswer             = []byte{ 0x00, 0x03 }
    var knownNS                 = []byte{ 0x00, 0x00 }
    var knownAdditional         = []byte{ 0x00, 0x01 }

    testHeaderSerialize(t, header, knownID, knownLowOpts, knownHighOpts, knownQuery, knownAnswer, knownNS, knownAdditional)
}

func testHeaderSerialize(t *testing.T, header MessageHeader, knownID, knownLow, knownHigh, knownQuery, knownAnswer, knownNS, knownAdditional []byte) {
    var known = make([]byte, 0)
    var buffer = bytes.NewBuffer(known)
    buffer.Write(knownID)
    buffer.Write(knownLow)
    buffer.Write(knownHigh)
    buffer.Write(knownQuery)
    buffer.Write(knownAnswer)
    buffer.Write(knownNS)
    buffer.Write(knownAdditional)

    var serialized = header.Serialize()
    if bytes.Compare(serialized, buffer.Bytes()) != 0 || len(serialized) != len(buffer.Bytes()) {
        serializeError(t, serialized, buffer.Bytes())
    }

    if bytes.Compare(serialized[:2], knownID) != 0 {
        serializeError(t, knownID, serialized[:2])
    }

    if bytes.Compare(serialized[2:3], knownLow) != 0 {
        serializeError(t, knownLow, serialized[2:3])
    }

    if bytes.Compare(serialized[3:4], knownHigh) != 0 {
        serializeError(t, knownHigh, serialized[3:4])
    }

    if bytes.Compare(serialized[4:6], knownQuery) != 0 {
        serializeError(t, knownQuery, serialized[4:6])
    }

    if bytes.Compare(serialized[6:8], knownAnswer) != 0 {
        serializeError(t, knownAnswer, serialized[6:8])
    }

    if bytes.Compare(serialized[8:10], knownNS) != 0 {
        serializeError(t, knownNS, serialized[8:10])
    }

    if bytes.Compare(serialized[10:12], knownAdditional) != 0 {
        serializeError(t, knownAdditional, serialized[10:12])
    }
}

func serializeError(t *testing.T, expected, got []byte) {
    t.Errorf("Incorrect Serialization:\n\tExpected: %+v\n\tGot: %+v\n", expected, got)
}


//----------------------------------------------
// Question Serializing Tests
//----------------------------------------------

var testQuestions = QuestionCollection{
    { // 0
        Name:        "zed.io",
        Type:        1,        // lookup A record
        Class:       1,        // IN class
    },
    { // 1
        Name:        "zed.io",
        Type:        255,      // lookup ANY record
        Class:       1,        // IN class
    },
    { // 2
        Name:        "app.production.zed.io",
        Type:        28,       // lookup AAAA record
        Class:       1,        // IN class
    },
}

var testQuestionsKnown = [][]byte{
    { // 0
        3, 0x7A, 0x65, 0x64, 2, 0x69, 0x6F, 0x00,        // zed.io
        0x00, 1,                                         // looking for A record
        0x00, 1},                                        // in class IN
    { // 1
        3, 0x7A, 0x65, 0x64, 2, 0x69, 0x6F, 0x00,        // zed.io
        0x00, 255,                                       // looking for ANY record
        0x00, 1},
    { // 2
        // app.production.zed.io
        3, 0x61, 0x70, 0x70, 10, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x69, 0x6f, 0x6e, 3, 0x7A, 0x65, 0x64, 2, 0x69, 0x6F, 0x00,
        0x00, 28,                                        // looking for ANY record
        0x00, 1},
}

func TestQuestionsSerialization(t *testing.T) {
    for i, question := range testQuestions {
        testQuestionSerialization(t, question, testQuestionsKnown[i])
    }
}


func testQuestionSerialization(t *testing.T, question Question, known []byte) {
    var serialized, err = question.Serialize()
    if err != nil {
        t.Error(err)
    }

    var qClassOffset = len(serialized) -1
    var qTypeOffset = len(serialized) -2

    var kClassOffset = len(known) -1
    var kTypeOffset = len(known) -2

    if bytes.Compare(serialized[qClassOffset:], known[kClassOffset:]) != 0 {
        t.Errorf("Incorrect Class Serialization:\n\tExpected: %+v\n\tGot: %+v\n\tRaw:\n\t\tKnown: %+v\n\t\tReceived: %+v\n", serialized[qClassOffset:], known[kClassOffset:], known, serialized)
    }

    if bytes.Compare(serialized[qTypeOffset:qClassOffset], known[kTypeOffset:kClassOffset]) != 0 {
        t.Errorf("Incorrect Type Serialization:\n\tExpected: %+v\n\tGot: %+v\n\tRaw:\n\t\tKnown: %+v\n\t\tReceived: %+v\n", serialized[qTypeOffset:qClassOffset], known[kTypeOffset:kClassOffset], known, serialized)
    }
}


//--------------------------------------------------------------
// Per-Record Serializing Tests included in record package
//--------------------------------------------------------------
