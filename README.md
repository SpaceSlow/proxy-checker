Формат конфигурационного файла `config.json`:
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
