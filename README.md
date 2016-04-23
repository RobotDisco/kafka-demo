Hack Day Project - Kafka and Go
===============================

For my work's April 2016 Hack day I played around with [Apache Kafka](http://kafka.apache.org) to see what these distributed queue systems are all about.

I had originally wanted to do something in the [Go](http://golang.org) programming language before deciding to also play around with distributed queues.

So here's a bunch of code I whipped up in a day. I wasn't trying to solve any grand business problem or anything (although I suppose there are a few places in which my work could take advantage of them!)

You may want to look elsewhere for a good intro into Kafka and the platform it will end up participating in. I really enjoyed the Kafka site's [introduction](http://kafka.apache.org/documentation.html#introduction) and for a bigger picture I highly recommend Nathan Marz and James Warren's [Big Data in Action](https://www.manning.com/books/big-data)

How to Run
==========

I didn't do much beyond what the Kafka documentation will quickstart you into. I just did what the introductio told me to do for version 0.9.0


The Code
========

I built three little programs in go. I'm sure they are awful and clobber existing binaries, I have no idea how Go programs mantain binary namespaces :)

* producer - A little program that uses goroutines to simulate multiple clients sending data to Kafka. It takes a single argument, indicating the number of workers it should spin up.

* filteredconsumer - A program that only displays messages coming from a specific worker id. It takes a single argument, an integer of the worker number to focus on (My producer basically created workers 1-N, where N was the number you gave it.)

* dashboard - A Dashboard that streams all messages in the top window, the latestest 10 messages (I used a random name generator) every 5 seconds in the middle window, and a sum of the different "Platforms" of each message in the bottom window (I pretended each message came from a different source type, mobile/browser/server)


Anyway that's it! This is all stuff I learned in one day so don't expect it to be any good :)
