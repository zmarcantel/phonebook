package server

import (
    "errors"
    "strings"

    "github.com/zmarcantel/phonebook/dns/record"
)

var ErrNotFound     error   = errors.New("ERROR: That record does not exist")


// naive cache...
// TODO: modular cache mechanisms
var cache map[string][]record.Record


//
// Add a record to the DNS server's lcoal cache
// This is the function responsible for registering new records so they can be queried
//
func AddRecord(rec record.Record) error {
    // input validation
    if rec == nil { return errors.New("ERROR: Cannot add nil record.") }
    if cache == nil { cache = make(map[string][]record.Record) }

    // check if a record's label has been registered in the cache
    // if not, we need to make an entry in the cache for that label
    var otherRecords = FindRecordsByLabel(rec.GetLabel())
    if otherRecords == nil {
        cache[rec.GetLabel()] = make([]record.Record, 0)
    }

    // whether it's a fresh list or one with existing records, append the
    // new record to the "cache line"
    otherRecords = append(otherRecords, rec)

    // give ability to handle per-record-type adding procedure
    switch (rec.GetType()) {
        default:
            // default -- trim any authoratative trailing "."
            cache[strings.TrimSuffix(rec.GetLabel(), ".")] = otherRecords
            break
    }

    return nil
}


//
// Return the full contents of the cache as a map[string][]Record
//
func GetCache() map[string][]record.Record {
    return cache
}
