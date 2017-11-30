package tweet

import (
	"strings"
	"gopkg.in/mgo.v2"
	"time"
	"net/http"
	"gopkg.in/mgo.v2/bson"
	"github.com/labstack/echo"
	"model"
	"github.com/dgrijalva/jwt-go"
	"math"
	"strconv"
)

type (
	// Handler : Handler for tweet related requests.
	Handler struct {
		db 				*mgo.Session
		dbName 			string
		key				string
		notifOperator	chan model.Notification
	}
)

// NewHandler : Create a Handler instance.
func NewHandler(dbURL string, key string, operator chan model.Notification) (h *Handler) {
	// Database connection
	session, err := mgo.Dial(dbURL)
	if err != nil {
		session.Close()
		panic(err)
	}

	return &Handler{session, strings.Split(dbURL, "/")[3], key, operator}
}

// FetchTweets : Handle requests asking for a list of tweets posted by a specific user.
//				 URL: "/api/v1/tweetlist/:user"
//				 Method: GET
//				 Return 200 OK on success.
//				 Return 404 Not Found if the user is not in the database.
func (h *Handler) FetchTweets (c echo.Context) (err error) {
	username := userNameFromToken(c)
	username = c.Param("username")
	/*
	t := new(model.Tweet)
	if err = c.Bind(t); err != nil {
		return
	}
	*/
	page, err := strconv.Atoi(c.QueryParam("page"))
	perpage, err := strconv.Atoi(c.QueryParam("perpage"))

	db := h.db.Clone()
	defer db.Close()

	// Check user is in the database
	
	user := model.User{}
	err = db.DB(h.dbName).C("users").Find(bson.M{"username": username}).One(&user)
	if err != nil {
		if err == mgo.ErrNotFound {
			return &echo.HTTPError{Code: http.StatusNotFound, Message: "User "+username+" does not exist."}
		}
		return
	}
	

	// Retrieve tweets from database
	tweets := []model.Tweet{}
	err = db.DB(h.dbName).C("tweets").Find(bson.M{"owner": username}).Sort("-timestamp").All(&tweets)

	var container struct {
		Page	string	`json:"page"`
		TotalPage	string	`json:"totalpage"`
		TotalTweets	string	`json:"totaltweets"`
		TweetList []model.Tweet `json:"tweetlist"`
	}

	if len(tweets) == 0 {
		container.Page = "0"
		container.TotalPage = "0"
		container.TotalTweets = "0"
		container.TweetList = []model.Tweet{}
		
		return c.JSON(http.StatusOK, container) 
	}

	totalTweets := len(tweets)
	totalPage := int(math.Ceil(float64(totalTweets)/float64(perpage)))

	var tweetList [] model.Tweet
	if page > totalPage{
		tweetList = []model.Tweet{}
	}else{
		if page == totalPage{
			tweetList = tweets[perpage*(page-1):]
		}else{
			tweetList = tweets[perpage*(page-1):perpage*page]
		}
	}

	// Change id to username
	/*
	for i := range tweetList {
		tweetList[i].Owner = user.Username
	}
	*/

	container.Page = strconv.Itoa(page)
	container.TotalPage = strconv.Itoa(totalPage)
	container.TotalTweets = strconv.Itoa(totalTweets)
	container.TweetList = tweetList
	

	return c.JSON(http.StatusOK, container)
	//return c.JSON(http.StatusOK, &tweets)
}

// NewTweet : Add one tweet for a specific user.
//			  URL: "/api/v1/newTweet"
//			  Method: POST
//			  Return 200 Created on success, along with the tweet data.
//			  Return 400 Bad Request if the content of the tweet is empty.
//			  Return 400 Bad Request if the image is larger than 10 MB.
func (h *Handler) NewTweet(c echo.Context) (err error) {
	userName := userNameFromToken(c)

	db := h.db.Clone()
	defer db.Close()

	tweet := &model.Tweet{ID: bson.NewObjectId(), Owner: string(userName), Numcomment: 0, Timestamp: time.Now(), Isretweet: false}
	if err = c.Bind(tweet); err != nil {
		return
	}

	// Retrieve user info from database
	user := model.User{}
	err = db.DB(h.dbName).C(model.UserCollection).Find(bson.M{"username": userName}).One(&user)
	if err != nil {
		if err == mgo.ErrNotFound {
			return &echo.HTTPError{Code: http.StatusNotFound, Message: "User does not exist."}
		}
		return
	}
	
	// Validation
	if tweet.Message == "" {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "Message cannot be empty."}
	}
	if len(tweet.Picture) > 10485760 {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "Image must be smaller than 10 MB."}
	}

	if user.Picture != ""{
		tweet.Picture = user.Picture
	}else{
		tweet.Picture = ""
	}
	
	// Save tweet in database
	err = db.DB(h.dbName).C("tweets").Insert(tweet)
	if err != nil {
		return
	}

	tweet.ID = ""
	var container struct {
		Owner	string	`json:"owner"`
		Message	string	`json:"message"`
		Picture string	`json:"picture"`
	}
	container.Owner = tweet.Owner
	container.Message = tweet.Message
	container.Picture = tweet.Picture

	h.notifOperator <- model.Notification{Timestamp: time.Now(), Type: model.NewTweetType, Detail: model.NewTweetNotif{Publisher: userName}}

	return c.JSON(http.StatusOK, container)
}

// DeleteTweet : Delete a specific tweet.
//				 URL: "/api/v1/deleteTweet/:tweet"
//				 Method: DELETE
//				 Return 204 No Content on success.
// 				 Return 400 Bad Request if tweet ID is malformed.
//				 Return 404 Not Found if the tweet is not in the database.
func (h *Handler) DeleteTweet(c echo.Context) (err error) {
	tweetID := c.Param("tweet")

	if !bson.IsObjectIdHex(tweetID) {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "Malformed tweet ID."}
	}

	db := h.db.Clone()
	defer db.Close()

	err = db.DB(h.dbName).C("tweets").RemoveId(bson.ObjectIdHex(tweetID))
	if err != nil {
		if err == mgo.ErrNotFound {
			return &echo.HTTPError{Code: http.StatusNotFound, Message: "Tweet does not exist."}
		}
		return
	}

	return c.NoContent(http.StatusNoContent)
}

// FetchTweetTimeLine : Handle requests asking for a list of tweets timeline of a specific user.
//				 	  	URL  : "/api/v1/tweettimeline/:username"
//				 		Method: GET
//				 		Return 200 OK on success.
//				 		Return 404 Not Found if the user is not in the database.
func (h *Handler) FetchTweetTimeLine (c echo.Context) (err error) {
	username := userNameFromToken(c)
	username = c.Param("username")
	
	page, err := strconv.Atoi(c.QueryParam("page"))
	perpage, err := strconv.Atoi(c.QueryParam("perpage"))

	db := h.db.Clone()
	defer db.Close()

	// Retrieve user info from database by username
	user := model.User{}
	err = db.DB(h.dbName).C("users").Find(bson.M{"username": username}).One(&user)
	if err != nil {
		if err == mgo.ErrNotFound {
			return &echo.HTTPError{Code: http.StatusNotFound, Message: "User "+username+" does not exist."}
		}
		return
	}
	// Get the user list of tweet timeline.
	following := user.Following
	timeLineUserList := append(following, user.Username)

	// Map between ID and username
	/*
	mapIDandUsername := make(map[string]string)
	for i := range timeLineUserList{
		tempUser := model.User{}
		err = db.DB(h.dbName).C("users").FindId(bson.ObjectIdHex(timeLineUserList[i])).One(&tempUser)
		mapIDandUsername[tempUser.ID.Hex()] = tempUser.Username
	}
	*/

	// Retrieve tweets from database
	tweets := []model.Tweet{}
	err = db.DB(h.dbName).C("tweets").Find(bson.M{"owner": bson.M{"$in": timeLineUserList}}).Sort("-timestamp").All(&tweets)
	
	var container struct {
		Page	string	`json:"page"`
		TotalPage	string	`json:"totalpage"`
		TotalTweets	string	`json:"totaltweets"`
		TweetList []model.Tweet `json:"tweetlist"`
	}

	if len(tweets) == 0 {
		container.Page = "0"
		container.TotalPage = "0"
		container.TotalTweets = "0"
		container.TweetList = []model.Tweet{}
		
		return c.JSON(http.StatusOK, container) 
	}

	totalTweets := len(tweets)
	totalPage := int(math.Ceil(float64(totalTweets)/float64(perpage)))

	var tweetList [] model.Tweet
	if page > totalPage{
		tweetList = []model.Tweet{}
	}else{
		if page == totalPage{
			tweetList = tweets[perpage*(page-1):]
		}else{
			tweetList = tweets[perpage*(page-1):perpage*page]
		}
	}

	// Change id to username
	/*
	for i := range tweetList {
		tweetList[i].Owner = mapIDandUsername[tweetList[i].Owner]
	}
	*/

	container.Page = strconv.Itoa(page)
	container.TotalPage = strconv.Itoa(totalPage)
	container.TotalTweets = strconv.Itoa(totalTweets)
	container.TweetList = tweetList

	return c.JSON(http.StatusOK, container)
	//return c.JSON(http.StatusOK, &tweets)
}

// ReTweet : Retweet a tweet.
//			  URL: "/api/v1/reTweet"
//			  Method: POST
//			  Return 200 Created on success, along with the tweet data.
//			  Return 400 Bad Request if the content of the tweet is empty.
//			  Return 400 Bad Request if the image is larger than 10 MB.
func (h *Handler) ReTweet(c echo.Context) (err error) {
	userName := userNameFromToken(c)

	db := h.db.Clone()
	defer db.Close()

	tweet := &model.Tweet{ID: bson.NewObjectId(), Owner: string(userName), Numcomment: 0, Timestamp: time.Now(), Isretweet: true}
	if err = c.Bind(tweet); err != nil {
		return
	}
	

	// Fetch the info of the tweet retweeted and check tweetID is valid
	retweet := new(model.Tweet)
	if bson.IsObjectIdHex(tweet.Idretweet){
		err = db.DB(h.dbName).C("tweets").FindId(bson.ObjectIdHex(tweet.Idretweet)).One(&retweet)
		if err != nil {
			if err == mgo.ErrNotFound {
				return &echo.HTTPError{Code: http.StatusNotFound, Message: "Invalid tweet ID"}
			}
			return &echo.HTTPError{Code: http.StatusNotFound, Message: "Invalid tweet ID"}
		}
	}else{
		return &echo.HTTPError{Code: http.StatusNotFound, Message: "Invalid tweet ID"}
	}

	// Validation
	if tweet.Message == "" {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "Message cannot be empty."}
	}
	if len(tweet.Picture) > 10485760 {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "Image must be smaller than 10 MB."}
	}
	
	tweet.Messageretweet = retweet.Message
	tweet.Ownerretweet = retweet.Owner

	// Save retweet in database
	err = db.DB(h.dbName).C("tweets").Insert(tweet)
	if err != nil {
		return
	}

	tweet.ID = ""
	var container struct {
		Owner	string	`json:"owner"`
		Message	string	`json:"message"`
		Ownerretweet	string	`json:"ownerretweet"`
		Messageretweet	string	`json:"messageretweet"`
		Picture string	`json:"picture"`
	}
	container.Owner = tweet.Owner
	container.Message = tweet.Message
	container.Ownerretweet = tweet.Ownerretweet
	container.Messageretweet = tweet.Messageretweet
	container.Picture = tweet.Picture

	//h.notifOperator <- model.Notification{Timestamp: time.Now(), Type: model.NewTweetType, Detail: model.NewTweetNotif{Publisher: userName}}

	return c.JSON(http.StatusOK, container)
}

func userNameFromToken(c echo.Context) string {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	return claims["username"].(string)
}

// Shutdown : Gracefully shutdown tweet handler.
func (h *Handler) Shutdown() {
	h.db.Close()
}
