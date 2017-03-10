package models

import "time"

type Book struct {
	Id 			string		`bson:"_id"`
	Type 			int		`bson:"type"`
	Title 			string		`bson:"title"`
	SubTitle 		string		`bson:"subtitle"`
	AddTitle 		string		`bson:"addtitle"`
	Authors 		string		`bson:"authors"`
	Year			int		`bson:"year"`
	AuthoringAssistants 	string		`bson:"authoringAssistants"`
	Institution		string		`bson:"institution"`
	PublishingLocation	string		`bson:"publishingLocation"`
	Publisher		string		`bson:"publisher"`
	NumberOfBooks		int		`bson:"numberOfBooks"`
	Edition			int		`bson:"edition"`
	LineTitle		string		`bson:"lineTitle"`
	LineBookNo		int		`bson:"lineBooknr"`
	ISBN			int		`bson:"isbn"`
	AdditionalUrl		string		`bson:"addUrl"`
	LastChecked		time.Time	`bson:"lastChecked"`
	FirstPublished		time.Time	`bson:"firstPublished"`
	Language		string		`bson:"language"`
	Price			float64		`bson:"price"`
	DateAdded		time.Time	`bson:"dateAdded"`
	DateModified		time.Time	`bson:"dateModified"`
	Notes			string		`bson:"notes"`
	Abstract		string		`bson:"abstract"`
	ToC			string		`bson:"toc"`
	Rating			int		`bson:"rating"`
	Tags			[]string	`bson:"tags"`
	Categories		[]string	`bson:"categories"`
	References		[]string	`bson:"references"`
	Quotations		[]Quotation	`bson:"quotes"`
	Tasks			[]Task		`bson:"tasks"`
	Locations		[]Location	`bson:"locations"`
}
