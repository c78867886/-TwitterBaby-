package model

import "gopkg.in/mgo.v2/bson"
import "time"

// TweetCollection : Tweet collection name in database.
const TweetCollection = "tweets"

// Tweet : Data structure that holds a single tweet.
type Tweet struct {
	ID			bson.ObjectId	`json:"id" bson:"_id"`
	Owner 		string			`json:"owner" bson:"owner"`
	//From 		string			`json:"from,omitempty" bson:"from,omitempty"`
	Message 	string			`json:"message" bson:"message"`
	Numcomment  int				`json:"numcomment" bson:"numcomment"`
	Picture		string			`json:"picture,omitempty" bson:"picture,omitempty"`
	Isretweet   bool			`json:"isretweet" bson:"isretweet"`
	Idretweet	string			`json:"idretweet" bson:"idretweet"`
	Ownerretweet 		string			`json:"ownerretweet" bson:"ownerretweet"`
	Messageretweet 	string			`json:"messageretweet" bson:"messageretweet"`
	Timestamp	time.Time		`json:"timestamp,omitempty" bson:"timestamp"`
}