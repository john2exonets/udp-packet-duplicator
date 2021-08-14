# udp-packet-duplicator
Pass UDP Packets to multiple downstream UDP servers.

![udpdup usage](/udp-dup-architecture.png)

In the `config.json` file, you define what port to open for incoming UDP connections. Then in the `dests` section, 
define all the nodes and ports to send all incoming UDP packets to.

Supports both IPv4 and IPv6.
