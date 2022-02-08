package manager

import (
	"log"
	"net/http"
	"os"
	"sync"

	"fakepokerserver/module/db"
)

const api_game = "/api/poker/"
const api_version = "v1"

//api 추가
const (
	// 회원 관련
	URI_USER_INFO  = api_game + api_version + "/user/manage" // POST : 가입 GET : 조회 (ID 찾기 또는 비번) PUT : 수정 DELETE : 탈퇴
	URI_USER_LOGIN = api_game + api_version + "/user/login"  // GET : 로그인 요청

	URI_JOIN_ROOM = api_game + api_version + "/rooms/random"
	URI_EXIT_ROOM = api_game + api_version + "/rooms/"

	// 게임관련
	URI_GAME_READY   = api_game + api_version + "/game/ready"         // PUT : 입장 후 레디.
	URI_GAME_STATUS  = api_game + api_version + "/game/status/rooms/" // GET : 상태를 얻어옴.
	URI_GAME_ENDTURN = api_game + api_version + "/game/endturn"       // PUT : 턴을 마무리 한다.
	URI_GAME_RESULT  = api_game + api_version + "/game/result/rooms/" // GET : 결과 정리.
)

func RunningServer(wg *sync.WaitGroup, done chan int) {
	if wg != nil {
		defer wg.Done()
	}

	db.InitDB()
	OpenServer()
}

func OpenServer() {

	//회원 관련
	http.HandleFunc(URI_USER_INFO, Handler_userInfo())
	http.HandleFunc(URI_USER_LOGIN, Handler_login())

	//대기실 관련
	http.HandleFunc(URI_JOIN_ROOM, Handle_joinRoom())
	http.HandleFunc(URI_EXIT_ROOM, Handle_ExitRoom())

	// 게임관련
	http.HandleFunc(URI_GAME_READY, Handler_Game_Ready())
	http.HandleFunc(URI_GAME_STATUS, Handler_Game_Status())
	http.HandleFunc(URI_GAME_ENDTURN, Handler_Game_Endturn())
	http.HandleFunc(URI_GAME_RESULT, Handler_Game_Result())

	log.Println("server is running...")
	err := http.ListenAndServe(":44447", nil)
	if err != nil {
		log.Println("Failed to listen : ", err)
		os.Exit(1)
	}

}
