package server

import (
    "errors"

    "github.com/zmarcantel/phonebook/dns/record"
)

var ErrNotFound = errors.New("ERROR: That record does not exist")

var cache map[string][]record.Record

func AddRecord(rec record.Record) error {
    if rec == nil { return errors.New("ERROR: Cannot add nil record.") }
    if cache == nil { cache = make(map[string][]record.Record) }

    var header = rec.Header()
    var otherRecords = FindRecordsByLabel(header.Name)
    if otherRecords == nil {
        cache[header.Name] = make([]record.Record, 0)
    }

    otherRecords = append(otherRecords, rec)
    cache[header.Name] = otherRecords

    return nil
}

func FindRecordsByLabel(label string) []record.Record {
    if records, exists := cache[label] ; exists {
        return records
    } else {
        return nil
    }
}

func FindRecord(label string, rType uint16, rClass uint16) (record.Record, error) {
    var records = FindRecordsByLabel(label)
    if records == nil {
        return nil, errors.New("ERROR: No records for label: " + label)
    }

    for _, rec := range records {
        var header = rec.Header()
        if header.Type == rType {
            return rec, nil
        }
    }

    return nil, ErrNotFound
}


func DumpCache() map[string][]record.Record {
    return cache
}
