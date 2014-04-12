package record

import (
    "fmt"
    "net"
    "time"
    "errors"
)

func RawToRecord(raw *RawRecord) (Record, error) {
    return nil, nil
}


//----------------------------------------------
//  A Record
//      Hostname -> IPV4
//----------------------------------------------

func A(hostname string, ttl time.Duration, target net.IP) (*ARecord, error) {
    if len(hostname) <= 0 {
        return nil, errors.New(fmt.Sprintf("The record must contain a hostname. Received: '%s'.", hostname))
    } else if ttl.Seconds() < 5 { // TODO: get actual max class int
        return nil, errors.New(fmt.Sprintf("TTL of <5s is not supported. Received: %d", ttl.Seconds))
    }
    // TODO: add checks on target -- are we remapping the current IP and some other security stuff

    return &ARecord{
        Header: RecordHeader{
            Name:            hostname,
            Class:           uint16( 1 ),                 // 'IN' class
            Type:            1,                           // A Records are class 1
            TTL:             ttl,
            RDataLength:     4,                           // A Records send the target IP in a 4-octet data section
        },
        IP:                  target,
    }, nil
}
