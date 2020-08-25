Here
===
![](https://github.com/alpn/here/workflows/Go/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/alpn/here)](https://goreportcard.com/report/github.com/alpn/here)

As-Simple-As-Possible HTTP server for **local** use.

## Build/Install
```bash
git clone https://github.com/alpn/here.git
cd here 
go build

# optionally, move the binary to a PATH directory, e.g
mv here /usr/local/bin
```
## Usage

```bash

# serve current directory at the default port (9898)
here

# serve current directory at the port of your choice
here -port 80

``` 

## License
MIT
