package tweet

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
	"encoding/json"
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

func TestNewTweet (t *testing.T) {
	// Setup
	e := echo.New()
	nh := notification.NewHandler("mongodb://SEavenger:SEavenger@ds121225.mlab.com:21225/se_avengers_test")
	h := NewHandler("mongodb://SEavenger:SEavenger@ds121225.mlab.com:21225/se_avengers_test", key, nh.Manager.Operator)

	const illegalSize int = 10485760+1
	oversizedImage := [illegalSize]byte{}
	for i := 0; i < illegalSize; i++ {
		oversizedImage[i] = 'A'
	}

	// Test cases
	testCases := []restHTTPTestCase {
		// test success
		restHTTPTestCase{
			true,
			restHTTPInput{"testNewTweet", `{"message":"testMessage"}`},
			restHTTPExpected{http.StatusOK, `{"owner":"testNewTweet","message":"testMessage","picture":""}`},
		},
		restHTTPTestCase{
			true,
			restHTTPInput{"testNewTweet", `{"message":"testMessage", "picture":""}`},
			restHTTPExpected{http.StatusOK, `{"owner":"testNewTweet","message":"testMessage","picture":""}`},
		},
		
		// test empty
		restHTTPTestCase{
			false,
			restHTTPInput{"testNewTweet", `{"message":"", "picture":"iVBORw0KGgoAAAANSUhEUgAAACAAAAAgCAMAAABEpIrGAAAAjVBMVEXp6en///+8vLwPZ66L1fPYiVBORw0KGgoAAAANSUhEUgAAACAAAAAgCAMAAABEpIrGAAAAjVBMVEXp6en///+8vLwPZ66L1fPYiVBORw0KGgoAAAANSUhEUgAAACAAAAAgCAMAAABEpIrGAAAAjVBMVEXp6en///+8vLwPZ66L1fPYiVBORw0KGgoAAAANSUhEUgAAACAAAAAgCAMAAABEpIrGAAAAjVBMVEXp6en///+8vLwPZ66L1fPYiVBORw0KGgoAAAANSUhEUgAAACAAAAAgCAMAAABEpIrGAAAAjVBMVEXp6en///+8vLwPZ66L1fPYiVBORw0KGgoAAAANSUhEUgAAACAAAAAgCAMAAABEpIrGAAAAjVBMVEXp6en///+8vLwPZ66L1fPYiVBORw0KGgoAAAANSUhEUgAAACAAAAAgCAMAAABEpIrGAAAAjVBMVEXp6en///+8vLwPZ66L1fPY"}`},
			restHTTPExpected{http.StatusBadRequest, "Message cannot be empty."},
		},
		
		// test illegal picture
		restHTTPTestCase{
			false,
			restHTTPInput{"testNewTweet", `{"message":"testMessage", "picture":"` + string(oversizedImage[:]) + `"}`},
			restHTTPExpected{http.StatusBadRequest, "Image must be smaller than 10 MB."},
		},
	}
	
	// Run
	for _, tc := range testCases {
		// Delete DB entry
		if tc.positive {
			_, err := h.db.DB(h.dbName).C(model.TweetCollection).RemoveAll(bson.M{"owner": "testSignup"})
			if err != nil {
				if err != mgo.ErrNotFound {
					t.Fatal(err)
				}
			}
		}

		// Setup
		req := httptest.NewRequest(echo.POST, "/", strings.NewReader(tc.input.target))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, "Bearer " + getToken(tc.input.user))
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/newTweet")
		processJWTToken(c)

		// Assertion
		if tc.positive {
			if assert.NoError(t, h.NewTweet(c), tc.input.target) {
				assert.Equal(t, tc.expected.code, rec.Code, tc.input.target)
				assert.Equal(t, tc.expected.message, rec.Body.String(), tc.input.target)
			}
		} else {
			if err := h.NewTweet(c); assert.Error(t, err, tc.input.target) {
				assert.Equal(t, tc.expected.code, err.(*echo.HTTPError).Code, tc.input.target)
				assert.Equal(t, tc.expected.message, err.(*echo.HTTPError).Message, tc.input.target)
			}
		}
	}
	h.Shutdown()
}

func TestFetchOwnTweets (t *testing.T) {
	e := echo.New()
	h := NewHandler("mongodb://SEavenger:SEavenger@ds121225.mlab.com:21225/se_avengers_test", key, make(chan model.Notification))

	testCases := []restHTTPTestCase {
		// test success
		restHTTPTestCase{
			true,
			restHTTPInput{"testFetchOenTweetsSuccess", ""},
			restHTTPExpected{http.StatusOK, `{"page":"1","totalpage":"4","totaltweets":"10","tweetlist":"[]"}`},
		},
		// test over page
		restHTTPTestCase{
			true,
			restHTTPInput{"testFetchOenTweetsOverPage", ""},
			restHTTPExpected{http.StatusOK, `{"page":"1","totalpage":"1","totaltweets":"1","tweetlist":"[]"}`},
		},
		// test the user with no tweet
		restHTTPTestCase{
			true,
			restHTTPInput{"testFetchOenTweetsWithNoTweet1", ""},
			restHTTPExpected{http.StatusOK, `{"page":"0","totalpage":"0","totaltweets":"0","tweetlist":"[]"}`},
		},
		restHTTPTestCase{
			true,
			restHTTPInput{"testFetchOenTweetsWithNoTweet2", ""},
			restHTTPExpected{http.StatusOK, `{"page":"0","totalpage":"0","totaltweets":"0","tweetlist":"[]"}`},
		},
		// tes username not in db
		restHTTPTestCase{
			false,
			restHTTPInput{"userNotExist", ""},
			restHTTPExpected{http.StatusNotFound, "User userNotExist does not exist."},
		},
	}
	// Run
	for _, tc := range testCases {
		// Setup
		req := httptest.NewRequest(echo.GET, "/?page=1&perpage=3", nil)
		req.Header.Set(echo.HeaderAuthorization, "Bearer " + getToken(tc.input.user))
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("api/v1/tweetlist")
		c.SetParamNames("username")
		c.SetParamValues(tc.input.user)
		processJWTToken(c)

		// Assertion
		if tc.positive {
			if assert.NoError(t, h.FetchTweets(c), tc.input.user + " fetching own tweet") {
				assert.Equal(t, tc.expected.code, rec.Code, tc.input.user + " fetching own tweet")
				assert.Equal(t, tc.expected.message, removeTweetList(rec.Body.String()), tc.input.user + " fetching own tweet")
			}
		} else {
			if err := h.FetchTweets(c); assert.Error(t, err, tc.input.user + " fetching own tweet") {
				assert.Equal(t, tc.expected.code, err.(*echo.HTTPError).Code, tc.input.user + " fetching own tweet")
				assert.Equal(t, tc.expected.message, err.(*echo.HTTPError).Message, tc.input.user + " fetching own tweet")
			}
		}
	}
}

func TestFetchTweetTimeLine (t *testing.T) {
	e := echo.New()
	h := NewHandler("mongodb://SEavenger:SEavenger@ds121225.mlab.com:21225/se_avengers_test", key, make(chan model.Notification))

	testCases := []restHTTPTestCase {
		// test success
		restHTTPTestCase{
			true,
			restHTTPInput{"TestFetchTweetTimeLineSuccess", ""},
			restHTTPExpected{http.StatusOK, `{"page":"1","totalpage":"3","totaltweets":"8","tweetlist":"[]"}`},
		},
		// test over page
		restHTTPTestCase{
			true,
			restHTTPInput{"TestFetchTweetTimeLineOverPage", ""},
			restHTTPExpected{http.StatusOK, `{"page":"0","totalpage":"0","totaltweets":"0","tweetlist":"[]"}`},
		},
		// test the user with no tweet
		restHTTPTestCase{
			true,
			restHTTPInput{"TestFetchTweetTimeLineNoTweet1", ""},
			restHTTPExpected{http.StatusOK, `{"page":"0","totalpage":"0","totaltweets":"0","tweetlist":"[]"}`},
		},
		restHTTPTestCase{
			true,
			restHTTPInput{"TestFetchTweetTimeLineNoTweet2", ""},
			restHTTPExpected{http.StatusOK, `{"page":"0","totalpage":"0","totaltweets":"0","tweetlist":"[]"}`},
		},
		// test username not in db
		restHTTPTestCase{
			false,
			restHTTPInput{"userNotExist", ""},
			restHTTPExpected{http.StatusNotFound, "User userNotExist does not exist."},
		},
	}
	// Run
	for _, tc := range testCases {
		// Setup
		req := httptest.NewRequest(echo.GET, "/?page=1&perpage=3", nil)
		req.Header.Set(echo.HeaderAuthorization, "Bearer " + getToken(tc.input.user))
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("api/v1/tweettimeline")
		c.SetParamNames("username")
		c.SetParamValues(tc.input.user)
		processJWTToken(c)

		// Assertion
		if tc.positive {
			if assert.NoError(t, h.FetchTweetTimeLine(c), tc.input.user + " fetching tweetTimeLine") {
				assert.Equal(t, tc.expected.code, rec.Code, tc.input.user + " fetching tweetTimeLine")
				assert.Equal(t, tc.expected.message, removeTweetList(rec.Body.String()), tc.input.user + " fetching tweetTimeLine")
			}
		} else {
			if err := h.FetchTweetTimeLine(c); assert.Error(t, err, tc.input.user + " fetching tweetTimeLine") {
				assert.Equal(t, tc.expected.code, err.(*echo.HTTPError).Code, tc.input.user + " fetching tweetTimeLine")
				assert.Equal(t, tc.expected.message, err.(*echo.HTTPError).Message, tc.input.user + " fetching tweetTimeLine")
			}
		}
	}
}

func getIDToken(id string) string {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = id
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	Token, err := token.SignedString([]byte(key))
	if err != nil {
		panic(err)
	}

	return Token
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

func removeTweetList(recorderString string) string {
	in := []byte(recorderString)
	var raw map[string]string
	json.Unmarshal(in, &raw)
	raw["tweetlist"] = "[]"
	out, _ := json.Marshal(raw)
	return string(out)
}
