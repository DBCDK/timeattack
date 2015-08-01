# timeattack

Small utility for loadtesting webservers.
Takes a list of (timestamp, url) as input, and performs the requests at the pace dictated by the timestamps. Requests can be sped up and down, or the timestamps can be ignored, which causes a flood of requests to be made as fast as possible.


## Example input

````
10.123 http://google.com
15.234 http://google.com
15.345 http://google.com
15.456 http://google.com
20.567 http://google.com
````
````
10 /
15 /
15 /
15 /
20 /
````


### Timestamps

Timestamps are parsed as floats with "seconds" as unit. The value of the first timestamp will be subtracted from all timestamps. That means that unix timestamps will work just fine, but it's also possible to have the first timestamp be '0' or whatever non-negative value floats your boat.


## Usage

 Input contains complete urls: ``$ timeattack run < urls``
 
 Input is missing part of the urls: ``$ timeattack run --prefix="http://google.com" < urls``

``--speedup`` can be used to change the pace of the test. A value of ``2`` will cut the time between requests in half, while a value of ``0.5`` will double the time between requests. Arbitrary positive values can be chosen.

``--flood`` disables the delay between requests. Use with care.

``--ramp-up=N`` ramps up trafic from 0% to 100% during a period of *n* seconds.


## Build-instructions

 1. install go
 2. run ``go build`` in the project root


## Why?

Once in a while we swap Solr-servers in and out of production. This utility was written to make it easy to ``tail -f`` a logfile on one server and forward the requests to a new server, thus testing how the new server would perform under the current load.


## Here be dragons

This is my first project written in Go. Pullrequests appreciated, but please don't learn from this code.
