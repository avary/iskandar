# Iskandar System Flow

A graph that helps understand the flow of requests/routines that are happening in the apps.

In summary, the most basic flow to understand the flow is the following:

**Setup**
On setting up the tunnel:

1. The CLI starts a WebSocket tunnel connection with the server.
2. The server creates and assigns a subdomain and sends it back to the client.
3. The server keeps a single blocking loop reading messages from the WebSocket.

**Request**
On a request to the tunnel domain:

1. The server catch-all handler sends the request information through the websocket connection to the CLI.
   On the same step creates a request-specific channel waiting for a response.
2. The CLI receives from the websocket the message and spawns a goroutine for each request and sends it on the target app.
3. When a response arrives from the target app to the request specific goroutine, the goroutine writes a message on the server websocket connection.
4. The blocking Setup.3 loop reads the message and routes the response information to the request specific channel waiting in Request.1 .
5. On the channel message arriving, the server sends back the appropriate response to the request.

## Connection Establishment

```
CLI                           Tunnel Server
 |                                   |
 |  HTTP GET /tunnel/connect         |
 | --------------------------------> |
 |                                   |
 |  WebSocket Upgrade (101)          |
 | <-------------------------------- |
 |                                   |
 |  Subdomain Registration           |
 | <-------------------------------- |
 |  {"subdomain": "abc123..."}       |
 |                                   |
 | ===== WebSocket Open ========     |
```

## Request Proxying Flow

```
HTTP Client          Tunnel Server              CLI                Local App
    |                      |                     |                     |
    | GET /api/data        |                     |                     |
    |--------------------> |                     |                     |
    |                      |                     |                     |
    |                      | Register Request    |                     |
    |                      | Channel (req-id)    |                     |
    |                      |                     |                     |
    |                      | WS: Request Message |                     |
    |                      | {"id": "req-id"...} |                     |
    |                      |-------------------> |                     |
    |                      |                     |                     |
    |                      |                     | [Spawn Goroutine]   |
    |                      |                     |                     |
    |                      |                     | HTTP Request        |
    |                      |                     |-------------------> |
    |                      |                     |                     |
    |                      |                     | HTTP Response       |
    |                      |                     | <------------------ |
    |                      |                     |                     |
    |                      | WS: Response Message|                     |
    |                      | {"id": "req-id"...} |                     |
    |                      | <------------------ |                     |
    |                      |                     |                     |
    |                      | Lookup Channel      |                     |
    |                      | Send to Channel     |                     |
    |                      |                     |                     |
    | HTTP Response        |                     |                     |
    | <------------------- |                     |                     |
    |                      |                     |                     |
    |                      | Remove Request      |                     |
    |                      | Channel (req-id)    |                     |
```

## Concurrent Request Handling

```
Tunnel Server                                     CLI
     |                                             |
     | WS: Request 1 (req-1)                       |
     |-------------------------------------------> |---> [Goroutine 1] ---> Local App
     |                                             |
     | WS: Request 2 (req-2)                       |
     |-------------------------------------------> |---> [Goroutine 2] ---> Local App
     |                                             |
     | WS: Request 3 (req-3)                       |
     |-------------------------------------------> |---> [Goroutine 3] ---> Local App
     |                                             |
     |                    WS: Response 2 (req-2)   |
     | <------------------------------------------ |<--- [Goroutine 2]
     |                                             |
     |                    WS: Response 1 (req-1)   |
     | <------------------------------------------ |<--- [Goroutine 1]
     |                                             |
     |                    WS: Response 3 (req-3)   |
     | <------------------------------------------ |<--- [Goroutine 3]
```

## Streaming Response Flow

```
HTTP Client          Tunnel Server              CLI                Local App
    |                      |                     |                     |
    | GET /stream          |                     |                     |
    |--------------------> |                     |                     |
    |                      | WS: Request         |                     |
    |                      |-------------------> |                     |
    |                      |                     | [Goroutine]         |
    |                      |                     |-------------------> |
    |                      |                     |                     |
    |                      | WS: Chunk 1         |    Stream Chunk 1   |
    |                      | done=false          | <------------------ |
    | HTTP Chunk 1         | <------------------ |                     |
    | <------------------- |                     |                     |
    |                      |                     |                     |
    |                      | WS: Chunk 2         |    Stream Chunk 2   |
    |                      | done=false          | <------------------ |
    | HTTP Chunk 2         | <------------------ |                     |
    | <------------------- |                     |                     |
    |                      |                     |                     |
    |                      | WS: Final Chunk     |    Stream End       |
    |                      | done=true           | <------------------ |
    | HTTP Complete        | <------------------ |                     |
    | <------------------- |                     |                     |
    |                      | Close Channel       |                     |
```
