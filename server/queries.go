package server

import (
    "fmt"

    "../dns"
    "../dns/record"
)

func AnswerQuestions(questions []dns.Question) ([]record.Record, error) {
    var result = make([]record.Record, 0)

    for _, question := range questions {
        var rec, err = FindRecord(question.Name, question.Type, question.Class)
        if err != nil {
            if err == ErrNotFound {
                fmt.Printf("%s|%d not found\n", question.Name, question.Type)
                continue
            } else {
                return nil, err
            }
        }

        result = append(result, rec)
    }

    return result, nil
}