## Concurrent http requests for Shield

go build -o client concurrent_get_requests.go
./client -u http://172.31.132.168:8080/urlinfo/1/google.com:8080/search?v=2  (Default to 500 concurrent requests. pass -n 1000 for 1000 requests in parallel).

```
SRKAUSHI-M-8A2X:examples srkaushi$ ./client -u http://172.31.132.168:8080/urlinfo/1/google.com:8080/search?v=2 -n 5
Response status code from http://172.31.132.168:8080/urlinfo/1/google.com:8080/search?v=2: 200
Response status code from http://172.31.132.168:8080/urlinfo/1/google.com:8080/search?v=2: 200
Response status code from http://172.31.132.168:8080/urlinfo/1/google.com:8080/search?v=2: 200
Response status code from http://172.31.132.168:8080/urlinfo/1/google.com:8080/search?v=2: 200
Response status code from http://172.31.132.168:8080/urlinfo/1/google.com:8080/search?v=2: 200
All requests completed.
SRKAUSHI-M-8A2X:examples srkaushi$ 
```
