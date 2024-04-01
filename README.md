# [ 환경 및 Command ]

## 설치 및 환경 변수

go 문서 참조.
문서:https://go.dev/doc/install

## go.mod 초기화

go mode init example/hello

## dependency 모듈 추가

- 현재 디렉토리 내부 코드 의존성을 가져온다.
  ```
  go get .
  ```

## go 실행

go run .

## Module install

go get rs.id/quote/v4 (모듈 이름)
문서: https://pkg.go.dev/

## go.sum 파일

checksum 파일을 포함하며 모듈의 모든 파일에 대한 해시값을 나타내며 이를 사용하여 모듈의 무결성으 확인.

### go 모듈 인증 과정

1. Checksums: 체크섬 확인
2. Module Signing: 모듈 소유자의 디지털 서명
3. Source Verification: HTTPS 모듈 다운로드

# [ Create a Go module ]

1. go mod init example.com/greetings
   - example.com/greetings: 모듈 이름
2. make file (Module)

```go
package greetings

import "fmt"

// Hello returns a greeting for the named person.
func Hello(name string) string {
    // Return a greeting that embeds the name in a message.
    message := fmt.Sprintf("Hi, %v. Welcome!", name)
    return message
}
```

4. go mod init example.com/hello

5. make file (Call your code from another module)

```go
package main

import (
    "fmt"

    "example.com/greetings"
)

func main() {
    // Get a greeting message and print it.
    message := greetings.Hello("Gladys")
    fmt.Println(message)
}
```

6. go mod edit -replace example.com/greetings=../greetings

   - The command specifies that example.com/greetings should be replaced with ../greetings for the purpose of locating the dependency.
   - the go.mod file in the hello directory should include a replace directive:

   ```
   module example.com/hello

   go 1.16

   replace example.com/greetings => ../greetings
   ```

7. go mod tidy

# [ Testing ]

1. Make file name: greetings_test.go
   - \_test.go: go test commands에게 이 파일이 test 함수들을 포함하다는 것을 전달한다.
2. Code:

   ```go
   package greetings

   import (
       "testing"
       "regexp"
   )

   // TestHelloName calls greetings.Hello with a name, checking
   // for a valid return value.
   func TestHelloName(t *testing.T) {
       name := "Gladys"
       want := regexp.MustCompile(`\b`+name+`\b`)
       msg, err := Hello("Gladys")
       if !want.MatchString(msg) || err != nil {
           t.Fatalf(`Hello("Gladys") = %q, %v, want match for %#q, nil`, msg, err, want)
       }
   }

   // TestHelloEmpty calls greetings.Hello with an empty string,
   // checking for an error.
   func TestHelloEmpty(t *testing.T) {
       msg, err := Hello("")
       if msg != "" || err == nil {
           t.Fatalf(`Hello("") = %q, %v, want "", error`, msg, err)
       }
   }
   ```

3. go test

   - 결과

   ```
   $ go test
   PASS
   ok      example.com/greetings   0.364s

   $ go test -v
   === RUN   TestHelloName
   --- PASS: TestHelloName (0.00s)
   === RUN   TestHelloEmpty
   --- PASS: TestHelloEmpty (0.00s)
   PASS
   ok      example.com/greetings   0.372s
   ```

# [ Build ]

## 빌드

1. 빌드 명령어: go build
2. 실행: ./hello

## 설치

1. 설치 경로 확인: go list -f '{{.Target}}'
   - Ex. /usr/local/go/bin/bin/hello -> 설치 경로: /usr/local/go/bin/bin
2. 설치 경로를 환경 변수에 추가
3. 컴파일 및 설치: go install
4. 실행: hello

# [ ADD ]

-
