# **Basic Server Operation**
	h for help
	q to shutdown server
	d to drop database
	i to reconstruct database to default (w/ some initial collections)

# **package model**
    import "model"

### TYPES

```go
type Tweet struct {
    ID        bson.ObjectId `json:"id" bson:"_id"`
    Owner     string        `json:"owner" bson:"owner"`
    From      string        `json:"from,omitempty" bson:"from,omitempty"`
    Message   string        `json:"message" bson:"message"`
    Timestamp time.Time     `json:"timestamp,omitempty" bson:"timestamp"`
}
```
> Tweet : Data structure that holds a single tweet.

```go
type User struct {
    ID        bson.ObjectId `json:"id" bson:"_id"`
    Username  string        `json:"username" bson:"username"`                       // unique
    FirstName string        `json:"firstname" bson:"firstname"`
    LastName  string        `json:"lastname,omitempty" bson:"lastname,omitempty"`
    Password  string        `json:"password,omitempty" bson:"password"`
    Email     string        `json:"email" bson:"email"`                             // unique
    Followers []string      `json:"followers,omitempty" bson:"followers,omitempty"`
    Following []string      `json:"following,omitempty" bson:"following,omitempty"`
    Bio       string        `json:"bio,omitempty" bson:"bio,omitempty"`
    Token     string        `json:"token,omitempty" bson:"-"`
}
```
> User : Data structure that holds a single user.

# **package handler**
    import "handler"

### TYPES

```go
type Handler struct {
    DB *mgo.Session
}
```
> Handler : Handler for dealing with requests.

### Functions

```go
func (h *Handler) DeleteTweet(c echo.Context) (err error)
```
> DeleteTweet : Delete a specific tweet.  
> URL: "/api/v1/deleteTweet/:tweet"  
> Method: DELETE  
> Return 204 No Content on success.  
> Return 400 Bad Request if∂ tweet ID is malformed.  
> Return 404 Not Found if the tweet is not in the database.  
> Sample call:  
> curl -X DELETE http://127.0.0.1:1323/api/v1/deleteTweet/59cbe3c1a54d75a4f514beb0 -H 'Content-Type: application/json' -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MDY1OTA3MTUsImlkIjoiNTljYjY5Y2NhNTRkNzU3YThlMzljOTc0In0.8A8hQQtbZeYBT3LDmOu_-OnrmRsfSby-KZw0eAMJ06s'

```go
func (h *Handler) FetchTweets(c echo.Context) (err error)
```
> FetchTweets : FetchTweets : Handle requests asking for a list of tweets posted by a specific user.  
> URL: "/api/v1/tweetlist/:user"  
> Method: GET  
> Return 200 OK on success.  
> Return 404 Not Found if the user is not in the database.  
> Sample call:  
> curl -X GET http://127.0.0.1:1323/api/v1/tweetlist/JasonHo -H 'Content-Type: application/json'

```go
func (h *Handler) FetchUserInfo(c echo.Context) (err error)
```
> FetchUserInfo : Return user info for a specific user, and whether it is followed by the current user.  
> \# Response does not include the full list of followers and following, only the counts.  
> URL: "/api/v1/userInfo/:user", "/api/v1/userInfo?username=:username"  
> Method: GET  
> Return 200 OK on success.  
> Return 404 Not Found if the user is not in the database.  
> Sample call:  
> curl -X GET http://127.0.0.1:1323/api/v1/userInfo/59cb69cca54d757a8e39c974 -H 'Content-Type: application/json' -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MDY1OTA3MTUsImlkIjoiNTljYjY5Y2NhNTRkNzU3YThlMzljOTc0In0.8A8hQQtbZeYBT3LDmOu_-OnrmRsfSby-KZw0eAMJ06s'

```go
func (h *Handler) Follow(c echo.Context) (err error)
```
> Follow : Add a specific user ID to the current user's following set, and add current user to that user's follower list.  
> URL: "/api/v1/follow/:id"  
> Method: POST  
> Return 200 OK on success, along with the user's following list.  
> Return 404 Not Found if the user is not in the database.  
> Sample call:  
> curl -X POST http://127.0.0.1:1323/api/v1/follow/59cb6e6ca54d757d12f0b0be -H 'Content-Type: application/json' -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MDY1OTA3MTUsImlkIjoiNTljYjY5Y2NhNTRkNzU3YThlMzljOTc0In0.8A8hQQtbZeYBT3LDmOu_-OnrmRsfSby-KZw0eAMJ06s'

```go
func (h *Handler) Login(c echo.Context) (err error)
```
> Login : Login to an account associated with the email address and the password.  
> URL: "/api/v1/login"  
> Method: POST  
> Return 200 OK on success, along with the user data, which now contains a token.  
> Return 401 Unauthorized if an account associated with the email address and password is not found.  
> Sample call:  
> curl -X POST http://127.0.0.1:1323/api/v1/login -H 'Content-Type: application/json' -d '{"email":"hojason117@gmail.com","password":"test1"}'

```go
func (h *Handler) NewTweet(c echo.Context) (err error)
```
> NewTweet : Add one tweet for a specific user.  
> URL: "/api/v1/newTweet"  
> Method: POST  
> Return 201 Created on success, along with the tweet data.  
> Return 404 Not Found if the user is not in the database.  
> Return 400 Bad Request if the content of the tweet is empty.  
> Sample call:  
> curl -X POST http://127.0.0.1:1323/api/v1/newTweet -H 'Content-Type: application/json' -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MDcwMTMxODYsImlkIjoiNTljYzhjOTVhNTRkNzVkNWM4YmUzYWE3In0.aiOsy8FbLcKV_EffHwOvfOZ-2rwg-kod5ov-v6xn98M' -d '{"message":"First tweet!"}'

```go
func (h *Handler) ShowFollower(c echo.Context) (err error)
```
> ShowFollower : Return the follower list for a specific user.  
> URL: "/api/v1/showFollower/:id"  
> Method: GET  
> Return 200 OK on success.  
> Return 404 Not Found if the user is not in the database.  
> Sample call:  
> curl -X GET http://127.0.0.1:1323/api/v1/showFollower/59cb69cca54d757a8e39c974 -H 'Content-Type: application/json' -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MDY1OTA3MTUsImlkIjoiNTljYjY5Y2NhNTRkNzU3YThlMzljOTc0In0.8A8hQQtbZeYBT3LDmOu_-OnrmRsfSby-KZw0eAMJ06s'

```go
func (h *Handler) ShowFollowing(c echo.Context) (err error)
```
> ShowFollowing : Return the following list for a specific user.  
> URL: "/api/v1/showFollowing/:id"  
> Method: GET  
> Return 200 OK on success.  
> Return 404 Not Found if the user is not in the database.  
> Sample call:  
> curl -X GET http://127.0.0.1:1323/api/v1/showFollowing/59cc8c95a54d75d5c8be3aa7 -H 'Content-Type: application/json' -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MDcwMDQ3ODYsImlkIjoiNTljYzhjOTVhNTRkNzVkNWM4YmUzYWE3In0.4m-C5cQ6wOgK_aPe8XiYXS5-CR3bZJgG-am1rGZUIgU'

```go
func (h *Handler) Signup(c echo.Context) (err error)
```
> Signup : Create an user instance.   
> URL: "/api/v1/signup"  
> Method: POST  
> Return 201 Created on success, along with the user data.  
> Return 400 Bad Request if one of username, firstname, email, password is empty, or username, email already used.  
> Sample call:  
> curl -X POST http://127.0.0.1:1323/api/v1/signup -H 'Content-Type: application/json' -d '{"firstname":"Evelyn","password":"pass","email":"eve@gmail.com"}'

# **Tips for Go**
* Don’t name a function starting with a capital letter if you don’t need to export it.
* Use godoc command to automatically generate documentation.

# **Tips for backend**
* [echo Twitter cookbook](https://github.com/labstack/echox/tree/master/cookbook/twitter)
* [Post v.s. Put](https://stackoverflow.com/questions/630453/put-vs-post-in-rest)
* [REST API Documentation](https://bocoup.com/blog/documenting-your-api)
* [QueryParam vs PathParam](https://stackoverflow.com/questions/11552248/when-to-use-queryparam-vs-pathparam)
* [path params vs. query params in a RESTful API](https://stackoverflow.com/questions/30967822/when-do-i-use-path-params-vs-query-params-in-a-restful-api)
