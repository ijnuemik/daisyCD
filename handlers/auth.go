package handlers

import (
	models "daisyCD/models"
	commons "daisyCD/commons"
	"gorm.io/gorm"
	"net/http"
	"time"
	"strings"
	"fmt"
	"strconv"
	"errors"
	"github.com/twinj/uuid"
	"github.com/gin-gonic/gin"
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
)

type TokenDetails struct {
	AccessToken    string
	RefreshToken   string
	AccessUuid     string
	RefreshUuid    string
	AtExpires      int64
	RtExpires      int64
}

type AccessDetails struct {
	AccessUuid string
	UserId     uint64
}

var ACCESS_SECRET = viper.GetString("token.ACCESS_SECRET")
var REFRESH_SECRET = viper.GetString("token.REFRESH_SECRET")

func Login(c *gin.Context) {
	// gorm
	db := c.MustGet("db").(*gorm.DB)

	var client models.User
	if err := c.ShouldBindJSON(&client); err != nil {
	   c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
	   return
	}
	
	// check valid username
	var user models.User
	if err := db.First(&user, "username = ?", client.Username).Error; err != nil{
		c.JSON(http.StatusUnauthorized, "Please provide valid login username")
		return
	}

	// check valid password
	if user.Password != client.Password {
		c.JSON(http.StatusUnauthorized, "Please provide valid login password")
	   return
	}

	ts, err := CreateToken(client.Id)
	if err != nil {
	c.JSON(http.StatusUnprocessableEntity, err.Error())
	  return
    }

	saveErr := CreateAuth(client.Id, ts)
	if saveErr != nil {
		c.JSON(http.StatusUnprocessableEntity, saveErr.Error())
	}

	tokens := map[string]string{
		"access_token":  ts.AccessToken,
		"refresh_token": ts.RefreshToken,
	 }
	 
	 c.JSON(http.StatusOK, tokens)
   }

func CreateToken(userid uint64) (*TokenDetails, error) {
	td := &TokenDetails{}
	td.AtExpires = time.Now().Add(time.Minute * 15).Unix()
	td.AccessUuid = uuid.NewV4().String()
	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	td.RefreshUuid = uuid.NewV4().String()
   
	var err error
	//Creating Access Token
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUuid
	atClaims["user_id"] = userid
	atClaims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(ACCESS_SECRET))
	if err != nil {
	   return nil, err
	}
	//Creating Refresh Token
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUuid
	rtClaims["user_id"] = userid
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(REFRESH_SECRET))
	if err != nil {
	   return nil, err
	}
	return td, nil
}

func CreateAuth(userid uint64, td *TokenDetails) error {
	client := commons.GetClient()

    at := time.Unix(td.AtExpires, 0) //converting Unix to UTC
    rt := time.Unix(td.RtExpires, 0)
    now := time.Now()
 
    errAccess := client.Set(td.AccessUuid, strconv.Itoa(int(userid)), at.Sub(now)).Err()
    if errAccess != nil {
        return errAccess
    }
    errRefresh := client.Set(td.RefreshUuid, strconv.Itoa(int(userid)), rt.Sub(now)).Err()
    if errRefresh != nil {
        return errRefresh
    }
    return nil
}

func ExtractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

// Parse, validate, and return a token.
// keyFunc will receive the parsed token and should return the key for validating.
func VerifyToken(r *http.Request) (token *jwt.Token, err error) {
	tokenString := ExtractToken(r)
	token, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(ACCESS_SECRET), nil
	})

	// if err != nil {
	// 	return nil, err
	// }
	// return token, nil
	return
}

func TokenValid(r *http.Request) error {
	token, err := VerifyToken(r)
	if err != nil {
		return err
	}

	if _, ok := token.Claims.(jwt.Claims); !ok || !token.Valid {
		return err
	}

	return nil
}

func ExtractTokenMetadata(r *http.Request) (*AccessDetails, error) {
	token, err := VerifyToken(r)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUuid, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, err
		}
		userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			return nil, err
		}
		return &AccessDetails{
			AccessUuid: accessUuid,
			UserId:     userId,
		}, nil
	}
	return nil, err
}

func ExtractTokenUserID(r *http.Request) (uint64, error) {
	metadata, err := ExtractTokenMetadata(r)
	if err != nil {
		return 0, err
	}
	return metadata.UserId, err
}

func FetchAuth(authD *AccessDetails) (uint64, error) {
	client := commons.GetClient()
	userid, err := client.Get(authD.AccessUuid).Result()
	if err != nil {
		return 0, err
	}

	userID, _ := strconv.ParseUint(userid, 10, 64)
	if authD.UserId != userID {
		return 0, errors.New("unauthorized")
	}
	return userID, nil
}

func DeleteAuth(givenUuid string) (uint64, error) {
	client := commons.GetClient()
	deleted, err := client.Del(givenUuid).Result()
	if err != nil {
		return 0, err
	}

	return uint64(deleted), nil
}

func Logout(c *gin.Context) {
	metadata, err := ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	delErr := DeleteTokens(metadata)
	if delErr != nil {
		c.JSON(http.StatusUnauthorized, delErr.Error())
		return
	}
	c.JSON(http.StatusOK, metadata)
}

func Refresh(c *gin.Context) {
	mapToken := map[string]string{}
	if err := c.ShouldBindJSON(&mapToken); err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	refreshToken := mapToken["refresh_token"]

	//verify the token
	// os.Setenv("REFRESH_SECRET", "mcmvmkmsdnfsdmfdsjf") //this should be in an env file
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(REFRESH_SECRET), nil
	})

	//if there is an error, the token must have expired
	if err != nil {
		c.JSON(http.StatusUnauthorized, "Refresh token expired")
		return
	}

	//is token valid?
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		c.JSON(http.StatusUnauthorized, err)
		return
	}

	//Since token is valid, get the uuid:
	claims, ok := token.Claims.(jwt.MapClaims) //the token claims should conform to MapClaims
	if ok && token.Valid {
		refreshUuid, ok := claims["refresh_uuid"].(string) //convert the interface to string
		if !ok {
			c.JSON(http.StatusUnprocessableEntity, err)
			return
		}
		userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, "Error occurred")
			return
		}
		//Delete the previous Refresh Token
		deleted, delErr := DeleteAuth(refreshUuid)
		if delErr != nil || deleted == 0 { //if any goes wrong
			c.JSON(http.StatusUnauthorized, "unauthorized")
			return
		}
		//Create new pairs of refresh and access tokens
		ts, createErr := CreateToken(userId)
		if createErr != nil {
			c.JSON(http.StatusForbidden, createErr.Error())
			return
		}
		//save the tokens metadata to redis
		saveErr := CreateAuth(userId, ts)
		if saveErr != nil {
			c.JSON(http.StatusForbidden, saveErr.Error())
			return
		}
		tokens := map[string]string{
			"access_token":  ts.AccessToken,
			"refresh_token": ts.RefreshToken,
		}
		c.JSON(http.StatusCreated, tokens)
	} else {
		c.JSON(http.StatusUnauthorized, "refresh expired")
	}
}

func DeleteTokens(authD *AccessDetails) error {
	client := commons.GetClient()
	//get the refresh uuid
	// refreshUuid := fmt.Sprintf("%s++%d", authD.AccessUuid, authD.UserId)

	//delete access token
	deletedAt, err := client.Del(authD.AccessUuid).Result()
	if err != nil {
		return err
	}

	// //delete refresh token
	// deletedRt, err := client.Del(refreshUuid).Result()
	// if err != nil {
	// 	return err
	// }

	//When the record is deleted, the return value is 1
	if deletedAt != 1 {
		return errors.New("something went wrong")
	}

	return nil
}
