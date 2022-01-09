package manager

import (
	"log"
	"net/http"
	"os"
)

//api 추가
const (
	URI_USER_INFO  = "/api/poker/v1/user"       // POST : 가입 GET : 조회 (ID 찾기 또는 비번) PUT : 수정 DELETE : 탈퇴
	URI_USER_LOGIN = "/api/poker/v1/user/login" // GET : 로그인 요청

)

func OpenServer() {
	http.HandleFunc(URI_USER_INFO, Handler_userInfo())

	err := http.ListenAndServe(":44447", nil)
	if err != nil {
		log.Println("Failed to listen : ", err)
		os.Exit(1)
	}
}
