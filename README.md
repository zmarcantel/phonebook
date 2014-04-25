phonebook
=========

A lightweight/minimal DNS server written in Go. [(godoc)](http://godoc.org/github.com/zmarcantel/phonebook)


Motivation
==========

The motivation for `phonebook` lies heavily in providing wrapper to a DNS server.

This initially sounds crazy, a wrapper to a server, but this offers two advantages:

1. Embeddable servers
2. Decoupling of cli/config interfaces from server routines

Allowing for domain-specific usage of DNS is key. The initial domain, for me, was service discovery ([whisper](https://gtihub.com/zmarcantel/whisper)). The use-cases for DNS in distributed systems is vast and largely uncharted.


Features
========

TODO
----

1. DNSSEC for verification of source


Records Supported
-----------------

1. `A` and `AAAA`
2. `SRV`
3. `CNAME`
4. `PTR`
5. `MX`
6. `TXT`


Server
------

1. Modular
    * The listener exists in its own thread separate from the calling context
    * Allow multiple listeners within the same process sharing a common error handler, pipeline, etc (if desired)
    * Even the data backing is pluggable! [modular storage](#modular-storage)
2. Fast
    * Every received packet/query is handled in an isolated thread
    * All operation are in memory so limited only by I/O speeds (network, task switching, memory latency)
3. Hackable
    * The codebase is designed to be extremely modular
    * Adding features, layers, extensions, etc become easier as the actual DNS wire protocol is entirely abstracted away while remaining reusable


Modular Storage
---------------

Storage of records will be domain specific.

Use redis, files, Mongo, Cassandra.... whatever you'd like.

The `server.Server` type includes a field `Store` that is of type `DNSStore`. This storage interface must support all the functions needed to query DNS records. However, this interface can be tweaked, expanded, or new ones created with no effect to the central server.

Similarly, an application can spin up two DNS servers querying records from two separate data sources within the same application (if you want).

For details on implementing your own `Store` check out the [(dnsstore godoc)](http://godoc.org/github.com/zmarcantel/phonebook/server/store) and the reference MapStore implementation.


Intentional Limitations
-----------------------

1. No recursion of DNS queries (security)
    * We will not query other nameservers to fullfil a query
        * Mitigate MITM attacks
        * Only serves records in local cache
    * `SRV`, `PTR`, and other relative records can be resolved to an IP using `A` or `AAAA` records
        * Those `A`/`AAAA` records must be loaded into the local cache


Usage
=====

The Golden Rule
---------------

DNS can be a serious security issue. __Always__ follow this rule:

1. Do not add `phonebook` to your machine's list of DNS servers
    * Inserting a malicious (overwriting) record could then cause a MITM
    * Always use a DNS client targeted at the single host machine

This is the case with any DNS server, but it __must__ be said. I cannot tell you enough how much you __should not__ do this.


Example
-------

Get a DNS server bound to `localhost` in 4 lines _(okay, adding records takes more)_:

````go
package main

import (
    "os"
    "fmt"
    "net"
    "time"

    dns    "github.com/zmarcantel/phonebook/dns/record"
    serve  "github.com/zmarcantel/phonebook/server"
)

func main() {
    // setup the server -- first, make a fatal error channel
    // then, start listening on localhost:53 (standard port)
    // the nil argument represents an optional [storage backing](#modular-storage) (defaults to map storage)
    // shorthand for the below is serve.Local(nil, lock)
    var lock = make(chan err)
    var server = serve.Start("localhost", 53, nil, lock)


    // create some records
    var productionRecords = recordSetA(server)
    var testingRecords = recordSetB(server)


    // erro handling
    var err := <-lock                      // blocks
    handleServerError(err)                 // not implemented here -- panic, print, whatever

    // dies on single error
    // wrap the above two lines in a "bare" (for{}) loop to go on forever
}

func recordSetA(server *serve.Server) []*dns.Record {
    // give the A record a label, TTL, and target IP
    a, err := dns.A("app.production", 10 * time.Second, net.ParseIP("10.0.8.15"))
    handleCreationErr("A", err)
    server.Store.Add(a)

    // give the AAAA record a label, TTL, and target IP
    b, err := dns.AAAA("ipv6.app.production", 10 * time.Second, net.ParseIP("2001:0db8:85a3:0042:1000:8a2e:0370:7334"))
    handleCreationErr("AAAA", err)
    server.Store.Add(b)

    // give the SRV record a label, target host, TTL, priority, weight, and port
    c, err := record.SRV("_logging._udp.app.production", "app.production", 10 * time.Second, 10, 5, 8053)
    handleCreationErr("SRV", err)
    server.Store.Add(c)

    return []*dns.Record { a, b, c }
}


func recordSetB(server *serve.Server) []*dns.Record {
    // give the A record a label, TTL, and target IP
    a, err := dns.A("app.test", 10 * time.Second, net.ParseIP("10.0.1.15"))
    handleCreationErr("A", err)
    server.Store.Add(a)

    // give the AAAA record a label, TTL, and target IP
    b, err := dns.AAAA("ipv6.app.production", 10 * time.Second, net.ParseIP("fe80:0000:0000:0000:0202:b3ff:fe1e:8329"))
    handleCreationErr("AAAA", err)
    server.Store.Add(b)

    // give the SRV record a label, target host, TTL, priority, weight, and port
    c, err := record.SRV("_logging._udp.app.test", "app.test", 10 * time.Second, 10, 5, 8053)
    handleCreationErr("SRV", err)
    server.Store.Add(c)

    return []*dns.Record { a, b, c }
}

func handleCreationErr(type string, err error) {
    if err != nil {
        fmt.Printf("ERROR: Could not create %s record: %s\n", type, err)
        os.Exit(1)
    }
}
````
