# HTTP from TCP

Servidor HTTP construido desde cero usando TCP en Go.

## 🚀 Cómo ejecutar

```bash
go run cmd/httpserver/main.go
```

## 🔧 Debugging de chunks

Para debuggear los chunks y ver los datos recibidos:

```bash
# HTML endpoint
echo -e "GET /httpbin/html HTTP/1.1\r\nHost: localhost:42069\r\nConnection: close\r\n\r\n" | nc localhost 42069

# Stream endpoint
echo -e "GET /httpbin/stream/100 HTTP/1.1\r\nHost: localhost:42069\r\nConnection: close\r\n\r\n" | nc localhost 42069
```

## 🌐 Endpoints disponibles

| Endpoint | Descripción |
|----------|-------------|
| `http://localhost:42069/video` | Ver video en el navegador |
| `http://localhost:42069/httpbin/stream/100` | Stream de 100 chunks |
| `http://localhost:42069/httpbin/html` | Proxy a httpbin HTML |
| `http://localhost:42069/yourproblem` | Error 400 Bad Request |
| `http://localhost:42069/myproblem` | Error 500 Internal Server Error |
| `http://localhost:42069/` | Respuesta exitosa 200 OK |

## 🧪 Probar con curl

```bash
curl -v http://127.0.0.1:42069/yourproblem
curl -v http://127.0.0.1:42069/httpbin/html
```

