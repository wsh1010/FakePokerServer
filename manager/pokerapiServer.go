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
	URI_GAME_READY   = api_game + api_version + "/game/ready"   // POST : 입장 후 레디.
	URI_GAME_STATUS  = api_game + api_version + "/game/status"  // GET : 상태를 얻어옴.
	URI_GAME_ENDTURN = api_game + api_version + "/game/endturn" // POST : 턴을 마무리 한다.
)

func OpenServer() {

	//회원 관련
	http.HandleFunc(URI_USER_INFO, Handler_userInfo())
	http.HandleFunc(URI_USER_INFO, Handler_login())

	// 게임관련
	/*http.HandleFunc(URI_GAME_READY, Handler_Game_Ready())
	http.HandleFunc(URI_GAME_STATUS, Handler_Game_Status())
	http.HandleFunc(URI_GAME_ENDTURN, Handler_Game_Endturn())*/

	err := http.ListenAndServe(":44447", nil)
	if err != nil {
		log.Println("Failed to listen : ", err)
		os.Exit(1)
	}
}
