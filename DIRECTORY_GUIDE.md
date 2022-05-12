# Project Structure
## Directory Structure
~~~
stockcontent-monitor-demo-back
    ├─ build
    ├─ core
    │   ├─ app
    │   ├─ config
    │   ├─ db
    │   ├─ echo
    │   │   └─ binder
    │   └─ lifecycle
    ├─ di
    │   └─ provides
    ├─ domain
    ├─ ent
    │   └─ schema
    ├─ util
    │   ├─ echox
    │   ├─ entx
    │   ├─ pointer
    │   ├─ safe
    └─ (domain_packages)
        ├─ handler
        ├─ repository
        └─ usecase
~~~

### Directory Description
#### build
빌드시에 셋팅될 ldflags 관리

#### core/app
어플리케이션 시작 포인트

#### core/config
어플리케이션 설정 값

#### core/db
데이터베이스 커넥션

#### core/echo
[Echo](https://echo.labstack.com/) 인스턴스 및 [Echo](https://echo.labstack.com/)의 기본적인 셋팅에 들어갈 코드추가

#### core/echo/binder
Controller 에 echo 인스턴스를 바인딩을 도와주는 부분

#### lifecycle
`core/app` 에서 어플리케이션 시작시 시작 직전의 `OnStart` 와 종료 직전의 `OnClose` 각각 호출 하는
라이프싸이클의 `OnStart`와 `OnClose`를 Provide 해주는 곳

#### di
의존성 주입

#### di/provides
기타 다른 패키지에 속하지 않은 상수, 스태틱 값 등을 provide 해주기 위한 곳

#### domain
도메인 모델을 정의하는 곳

#### ent/schema
테이블 정의

#### util/echox
[Echo](https://echo.labstack.com/) 관련 유틸리티 함수들

#### util/entx
[ent](https://entgo.io/) 관련 유틸리티 함수들

#### util/pointer
빌트인 타입들을 포인터로 캐스팅 해주는 함수들

- Example

```go
package main

import (
	"stockcontent-monitor-demo-back/util/pointer"
	"time"
)

func main() {
	pointer.Int64(1)
	pointer.String("hello world")
	pointer.Time(time.Now())
}
```

#### util/safe
빌트인 타입들이 포인터인 경우 `nil` 일 때 다른 안전한 값으로 가져오게 해주는 유틸리티 함수들 

- Example

```go
package main

import "stockcontent-monitor-demo-back/util/safe"

func main() {
	var i *int
	safe.IOrZero(i) // if i nil return 0
}
```

#### (domain_packages)/handler
컨트롤러 / 라우팅 / Presentation 영역

#### (domain_packages)/repository
domain model의 영속성 담당

#### (domain_packages)/usecase
domain model 실제 사용 및 비즈니스 로직을 실행하는 곳