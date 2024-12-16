package models

type LeaderboardMemberChange struct {
	QuizID   string `json:"quizID"`
	UserID   string `json:"userID"`
	NewRank  int    `json:"newRank"`
	NewScore int    `json:"newScore"`
}

