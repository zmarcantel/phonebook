phonebook
=========

A lightweight/minimal DNS server written in Go.


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

1. Allow address binding other than `localhost`
2. Support all standard (non-NS) record types
3. Modular cache backing
4. DNSSEC and signing of packets


Records Supported
-----------------

1. `A` and `AAAA`
    * Resolve a top-level name (example.com) to an IP (or IPv6 with `AAAA`)
2. `SRV`
    * Report an available service by pointing to an `A` or `AAAA` record to resolve IP of machine


Server
------

1. Modular
    * The listener exists in its own thread separate from the calling context
    * Allow multiple listeners within the same process sharing a common error handler, pipeline, etc (if desired)
2. Fast
    * Every received packet/query is handled in an isolated thread
    * All operation are in memory so limited only by I/O sppeds (network, task switching, memory latency)
3. Hackable
    * The codebase is designed to be extremely modular
    * Adding features, layers, extensions, etc become easier as the actual DNS wire protocol is entirely abstracted away while remaining reusable


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
    // create some records
    var productionRecords = recordSetA()
    var testingRecords = recordSetB()

    var lock = make(chan err)      // make an error channel
    serve.Start(lock)              // start listening on localhost

    var err := <-lock              // blocks
    handleServerError(err)         // not implemented here -- panic, print, whatever

    // dies on single error
    // wrap the above two lines in a "bare" (for{}) loop to go on forever
}

func recordSetA() []*dns.Record {
    // give the A record a label, TTL, and target IP
    a, err := dns.A("app.production", 10 * time.Second, net.ParseIP("10.0.8.15"))
    handleCreationErr("A", err)
    serve.AddRecord(a)

    // give the AAAA record a label, TTL, and target IP
    b, err := dns.AAAA("ipv6.app.production", 10 * time.Second, net.ParseIP("2001:0db8:85a3:0042:1000:8a2e:0370:7334"))
    handleCreationErr("AAAA", err)
    serve.AddRecord(b)

    // give the SRV record a label, target host, TTL, priority, weight, and port
    c, err := record.SRV("_logging._udp.app.production", "app.production", 10 * time.Second, 10, 5, 8053)
    handleCreationErr("SRV", err)
    serve.AddRecord(c)

    return []*dns.Record { a, b, c }
}


func recordSetB() []*dns.Record {
    // give the A record a label, TTL, and target IP
    a, err := dns.A("app.test", 10 * time.Second, net.ParseIP("10.0.1.15"))
    handleCreationErr("A", err)
    serve.AddRecord(a)

    // give the AAAA record a label, TTL, and target IP
    b, err := dns.AAAA("ipv6.app.production", 10 * time.Second, net.ParseIP("fe80:0000:0000:0000:0202:b3ff:fe1e:8329"))
    handleCreationErr("AAAA", err)
    serve.AddRecord(b)

    // give the SRV record a label, target host, TTL, priority, weight, and port
    c, err := record.SRV("_logging._udp.app.test", "app.test", 10 * time.Second, 10, 5, 8053)
    handleCreationErr("SRV", err)
    serve.AddRecord(c)

    return []*dns.Record { a, b, c }
}

func handleCreationErr(type string, err error) {
    if err != nil {
        fmt.Printf("ERROR: Could not create %s record: %s\n", type, err)
        os.Exit(1)
    }
}
````
