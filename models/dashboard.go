package models

type AssigneeTaskCount struct {
	Assignee  string `json:"assignee"`
	TaskCount int64  `json:"taskCount"`
}

type TaskPercentage struct {
	Name  string `json:"name"`
	Value int64  `json:"value"`
	Count int64  `json:"count"`
}
