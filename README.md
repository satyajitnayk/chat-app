# chat-app

A full-fledged chat application using golang &amp; js

```
Client                            Server
  |                                 |
  |    HTTP (Upgrade Request)       |
  | ------------------------------> |
  |                                 |
  |   WebSocket Supported           |
  |   HTTP (101 Switching Protocol) |
  | <-------------------------------|
  |                                 |
  |  WebSocket Communication        |
  | <---------------------------->  |
  |   WebSocket Close               |
  | ------------------------------> |

```

![websockt_connection](/assets/websocket-coonection.png)
