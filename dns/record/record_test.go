package record

import (
    "net"
    "time"
    "bytes"
    "testing"
)

//----------------------------------------------
// A Tests
//----------------------------------------------

func TestA_CreateValid(t *testing.T) {
    var label = "zed.io"
    var TTL = 10 * time.Second
    var ip = net.ParseIP("127.0.0.1")

    var record, err = A(label, TTL, ip)
    if err != nil {
        t.Error(err)
    }

    if record.Name != label || record.GetLabel() != label {
        t.Errorf("Incorrect Name:\n\tExpected: %s\n\tGot: %s\n", label, record.Name)
    }

    if record.TTL != TTL {
        t.Errorf("Incorrect TTL:\n\tExpected: %d\n\tGot: %d\n", TTL, record.TTL)
    }

    if record.Class != 1 {
        t.Errorf("Incorrect Class:\n\tExpected: %d\n\tGot: %d\n", 1, record.Class)
    }

    if record.Type != A_RECORD || record.GetType() != A_RECORD {
        t.Errorf("Incorrect Type:\n\tExpected: %d\n\tGot: %d\n", A_RECORD, record.Type)
    }

    if record.RDataLength != 4 {
        t.Errorf("Incorrect Data Length:\n\tExpected: %d\n\tGot: %d\n", 4, record.RDataLength)
    }
}

func TestA_CreateInvalid_EmptyLabel(t *testing.T) {
    var label = ""
    var TTL = 10 * time.Second
    var ip = net.ParseIP("127.0.0.1")

    var _, err = A(label, TTL, ip)
    if err == nil{
        t.Errorf("Didn't catch empty name error:\n\tExpected: %s\n\tGot: %+v\n", "non-nil", err)
    }
}

func TestA_CreateInvalid_ShortTTL(t *testing.T) {
    var label = "zed.io"
    var TTL = 4 * time.Second
    var ip = net.ParseIP("127.0.0.1")

    var _, err = A(label, TTL, ip)
    if err == nil{
        t.Errorf("Didn't catch <5s TTL error:\n\tExpected: %s\n\tGot: %+v\n", "non-nil", err)
    }
}

func TestA_CreateInvalid_IPv6(t *testing.T) {
    var label = "zed.io"
    var TTL = 10 * time.Second
    var ip = net.ParseIP("::1")

    var _, err = A(label, TTL, ip)
    if err == nil || err != ErrInvalidIP {
        t.Errorf("Didn't catch invalid IP error:\n\tExpected: %+v\n\tGot: %+v\n", ErrInvalidIP, err)
    }
}

func TestA_CreateInvalid_NilIP(t *testing.T) {
    var label = "zed.io"
    var TTL = 10 * time.Second

    var _, err = A(label, TTL, nil)
    if err == nil {
        t.Errorf("Allowed nil IP:\n\tExpected: %s\n\tGot: %+v\n", "non-nil", ErrInvalidIP)
    }
}

func TestA_Serialize(t *testing.T) {
    var label = "zed.io"
    var TTL = 10 * time.Second
    var ip = net.ParseIP("127.0.0.1")

    testSerializeA(t, label, TTL, ip, []byte{
        3, 0x7A, 0x65, 0x64, 2, 0x69, 0x6F, 0x00,        // zed.io
        0x00, 0x01,                                      // type
        0x00, 0x01,                                      // class
        0x00, 0x00, 0x00, 0xA,                           // ttl
        0x00, 0x04,                                      // data length
        127, 0, 0, 1})
}

func TestA_SerializeSubdomain(t *testing.T) {
    var label = "app.production.zed.io"
    var TTL = 10 * time.Second
    var ip = net.ParseIP("127.0.0.1")

    testSerializeA(t, label, TTL, ip, []byte{
        // app.production.zed.io
        3, 0x61, 0x70, 0x70, 10, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x69, 0x6f, 0x6e, 3, 0x7A, 0x65, 0x64, 2, 0x69, 0x6F, 0x00,
        0x00, 0x01,                                      // type
        0x00, 0x01,                                      // class
        0x00, 0x00, 0x00, 0xA,                           // ttl
        0x00, 0x04,                                      // data length
        127, 0, 0, 1})
}

func testSerializeA(t *testing.T, label string, ttl time.Duration, ip net.IP, known []byte) {
    var record, err = A(label, ttl, ip)
    if err != nil {
        t.Error(err)
    }

    serialized, err := record.Serialize()
    if err != nil {
        t.Errorf("Error while serializing:\n\t%s\n", err)
    }

    if bytes.Compare(serialized, known) != 0 {
        t.Errorf("Incorrect Record Serialization:\n\tExpected: %+v\n\t     Got: %+v\n", known, serialized)
    }
}


//----------------------------------------------
// AAAA Tests
//----------------------------------------------

func TestAAAA_CreateValid(t *testing.T) {
    var label = "zed.io"
    var TTL = 10 * time.Second
    var ip = net.ParseIP("::1")

    var record, err = AAAA(label, TTL, ip)
    if err != nil {
        t.Error(err)
    }

    if record.Name != label || record.GetLabel() != label {
        t.Errorf("Incorrect Name:\n\tExpected: %s\n\tGot: %s\n", label, record.Name)
    }

    if record.TTL != TTL {
        t.Errorf("Incorrect TTL:\n\tExpected: %d\n\tGot: %d\n", TTL, record.TTL)
    }

    if record.Class != 1 {
        t.Errorf("Incorrect Class:\n\tExpected: %d\n\tGot: %d\n", 1, record.Class)
    }

    if record.Type != AAAA_RECORD || record.GetType() != AAAA_RECORD {
        t.Errorf("Incorrect Type:\n\tExpected: %d\n\tGot: %d\n", AAAA_RECORD, record.Type)
    }

    if record.RDataLength != 16 {
        t.Errorf("Incorrect Data Length:\n\tExpected: %d\n\tGot: %d\n", 16, record.RDataLength)
    }
}

func TestAAAA_CreateInvalid_EmptyLabel(t *testing.T) {
    var label = ""
    var TTL = 10 * time.Second
    var ip = net.ParseIP("::1")

    var _, err = AAAA(label, TTL, ip)
    if err == nil{
        t.Errorf("Didn't catch empty name error:\n\tExpected: %s\n\tGot: %+v\n", "non-nil", err)
    }
}

func TestAAAA_CreateInvalid_ShortTTL(t *testing.T) {
    var label = "zed.io"
    var TTL = 4 * time.Second
    var ip = net.ParseIP("::1")

    var _, err = AAAA(label, TTL, ip)
    if err == nil{
        t.Errorf("Didn't catch <5s TTL error:\n\tExpected: %s\n\tGot: %+v\n", "non-nil", err)
    }
}

func TestAAAA_CreateInvalid_IPv4(t *testing.T) {
    var label = "zed.io"
    var TTL = 10 * time.Second
    var ip = net.ParseIP("127.0.0.1")

    var _, err = AAAA(label, TTL, ip)
    if err == nil || err != ErrInvalidIP {
        t.Errorf("Didn't catch invalid IP error:\n\tExpected: %+v\n\tGot: %+v\n", ErrInvalidIP, err)
    }
}


func TestAAAA_CreateInvalid_NilIP(t *testing.T) {
    var label = "zed.io"
    var TTL = 10 * time.Second

    var _, err = AAAA(label, TTL, nil)
    if err == nil {
        t.Errorf("Allowed nil IP:\n\tExpected: %s\n\tGot: %+v\n", "non-nil", err)
    }
}

func TestAAAA_Serialize(t *testing.T) {
    var label = "zed.io"
    var TTL = 10 * time.Second
    var ip = net.ParseIP("::1")

    testSerializeAAAA(t, label, TTL, ip, []byte{
        3, 0x7A, 0x65, 0x64, 2, 0x69, 0x6F, 0x00,        // zed.io
        0x00, 0x1C,                                      // type
        0x00, 0x01,                                      // class
        0x00, 0x00, 0x00, 0xA,                           // ttl
        0x00, 0x10,                                      // data length
        0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,})
}

func TestAAAA_SerializeSubdomainIPv6(t *testing.T) {
    var label = "app.production.zed.io"
    var TTL = 10 * time.Second
    var ip = net.ParseIP("2001:0db8:85a3:0042:1000:8a2e:0370:7334")

    testSerializeAAAA(t, label, TTL, ip, []byte{
        // app.production.zed.io
        3, 0x61, 0x70, 0x70, 10, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x69, 0x6f, 0x6e, 3, 0x7A, 0x65, 0x64, 2, 0x69, 0x6F, 0x00,
        0x00, 0x1C,                                      // type
        0x00, 0x01,                                      // class
        0x00, 0x00, 0x00, 0xA,                           // ttl
        0x00, 0x10,                                      // data length
        0x20, 0x01, 0x0d, 0xb8, 0x85, 0xa3, 0x00, 0x42, 0x10, 0x00, 0x8a, 0x2e, 0x03, 0x70, 0x73, 0x34})
}

func testSerializeAAAA(t *testing.T, label string, ttl time.Duration, ip net.IP, known []byte) {
    var record, err = AAAA(label, ttl, ip)
    if err != nil {
        t.Error(err)
    }

    serialized, err := record.Serialize()
    if err != nil {
        t.Errorf("Error while serializing:\n\t%s\n", err)
    }

    if bytes.Compare(serialized, known) != 0 {
        t.Errorf("Incorrect Record Serialization:\n\tExpected: %+v\n\t     Got: %+v\n", known, serialized)
    }
}


//----------------------------------------------
// SRV Tests
//----------------------------------------------

func TestSRV_CreateValid(t *testing.T) {
    var label = "_dns._udp.zed.io"
    var TTL = 10 * time.Second
    var priority = uint16(10)
    var weight = uint16(5)
    var port = uint16(8053)
    var target = "zed.io"
    var knownLength = uint16(8 + len(target))

    testSRV(t, false, label, target, TTL, priority, weight, port, knownLength)
}


func TestSRV_CreateValidSubdomain(t *testing.T) {
    var label = "_dns._udp.production.zed.io"
    var TTL = 10 * time.Second
    var priority = uint16(10)
    var weight = uint16(5)
    var port = uint16(8053)
    var target = "zed.io"
    var knownLength = uint16(8 + len(target))

    testSRV(t, false, label, target, TTL, priority, weight, port, knownLength)
}

func TestSRV_CreateValidMultiSubdomain(t *testing.T) {
    var label = "_dns._udp.mongo-1.east-1.production.zed.io"
    var TTL = 10 * time.Second
    var priority = uint16(10)
    var weight = uint16(5)
    var port = uint16(8053)
    var target = "zed.io"
    var knownLength = uint16(8 + len(target))

    testSRV(t, false, label, target, TTL, priority, weight, port, knownLength)
}

// ----- error testing

// time.Duration corrects negative TTL with abs(x)
// uint16(x) can overflow (not my problem) but is unsigned so negative is ok
// a port, weight, or priority can be 0 without issue per spec

func TestSRV_CreateInvalid_EmptyLabel(t *testing.T) {
    var label = ""
    var TTL = 10 * time.Second
    var priority = uint16(10)
    var weight = uint16(5)
    var port = uint16(8053)
    var target = "zed.io"
    var knownLength = uint16(8 + len(target))

    testSRV(t, true, label, target, TTL, priority, weight, port, knownLength)
}

func TestSRV_CreateInvalid_ShortTTL(t *testing.T) {
    var label = "_phonebook._tcp.zed.io"
    var TTL = 4 * time.Second
    var priority = uint16(10)
    var weight = uint16(5)
    var port = uint16(8053)
    var target = "zed.io"
    var knownLength = uint16(8 + len(target))

    testSRV(t, true, label, target, TTL, priority, weight, port, knownLength)
}

func TestSRV_CreateInvalid_SubdomainTarget(t *testing.T) {
    var label = "_test._tcp.zed.io"
    var TTL = 10 * time.Second
    var priority = uint16(10)
    var weight = uint16(5)
    var port = uint16(8053)
    var target = "east.production.zed.io"
    var knownLength = uint16(8 + len(target))

    testSRV(t, false, label, target, TTL, priority, weight, port, knownLength)
}

// ----- helper

func testSRV(t *testing.T, erroneous bool, label, target string, TTL time.Duration, priority, weight, port, length uint16) {
    var record, err = SRV(label, target, TTL, priority, weight, port)
    if err != nil && !erroneous {
        t.Error(err)
        return
    }

    if erroneous {
        if err == nil {
            t.Errorf("Failed to catch error:\n\tExpected: %s\n\tGot: %s\n", "non-nil", err)
        }
        return
    }

    if record.Name != label || record.GetLabel() != label {
        t.Errorf("Incorrect Name:\n\tExpected: %s\n\tGot: %s\n", label, record.Name)
    }

    if record.TTL != TTL {
        t.Errorf("Incorrect TTL:\n\tExpected: %d\n\tGot: %d\n", TTL, record.TTL)
    }

    if record.Class != 1 {
        t.Errorf("Incorrect Class:\n\tExpected: %d\n\tGot: %d\n", 1, record.Class)
    }

    if record.Type != SRV_RECORD || record.GetType() != SRV_RECORD {
        t.Errorf("Incorrect Type:\n\tExpected: %d\n\tGot: %d\n", SRV_RECORD, record.Type)
    }

    if record.RDataLength != length {
        t.Errorf("Incorrect Data Length:\n\tExpected: %d\n\tGot: %d\n", length, record.RDataLength)
    }

    if record.Target != target {
        t.Errorf("Incorrect Target Host:\n\tExpected: %s\n\tGot: %s\n", target, record.Target)
    }
}

func TestSRV_Serialize(t *testing.T) {
    var label = "_phonebook._tcp.zed.io"
    var TTL = 10 * time.Second
    var priority = uint16(10)
    var weight = uint16(5)
    var port = uint16(8053)
    var target = "zed.io"

    testSerializeSRV(t, label, target, TTL, priority, weight, port, []byte{
        0x5f, 0x70, 0x68, 0x6f, 0x6e, 0x65, 0x62, 0x6f, 0x6f, 0x6b,
        4, 0x5f, 0x74, 0x63, 0x70,
        3, 0x7a, 0x65, 0x64, 2, 0x69, 0x6f, 0x00,
        0x00, 0x21,                                      // type
        0x00, 0x01,                                      // class
        0x00, 0x00, 0x00, 0xA,                           // ttl
        0x00, 0x0E,                                      // data length
        0x00, 10,
        0x00, 5,
        0x1F, 0x75,
        3, 0x7a, 0x65, 0x64, 2, 0x69, 0x6f, 0x00,
        })
}

func TestSRV_SerializeSubdomain(t *testing.T) {
    var label = "_phonebook._tcp.west-1a.production.zed.io"
    var TTL = 10 * time.Second
    var priority = uint16(10)
    var weight = uint16(5)
    var port = uint16(8053)
    var target = "zed.io"

    testSerializeSRV(t, label, target, TTL, priority, weight, port, []byte{
        0x5f, 0x70, 0x68, 0x6f, 0x6e, 0x65, 0x62, 0x6f, 0x6f, 0x6b,
        4, 0x5f, 0x74, 0x63, 0x70,
        7, 0x77, 0x65, 0x73, 0x74, 0x2d, 0x31, 0x61,
        10, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x69, 0x6f, 0x6e,
        3, 0x7a, 0x65, 0x64, 2, 0x69, 0x6f, 0x00,
        0x00, 0x21,                                      // type
        0x00, 0x01,                                      // class
        0x00, 0x00, 0x00, 0xA,                           // ttl
        0x00, 0x0E,                                      // data length
        0x00, 10,
        0x00, 5,
        0x1F, 0x75,
        3, 0x7a, 0x65, 0x64, 2, 0x69, 0x6f, 0x00,
        })
}

func testSerializeSRV(t *testing.T, label, target string, TTL time.Duration, priority, weight, port uint16, known []byte) {
    var record, err = SRV(label, target, TTL, priority, weight, port)
    if err != nil {
        t.Error(err)
    }

    serialized, err := record.Serialize()
    if err != nil {
        t.Errorf("Error while serializing:\n\t%s\n", err)
    }

    if bytes.Compare(serialized, known) != 0 {
        t.Errorf("Incorrect Record Serialization:\n\tExpected: %+v\n\t     Got: %+v\n", known, serialized)
    }
}


//----------------------------------------------
// CNAME Tests
//----------------------------------------------

func TestCNAME_CreateValid(t *testing.T) {
    var label = "mongo-1.testing.zed.io"
    var TTL = 10 * time.Second
    var target = "10-0-2-15.east-1b.zed.io"
    var knownLength = uint16(len(target) + 1)

    testCNAME(t, false, label, target, TTL, knownLength)
}

func TestCNAME_CreateInvalid_EmptyLabel(t *testing.T) {
    var label = ""
    var TTL = 10 * time.Second
    var target = "10-0-2-15.east-1b.zed.io"
    var knownLength = uint16(len(target) + 1)

    testCNAME(t, true, label, target, TTL, knownLength)
}

// ----- error testing

func TestCNAME_CreateInvalid_ShortTTL(t *testing.T) {
    var label = "mongo-1.testing.zed.io"
    var TTL = 4 * time.Second
    var target = "10-0-2-15.east-1b.zed.io"
    var knownLength = uint16(len(target) + 1)

    testCNAME(t, true, label, target, TTL, knownLength)
}

func TestCNAME_CreateInvalid_EmptyTarget(t *testing.T) {
    var label = "mongo-1.testing.zed.io"
    var TTL = 4 * time.Second
    var target = ""
    var knownLength = uint16(len(target) + 1)

    testCNAME(t, true, label, target, TTL, knownLength)
}

// ----- helper

func testCNAME(t *testing.T, erroneous bool, label, target string, TTL time.Duration, length uint16) {
    var record, err = CNAME(label, target, TTL)
    if err != nil && !erroneous {
        t.Error(err)
        return
    }

    if erroneous {
        if err == nil {
            t.Errorf("Failed to catch error:\n\tExpected: %s\n\tGot: %s\n", "non-nil", err)
        }
        return
    }

    if record.Name != label || record.GetLabel() != label {
        t.Errorf("Incorrect Name:\n\tExpected: %s\n\tGot: %s\n", label, record.Name)
    }

    if record.TTL != TTL {
        t.Errorf("Incorrect TTL:\n\tExpected: %d\n\tGot: %d\n", TTL, record.TTL)
    }

    if record.Class != 1 {
        t.Errorf("Incorrect Class:\n\tExpected: %d\n\tGot: %d\n", 1, record.Class)
    }

    if record.Type != CNAME_RECORD || record.GetType() != CNAME_RECORD {
        t.Errorf("Incorrect Type:\n\tExpected: %d\n\tGot: %d\n", CNAME_RECORD, record.Type)
    }

    if record.RDataLength != length {
        t.Errorf("Incorrect Data Length:\n\tExpected: %d\n\tGot: %d\n", length, record.RDataLength)
    }

    if record.Target != target {
        t.Errorf("Incorrect Target Host:\n\tExpected: %s\n\tGot: %s\n", target, record.Target)
    }
}

// ----- serializing tests

func TestCNAME_SerializeNonHost(t *testing.T) {
    var label = "www"
    var TTL = 10 * time.Second
    var target = "zed.io"

    testSerializeCNAME(t, label, target, TTL, []byte{
        3, 0x77, 0x77, 0x77, 0x00,
        0x00, 0x05,                                      // type
        0x00, 0x01,                                      // class
        0x00, 0x00, 0x00, 0xA,                           // ttl
        0x00, 0x08,                                      // data length
        3, 0x7a, 0x65, 0x64, 2, 0x69, 0x6f, 0x00,
        })
}

func TestCNAME_SerializePlain(t *testing.T) {
    var label = "www.zed.io"
    var TTL = 10 * time.Second
    var target = "zed.io"

    testSerializeCNAME(t, label, target, TTL, []byte{
        3, 0x77, 0x77, 0x77, 3, 0x7a, 0x65, 0x64, 2, 0x69, 0x6f, 0x00,
        0x00, 0x05,                                      // type
        0x00, 0x01,                                      // class
        0x00, 0x00, 0x00, 0xA,                           // ttl
        0x00, 0x08,                                      // data length
        3, 0x7a, 0x65, 0x64, 2, 0x69, 0x6f, 0x00,
        })
}

func TestCNAME_SerializeSubdomain(t *testing.T) {
    var label = "smtp.zed.io"
    var TTL = 10 * time.Second
    var target = "east-1b.mail.zed.io"

    testSerializeCNAME(t, label, target, TTL, []byte{
        4, 0x73, 0x6d, 0x74, 0x70, 3, 0x7a, 0x65, 0x64, 2, 0x69, 0x6f, 0x00,
        0x00, 0x05,                                      // type
        0x00, 0x01,                                      // class
        0x00, 0x00, 0x00, 0xA,                           // ttl
        0x00, 0x15,                                      // data length
        7, 0x65, 0x61, 0x73, 0x74, 0x2d, 0x31, 0x62, 4, 0x6d, 0x61, 0x69, 0x6c, 3, 0x7a, 0x65, 0x64, 2, 0x69, 0x6f, 0x00,
        })
}

func testSerializeCNAME(t *testing.T, label, target string, TTL time.Duration, known []byte) {
    var record, err = CNAME(label, target, TTL)
    if err != nil {
        t.Error(err)
    }

    serialized, err := record.Serialize()
    if err != nil {
        t.Errorf("Error while serializing:\n\t%s\n", err)
    }

    if bytes.Compare(serialized, known) != 0 {
        t.Errorf("Incorrect Record Serialization:\n\tExpected: %+v\n\t     Got: %+v\n", known, serialized)
    }
}


//----------------------------------------------
// PTR Tests
//----------------------------------------------

func TestPTR_CreateValid_IPv4(t *testing.T) {
    var label = "127.0.0.1"
    var TTL = 10 * time.Second
    var target = "zed.io"

    testPTR(t, false, label, target, TTL)
}


func TestPTR_CreateValid_Cidr(t *testing.T) {
    var label = "10.0.8/24-west-1b.zed.io"
    var TTL = 10 * time.Second
    var target = "long.production.zed.io"

    testPTR(t, false, label, target, TTL)
}

func TestPTR_CreateValid_SubdomainTarget(t *testing.T) {
    var label = "127.0.0.1"
    var TTL = 10 * time.Second
    var target = "10-0-2-15.east-1b.zed.io"

    testPTR(t, false, label, target, TTL)
}

// ----- error testing

func TestPTR_CreateInvalid_EmptyLabel(t *testing.T) {
    var label = ""
    var TTL = 10 * time.Second
    var target = "zed.io"

    testPTR(t, true, label, target, TTL)
}

func TestPTR_CreateInvalid_EmptyTarget(t *testing.T) {
    var label = "127.0.0.1"
    var TTL = 10 * time.Second
    var target = ""

    testPTR(t, true, label, target, TTL)
}

// ----- helper

func testPTR(t *testing.T, erroneous bool, label, target string, TTL time.Duration) {
    var record, err = PTR(label, target, TTL)
    if err != nil && !erroneous {
        t.Error(err)
        return
    }

    if erroneous {
        if err == nil {
            t.Errorf("Failed to catch error:\n\tExpected: %s\n\tGot: %s\n", "non-nil", err)
        }
        return
    }

    if record.Name != label || record.GetLabel() != label {
        t.Errorf("Incorrect Name:\n\tExpected: %s\n\tGot: %s\n", label, record.Name)
    }

    if record.TTL != TTL {
        t.Errorf("Incorrect TTL:\n\tExpected: %d\n\tGot: %d\n", TTL, record.TTL)
    }

    if record.Class != 1 {
        t.Errorf("Incorrect Class:\n\tExpected: %d\n\tGot: %d\n", 1, record.Class)
    }

    if record.Type != PTR_RECORD || record.GetType() != PTR_RECORD {
        t.Errorf("Incorrect Type:\n\tExpected: %d\n\tGot: %d\n", PTR_RECORD, record.Type)
    }

    if record.Target != target {
        t.Errorf("Incorrect Target Host:\n\tExpected: %s\n\tGot: %s\n", target, record.Target)
    }
}


// ----- serializing tests

func TestPTR_SerializeIPv4Localhost(t *testing.T) {
    var label = "127.0.0.1"
    var TTL = 10 * time.Second
    var target = "zed.io"

    testSerializePTR(t, label, target, TTL, []byte{
        3, 0x31, 0x32, 0x37, 1, 0x30, 1, 0x30, 1, 0x31, 0x00,
        0x00, 0x0C,                                      // type
        0x00, 0x01,                                      // class
        0x00, 0x00, 0x00, 0xA,                           // ttl
        0x00, 0x08,                                      // data length
        3, 0x7a, 0x65, 0x64, 2, 0x69, 0x6f, 0x00,
        })
}

func TestPTR_SerializeIPv4(t *testing.T) {
    var label = "192.168.33.50"
    var TTL = 10 * time.Second
    var target = "zed.io"

    testSerializePTR(t, label, target, TTL, []byte{
        3, 0x31, 0x39, 0x32, 3, 0x31, 0x36, 0x38, 2, 0x33, 0x33, 2, 0x35, 0x30, 0x00,
        0x00, 0x0C,                                      // type
        0x00, 0x01,                                      // class
        0x00, 0x00, 0x00, 0xA,                           // ttl
        0x00, 0x08,                                      // data length
        3, 0x7a, 0x65, 0x64, 2, 0x69, 0x6f, 0x00,
        })
}

func TestPTR_SerializeCidr(t *testing.T) {
    var label = "10.0.8/24-west-1b.zed.io"
    var TTL = 10 * time.Second
    var target = "east-1b.mail.zed.io"

    testSerializePTR(t, label, target, TTL, []byte{
        2, 0x31, 0x30, 1, 0x30, 12, 0x38, 0x2f, 0x32, 0x34, 0x2d, 0x77, 0x65, 0x73, 0x74, 0x2d, 0x31, 0x62, 3, 0x7a, 0x65, 0x64, 2, 0x69, 0x6f, 0x00,
        0x00, 0x0C,                                      // type
        0x00, 0x01,                                      // class
        0x00, 0x00, 0x00, 0xA,                           // ttl
        0x00, 0x15,                                      // data length
        7, 0x65, 0x61, 0x73, 0x74, 0x2d, 0x31, 0x62, 4, 0x6d, 0x61, 0x69, 0x6c, 3, 0x7a, 0x65, 0x64, 2, 0x69, 0x6f, 0x00,
        })
}

func testSerializePTR(t *testing.T, label, target string, TTL time.Duration, known []byte) {
    var record, err = PTR(label, target, TTL)
    if err != nil {
        t.Error(err)
    }

    serialized, err := record.Serialize()
    if err != nil {
        t.Errorf("Error while serializing:\n\t%s\n", err)
    }

    if bytes.Compare(serialized, known) != 0 {
        t.Errorf("Incorrect Record Serialization:\n\tExpected: %+v\n\t     Got: %+v\n", known, serialized)
    }
}



//----------------------------------------------
// MX Tests
//----------------------------------------------

func TestMX_CreateValid(t *testing.T) {
    var label = "mail.zed.io"
    var TTL = 10 * time.Second
    var priority = uint16(10)
    var target = "zed.io"

    testMX(t, false, label, target, priority, TTL)
}


func TestMX_CreateValid_Subdomains(t *testing.T) {
    var label = "smtp.west-1b.zed.io"
    var TTL = 10 * time.Second
    var priority = uint16(10)
    var target = "long.production.zed.io"

    testMX(t, false, label, target, priority, TTL)
}

// ----- error testing

func TestMX_CreateInvalid_EmptyLabel(t *testing.T) {
    var label = ""
    var TTL = 10 * time.Second
    var priority = uint16(10)
    var target = "zed.io"

    testMX(t, true, label, target, priority, TTL)
}

func TestMX_CreateInvalid_EmptyTarget(t *testing.T) {
    var label = "mail.zed.io"
    var TTL = 10 * time.Second
    var priority = uint16(10)
    var target = ""

    testMX(t, true, label, target, priority, TTL)
}

// ----- helper

func testMX(t *testing.T, erroneous bool, label, target string, priority uint16, TTL time.Duration) {
    var record, err = MX(label, target, priority, TTL)
    if err != nil && !erroneous {
        t.Error(err)
        return
    }

    if erroneous {
        if err == nil {
            t.Errorf("Failed to catch error:\n\tExpected: %s\n\tGot: %s\n", "non-nil", err)
        }
        return
    }

    if record.Name != label || record.GetLabel() != label {
        t.Errorf("Incorrect Name:\n\tExpected: %s\n\tGot: %s\n", label, record.Name)
    }

    if record.TTL != TTL {
        t.Errorf("Incorrect TTL:\n\tExpected: %d\n\tGot: %d\n", TTL, record.TTL)
    }

    if record.Class != 1 {
        t.Errorf("Incorrect Class:\n\tExpected: %d\n\tGot: %d\n", 1, record.Class)
    }

    if record.Type != MX_RECORD || record.GetType() != MX_RECORD {
        t.Errorf("Incorrect Type:\n\tExpected: %d\n\tGot: %d\n", MX_RECORD, record.Type)
    }

    if record.Priority != priority {
        t.Errorf("Incorrect Priority:\n\tExpected: %d\n\tGot: %d\n", priority, record.Priority)
    }

    if record.Target != target {
        t.Errorf("Incorrect Target Host:\n\tExpected: %s\n\tGot: %s\n", target, record.Target)
    }
}

// ----- serializing tests

func TestMX_Serialize(t *testing.T) {
    var label = "mail.zed.io"
    var TTL = 10 * time.Second
    var priority = uint16(10)
    var target = "zed.io"

    testSerializeMX(t, label, target, TTL, priority, []byte{
        4, 0x6d, 0x61, 0x69, 0x6c, 3, 0x7a, 0x65, 0x64, 2, 0x69, 0x6f, 0x00,
        0x00, 0x0F,                                      // type
        0x00, 0x01,                                      // class
        0x00, 0x00, 0x00, 0xA,                           // ttl
        0x00, 0x0A,                                      // data length
        0x00, 0x0A,                                      // priority
        3, 0x7a, 0x65, 0x64, 2, 0x69, 0x6f, 0x00,
        })
}

func TestMX_SerializeSubdomain(t *testing.T) {
    var label = "mail.zed.io"
    var TTL = 10 * time.Second
    var priority = uint16(10)
    var target = "east-1b.smtp.zed.io"

    testSerializeMX(t, label, target, TTL, priority, []byte{
        4, 0x6d, 0x61, 0x69, 0x6c, 3, 0x7a, 0x65, 0x64, 2, 0x69, 0x6f, 0x00,
        0x00, 0x0F,                                      // type
        0x00, 0x01,                                      // class
        0x00, 0x00, 0x00, 0xA,                           // ttl
        0x00, 0x17,                                      // data length
        0x00, 0x0A,                                      // priority
        7, 0x65, 0x61, 0x73, 0x74, 0x2d, 0x31, 0x62, 4, 0x73, 0x6d, 0x74, 0x70, 3, 0x7a, 0x65, 0x64, 2, 0x69, 0x6f, 0x00,
        })
}

func testSerializeMX(t *testing.T, label, target string, TTL time.Duration, priority uint16, known []byte) {
    var record, err = MX(label, target, priority, TTL)
    if err != nil {
        t.Error(err)
    }

    serialized, err := record.Serialize()
    if err != nil {
        t.Errorf("Error while serializing:\n\t%s\n", err)
    }

    if bytes.Compare(serialized, known) != 0 {
        t.Errorf("Incorrect Record Serialization:\n\tExpected: %+v\n\t     Got: %+v\n", known, serialized)
    }
}


//----------------------------------------------
// TXT Tests
//----------------------------------------------

func TestTXT_CreateValid(t *testing.T) {
    var label = "mail.production"
    var TTL = 10 * time.Second
    var text = "admin user is Zach Marcantel"

    testTXT(t, false, label, text, TTL)
}


func TestTXT_CreateValid_SubdomainWithJSON(t *testing.T) {
    var label = "app.east-1b.production"
    var TTL = 10 * time.Second
    var text = "{ \"version\": 0.1, \"upSince\": \"2014-04-25T06:19:59.085Z\" }"

    testTXT(t, false, label, text, TTL)
}

// ----- error testing

func TestTXT_CreateInvalid_EmptyLabel(t *testing.T) {
    var label = ""
    var TTL = 10 * time.Second
    var text = "{ \"version\": 0.1, \"upSince\": \"2014-04-25T06:19:59.085Z\" }"

    testTXT(t, true, label, text, TTL)
}

func TestTXT_CreateInvalid_EmptyText(t *testing.T) {
    var label = "zed.io"
    var TTL = 10 * time.Second
    var text = ""

    testTXT(t, true, label, text, TTL)
}

// ----- helper

func testTXT(t *testing.T, erroneous bool, label, text string, TTL time.Duration) {
    var record, err = TXT(label, TTL, text)
    if err != nil && !erroneous {
        t.Error(err)
        return
    }

    if erroneous {
        if err == nil {
            t.Errorf("Failed to catch error:\n\tExpected: %s\n\tGot: %s\n", "non-nil", err)
        }
        return
    }

    if record.Name != label || record.GetLabel() != label {
        t.Errorf("Incorrect Name:\n\tExpected: %s\n\tGot: %s\n", label, record.Name)
    }

    if record.TTL != TTL {
        t.Errorf("Incorrect TTL:\n\tExpected: %d\n\tGot: %d\n", TTL, record.TTL)
    }

    if record.Class != 1 {
        t.Errorf("Incorrect Class:\n\tExpected: %d\n\tGot: %d\n", 1, record.Class)
    }

    if record.Type != TXT_RECORD || record.GetType() != TXT_RECORD {
        t.Errorf("Incorrect Type:\n\tExpected: %d\n\tGot: %d\n", TXT_RECORD, record.Type)
    }

    if record.Text != text {
        t.Errorf("Incorrect Text Data:\n\tExpected: %s\n\tGot: %s\n", text, record.Text)
    }
}


// ----- serializing tests

func TestTXT_Serialize(t *testing.T) {
    var label = "mail.production"
    var TTL = 10 * time.Second
    var text = "admin user is Zach Marcantel"

    testSerializeTXT(t, label, text, TTL, []byte{
        4, 0x6d, 0x61, 0x69, 0x6c, 10, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x00,
        0x00, 0x10,                                      // type
        0x00, 0x01,                                      // class
        0x00, 0x00, 0x00, 0xA,                           // ttl
        0x00, 0x1C,                                      // data length
        0x61, 0x64, 0x6d, 0x69, 0x6e, 0x20, 0x75, 0x73, 0x65, 0x72, 0x20, 0x69, 0x73, 0x20, 0x5a, 0x61, 0x63, 0x68, 0x20, 0x4d, 0x61, 0x72, 0x63, 0x61, 0x6e, 0x74, 0x65, 0x6c,
        })
}

func testSerializeTXT(t *testing.T, label, text string, TTL time.Duration, known []byte) {
    var record, err = TXT(label, TTL, text)
    if err != nil {
        t.Error(err)
    }

    serialized, err := record.Serialize()
    if err != nil {
        t.Errorf("Error while serializing:\n\t%s\n", err)
    }

    if bytes.Compare(serialized, known) != 0 {
        t.Errorf("Incorrect Record Serialization:\n\tExpected: %+v\n\t     Got: %+v\n", known, serialized)
    }
}
