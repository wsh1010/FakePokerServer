package manager

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"syscall"
	"time"

	"fakepokerserver/module/db"
)

const api_game = "/api/poker/"
const api_version = "v1"

//api 추가
const (
	// 회원 관련
	URI_USER_INFO = api_game + api_version + "/user/info"     // POST : 가입 GET : 조회 (ID 찾기 또는 비번) PUT : 수정 DELETE : 탈퇴
	URI_ADD_COIN  = api_game + api_version + "/user/addCoins" // GET : 코인 추가.

	URI_JOIN_ROOM = api_game + api_version + "/rooms/random"
	URI_EXIT_ROOM = api_game + api_version + "/rooms/"

	// 게임관련
	URI_GAME_READY   = api_game + api_version + "/game/ready"         // PUT : 입장 후 레디.
	URI_GAME_STATUS  = api_game + api_version + "/game/status/rooms/" // GET : 상태를 얻어옴.
	URI_GAME_ENDTURN = api_game + api_version + "/game/endturn"       // PUT : 턴을 마무리 한다.
	URI_GAME_RESULT  = api_game + api_version + "/game/result/rooms/" // GET : 결과 정리.
)

func RunningServer(wg *sync.WaitGroup, done chan int) {
	defer log.Println("end RunningServer")
	if wg != nil {
		defer wg.Done()

	}

	db.InitDB()
	OpenServer()
}
func CheckDB(wg *sync.WaitGroup, done chan int) {
	defer log.Println("end CheckDB")
	if wg != nil {
		defer wg.Done()
	}
	for {
		timer1 := time.NewTimer(time.Minute * 10)
		db.CheckPing()
		query := "SELECT t_rooms_info.room_id, COUNT(DISTINCT t_users_gameinfo.id) FROM t_rooms_info LEFT JOIN t_users_gameinfo ON t_rooms_info.room_id = t_users_gameinfo.room_id GROUP BY t_rooms_info.room_id ;"
		rows, err := db.SelectQueryRows(query)
		if err != nil {
			log.Println(err)
			continue
		}
		defer rows.Close()

		for rows.Next() {
			var roomID string
			var count int
			rows.Scan(&roomID, &count)
			if count == 0 {
				query := fmt.Sprintf("DELETE FROM t_rooms_info WHERE room_id = '%s';", roomID)
				db.ExecuteQuery(query)
			}
		}
		count := 0
		for {
			timer2 := time.NewTimer(time.Second * 10)
			<-timer2.C
			count++
			if count == 59 {
				break
			}
			if len(done) == 0 {
				log.Println("return")
				return
			}
		}

		<-timer1.C
	}

}

func OpenServer() {

	//회원 관련
	http.HandleFunc(URI_USER_INFO, Handler_userInfo())

	//대기실 관련
	http.HandleFunc(URI_JOIN_ROOM, Handle_joinRoom())
	http.HandleFunc(URI_EXIT_ROOM, Handle_ExitRoom())
	http.HandleFunc(URI_ADD_COIN, handle_AddCoin())

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
	log.Println("end openServer ")
	syscall.Kill(syscall.Getpid(), syscall.SIGTERM)

}
func CloseServer() {

}
