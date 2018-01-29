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

### MAL/GO QUICK INSTALLATION

### MAL/GO TEST

```
options: -v -timeout 1m

go test ./src/mal/encoding/binary
go test ./src/mal/encoding/splitbinary
go test ./src/mal/transport/invm
go test ./src/mal/transport/tcp
go test ./src/mal/api
```
