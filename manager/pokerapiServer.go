package manager

import (
	"log"
	"net/http"
	"os"
)

const api_game = "/api/porker/"
const api_version = "v1"

//api 추가
const (
	// 회원 관련
	URI_USER_INFO  = api_game + api_version + "/user/manage" // POST : 가입 GET : 조회 (ID 찾기 또는 비번) PUT : 수정 DELETE : 탈퇴
	URI_USER_LOGIN = api_game + api_version + "/user/login"  // GET : 로그인 요청

	// 게임관련
)

func OpenServer() {
	http.HandleFunc(URI_USER_INFO, Handler_userInfo())
	http.HandleFunc(URI_USER_INFO, Handler_login())

	err := http.ListenAndServe(":44447", nil)
	if err != nil {
		log.Println("Failed to listen : ", err)
		os.Exit(1)
	}
}
