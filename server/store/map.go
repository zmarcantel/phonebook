package store

import (
    "fmt"
    "strings"

    "github.com/zmarcantel/phonebook/dns/record"
)

type MapStore struct {
    Backing          map[string][]record.Record
    Labels           int64
    Records          int64
}

func Map() *MapStore {
    return &MapStore{
        make(map[string][]record.Record, 0),
        0,
        0,
    }
}

//
// Add a record to the naive map implementation so the record can be queried
//
func (self *MapStore) Add(rec record.Record) error {
    // input validation
    if rec == nil { return ErrNilRecord }

    // check if there are any other records sharing the label
    // if so, there is a map entry all ready so just add it
    var cleanLabel = strings.TrimSuffix(rec.GetLabel(), ".")
    if collection, exists := self.Backing[cleanLabel] ; exists {
        self.Backing[cleanLabel] = append(collection, rec)
    } else {
        self.Labels += 1
        self.Backing[cleanLabel] = []record.Record{ rec }
    }

    self.Records += 1
    return nil
}

//
// Delete a record from the map given the structural record
//
func (self *MapStore) Delete(rec record.Record) error {
    // input validation
    if rec == nil { return ErrNilRecord }

    // check if the record exists and if so, delete it
    var cleanLabel = strings.TrimSuffix(rec.GetLabel(), ".")
    if collection, exists := self.Backing[cleanLabel] ; exists {
        for i, curr := range collection {
            // if the labels and types match
            if curr.GetLabel() == rec.GetLabel() && curr.GetType() == rec.GetType() {
                // do an "in-place remove" of the record -- still O(n-i), but preserves added order
                // but first, null out the position in the slice the record was
                // records are pointers and will not get garbage collected if the entry is not nil'd
                self.Backing[cleanLabel][i] = nil
                self.Backing[cleanLabel] = append(self.Backing[cleanLabel][:i], self.Backing[cleanLabel][i + 1:]...)
                self.Records -= 1
                return nil
            }
        }
    }

    // either there was no collection at the label,
    // or the record did not exist in the collection.... either way 404
    return ErrNotFound
}


//
// Find a record from the map given the structural record
//
func (self *MapStore) Find(rLabel string, rType uint16) (record.Record, error) {
    // input validation
    if rLabel == "" { return nil, ErrNilRecord }

    // check if the record exists and if so, return it
    var cleanLabel = strings.TrimSuffix(rLabel, ".")
    if collection, exists := self.Backing[cleanLabel] ; exists {
        for _, curr := range collection {
            // if the labels and types match
            if curr.GetLabel() == rLabel && curr.GetType() == rType {
                return curr, nil
            }
        }
    }

    // either there was no collection at the label,
    // or the record did not exist in the collection.... either way 404
    return nil, ErrNotFound
}

//
// Find a collection of records from the map given the label
//
func (self *MapStore) FindLabel(rLabel string) ([]record.Record, error) {
    // input validation
    if rLabel == "" { return nil, ErrNilRecord }

    // check if the label exists and if so, return it
    var cleanLabel = strings.TrimSuffix(rLabel, ".")
    if collection, exists := self.Backing[cleanLabel] ; exists {
        return collection, nil
    }

    // either there was no collection at the label,
    // or the record did not exist in the collection.... either way 404
    return nil, ErrNotFound
}


//
// Delete a record from the map given the label and type
//
func (self *MapStore) FindAndDelete(rLabel string, rType uint16) error {
    // input validation
    if rLabel == "" { return ErrNilRecord }
    if rType == 0 { return ErrInvalidType }

    // check if the record exists and if so, delete it
    var cleanLabel = strings.TrimSuffix(rLabel, ".")
    if collection, exists := self.Backing[cleanLabel] ; exists {
        for i, curr := range collection {
            // if the labels and types match
            if curr.GetLabel() == rLabel && curr.GetType() == rType {
                // do an "in-place remove" of the record -- still O(n-i), but preserves added order
                // but first, null out the position in the slice the record was
                // records are pointers and will not get garbage collected if the entry is not nil'd
                self.Backing[cleanLabel][i] = nil
                self.Backing[cleanLabel] = append(self.Backing[cleanLabel][:i], self.Backing[cleanLabel][i + 1:]...)
                self.Records -= 1
                return nil
            }
        }
    }

    // either there was no collection at the label,
    // or the record did not exist in the collection.... either way 404
    return ErrNotFound
}


//
// Replace a record from the map given the label and type with a newer version
//
func (self *MapStore) FindAndReplace(rLabel string, rType uint16, newer record.Record) error {
    // input validation
    if rLabel == "" { return ErrNilRecord }
    if rType == 0 { return ErrInvalidType }

    // check if the record exists and if so, delete it
    var cleanLabel = strings.TrimSuffix(rLabel, ".")
    if collection, exists := self.Backing[cleanLabel] ; exists {
        for i, curr := range collection {
            // if the labels and types match
            if curr.GetLabel() == rLabel && curr.GetType() == rType {
                // do an "in-place replace" of the record -- still O(1)
                // first, null out the position in the slice the record was
                // records are pointers and will not get garbage collected if the entry is not nil'd
                self.Backing[cleanLabel][i] = nil
                self.Backing[cleanLabel][i] = newer
                return nil
            }
        }
    }

    // either there was no collection at the label,
    // or the record did not exist in the collection.... either way 404
    return ErrNotFound
}

//
// Find records recursively from the local collection
// This primarily applies to CNAME records
//
func (self *MapStore) FindRecursively(rLabel string, rType uint16) ([]record.Record, error) {
    fmt.Printf("Recursively looking for: %s (%d)\n", rLabel, rType)
    // input validation
    if rLabel == "" { return nil, ErrNilRecord }
    if rType == 0 { return nil, ErrInvalidType }

    // check if the record exists and if so, return it
    var cleanLabel = strings.TrimSuffix(rLabel, ".")
    if collection, exists := self.Backing[cleanLabel] ; exists {
        var result = make([]record.Record, 0)

        for _, curr := range collection {
            // if the labels match
            if curr.GetLabel() == rLabel {

                // if the current record is a CNAME -- lookup any A/AAAA records at the target
                if curr.GetType() == record.CNAME_RECORD {
                    // the CNAME always comes first
                    // also reflect it for convenience
                    result = append(result, curr)
                    var cname = curr.(*record.CNAMERecord)

                    // but we were looking for and A/AAAA record... recurse
                    if rType == record.A_RECORD || rType == record.AAAA_RECORD {
                        // lookup A records and append them
                        aRecord, err := self.Find(cname.Target, record.A_RECORD)
                        if err != nil && err != ErrNotFound { return nil, err }
                        result = append(result, aRecord)

                        // lookup AAAA records and append them
                        aaaaRecord, err := self.Find(cname.Target, record.AAAA_RECORD)
                        if err != nil && err != ErrNotFound { return nil, err }
                        result = append(result, aaaaRecord)
                    }
                } else if curr.GetType() == rType {
                    result = append(result, curr)
                }
            }
        }

        return result, nil
    }

    // either there was no collection at the label,
    // or the record did not exist in the collection.... either way 404
    return nil, ErrNotFound
}

func (self *MapStore) Size() int64 {
    return self.Records
}

func (self *MapStore) LabelSize(label string) int {
    // input validation
    if label == "" { return 0 }

    // check if the record exists and if so, return it
    var cleanLabel = strings.TrimSuffix(label, ".")
    if collection, exists := self.Backing[cleanLabel] ; exists {
        return len(collection)
    } else {
        return 0
    }
}


//
// Print the contents of the map to stdout
//
func (self *MapStore) Print() {
    fmt.Printf("\n\nMapStore Data:\n%+v\n\n", self.Backing)
}

