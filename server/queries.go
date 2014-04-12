package server

import (
    "fmt"

    "github.com/zmarcantel/phonebook/dns"
    "github.com/zmarcantel/phonebook/dns/record"
)

const (
    DNS_QUERY_ALL int        = 255
)

func AnswerQuestions(questions []dns.Question) ([]record.Record, error) {
    var result = make([]record.Record, 0)

    for _, question := range questions {
        if int(+question.Type) == DNS_QUERY_ALL {
            var records = FindRecordsByLabel(question.Name)
            result = append(result, records...)
            continue
        }

        var rec, err = FindRecord(question.Name, question.Type, question.Class)
        if err != nil {
            if err == ErrNotFound {
                fmt.Printf("%s (type: %d) | not found\n", question.Name, question.Type)
                continue
            } else {
                return nil, err
            }
        }

        result = append(result, rec)
    }

    return result, nil
}
