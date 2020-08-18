# tcprobe

A simple tcp port scanner written in golang.

## Install 

```
go get github.com/13ph03nix/tcprobe
```

## Help

```
➜ ./tcprobe -h
Usage of ./tcprobe:
  -c int
    	set the concurrency level (default 128)
  -p value
    	add additional port probe
  -s	skip the address built-in port check
  -t int
    	timeout (milliseconds) (default 10000)

```

## Basic Usage

tcprobe accepts line-delimited addresses (host:port) on `stdin`:

```
➜ cat /tmp/addresses.txt
localhost:22
localhost:9001
192.168.50.1:80
192.168.50.1:22
192.168.50.1:23
192.168.50.97:80
192.168.50.97:443
192.168.50.121:22
192.168.50.121:80
192.168.50.177:22
192.168.50.179:22

➜ cat /tmp/addresses.txt | tcprobe
localhost:9001
192.168.50.1:22
192.168.50.1:80
192.168.50.121:22
192.168.50.97:80
192.168.50.179:22
192.168.50.177:22
192.168.50.121:80
```

## Advanced Usage

```
➜ echo '192.168.50.1/24' | cidr_to_ips | tcprobe -p 22 -p 80 -c 1000 -t 20000
192.168.50.1:80
192.168.50.1:22
192.168.50.97:80
192.168.50.177:22
192.168.50.121:80
192.168.50.121:22
192.168.50.179:22


➜ echo '192.168.50.1/24' | cidr_to_ips | tcprobe -p large -c 1000
192.168.50.1:445
192.168.50.1:80
192.168.50.1:22
192.168.50.1:139
192.168.50.97:23
192.168.50.97:80
192.168.50.121:80
192.168.50.121:22
192.168.50.177:22
192.168.50.179:22
...
```

## Credits

* [httprobe](https://github.com/tomnomnom/httprobe)
* [cidr_to_ips](https://github.com/13ph03nix/cidr_to_ips)
