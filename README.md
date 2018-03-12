# malgo
CCSDS MO MAL GO API

This project is an implementation of the [CCSDS MO Message Abstraction Layer (MAL) Standard](https://en.wikipedia.org/wiki/CCSDS_Mission_Operations) in Go language.

CCSDS Mission Operation implementations for other languages (e.g. Go, Java, etc.) can be found on the [CCSDS MO WebSite](http://ccsdsmo.github.io/)

Complete MAL specification can be found on the [CCSDS website](http://public.ccsds.org/publications/BlueBooks.aspx) in the *published documents* section.

In particular:

- [CCSDS 521.0-B-2, Mission Operations Message Abstraction Layer](https://public.ccsds.org/Pubs/521x0b2e1.pdf)
- [CCSDS 524.1-B-1, Mission Operations--MAL Space Packet Transport Binding and Binary Encoding](https://public.ccsds.org/Pubs/524x1b1.pdf)
- [CCSDS 524.2-B-1, Mission Operations--MAL Binding to TCP/IP Transport and Split Binary Encoding](https://public.ccsds.org/Pubs/524x2b1.pdf)

## ABOUT

This CCSDS MO MAL Go API was originally developed for the [CNES](http://cnes.fr), the French Space Agency, by [ScalAgent](http://www.scalagent.com/en/), a french company specialized in distributed technologies. All contributions are welcome.

## PROJECT DOCUMENTATION

A MAL/GO description and user's guide is available in the doc directory.

### MAL/GO Description

This GO API basically includes 4 packages:

  - **mal** package defines all MAL Concepts: message, data types, etc.
  - **mal/encoding** package includes encoding technologies.
  - **mal/transport** package includes transport technologies.
  - **mal/api** defines the high level consumer and provider APIs.

### MAL/GO QUICK INSTALLATION

### MAL/GO TEST

```
options: -v -timeout 1m

cd src
go test github.com/ccsdsmo/malgo/mal/encoding/binary
go test github.com/ccsdsmo/malgo/mal/encoding/splitbinary
go test github.com/ccsdsmo/malgo/mal/transport/invm
go test github.com/ccsdsmo/malgo/mal/transport/tcp
go test github.com/ccsdsmo/malgo/mal/api
go test github.com/ccsdsmo/malgo/tests/encoding
go test github.com/ccsdsmo/malgo/tests/issue1
```
