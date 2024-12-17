package gatewayserver

import (
	"context"
	"encoding/json"
	"log"
	pb "quiz/api/gen/quiz"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var upgrader = websocket.Upgrader{}

func (srv *Server) RegisterRoute(e *echo.Group) error {
	e.GET("/quiz/ws", srv.joinRoomWS)

	return nil
}

func (srv *Server) joinRoomWS(c echo.Context) error {
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer conn.Close()

	// assume this will read authentication info and get quizID and userID from quiz service
	// in this demo i will only read user_id and quiz id direct from header
	userID := c.Request().Header.Get("user_id")
	quizID := c.Request().Header.Get("quiz_id")

	srv.hub.joinRoom(userID, quizID, conn)
	defer srv.hub.leaveRoom(userID, quizID)

	for {
		var data WebsocketMessage[any]
		// Read
		err := conn.ReadJSON(&data)
		if err != nil {
			c.Logger().Error(err)
			conn.WriteMessage(websocket.TextMessage, []byte("invalid request!"))
			continue
		}

		srv.handleMessage(userID, quizID, data, conn)
	}
}

func (srv *Server) handleMessage(userID, quizID string, message WebsocketMessage[any], conn *websocket.Conn) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	switch message.Type {
	case SubmitAnswer:
		// parse msg
		msgByte, err := json.Marshal(message.Message)
		if err != nil {
			log.Println(err)
			return
		}

		var msg SubmitAnwserMessage
		if err := json.Unmarshal(msgByte, &msg); err != nil {
			log.Println(err)
			return
		}

		resp, err := srv.quizService.SubmitAnswer(ctx, &pb.AnswerRequest{
			UserID:     userID,
			QuizID:     quizID,
			QuestionID: msg.QuestionID,
			Answer:     msg.Anwser,
		})
		if err != nil {
			log.Println(err)
			return
		}

		if err := conn.WriteJSON(WebsocketMessage[*pb.AnswerResponse]{
			Type:    SubmitAnswerResp,
			Message: resp,
		}); err != nil {
			log.Println(err)
			return
		}
	case ListLeaderboardRequest:
		resp, err := srv.quizService.GetLeaderboard(ctx, &pb.LeaderboardRequest{
			QuizID: quizID,
		})
		if err != nil {
			log.Println(err)
			return
		}

		if err := conn.WriteJSON(WebsocketMessage[*pb.LeaderboardResponse]{
			Type:    ListLeaderboardResponse,
			Message: resp,
		}); err != nil {
			log.Println(err)
			return
		}
	default:
		conn.WriteMessage(websocket.TextMessage, []byte("unsupprt request type!"))
	}
}
