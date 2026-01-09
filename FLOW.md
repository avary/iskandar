# Iskandar System Flow

A graph that helps understand the flow of requests/routines that are happening in the apps.

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
