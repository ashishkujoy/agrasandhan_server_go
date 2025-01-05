package requests

import "ashishkujoy/agrasandhan/repositories/models"

type AssignMentorRequest struct {
	Id          string                  `json:"id"`
	Permissions models.MentorPermission `json:"permissions"`
}
