### Build an executable file
#### for Windows:
`go build -ldflags -H=windowsgui`

### Configuration file `config.json` with the following contents:
```
[
  {
    "name": "Without proxy"
  },
  {
    "name": "Anonymous proxy",
    "host": "10.10.10.10:8080"
  },
  {
    "name": "Proxy with authentication",
    "host": "10.10.10.11:8080",
    "user": {
      "username": "user",
      "password": "pass"
    }
  }
]
```
