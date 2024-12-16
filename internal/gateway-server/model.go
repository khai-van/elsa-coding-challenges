package gatewayserver

type MessageType int

const (
	LeaderboardChangeType MessageType = iota + 1
	SubmitAnswer
	SubmitAnswerResp
)

type WebsocketMessage[T any] struct {
	Type    MessageType `json:"type"`
	Message T           `json:"msg"`
}

type SubmitAnwserMessage struct {
	QuestionID string `json:"questionID"`
	Anwser     string `json:"anwser"`
}
