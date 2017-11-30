package user

import (
	"strings"
	"time"
	"net/http"
	"gopkg.in/mgo.v2"
	"github.com/dgrijalva/jwt-go"
	"gopkg.in/mgo.v2/bson"
	"github.com/labstack/echo"
	"model"
)

type (
	// Handler : Handler for user related requests.
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

// Signup : Create an user instance.
//			URL: "/api/v1/signup"
//			Method: POST
//			Return 201 Created on success, along with the user data.
//			Return 400 Bad Request if one of username, firstname, email, password is empty, or username, email already used.
func (h *Handler) Signup(c echo.Context) (err error) {
	// Bind
	u := &model.User{}
	if err = c.Bind(u); err != nil {
		return
	}

	db := h.db.Clone()
	defer db.Close()

	u.ID = bson.NewObjectId()
	u.Followers = []string{}
	u.Following = []string{}

	// Validate
	if u.Username == "" || u.FirstName == "" || u.Email == "" || u.Password == "" {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "Invalid username, firstname, email or password."}
	}

	varify := &model.User{}
	err = db.DB(h.dbName).C(model.UserCollection).Find(bson.M{"$or": []bson.M{bson.M{"username": u.Username}, bson.M{"email": u.Email}}}).One(varify)
	if err == nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "Username or email already used."}
	}

	// Save user
	err = db.DB(h.dbName).C(model.UserCollection).Insert(u);
	if err != nil {
		return
	}
	err = db.DB(h.dbName).C(model.NotificationCollection).Insert(model.Individual{ID: bson.NewObjectId(), Username: u.Username, Notifications: make([]model.Notification, 0)})
	if err != nil {
		return
	}

	// Don't send password
	u.Password = ""

	return c.JSON(http.StatusCreated, u)
}

// Login : Login to an account associated with the email address and the password.
//		   URL: "/api/v1/login"
//		   Method: POST
//		   Return 200 OK on success, along with the user data, which now contains a token.
//		   Return 401 Unauthorized if an account associated with the email address and password is not found.
func (h *Handler) Login(c echo.Context) (err error) {
	// Bind
	u := new(model.User)
	if err = c.Bind(u); err != nil {
		return
	}

	// Find user
	db := h.db.Clone()
	defer db.Close()

	err = db.DB(h.dbName).C(model.UserCollection).Find(bson.M{"email": u.Email, "password": u.Password}).One(u)
	if err != nil {
		if err == mgo.ErrNotFound {
			return &echo.HTTPError{Code: http.StatusUnauthorized, Message: "Invalid email or password."}
		}
		return
	}

	// JWT

	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = u.Username
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Generate encoded token and send it as response
	u.Token, err = token.SignedString([]byte(h.key))
	if err != nil {
		return err
	}

	// Don't send password
	u.Password = ""

	return c.JSON(http.StatusOK, u)
}

// FetchUserInfo : Return user info for a specific user, and whether it is followed by the current user.
//				   # Response does not include the full list of followers and following, only the counts.
//				   URL: "/api/v1/userInfo/:username"
//				   Method: GET
//				   Return 200 OK on success.
//				   Return 404 Not Found if the user is not in the database.
func (h *Handler) FetchUserInfo (c echo.Context) (err error) {
	selfUsername := usernameFromToken(c)

	db := h.db.Clone()
	defer db.Close()

	// Retrieve user info from database
	user := model.User{}
	err = db.DB(h.dbName).C(model.UserCollection).Find(bson.M{"username": c.Param("username")}).One(&user)
	if err != nil {
		if err == mgo.ErrNotFound {
			return &echo.HTTPError{Code: http.StatusNotFound, Message: "User does not exist."}
		}
		return
	}

	// Don't send password
	user.Password = ""

	var container struct {
		U				model.User	`json:"userinfo"`
		FollowerCount	int			`json:"followercount"`
		FollowingCount	int			`json:"followingcount"`
		Followed		bool		`json:"followed"`
	}
	container.U = user
	container.FollowerCount = len(user.Followers)
	container.FollowingCount = len(user.Following)

	container.Followed = false
	for i := range user.Followers {
		if user.Followers[i] == selfUsername {
			container.Followed = true
			break
		}
	}

	container.U.Followers = nil
	container.U.Following = nil

	return c.JSON(http.StatusOK, container)
}

// Follow : Add a specific username to the current user's following set, and add current user to that user's follower list.
//			URL: "/api/v1/follow/:username"
//			Method: POST
//			Return 200 OK on success, along with the user's following list.
//			Return 404 Not Found if the followee is not in the database.
func (h *Handler) Follow(c echo.Context) (err error) {
	username := usernameFromToken(c)
	followee := c.Param("username")

	db := h.db.Clone()
	defer db.Close()

	// Add self to followee's follower list
	err = db.DB(h.dbName).C(model.UserCollection).Update(bson.M{"username": followee}, bson.M{"$addToSet": bson.M{"followers": username}})
	if err != nil {
		if err == mgo.ErrNotFound {
			return &echo.HTTPError{Code: http.StatusNotFound, Message: "User does not exist."}
		}
		return
	}

	// Add followee to self following list
	err = db.DB(h.dbName).C(model.UserCollection).Update(bson.M{"username": username}, bson.M{"$addToSet": bson.M{"following": followee}})
	if err != nil {
		return
	}
	
	user := model.User{}
	err = db.DB(h.dbName).C(model.UserCollection).Find(bson.M{"username": username}).One(&user)

	h.notifOperator <- model.Notification{Timestamp: time.Now(), Type: model.FollowType, Detail: model.FollowNotif{Followee: followee, Follower: username}}

	return c.JSON(http.StatusOK, user.Following)
}

// Unfollow : Remove a specific username from the current user's following set, and remove current user from that user's follower list.
//			URL: "/api/v1/unfollow/:username"
//			Method: POST
//			Return 200 OK on success, along with the user's following list.
//			Return 404 Not Found if the followee is not in the database.
func (h *Handler) Unfollow(c echo.Context) (err error) {
	username := usernameFromToken(c)
	followee := c.Param("username")

	db := h.db.Clone()
	defer db.Close()

	// Remove self from followee's follower list
	err = db.DB(h.dbName).C(model.UserCollection).Update(bson.M{"username": followee}, bson.M{"$pull": bson.M{"followers": username}})
	if err != nil {
		if err == mgo.ErrNotFound {
			return &echo.HTTPError{Code: http.StatusNotFound, Message: "User does not exist."}
		}
		return
	}

	// Remove followee from self following list
	err = db.DB(h.dbName).C(model.UserCollection).Update(bson.M{"username": username}, bson.M{"$pull": bson.M{"following": followee}})
	if err != nil {
		return
	}
	
	user := model.User{}
	err = db.DB(h.dbName).C(model.UserCollection).Find(bson.M{"username": username}).One(&user)

	return c.JSON(http.StatusOK, user.Following)
}

// ShowFollower : Return the follower list for a specific user, along with some followers info.
//				  URL: "/api/v1/showFollower/:username"
//				  Method: GET
//				  Return 200 OK on success.
//				  Return 404 Not Found if the user is not in the database.
func (h *Handler) ShowFollower(c echo.Context) (err error) {
	username := c.Param("username")

	db := h.db.Clone()
	defer db.Close()

	// Retrieve user info from database
	user := model.User{}
	err = db.DB(h.dbName).C(model.UserCollection).Find(bson.M{"username": username}).One(&user)
	if err != nil {
		if err == mgo.ErrNotFound {
			return &echo.HTTPError{Code: http.StatusNotFound, Message: "User does not exist."}
		}
		return
	}

	type followerData struct {
		Username	string			`json:"username" bson:"username"`
		FirstName	string			`json:"firstname" bson:"firstname"`
		LastName	string			`json:"lastname,omitempty" bson:"lastname,omitempty"`
		Bio			string			`json:"bio,omitempty" bson:"bio,omitempty"`
	}
	container := []followerData{}

	for _, f := range user.Followers {
		follower := followerData{}
		err = db.DB(h.dbName).C(model.UserCollection).Find(bson.M{"username": f}).One(&follower)
		if err != nil {
			return
		}
		container = append(container, follower)
	}

	return c.JSON(http.StatusOK, &container)
}

// ShowFollowing : Return the following list for a specific user, along with some followings info.
//				   URL: "/api/v1/showFollowing/:username"
//				   Method: GET
//				   Return 200 OK on success.
//				   Return 404 Not Found if the user is not in the database.
func (h *Handler) ShowFollowing(c echo.Context) (err error) {
	username := c.Param("username")
	
	db := h.db.Clone()
	defer db.Close()

	// Retrieve user info from database
	user := model.User{}
	err = db.DB(h.dbName).C(model.UserCollection).Find(bson.M{"username": username}).One(&user)
	if err != nil {
		if err == mgo.ErrNotFound {
			return &echo.HTTPError{Code: http.StatusNotFound, Message: "User does not exist."}
		}
		return
	}

	type followingData struct {
		Username	string			`json:"username" bson:"username"`
		FirstName	string			`json:"firstname" bson:"firstname"`
		LastName	string			`json:"lastname,omitempty" bson:"lastname,omitempty"`
		Bio			string			`json:"bio,omitempty" bson:"bio,omitempty"`
	}
	container := []followingData{}

	for _, f := range user.Following {
		following := followingData{}
		err = db.DB(h.dbName).C(model.UserCollection).Find(bson.M{"username": f}).One(&following)
		if err != nil {
			return
		}
		container = append(container, following)
	}

	return c.JSON(http.StatusOK, &container)
}

// UpdateUserInfo : Update a user's info, only firstname, lastname, bio, and tag are allowed to be modified.
//					# Request body must contain all four fields, even if they are empty or not modified.
//					URL: "/api/v1/updateUserInfo"
//					Method: POST
//					Return 200 OK on success, along with the user data.
//					Return 400 Bad Request if one or more fields are empty.
//					Return 404 Not Found if the user is not in the database.
func (h *Handler) UpdateUserInfo(c echo.Context) (err error) {
	username := usernameFromToken(c)

	db := h.db.Clone()
	defer db.Close()

	type userUpdate struct {
		FirstName 	string	`json:"firstname"`
		LastName	string	`json:"lastname"`
		Bio			string	`json:"bio"`
		Tag			string	`json:"tag"`
	}

	update := &userUpdate{}
	if err = c.Bind(update); err != nil {
		return
	}

	if update.FirstName == "" {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "Firstname must not be empty."}
	}

	err = db.DB(h.dbName).C(model.UserCollection).Update(bson.M{"username": username}, bson.M{"$set": bson.M{"firstname": update.FirstName, "lastname": update.LastName, "bio": update.Bio, "tag": update.Tag}})
	if err != nil {
		if err == mgo.ErrNotFound {
			return &echo.HTTPError{Code: http.StatusNotFound, Message: "User does not exist."}
		}
		return
	}

	user := &model.User{}
	err = db.DB(h.dbName).C(model.UserCollection).Find(bson.M{"username": username}).One(user)
	if err != nil {
		return
	}

	// Don't send password
	user.Password = ""
	
	//return c.JSON(http.StatusOK, user)
	return c.JSON(http.StatusOK, user)
}

// UpdateProfilePicture : Update a user's profile picture, or remove profile picture by sending an empty string.
//						  # Image is to be encoded into base 64 string, and cannot be larger than 10 MB.
//						  URL: "/api/v1/updateProfilePic"
//						  Method: POST
//						  Return 201 Created on success, along with the user data.
//						  Return 400 Bad Request if the image is larger than 10 MB.
//						  Return 404 Not Found if the user is not in the database.
func (h *Handler) UpdateProfilePicture(c echo.Context) (err error) {
	username := usernameFromToken(c)
	
	db := h.db.Clone()
	defer db.Close()

	type base64Image struct {
		Base64String	string	`json:"picture"`
	}
	
	image := base64Image{}
	if err = c.Bind(&image); err != nil {
		return
	}

	// Do not accept image larger than 10 MB
	if len(image.Base64String) > 10485760 {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "Image must be smaller than 10 MB."}
	}
	
	err = db.DB(h.dbName).C(model.UserCollection).Update(bson.M{"username": username}, bson.M{"$set": bson.M{"picture": &image.Base64String}})
	if err != nil {
		if err == mgo.ErrNotFound {
			return &echo.HTTPError{Code: http.StatusNotFound, Message: "User does not exist."}
		}
		return
	}

	db.DB(h.dbName).C(model.TweetCollection).UpdateAll(bson.M{"owner": username}, bson.M{"$set": bson.M{"picture": &image.Base64String}})

	user := &model.User{}
	err = db.DB(h.dbName).C(model.UserCollection).Find(bson.M{"username": username}).One(user)
	if err != nil {
		return
	}

	// Don't send password
	user.Password = ""
	
	//return c.JSON(http.StatusOK, user)
	return c.JSON(http.StatusCreated, user)
}

func usernameFromToken(c echo.Context) string {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	return claims["username"].(string)
}

// Shutdown : Gracefully shutdown user handler.
func (h *Handler) Shutdown() {
	h.db.Close()
}