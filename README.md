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

```
$ timeattack --help-long
usage: timeattack [<flags>] <command> [<args> ...]

Flags:
  --help           Show help (also see --help-long and --help-man).
  --prefix=http://example.com  
                   String to prepend to urls.
  --ramp-up=0      Increase the amount of requests let through over a number of seconds.
  --concurrency=0  Allowed concurrent requests. 0 is unlimited.
  --limit=0        Maximum number of requests that will be sent. 0 is unlimited.

Commands:
  help [<command>...]
    Show help.


  run timed [<flags>]
    Schedule requests based on timestamps.

    --speedup=1  Change replay speed; 1 is 100%.

  run flood
    Ignore timestamps (send requests as fast as possible).


  run ticker [<flags>]
    Schedule requests with a certain frequency.

    --freq=1  Ticker frequency [1/s].

  parse solr
    Parses solr logs into suitable format.
```

## Build-instructions

 1. install go
 2. run ``go build`` in the project root

