[![Build Status](https://travis-ci.com/apostergiou/arachne.svg?branch=master)](https://travis-ci.com/apostergiou/arachne)

# arachne

                             _
        /\                  | |
       /  \   _ __ __ _  ___| |__  _ __   ___
      / /\ \ | '__/ _` |/ __| '_ \| '_ \ / _ \
     / ____ \| | | (_| | (__| | | | | | |  __/
    /_/    \_\_|  \__,_|\___|_| |_|_| |_|\___|

arachne is a DNS over TLS proxy.

It listens for conventional DNS requests and resolves them using TLS.

         +------------------+          +-----------------+           +-----------------+
         |                  |          |                 |           |                 |
         |   DNS Reqeust    +--------->+    Arachne      +---------->+  DNS/TLS server |
         |                  |          |                 |           |                 |
         +------------------+          +-----------------+           +-----------------+

                             DNS tcp/53                     DNS/TLS

## Quick start

```
$ make start

Building image...

2018/12/09 20:04:52 Starting arachne
[server] 2018/12/09 20:04:52 Listening on 127.0.0.1:5300 (tcp)
[server] 2018/12/09 20:04:52 Forwarding to 1.1.1.1:853
[server] 2018/12/09 20:04:52 Configuration: &main.Config{Listen:"127.0.0.1:5300", Upstream:"1.1.1.1:853", Network:"tcp"}
```

This will start arachne with the resolvers set to Cloudflare's 1.1.1.1.

Send a DNS query as normal (from another shell):

```
$ kdig -d @localhost -p 5300 +tcp apostergiou.com
;; DEBUG: Querying for owner(apostergiou.com.), class(1), type(1), server(localhost), port(5300), protocol(TCP)
;; ->>HEADER<<- opcode: QUERY; status: NOERROR; id: 63742
;; Flags: qr rd ra; QUERY: 1; ANSWER: 1; AUTHORITY: 0; ADDITIONAL: 0

;; QUESTION SECTION:
;; apostergiou.com.             IN      A

;; ANSWER SECTION:
apostergiou.com.        2947    IN      A       46.4.121.137

;; Received 64 B
;; Time 2018-12-09 23:50:47 CET
;; From 127.0.0.1@5300(TCP) in 157.7 ms
```

Capture traffic and notice TLS usage:

```
$ tshark -f "tcp port 853"

1 0.000000000 192.168.2.101 → 1.1.1.1      TCP 74 56166 → 853 [SYN] Seq=0 Win=29200 Len=0 MSS=1460 SACK_PERM=1 TSval=3530330615 TSecr=0 WS=128
2 0.023942864      1.1.1.1 → 192.168.2.101 TCP 66 853 → 56166 [SYN, ACK] Seq=0 Ack=1 Win=29200 Len=0 MSS=1412 SACK_PERM=1 WS=1024
3 0.024022667 192.168.2.101 → 1.1.1.1      TCP 54 56166 → 853 [ACK] Seq=1 Ack=1 Win=29312 Len=0
4 0.024191502 192.168.2.101 → 1.1.1.1      TLSv1 196 Client Hello
5 0.044670267      1.1.1.1 → 192.168.2.101 TCP 54 853 → 56166 [ACK] Seq=1 Ack=143 Win=30720 Len=0
6 0.047519885      1.1.1.1 → 192.168.2.101 TLSv1.2 2254 Server Hello, Certificate, Server Key Exchange, Server Hello Done
7 0.047548909      1.1.1.1 → 192.168.2.101 TCP 54 [TCP Dup ACK 5#1] 853 → 56166 [ACK] Seq=2201 Ack=143 Win=30720 Len=0
8 0.047587281 192.168.2.101 → 1.1.1.1      TCP 54 56166 → 853 [ACK] Seq=143 Ack=2201 Win=33664 Len=0
9 0.062991329 192.168.2.101 → 1.1.1.1      TLSv1.2 147 Client Key Exchange, Change Cipher Spec, Encrypted Handshake Message
10 0.105217263      1.1.1.1 → 192.168.2.101 TLSv1.2 105 Change Cipher Spec, Encrypted Handshake Message
11 0.105405229 192.168.2.101 → 1.1.1.1      TLSv1.2 118 Application Data
12 0.156837233      1.1.1.1 → 192.168.2.101 TLSv1.2 134 Application Data
13 0.157012031 192.168.2.101 → 1.1.1.1      TLSv1.2 85 Encrypted Alert
14 0.157049549 192.168.2.101 → 1.1.1.1      TCP 54 56166 → 853 [FIN, ACK] Seq=331 Ack=2332 Win=33664 Len=0
15 0.207019388      1.1.1.1 → 192.168.2.101 TCP 66 [TCP Dup ACK 12#1] 853 → 56166 [ACK] Seq=2332 Ack=300 Win=30720 Len=0 SLE=331 SRE=332
16 0.220087808      1.1.1.1 → 192.168.2.101 TCP 54 853 → 56166 [ACK] Seq=2332 Ack=332 Win=30720 Len=0
17 0.220120682      1.1.1.1 → 192.168.2.101 TCP 54 853 → 56166 [FIN, ACK] Seq=2332 Ack=332 Win=30720 Len=0
18 0.220169348 192.168.2.101 → 1.1.1.1      TCP 54 56166 → 853 [ACK] Seq=332 Ack=2333 Win=33664 Len=0
^C18 packets captured
```

Server logs:

```
[server] 2018/12/09 21:29:46 Received DNS query for: apostergiou.com
[server] 2018/12/09 21:29:46 Forwarding the query to: 1.1.1.1:853
[server] 2018/12/09 21:29:46 Upstream answer: [apostergiou.com. 7809    IN      A       46.4.121.137]
```

## Tuning arachne

List available commands:

```
$ make

Available commands:
  build    - build the image
  arachned - start the daemon
  start    - start the arachne server locally
  test     - run the tests in the container
  install  - install the library as binary
  deps     - check dependencies
  lint     - run the golint tool
```

To use other options:

```
$ make start LISTEN="0.0.0.0:53" UPSTREAM="1.1.1.0:853" HOST_PORT="53"

2018/12/09 23:25:54 Starting arachne
[server] 2018/12/09 23:25:54 Listening on 0.0.0.0:53 (tcp)
[server] 2018/12/09 23:25:54 Forwarding to 1.1.1.0:853
[server] 2018/12/09 23:25:54 Configuration: &main.Config{Listen:"0.0.0.0:53", Upstream:"1.1.1.0:853", Network:"tcp"}
```

Start in the background:

```
$ make arachned

$ docker ps
CONTAINER ID        IMAGE               COMMAND             CREATED             STATUS              PORTS                  NAMES
432b285273aa        arachne             "/app/arachne"      23 seconds ago      Up 22 seconds       0.0.0.0:5300->53/tcp   serene_gauss

$ kdig -d @localhost -p 5300 +tcp apostergiou.com
;; DEBUG: Querying for owner(apostergiou.com.), class(1), type(1), server(localhost), port(5300), protocol(TCP)
;; ->>HEADER<<- opcode: QUERY; status: NOERROR; id: 63742
;; Flags: qr rd ra; QUERY: 1; ANSWER: 1; AUTHORITY: 0; ADDITIONAL: 0

;; QUESTION SECTION:
;; apostergiou.com.             IN      A

;; ANSWER SECTION:
apostergiou.com.        2947    IN      A       46.4.121.137

;; Received 64 B
;; Time 2018-12-09 23:50:47 CET
;; From 127.0.0.1@5300(TCP) in 157.7 ms

$ docker logs -f 432b285273aa
2018/12/09 22:33:42 Starting arachne
[server] 2018/12/09 22:33:42 Listening on 0.0.0.0:53 (tcp)
[server] 2018/12/09 22:33:42 Forwarding to 1.1.1.1:853
[server] 2018/12/09 22:33:42 Configuration: &main.Config{Listen:"0.0.0.0:53", Upstream:"1.1.1.1:853", Network:"tcp"}
[server] 2018/12/09 22:35:08 Received DNS query for: apostergiou.com
[server] 2018/12/09 22:35:08 Forwarding the query to: 1.1.1.1:853
[server] 2018/12/09 22:35:08 Upstream answer: [apostergiou.com. 3887    IN      A       46.4.121.137]
```

## Test suite

To run the tests execute:

```
$ make test

=== RUN   TestSetupConfig
--- PASS: TestSetupConfig (0.00s)
=== RUN   TestMissingConfig
--- PASS: TestMissingConfig (0.00s)
PASS
ok      app     0.003s
```

## Security Concerns

Quoting from RFC7858:

>Use of DNS over TLS is designed to address the privacy risks that
arise out of the ability to eavesdrop on DNS messages.  It does not
address other security issues in DNS.

1. Person-in-the-middle and downgrade attack
2. Middleboxes
3. Traffic analysis and side-channel leaks
4. Server certificate may get edited by a malicious user

## Future

* Handle UDP requests
* Support caching queries
* Display DNS requests in a Web GUI
* Use Prometheus for metrics
* DNSSEC implementation
* Allow usage of multiple upstream resolvers
* Add validations
* Add domain specific error handling

## References

* [cloudflare docs](https://developers.cloudflare.com/1.1.1.1/dns-over-tls/)
* [domain names RFC](https://tools.ietf.org/html/rfc1035)
* [DNS over TLS RFC](https://tools.ietf.org/html/rfc7858)

## Author

Apostolis Stergiou, apostergiou@gmail.com
