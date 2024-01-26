# chat-app

A full-fledged chat application using golang &amp; js

- use username: satya & password: 1234 for login

### Generate self signing server certificates for https:// and wss://

- use `gencert.bash` script to generate

![generate_server_certificate](/assets/generate-certificates.png)

- When you use self-sign certificate then you will get alert `Your Connection is not private` as below screenshot
  ![connection_not_private](/assets/connection_not_private.png)

- Just click on advanced --> proceed to unsafe website

- You will see following things in terminal
  ![terminal_output](/assets/terminal_output.png)

> Alert: Do not push your key & certificate to repo. Pull them from somewhere secure.

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

### CORS Requirements in WebSocket:

- **Origin Policy**: WebSocket connections follow the browser's same-origin policy.
- **Cross-Origin Connections**: WebSocket connections to different domains are blocked by default.
- **CORS Headers**: Servers must include CORS headers like Access-Control-Allow-Origin to authorize cross-origin WebSocket connections.
- **Server Configuration**: WebSocket servers need proper configuration to add CORS headers and comply with browser security policies.

![cors_error](/assets/cors-error.png)

### Authentication in WebSocket applications

- **Token-Based Authentication**: Users provide JWT or OAuth tokens for validation during WebSocket connections.
- **Cookie-Based Authentication**: Authentication tokens or session identifiers are transmitted via cookies during WebSocket sessions.
- **HTTP Authentication Headers**: Clients send authentication credentials using standard HTTP authentication headers during the WebSocket handshake.
- **Custom Authentication Headers**: Clients transmit custom headers containing authentication information during WebSocket connections.
- **IP-Based Authentication**: Authentication is based on client IP addresses, useful for trusted environments but limited in dynamic settings.
- **Certificate-Based Authentication**: Clients present X.509 certificates during the WebSocket handshake for server verification, ensuring strong mutual authentication.
