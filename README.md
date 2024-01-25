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

- Read more about websocket [RFC6455](https://datatracker.ietf.org/doc/html/rfc6455)

### Optimizing Gorilla WebSocket Handling with Unbuffered Channels

- Reading and writing messages on Gorilla WebSocket is easy but limited to one concurrent writer, causing potential issues with heavy traffic or spam.
- Gorilla WebSocket's single concurrent writer limitation can lead to problems when handling a high volume of clients or messages simultaneously.
- We can address this limitation by using an unbuffered channel, preventing the connection from being overloaded and ensuring controlled message flow in challenging scenarios.
- for more see [documentation](https://github.com/gorilla/websocket/blob/main/examples/chat/client.go#L47)

### Message types in websocket

| Opcode | Meaning                | Description                                |
| ------ | ---------------------- | ------------------------------------------ |
| 0      | Continuation Frame     | Fragment of a message from the server      |
| 1      | Text Frame             | UTF-8 encoded text data                    |
| 2      | Binary Frame           | Frames carrying binary data                |
| 8      | Connection Close Frame | Indicates a desire to close the connection |
| 9      | Ping Frame             | Sent to check the health of the connection |
| 10     | Pong Frame             | Response confirming a healthy connection   |

## Connection States in websocket:

| State      | Description                                                     |
| ---------- | --------------------------------------------------------------- |
| CONNECTING | Initial state when a WebSocket connection is being established. |
| OPEN       | The WebSocket connection is open and ready for communication.   |
| CLOSING    | The connection is in the process of closing.                    |
| CLOSED     | The WebSocket connection has been closed.                       |
