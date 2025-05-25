## HTTP Server from scratch

### Introduction

A project that implements a small part of the HTTP protocol. It uses TCP as its base and then implements the protocol on top of it. At its current state, it can serve static files such as HTML and CSS. 

### Features

Since the project was experimental, I haven't really implemented a lot of features you'd find in a complete HTTP library. 
You have a ```server.Get(route, path-to-file)``` function which you call in the main func which routes your path to a file on your server. See ```main.go```.

You can provide a ```notfound.html``` for routes that don't exist, but the server automatically handles that as well if none is found.

Every file you want to be served should have a Get() func for it.

### Metrics

Using wrk to benchmark the performance.
I managed to get
```
~194519 Requests/sec
~115.57 Transfer/sec

Socket Errors: timeout 981
Total Requests: 6,546,526 in 33.65s

Cmd for wrk: wrk -t12 -c1000 -d30s --latency http://localhost:8080/

12 threads
1000 concurrent connections
30 secs test time
```

CPU usage does get high, so i may improve that. 