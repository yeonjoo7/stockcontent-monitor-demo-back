# stockcontent-monitor-demo-back
# How to Start

## requires
- [make](https://www.gnu.org/software/make/)
- [golang 1.17 🔺](https://golang.org/)
- [mysql 5 🔺](https://www.mysql.com/)

## 1. Initialize Command
```bash
$ pwd
(~someDirectoryPath~)/stockcontent-monitor-demo-back
$ make init
...
$ make gen
...
```

## 2. Add Config file

프로젝트 루트 폴더에 `config.local.json` 파일 추가

```bash
$ make cfg
```

### `config.local.json` data structure
```json
{
  "db": {
    "user": "root",       // string, 디비유저
    "pass": "1234",       // string, 디비비밀번호
    "host": "localhost",  // string, 디비호스트
    "port": 3306,         // uint16, 디비포트
    "name": "demo",       // string, 디비이름
    "query_values": {     // fixed
      "charset": ["utf8mb4"], // 문자/문자열 포맷
      "parseTime": ["true"],  // 시간 파싱
      "loc": ["UTC"]          // Timezone
    }
  },
  "serve_addr": ":8000",    // string, listen 주소
  "use_case_timeout": "3s"  // string - 3s == 3초
}
```

## 3. Run
```bash
// Type-1
$ make go-run

// Type-2
$ go run .
```

## Used
### HTTP Router
[Echo Framework](https://echo.labstack.com/)

### ORM
[ent](https://entgo.io/)

### Etc
- [google/wire](https://github.com/google/wire) - Compile-time Dependency Injection for Go
- ~~[go-playground/validator](https://github.com/go-playground/validator) - About 💯Go Struct and Field validation, including Cross Field, Cross Struct, Map, Slice and Array diving~~ 준비중
    - ~~[document](https://pkg.go.dev/github.com/go-playground/validator/v10)~~
- ~~[swagger](https://swagger.io/) - Automatically generate RESTful API documentation with Swagger 2.0~~ 준비중
    - ~~[swaggo/swag](https://github.com/swaggo/swag#declarative-comments-format)~~
    - ~~[swaggo/echo-swagger](https://github.com/swaggo/echo-swagger)~~
- ~~[sirupsen/logrus](https://github.com/sirupsen/logrus) - Structured, pluggable logging for Go.~~ 준비중

## Contributors Guide
[보러가기](./CONTRIBUTING.md)

# License
[`MIT License`](./LICENSE)