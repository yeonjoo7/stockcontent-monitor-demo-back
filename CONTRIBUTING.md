# 기여자 가이드
## 디렉토리 가이드
[보러가기](./DIRECTORY_GUIDE.md)

## 도메인 생성 및 관리
### 1. 엔티티 생성
```bash
$ make entity name={도메인이름}
```

- 필드 정의
- 릴레이션 정의
- **[Annotation](https://entgo.io/docs/schema-annotations#custom-table-name), 테이블명 명시 필수**

```bash
$ make entity-gen
```

### 2. 도메인 생성
```bash
$ make domain name={도메인이름}
```

생성된 파일 리스트 확인
- `domain/{도메인이름}.go`
- `{도메인이름}/handler/http.go`
- `{도메인이름}/usecase/usecase.go`
- `{도메인이름}/repository/mysql_repo.go`

## Code Convention
작성중