package models

import "time"

type Quotation struct {
	Id		string		`bson:"_id"`
	Type		int		`bson:"type"`
	CoreStatement	string		`bson:"coreStatement"`
	Files		[]string	`bson:"files"`
	Legend		string		`bson:"legend"`
	Description	string		`bson:"description"`
	Text		string		`bson:"text"`
	PageNrStart	int		`bson:"pagenrStart"`
	PageNrStop	int		`bson:"pagenrStop"`
	Tags		[]string	`bson:"tags"`
	Categories	[]string	`bson:"categories"`
	CreatorId	string		`bson:"creatorId"`
	Permissions	[]Permission	`bson:"permissions"`
	CreationDate	time.Time	`bson:"creationDate"`
	LastModified	time.Time	`bson:"lastModified"`
	History		[]Quotation	`bson:"history"`
}
