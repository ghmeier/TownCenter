package handlers

import (
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/alexcesaro/statsd.v2"
	"gopkg.in/gin-gonic/gin.v1"

	"github.com/dgrijalva/jwt-go"

	"github.com/ghmeier/bloodlines/gateways"
	"github.com/ghmeier/bloodlines/handlers"
	"github.com/jakelong95/TownCenter/helpers"
	"github.com/jakelong95/TownCenter/models"
)

type UserI interface {
	New(ctx *gin.Context)
	ViewAll(ctx *gin.Context)
	View(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
	Login(ctx *gin.Context)
	Upload(ctx *gin.Context)
	Time() gin.HandlerFunc
	GetJWT() gin.HandlerFunc
}

type User struct {
	*handlers.BaseHandler
	Helper     helpers.UserI
	Bloodlines gateways.Bloodlines
}

func NewUser(ctx *handlers.GatewayContext) UserI {
	stats := ctx.Stats.Clone(statsd.Prefix("api.user"))
	return &User{
		BaseHandler: &handlers.BaseHandler{Stats: stats},
		Helper:      helpers.NewUser(ctx.Sql, ctx.S3),
		Bloodlines:  ctx.Bloodlines,
	}
}

func (u *User) New(ctx *gin.Context) {
	//Bind the json to a user object
	var json models.User
	err := ctx.BindJSON(&json)
	if err != nil {
		u.UserError(ctx, "Error: Unable to parse json", err)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(json.PassHash), bcrypt.DefaultCost)
	if err != nil {
		u.ServerError(ctx, err, nil)
	}

	//Create the new user in the database
	user := models.NewUser(string(hashedPassword), json.FirstName, json.LastName, json.Email, json.Phone,
		json.AddressLine1, json.AddressLine2, json.AddressCity, json.AddressState, json.AddressZip,
		json.AddressCountry)
	err = u.Helper.Insert(user)
	if err != nil {
		u.ServerError(ctx, err, json)
		return
	}

	//Don't need to pass the password hash back
	user.PassHash = ""

	_, err = u.Bloodlines.NewPreference(user.ID)
	if err != nil {
		u.ServerError(ctx, err, json)
		return
	}

	signedToken, _ := CreateJWT()

	ctx.Header("Auth", signedToken)
	u.Success(ctx, user)
}

func (u *User) ViewAll(ctx *gin.Context) {
	//Use paging when getting lists of users
	offset, limit := u.GetPaging(ctx)

	//Query the database for all users
	users, err := u.Helper.GetAll(offset, limit)
	if err != nil {
		u.ServerError(ctx, err, users)
		return
	}

	//Don't pass the password hashes back
	for _, user := range users {
		user.PassHash = ""
	}

	u.Success(ctx, users)
}

func (u *User) View(ctx *gin.Context) {
	userId := ctx.Param("userId")

	//Query the database for the user
	user, err := u.Helper.GetByID(userId)
	if err != nil {
		u.ServerError(ctx, err, userId)
		return
	}

	if user == nil {
		u.NotFoundError(ctx, "Error: User with ID "+userId+" does not exist")
		return
	}

	//Don't pass the password hash back
	user.PassHash = ""

	u.Success(ctx, user)
}

func (u *User) Update(ctx *gin.Context) {
	userId := ctx.Param("userId")

	//Bind the json to a user object
	var json models.User
	err := ctx.BindJSON(&json)
	if err != nil {
		u.UserError(ctx, "Error: Unable to parse json", err)
		return
	}

	if json.PassHash != "" {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(json.PassHash), bcrypt.DefaultCost)
		json.PassHash = string(hashedPassword)
	}

	//Update the user in the database
	err = u.Helper.Update(&json, userId)
	if err != nil {
		u.ServerError(ctx, err, userId)
		return
	}

	//Don't pass the password hash bash
	json.PassHash = ""

	u.Success(ctx, json)
}

func (u *User) Delete(ctx *gin.Context) {
	userId := ctx.Param("userId")

	//Delete the user from the database
	err := u.Helper.Delete(userId)
	if err != nil {
		u.ServerError(ctx, err, userId)
		return
	}

	u.Success(ctx, nil)
}

func (u *User) Login(ctx *gin.Context) {
	//Bind the json to a user object
	var json models.User
	err := ctx.BindJSON(&json)
	if err != nil {
		u.UserError(ctx, "Error: Unable to parse json", err)
		return
	}

	//Get the user from the database
	user, err := u.Helper.GetByEmail(json.Email)
	if err != nil {
		u.ServerError(ctx, err, json.Email)
		return
	}

	if user == nil {
		u.NotFoundError(ctx, "Error: User with email "+json.Email+" not found")
		return
	}

	//Don't pass the password hash back
	tmpHash := user.PassHash
	user.PassHash = ""

	if err != nil {
		u.ServerError(ctx, err, json.Email)
	}

	err = bcrypt.CompareHashAndPassword([]byte(tmpHash), []byte(json.PassHash))

	if err == nil {
		signedToken, _ := CreateJWT()

		ctx.Header("X-Auth", signedToken)
		u.Success(ctx, user)
	} else {
		u.UserError(ctx, "Incorrect login credentials", nil)
	}
}

func (u *User) Upload(ctx *gin.Context) {
	id := ctx.Param("userId")
	file, headers, err := ctx.Request.FormFile("profile")
	if err != nil {
		u.UserError(ctx, "ERROR: unable to find body", nil)
	}
	defer file.Close()

	err = u.Helper.Profile(id, headers.Filename, file)
	if err != nil {
		u.ServerError(ctx, err, id)
		return
	}

	u.Success(ctx, nil)
}

/*CreateJWT creates a new JSON Web Token that expires in 30 days*/
func CreateJWT() (string, error) {
	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * 24 * 30).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT")))
}
