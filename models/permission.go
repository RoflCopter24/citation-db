package models

type Permission struct {
	Group	bool	`bson:"group"`
	GroupId int	`bson:"groupId"`
	UserId	string	`bson:"userId"`
	Read	bool	`bson:"read"`
	Write	bool	`bson:"write"`
}
