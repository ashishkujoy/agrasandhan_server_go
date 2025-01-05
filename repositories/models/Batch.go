package models

import "time"

type MentorPermission struct {
	AllowProvideObservations bool `json:"allowProvideObservations" bson:"allowProvideObservations"`
	AllowReleaseIntern       bool `json:"allowReleaseIntern" bson:"allowReleaseIntern"`
	AllowProvideFeedback     bool `json:"allowProvideFeedback" bson:"allowProvideFeedback"`
	AllowDeliverFeedback     bool `json:"allowDeliverFeedback" bson:"allowDeliverFeedback"`
}

type Mentor struct {
	ID          string           `json:"id" bson:"id"`
	Permissions MentorPermission `json:"permissions" bson:"permissions"`
}

type Batch struct {
	ID        int       `json:"id" bson:"id"`
	StartDate time.Time `json:"startDate" bson:"startDate"`
	EndDate   time.Time `json:"endDate" bson:"endDate"`
	Interns   []string  `json:"interns" bson:"interns"`
	Mentors   []Mentor  `json:"mentors" bson:"mentors"`
	Name      string    `json:"name" bson:"name"`
}
