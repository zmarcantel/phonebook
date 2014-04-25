package server

import (
    "github.com/zmarcantel/phonebook/dns"
    "github.com/zmarcantel/phonebook/dns/record"
)

const (
    DNS_QUERY_ALL  int       = 255
    DNS_QUERY_A    int       = 1
    DNS_QUERY_AAAA int       = 28
)

//
// Given a set of DNS queries, look into the local cache and get answers
//
func (self *Server) Answer(questions []dns.Question) ([]record.Record, error) {
    var result = make([]record.Record, 0)

    for _, question := range questions {

        switch (int(+question.Type)) { // cast to positive integer
            // if we are querying for ANY (255) then just lookup the label
            case DNS_QUERY_ALL:
                var collection, err = self.Store.FindLabel(question.Name)
                if err != nil { return nil, err }
                result = append(result, collection...)
                break

            case DNS_QUERY_A, DNS_QUERY_AAAA:
                var collection, err = self.Store.FindRecursively(question.Name, question.Type)
                if err != nil { return nil, err }
                result = append(result, collection...)
                break

            default:
                var single, err = self.Store.Find(question.Name, question.Type)
                if err != nil { return nil, err }
                result = append(result, single)
                break
        }
    }

    return result, nil
}
