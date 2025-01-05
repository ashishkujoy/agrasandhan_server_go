package models

import "time"

type MentorPermission struct {
	AllowProvideObservations bool
	AllowReleaseIntern       bool
	AllowProvideFeedback     bool
	AllowDeliverFeedback     bool
}

type Mentor struct {
	ID          int              `json:"id" bson:"id"`
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
