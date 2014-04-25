package store

import (
    "errors"

    "github.com/zmarcantel/phonebook/dns/record"
)

var ErrNotFound     error   = errors.New("ERROR: That record does not exist")
var ErrNilRecord    error   = errors.New("ERROR: Cannot operate on nil record.")
var ErrInvalidType  error   = errors.New("ERROR: Invalid record type")

type DNSStore interface {
    // record interaction operations
    Add(record.Record) error
    Delete(record.Record) error
    Find(rLabel string, rType uint16) (record.Record, error)
    FindLabel(rLabel string) ([]record.Record, error)
    FindAndDelete(rLabel string, rType uint16) error
    FindAndReplace(rLabel string, rType uint16, newer record.Record) error
    FindRecursively(rLabel string, rType uint16) ([]record.Record, error)

    // statistics
    Size() int64
    LabelSize(label string) int
}