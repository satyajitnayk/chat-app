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

Certainly! Below is a crisp note accompanied by a simple dotted line diagram illustrating the ping-pong mechanism in WebSocket communication:

---

## WebSocket Ping-Pong/Heartbeat Mechanism:

WebSocket communication employs a ping-pong mechanism to sustain connection health between client and server:

- **Server Sends Ping**: Periodically, the server sends a ping frame to the client.
- **Client Responds with Pong**: Upon receiving the ping frame, the client automatically replies with a pong frame.
- **Acknowledgment by Server**: The server acknowledges the pong frame, ensuring the connection remains active.

## Security Concerns with websocket

### Security Considerations with Jumbo Frames in WebSocket:

Reason :larger frame size that exceeds the typical MTU (Maximum Transmission Unit) size of the network.

- **DDoS Amplification:** Jumbo frames can be exploited in DDoS attacks, amplifying the impact by overwhelming server resources with large payloads.
- **Resource Strain:** Handling jumbo frame payloads can strain server resources, potentially leading to denial of service if servers become overwhelmed.
- **Protocol Vulnerabilities:** Vulnerabilities in WebSocket implementations related to jumbo frame handling can be exploited for arbitrary code execution or server crashes.
- **Mitigation Measures:** Enforce maximum frame size limits, implement rate limiting, and regularly update software to mitigate risks associated with jumbo frames.
- **Best Practices:** Configure WebSocket servers securely, considering network settings, parameters, and access controls to minimize attack surface and enhance security.
