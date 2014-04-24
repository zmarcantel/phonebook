package record

import (
    "net"
    "time"
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
    var label = ""
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



//----------------------------------------------
// MX Tests
//----------------------------------------------

func TestMX_CreateValid_IPv4(t *testing.T) {
    var label = "127.0.0.1"
    var TTL = 10 * time.Second
    var priority = uint16(10)
    var target = "zed.io"

    testMX(t, false, label, target, priority, TTL)
}


func TestMX_CreateValid_Cidr(t *testing.T) {
    var label = "10.0.8/24-west-1b.zed.io"
    var TTL = 10 * time.Second
    var priority = uint16(10)
    var target = "long.production.zed.io"

    testMX(t, false, label, target, priority, TTL)
}

func TestMX_CreateValid_SubdomainTarget(t *testing.T) {
    var label = "127.0.0.1"
    var TTL = 10 * time.Second
    var priority = uint16(10)
    var target = "10-0-2-15.east-1b.zed.io"

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
    var label = "127.0.0.1"
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
