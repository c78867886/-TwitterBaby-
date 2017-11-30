package user

import (
	"model"
	"notification"
	"fmt"
	"strings"
	"time"
	"net/http"
	"net/http/httptest"
	"testing"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"github.com/labstack/echo"
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
)

type (
	restHTTPTestCase struct {
		positive	bool
		input		restHTTPInput
		expected	restHTTPExpected
	}

	restHTTPInput struct {
		user	string
		target	string
	}

	restHTTPExpected struct {
		code 		int
		message 	string
	}
)

const key  = "UYrtPaa0Pky7QZyVrkIfnouatG7LjTKystf0FGdOuDiXCZyCSuVz8YdK7OBeSrC"

func TestSignup(t *testing.T) {
	// Setup
	e := echo.New()
	h := NewHandler("mongodb://SEavenger:SEavenger@ds121225.mlab.com:21225/se_avengers_test", key, make(chan model.Notification))

	// Test cases
	testCases := []restHTTPTestCase {
		// test success
		restHTTPTestCase{
			true,
			restHTTPInput{"", `{"username":"testSignup","firstname":"test","lastname":"signup","password":"test","email":"testSignup@gmail.com","bio":"testtest"}`},
			restHTTPExpected{http.StatusCreated, `{"username":"testSignup","firstname":"test","lastname":"signup","email":"testSignup@gmail.com","bio":"testtest"}`},
		},
		restHTTPTestCase{
			true,
			restHTTPInput{"", `{"username":"testSignup","firstname":"test","lastname":"signup","password":"test","email":"testSignup@gmail.com"}`,},
			restHTTPExpected{http.StatusCreated, `{"username":"testSignup","firstname":"test","lastname":"signup","email":"testSignup@gmail.com"}`},
		},
		restHTTPTestCase{
			true,
			restHTTPInput{"", `{"username":"testSignup","firstname":"test","password":"test","email":"testSignup@gmail.com","bio":"testtest"}`},
			restHTTPExpected{http.StatusCreated, `{"username":"testSignup","firstname":"test","email":"testSignup@gmail.com","bio":"testtest"}`},
		},
		restHTTPTestCase{
			true,
			restHTTPInput{"", `{"username":"testSignup","firstname":"test","password":"test","email":"testSignup@gmail.com","bio":"testtest","tag":"Einstein"}`},
			restHTTPExpected{http.StatusCreated, `{"username":"testSignup","firstname":"test","email":"testSignup@gmail.com","bio":"testtest","tag":"Einstein"}`},
		},
		// test empty
		restHTTPTestCase{
			false,
			restHTTPInput{"", `{"firstname":"test","lastname":"signup","password":"test","email":"testSignup@gmail.com"}`},
			restHTTPExpected{http.StatusBadRequest, "Invalid username, firstname, email or password."},
		},
		restHTTPTestCase{
			false,
			restHTTPInput{"", `{"username":"testSignup","lastname":"signup","password":"test","email":"testSignup@gmail.com"}`,},
			restHTTPExpected{http.StatusBadRequest, "Invalid username, firstname, email or password."},
		},
		restHTTPTestCase{
			false,
			restHTTPInput{"", `{"username":"testSignup","firstname":"test","lastname":"signup","password":"test"}`},
			restHTTPExpected{http.StatusBadRequest, "Invalid username, firstname, email or password."},
		},
		restHTTPTestCase{
			false,
			restHTTPInput{"", `{"username":"testSignup","firstname":"test","lastname":"signup","email":"testSignup@gmail.com","bio":"testtest"}`},
			restHTTPExpected{http.StatusBadRequest, "Invalid username, firstname, email or password."},
		},
		// test duplicate
		restHTTPTestCase{
			false,
			restHTTPInput{"", `{"username":"testSignup_dup","firstname":"test","lastname":"signup_dup","password":"test","email":"testSignup_dup@gmail.com","bio":"testtest"}`},
			restHTTPExpected{http.StatusBadRequest, "Username or email already used."},
		},
	}

	// Run
	for _, tc := range testCases {
		// Delete DB entry
		if tc.positive {
			err := h.db.DB(h.dbName).C(model.UserCollection).Remove(bson.M{"username": "testSignup"})
			if err != nil {
				if err != mgo.ErrNotFound {
					t.Fatal(err)
				}
			}
		}

		// Setup
		req := httptest.NewRequest(echo.POST, "/", strings.NewReader(tc.input.target))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		// Set the registered path for the handler.
		c.SetPath("/api/v1/signup")

		// Assertion
		if tc.positive {
			if assert.NoError(t, h.Signup(c), tc.input.target) {
				assert.Equal(t, tc.expected.code, rec.Code, tc.input.target)
				assert.Equal(t, tc.expected.message, rec.Body.String(), tc.input.target)
			}
		} else {
			if err := h.Signup(c); assert.Error(t, err, tc.input.target) {
				assert.Equal(t, tc.expected.code, err.(*echo.HTTPError).Code, tc.input.target)
				assert.Equal(t, tc.expected.message, err.(*echo.HTTPError).Message, tc.input.target)
			}
		}
	}
	h.Shutdown()
}

func TestLogin(t *testing.T) {
	e := echo.New()
	h := NewHandler("mongodb://SEavenger:SEavenger@ds121225.mlab.com:21225/se_avengers_test", key, make(chan model.Notification))

	testCases := []restHTTPTestCase {
		// test success
		restHTTPTestCase{
			true,
			restHTTPInput{"", `{"email":"testLogin@gmail.com","password":"test"}`},
			restHTTPExpected{http.StatusOK, `{"username":"testLogin","firstname":"testLogin","email":"testLogin@gmail.com"`},
		},
		// tes fail
		restHTTPTestCase{
			false,
			restHTTPInput{"", `{"email":"testLogin@gmail.com","password":""}`},
			restHTTPExpected{http.StatusUnauthorized, "Invalid email or password."},
		},
		restHTTPTestCase{
			false,
			restHTTPInput{"", `{"email":"","password":"test"}`},
			restHTTPExpected{http.StatusUnauthorized, "Invalid email or password."},
		},
		restHTTPTestCase{
			false,
			restHTTPInput{"", `{"email":"testLogin","password":"test"}`},
			restHTTPExpected{http.StatusUnauthorized, "Invalid email or password."},
		},
	}

	for _, tc := range testCases {
		req := httptest.NewRequest(echo.POST, "/", strings.NewReader(tc.input.target))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/login")

		if tc.positive {
			if assert.NoError(t, h.Login(c), tc.input.target) {
				assert.Equal(t, tc.expected.code, rec.Code, tc.input.target)
				assert.Equal(t, tc.expected.message, strings.Split(rec.Body.String(), `,"token`)[0], tc.input.target)
			}
		} else {
			if err := h.Login(c); assert.Error(t, err, tc.input.target) {
				assert.Equal(t, tc.expected.code, err.(*echo.HTTPError).Code, tc.input.target)
				assert.Equal(t, tc.expected.message, err.(*echo.HTTPError).Message, tc.input.target)
			}
		}
	}
	h.Shutdown()
}

func TestFetchUserInfo(t *testing.T) {
	e := echo.New()
	h := NewHandler("mongodb://SEavenger:SEavenger@ds121225.mlab.com:21225/se_avengers_test", key, make(chan model.Notification))

	testCases := []restHTTPTestCase {
		// test success
		restHTTPTestCase{
			true,
			restHTTPInput{"testUserInfo_1", "testUserInfo_2"},
			restHTTPExpected{http.StatusOK, `{"userinfo":{"username":"testUserInfo_2","firstname":"testUserInfo_2","email":"testUserInfo_2@gmail.com"},"followercount":1,"followingcount":0,"followed":true}`},
		},
		restHTTPTestCase{
			true,
			restHTTPInput{"testUserInfo_2", "testUserInfo_1"},
			restHTTPExpected{http.StatusOK, `{"userinfo":{"username":"testUserInfo_1","firstname":"testUserInfo_1","email":"testUserInfo_1@gmail.com"},"followercount":0,"followingcount":1,"followed":false}`},
		},
		// tes fail
		restHTTPTestCase{
			false,
			restHTTPInput{"testUserInfo_1", "userNotExist"},
			restHTTPExpected{http.StatusNotFound, "User does not exist."},
		},
	}

	for _, tc := range testCases {
		req := httptest.NewRequest(echo.GET, "/", nil)
		req.Header.Set(echo.HeaderAuthorization, "Bearer " + getToken(tc.input.user))
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/userInfo")
		c.SetParamNames("username")
		c.SetParamValues(tc.input.target)
		processJWTToken(c)

		if tc.positive {
			if assert.NoError(t, h.FetchUserInfo(c), tc.input.user + " fetching " + tc.input.target) {
				assert.Equal(t, tc.expected.code, rec.Code, tc.input.user + " fetching " + tc.input.target)
				assert.Equal(t, tc.expected.message, rec.Body.String(), tc.input.user + " fetching " + tc.input.target)
			}
		} else {
			if err := h.FetchUserInfo(c); assert.Error(t, err, tc.input.user + " fetching " + tc.input.target) {
				assert.Equal(t, tc.expected.code, err.(*echo.HTTPError).Code, tc.input.user + " fetching " + tc.input.target)
				assert.Equal(t, tc.expected.message, err.(*echo.HTTPError).Message, tc.input.user + " fetching " + tc.input.target)
			}
		}
	}
	h.Shutdown()
}

func TestFollow(t *testing.T) {
	e := echo.New()
	nh := notification.NewHandler("mongodb://SEavenger:SEavenger@ds121225.mlab.com:21225/se_avengers_test")
	h := NewHandler("mongodb://SEavenger:SEavenger@ds121225.mlab.com:21225/se_avengers_test", key, nh.Manager.Operator)

	testCases := []restHTTPTestCase {
		// test success
		restHTTPTestCase{
			true,
			restHTTPInput{"testFollow", "testFollow_2"},
			restHTTPExpected{http.StatusOK, `["testFollow_1","testFollow_2"]`},
		},
		// test fail
		restHTTPTestCase{
			false,
			restHTTPInput{"testFollow", "testFollow_not_exist"},
			restHTTPExpected{http.StatusNotFound, "User does not exist."},
		},
	}

	for _, tc := range testCases {
		req := httptest.NewRequest(echo.POST, "/", nil)
		req.Header.Set(echo.HeaderAuthorization, "Bearer " + getToken(tc.input.user))
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/follow")
		c.SetParamNames("username")
		c.SetParamValues(tc.input.target)
		processJWTToken(c)

		if tc.positive {
			if assert.NoError(t, h.Follow(c), tc.input.user + " follows " + tc.input.target) {
				assert.Equal(t, tc.expected.code, rec.Code, tc.input.user + " follows " + tc.input.target)
				assert.Equal(t, tc.expected.message, rec.Body.String(), tc.input.user + " follows " + tc.input.target)

				user := model.User{}
				h.db.DB(h.dbName).C(model.UserCollection).Find(bson.M{"username": tc.input.target}).One(&user)
				assert.Equal(t, []string{tc.input.user}, user.Followers, tc.input.user + " follows " + tc.input.target)

				h.db.DB(h.dbName).C(model.UserCollection).Update(bson.M{"username": tc.input.target}, bson.M{"$pull": bson.M{"followers": tc.input.user}})
				h.db.DB(h.dbName).C(model.UserCollection).Update(bson.M{"username": tc.input.user}, bson.M{"$pull": bson.M{"following": tc.input.target}})
			}
		} else {
			if err := h.Follow(c); assert.Error(t, err, tc.input.user + " follows " + tc.input.target) {
				assert.Equal(t, tc.expected.code, err.(*echo.HTTPError).Code, tc.input.user + " follows " + tc.input.target)
				assert.Equal(t, tc.expected.message, err.(*echo.HTTPError).Message, tc.input.user + " follows " + tc.input.target)
			}
		}
	}
	h.Shutdown()
}

func TestUnfollow(t *testing.T) {
	e := echo.New()
	h := NewHandler("mongodb://SEavenger:SEavenger@ds121225.mlab.com:21225/se_avengers_test", key, make(chan model.Notification))

	testCases := []restHTTPTestCase {
		// test success
		restHTTPTestCase{
			true,
			restHTTPInput{"testUnfollow", "testUnfollow_2"},
			restHTTPExpected{http.StatusOK, `["testUnfollow_1"]`},
		},
		// test fail
		restHTTPTestCase{
			false,
			restHTTPInput{"testUnfollow", "testUnfollow_not_exist"},
			restHTTPExpected{http.StatusNotFound, "User does not exist."},
		},
	}

	for _, tc := range testCases {
		req := httptest.NewRequest(echo.POST, "/", nil)
		req.Header.Set(echo.HeaderAuthorization, "Bearer " + getToken(tc.input.user))
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/unfollow")
		c.SetParamNames("username")
		c.SetParamValues(tc.input.target)
		processJWTToken(c)

		if tc.positive {
			if assert.NoError(t, h.Unfollow(c), tc.input.user + " unfollows " + tc.input.target) {
				assert.Equal(t, tc.expected.code, rec.Code, tc.input.user + " unfollows " + tc.input.target)
				assert.Equal(t, tc.expected.message, rec.Body.String(), tc.input.user + " unfollows " + tc.input.target)

				user := model.User{}
				h.db.DB(h.dbName).C(model.UserCollection).Find(bson.M{"username": tc.input.target}).One(&user)
				assert.Equal(t, []string{}, user.Followers, tc.input.user + " unfollows " + tc.input.target)

				h.db.DB(h.dbName).C(model.UserCollection).Update(bson.M{"username": tc.input.target}, bson.M{"$addToSet": bson.M{"followers": tc.input.user}})
				h.db.DB(h.dbName).C(model.UserCollection).Update(bson.M{"username": tc.input.user}, bson.M{"$addToSet": bson.M{"following": tc.input.target}})
			}
		} else {
			if err := h.Unfollow(c); assert.Error(t, err, tc.input.user + " unfollows " + tc.input.target) {
				assert.Equal(t, tc.expected.code, err.(*echo.HTTPError).Code, tc.input.user + " unfollows " + tc.input.target)
				assert.Equal(t, tc.expected.message, err.(*echo.HTTPError).Message, tc.input.user + " unfollows " + tc.input.target)
			}
		}
	}
	h.Shutdown()
}

func TestShowFollower(t *testing.T) {
	e := echo.New()
	h := NewHandler("mongodb://SEavenger:SEavenger@ds121225.mlab.com:21225/se_avengers_test", key, make(chan model.Notification))

	testCases := []restHTTPTestCase {
		// test success
		restHTTPTestCase{
			true,
			restHTTPInput{"", "testShowFollower_1"},
			restHTTPExpected{http.StatusOK, "[]"},
		},
		restHTTPTestCase{
			true,
			restHTTPInput{"", "testShowFollower_2"},
			restHTTPExpected{http.StatusOK, `[{"username":"testShowFollower_1","firstname":"testShowFollower_1"},{"username":"testShowFollower_3","firstname":"testShowFollower_3"}]`},
		},
		restHTTPTestCase{
			true,
			restHTTPInput{"", "testShowFollower_3"},
			restHTTPExpected{http.StatusOK, `[{"username":"testShowFollower_2","firstname":"test","lastname":"ShowFollower_2","bio":"testtest"}]`},
		},
		// test fail
		restHTTPTestCase{
			false,
			restHTTPInput{"", "testShowFollower_not_exist"},
			restHTTPExpected{http.StatusNotFound, "User does not exist."},
		},
	}

	for _, tc := range testCases {
		req := httptest.NewRequest(echo.GET, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/showFollower")
		c.SetParamNames("username")
		c.SetParamValues(tc.input.target)

		if tc.positive {
			if assert.NoError(t, h.ShowFollower(c), tc.input.target) {
				assert.Equal(t, tc.expected.code, rec.Code, tc.input.target)
				assert.Equal(t, tc.expected.message, rec.Body.String(), tc.input.target)
			}
		} else {
			if err := h.ShowFollower(c); assert.Error(t, err, tc.input.target) {
				assert.Equal(t, tc.expected.code, err.(*echo.HTTPError).Code, tc.input.target)
				assert.Equal(t, tc.expected.message, err.(*echo.HTTPError).Message, tc.input.target)
			}
		}
	}
	h.Shutdown()
}

func TestShowFollowing(t *testing.T) {
	e := echo.New()
	h := NewHandler("mongodb://SEavenger:SEavenger@ds121225.mlab.com:21225/se_avengers_test", key, make(chan model.Notification))

	testCases := []restHTTPTestCase {
		// test success
		restHTTPTestCase{
			true,
			restHTTPInput{"", "testShowFollowing_1"},
			restHTTPExpected{http.StatusOK, "[]"},
		},
		restHTTPTestCase{
			true,
			restHTTPInput{"", "testShowFollowing_2"},
			restHTTPExpected{http.StatusOK, `[{"username":"testShowFollowing_1","firstname":"testShowFollowing_1"},{"username":"testShowFollowing_3","firstname":"testShowFollowing_3"}]`},
		},
		restHTTPTestCase{
			true,
			restHTTPInput{"", "testShowFollowing_3"},
			restHTTPExpected{http.StatusOK, `[{"username":"testShowFollowing_2","firstname":"test","lastname":"ShowFollowing_2","bio":"testtest"}]`},
		},
		// test fail
		restHTTPTestCase{
			false,
			restHTTPInput{"", "testShowFollowing_not_exist"},
			restHTTPExpected{http.StatusNotFound, "User does not exist."},
		},
	}

	for _, tc := range testCases {
		req := httptest.NewRequest(echo.GET, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/showFollowing")
		c.SetParamNames("username")
		c.SetParamValues(tc.input.target)

		if tc.positive {
			if assert.NoError(t, h.ShowFollowing(c), tc.input.target) {
				assert.Equal(t, tc.expected.code, rec.Code, tc.input.target)
				assert.Equal(t, tc.expected.message, rec.Body.String(), tc.input.target)
			}
		} else {
			if err := h.ShowFollowing(c); assert.Error(t, err, tc.input.target) {
				assert.Equal(t, tc.expected.code, err.(*echo.HTTPError).Code, tc.input.target)
				assert.Equal(t, tc.expected.message, err.(*echo.HTTPError).Message, tc.input.target)
			}
		}
	}
	h.Shutdown()
}

func TestUpdateUserInfo(t *testing.T) {
	e := echo.New()
	h := NewHandler("mongodb://SEavenger:SEavenger@ds121225.mlab.com:21225/se_avengers_test", key, make(chan model.Notification))

	testCases := []restHTTPTestCase {
		// test success
		restHTTPTestCase{
			true,
			restHTTPInput{"testUpdateUserInfo", `{"firstname":"testUpdateUserInfo_mod","lastname":"mod","bio":"mod","tag":"mod"}`},
			restHTTPExpected{http.StatusOK, `{"username":"testUpdateUserInfo","firstname":"testUpdateUserInfo_mod","lastname":"mod","email":"testUpdateUserInfo@gmail.com","bio":"mod","tag":"mod"}`},
		},
		restHTTPTestCase{
			true,
			restHTTPInput{"testUpdateUserInfo", `{"firstname":"testUpdateUserInfo","lastname":"","bio":"","tag":""}`},
			restHTTPExpected{http.StatusOK, `{"username":"testUpdateUserInfo","firstname":"testUpdateUserInfo","email":"testUpdateUserInfo@gmail.com"}`},
		},
		// test fail
		restHTTPTestCase{
			false,
			restHTTPInput{"testUpdateUserInfo_empty_firstname", `{"firstname":"","lastname":"","bio":"","tag":"testUpdate"}`},
			restHTTPExpected{http.StatusBadRequest, "Firstname must not be empty."},
		},
		restHTTPTestCase{
			false,
			restHTTPInput{"testUpdateUserInfo_not_exist", `{"firstname":"testUpdateUserInfo_not_exist","lastname":"","bio":"","tag":""}`},
			restHTTPExpected{http.StatusNotFound, "User does not exist."},
		},
	}

	for _, tc := range testCases {
		req := httptest.NewRequest(echo.POST, "/", strings.NewReader(tc.input.target))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, "Bearer " + getToken(tc.input.user))
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/updateUserInfo")
		processJWTToken(c)

		if tc.positive {
			if assert.NoError(t, h.UpdateUserInfo(c), tc.input.target) {
				assert.Equal(t, tc.expected.code, rec.Code, tc.input.target)
				assert.Equal(t, tc.expected.message, rec.Body.String(), tc.input.target)
			}
		} else {
			if err := h.UpdateUserInfo(c); assert.Error(t, err, tc.input.target) {
				assert.Equal(t, tc.expected.code, err.(*echo.HTTPError).Code, tc.input.target)
				assert.Equal(t, tc.expected.message, err.(*echo.HTTPError).Message, tc.input.target)
			}
		}
	}
	h.Shutdown()
}

func TestUpdateProfilePicture(t *testing.T) {
	e := echo.New()
	h := NewHandler("mongodb://SEavenger:SEavenger@ds121225.mlab.com:21225/se_avengers_test", key, make(chan model.Notification))

	oversizedImage := [10485760 + 1]byte{}
	for i := 0; i < 10485760 + 1; i++ {
		oversizedImage[i] = 'A'
	}

	testCases := []restHTTPTestCase {
		// test success
		restHTTPTestCase{
			true,
			restHTTPInput{"testUpdateProfilePicture", `{"picture":"iVBORw0KGgoAAAANSUhEUgAAACAAAAAgCAMAAABEpIrGAAAAjVBMVEXp6en///+8vLwPZ66L1fPY` + 
														  `//9Xo95eqeH4/f8wg8AQaK6U2flirOJkruJ4wutnseTf8/+F0PF0v+p5y+7l9/+CzO9puefw+//2` + 
														  `9vZ8xux/ye3X8f/b29trteWt6f+V0/O48P8YbbCo3/WVrcC1zeHFxcWo5/+h4f/K+/++9P+FwuzM` + 
														  `/f/H6fgjd7hnvekVvPy4AAABSElEQVQ4jXXT23KDIBCAYbY5WaNiqEGJWoOao0ne//G6gBjphH/G` + 
														  `GVk/1ysJQC/DxcdC2QMQ6J/tylP77IEUsm3WnppWFoQsVuuNp/VqQRAMm8jTZhhBjEX4RhS7GRAO` + 
														  `ah7dy7S8/xPRECpwjSmlUamKqFN8taCmcYbP0zeosQmo0wjscg1qasBNnWiWIshsKdPgZkBtwLys` + 
														  `1Dss4JzXuQtyPRwB10CtyN7fUIAbcOSMMZ6rcW7L7np4NIBZkPMpppuDQAPmZsBvgvEAFwc8cfsd` + 
														  `gRAJ0+A1PRIqC/Q9aPBygDDgbACuCGCq0eBswAETgRMEQk1HsMcOD3CFOODUgIsG+xzmBUJNLwbs` + 
														  `8Bavxxw0arKzwLT/mWVGI1juPC01kMl26WmbSAQdCk+J7Agpqk4+vz8mZVcVBIrq9OXpVBX4dwP+` + 
														  `n74KgD9PVzqtRZGhswAAAABJRU5ErkJggg=="}`},
			restHTTPExpected{http.StatusCreated, `{"username":"testUpdateProfilePicture","firstname":"testUpdateProfilePicture","email":"testUpdateProfilePicture@gmail.com",` + 
											   `"picture":"iVBORw0KGgoAAAANSUhEUgAAACAAAAAgCAMAAABEpIrGAAAAjVBMVEXp6en///+8vLwPZ66L1fPY` + 
														  `//9Xo95eqeH4/f8wg8AQaK6U2flirOJkruJ4wutnseTf8/+F0PF0v+p5y+7l9/+CzO9puefw+//2` + 
														  `9vZ8xux/ye3X8f/b29trteWt6f+V0/O48P8YbbCo3/WVrcC1zeHFxcWo5/+h4f/K+/++9P+FwuzM` + 
														  `/f/H6fgjd7hnvekVvPy4AAABSElEQVQ4jXXT23KDIBCAYbY5WaNiqEGJWoOao0ne//G6gBjphH/G` + 
														  `GVk/1ysJQC/DxcdC2QMQ6J/tylP77IEUsm3WnppWFoQsVuuNp/VqQRAMm8jTZhhBjEX4RhS7GRAO` + 
														  `ah7dy7S8/xPRECpwjSmlUamKqFN8taCmcYbP0zeosQmo0wjscg1qasBNnWiWIshsKdPgZkBtwLys` + 
														  `1Dss4JzXuQtyPRwB10CtyN7fUIAbcOSMMZ6rcW7L7np4NIBZkPMpppuDQAPmZsBvgvEAFwc8cfsd` + 
														  `gRAJ0+A1PRIqC/Q9aPBygDDgbACuCGCq0eBswAETgRMEQk1HsMcOD3CFOODUgIsG+xzmBUJNLwbs` + 
														  `8Bavxxw0arKzwLT/mWVGI1juPC01kMl26WmbSAQdCk+J7Agpqk4+vz8mZVcVBIrq9OXpVBX4dwP+` + 
											  			  `n74KgD9PVzqtRZGhswAAAABJRU5ErkJggg=="}`},
		},
		restHTTPTestCase{
			true,
			restHTTPInput{"testUpdateProfilePicture", `{"picture":""}`},
			restHTTPExpected{http.StatusCreated, `{"username":"testUpdateProfilePicture","firstname":"testUpdateProfilePicture","email":"testUpdateProfilePicture@gmail.com"}`},
		},
		// test fail
		restHTTPTestCase{
			false,
			restHTTPInput{"testUpdateProfilePicture_oversized_image", `{"picture":"` + string(oversizedImage[:]) + `"}`},
			restHTTPExpected{http.StatusBadRequest, "Image must be smaller than 10 MB."},
		},
		restHTTPTestCase{
			false,
			restHTTPInput{"testUpdateProfilePicture_not_exist", `{"picture":""}`},
			restHTTPExpected{http.StatusNotFound, "User does not exist."},
		},
	}

	for _, tc := range testCases {
		req := httptest.NewRequest(echo.POST, "/", strings.NewReader(tc.input.target))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, "Bearer " + getToken(tc.input.user))
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/updateProfilePic")
		processJWTToken(c)

		if tc.positive {
			if assert.NoError(t, h.UpdateProfilePicture(c), tc.input.target) {
				assert.Equal(t, tc.expected.code, rec.Code, tc.input.target)
				assert.Equal(t, tc.expected.message, rec.Body.String(), tc.input.target)
			}
		} else {
			if err := h.UpdateProfilePicture(c); assert.Error(t, err, tc.input.target) {
				assert.Equal(t, tc.expected.code, err.(*echo.HTTPError).Code, tc.input.target)
				assert.Equal(t, tc.expected.message, err.(*echo.HTTPError).Message, tc.input.target)
			}
		}
	}
	h.Shutdown()
}

func getToken(username string) string {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = username
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	Token, err := token.SignedString([]byte(key))
	if err != nil {
		panic(err)
	}

	return Token
}

func processJWTToken(c echo.Context) {
	type JWTConfig struct {
		SigningKey 		interface{}
		SigningMethod	string
		ContextKey		string
		TokenLookup		string
		AuthScheme 		string
		Claims 			jwt.Claims
		keyFunc 		jwt.Keyfunc
	}
	var DefaultJWTConfig = JWTConfig {
		SigningKey: 	[]byte(key),
		SigningMethod: 	"HS256",
		ContextKey:    	"user",
		TokenLookup:   	"header:" + echo.HeaderAuthorization,
		AuthScheme:    	"Bearer",
		Claims:       	jwt.MapClaims{},
	}
	DefaultJWTConfig.keyFunc = func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != DefaultJWTConfig.SigningMethod {
			return nil, fmt.Errorf("Unexpected jwt signing method=%v", t.Header["alg"])
		}
		return DefaultJWTConfig.SigningKey, nil
	}

	parts := strings.Split(DefaultJWTConfig.TokenLookup, ":")
	extractor := jwtFromHeader(parts[1], DefaultJWTConfig.AuthScheme)

	auth, err := extractor(c)
	token := new(jwt.Token)
	
	if _, ok := DefaultJWTConfig.Claims.(jwt.MapClaims); ok {
		token, err = jwt.Parse(auth, DefaultJWTConfig.keyFunc)
	}

	if err == nil && token.Valid {
		c.Set(DefaultJWTConfig.ContextKey, token)
	}
}

var ErrJWTMissing = echo.NewHTTPError(http.StatusBadRequest, "Missing or malformed jwt")
type jwtExtractor func(echo.Context) (string, error)

func jwtFromHeader(header string, authScheme string) jwtExtractor {
	return func(c echo.Context) (string, error) {
		auth := c.Request().Header.Get(header)
		l := len(authScheme)
		if len(auth) > l+1 && auth[:l] == authScheme {
			return auth[l+1:], nil
		}
		return "", ErrJWTMissing
	}
}
