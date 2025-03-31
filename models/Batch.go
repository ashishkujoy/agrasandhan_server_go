package models

import "time"

type InternWithFeedbacks struct {
	ColorCode           int       `json:"colorCode"`
	ProfilePicId        int       `json:"profilePicId"`
	BatchId             int       `json:"batchId"`
	IsReleased          bool      `json:"isReleased"`
	IsOnNotice          bool      `json:"isOnNotice"`
	Id                  string    `json:"id"`
	Name                string    `json:"name"`
	Gender              string    `json:"gender"`
	GithubUsername      string    `json:"githubUsername"`
	DiplomaBranch       string    `json:"diplomaBranch"`
	TwEmailId           string    `json:"twEmailId"`
	LastFeedbackDate    time.Time `json:"lastFeedbackDate"`
	LastObservationDate time.Time `json:"lastObservationDate"`
}

type Batch struct {
	Id        int                   `json:"id"`
	Name      string                `json:"name"`
	StartDate time.Time             `json:"startDate"`
	EndDate   time.Time             `json:"endDate"`
	Interns   []InternWithFeedbacks `json:"interns"`
}
