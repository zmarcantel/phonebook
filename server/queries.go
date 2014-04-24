package server

import (
    "fmt"
    "errors"
    "strings"

    "github.com/zmarcantel/phonebook/dns"
    "github.com/zmarcantel/phonebook/dns/record"
)

const (
    DNS_QUERY_ALL int        = 255
)

//
// Given a set of DNS queries, look into the local cache and get answers
//
func AnswerQuestions(questions []dns.Question) ([]record.Record, error) {
    var result = make([]record.Record, 0)

    for _, question := range questions {
        // handle the 'ANY' search
        if int(+question.Type) == DNS_QUERY_ALL {
            var records = FindRecordsByLabel(question.Name)
            result = append(result, records...)
            continue
        }

        // handle a (label x type) search
        var recs, err = FindRecords(question.Name, question.Type, true)
        if err != nil {
            // handle not founds differently
            if err == ErrNotFound {
                fmt.Printf("%s (type: %d) | not found\n", question.Name, question.Type)
                continue
            } else {
                return nil, err
            }
        }

        // append the record we found to any previously found
        result = append(result, recs...)
    }

    return result, nil
}

//
// Returns the records in the cache that match the label queried
// With the naive map implementation, label lookup is O(1)
//
func FindRecordsByLabel(label string) []record.Record {
    if records, exists := cache[strings.TrimSuffix(label, ".")] ; exists {
        return records
    } else {
        return nil
    }
}

//
// Returns the records in the cache that match the label and type queried
//
func FindRecords(label string, rType uint16, recursively bool) ([]record.Record, error) {
    // naive iterative search dependent on FindRecordsByLabel
    var records = FindRecordsByLabel(label)
    if records == nil {
        // welp, 404
        return nil, errors.New("ERROR: No records for label: " + label)
    }

    // we now have all the records for that label in O(1) time
    // so filter out any non-matching record types
    var result = make([]record.Record, 0)
    for _, rec := range records {
        if rec.GetType() == rType {
            result = append(result, rec)
        }

        // handle any locally-recursive queries needed

        var err error
        var recursiveRecords []record.Record

        // if the current record is a CNAME
        if rec.GetType() == record.CNAME_RECORD && recursively {

            // but we were looking for and A/AAAA record... recurse
            if rType == record.A_RECORD || rType == record.AAAA_RECORD {
                // lookup A records
                ipv4Records, ipv4Err := FindRecords(rec.(*record.CNAMERecord).Target, record.A_RECORD, true)
                if ipv4Err != nil && ipv4Err != ErrNotFound { return nil, ipv4Err }

                // lookup AAAA records
                ipv6Records, ipv6Err := FindRecords(rec.(*record.CNAMERecord).Target, record.AAAA_RECORD, true)
                if ipv6Err != nil && ipv6Err != ErrNotFound { return nil, ipv6Err }

                // the CNAME comes first, then A, then AAAA
                recursiveRecords = []record.Record{rec}
                recursiveRecords = append(recursiveRecords, ipv4Records...)
                recursiveRecords = append(recursiveRecords, ipv6Records...)
            }

        }

        // cleanup recursive errors and append their results
        if err != nil && err != ErrNotFound { return nil, err }
        result = append(result, recursiveRecords...)
    }

    // if we found none, 404
    if len(result) < 1 {
        return nil, ErrNotFound
    }

    // yay, return records
    return result, nil
}
