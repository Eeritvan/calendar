package api

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/eeritvan/calendar/internal/sqlc"
	"github.com/eeritvan/calendar/internal/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
)

// (POST /signup)
func (s *Server) PostSignup(c echo.Context) error {
	body := new(Signup)

	if err := c.Bind(&body); err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	if body.Password != body.PasswordConfirmation {
		fmt.Println("passwords did not match")
		return c.JSON(http.StatusInternalServerError, nil)
	}

	hashedPW, err := bcrypt.GenerateFromPassword([]byte(body.Password), 12)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}

	ctx := c.Request().Context()
	queryResp, err := s.queries.Signup(ctx, sqlc.SignupParams{
		Name:         body.Name,
		PasswordHash: string(hashedPW),
	})

	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	jwtToken, err := utils.GenerateJWT(queryResp.Name, queryResp.ID.String())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}

	resp := UserCredentials{
		Name: queryResp.Name,
		JWT:  jwtToken,
	}

	return c.JSON(http.StatusOK, resp)
}

// (POST /login)
func (s *Server) PostLogin(c echo.Context) error {
	body := new(Login)

	if err := c.Bind(&body); err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	ctx := c.Request().Context()
	queryResp, err := s.queries.Login(ctx, body.Name)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(queryResp.PasswordHash), []byte(body.Password)); err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}

	if queryResp.Totp != "" {
		JWTkey := os.Getenv("JWT_KEY")
		secretKey := []byte(JWTkey)
		token := jwt.NewWithClaims(jwt.SigningMethodHS256,
			jwt.MapClaims{
				"userId": queryResp.ID,
				"exp":    time.Now().Add(time.Hour * 1).Unix(),
			})
		returnToken, err := token.SignedString(secretKey)
		if err != nil {
			fmt.Println(err)
			return c.JSON(http.StatusInternalServerError, nil)
		}
		resp := TotpRequired{
			VerificationToken: returnToken,
		}

		return c.JSON(http.StatusOK, resp)
	}

	jwtToken, err := utils.GenerateJWT(queryResp.Name, queryResp.ID.String())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}

	resp := UserCredentials{
		Name: queryResp.Name,
		JWT:  jwtToken,
	}

	return c.JSON(http.StatusOK, resp)
}

// (POST /totp/enable)
func (s *Server) PostTotpEnable(c echo.Context) error {
	username := c.Get("username").(string)
	userIdStr := c.Get("userId").(string)

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "void",
		AccountName: username,
	})
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	userUUID, err := uuid.Parse(userIdStr)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	// TODO: userId or not ???
	JWTkey := os.Getenv("JWT_KEY")
	secretKey := []byte(JWTkey)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"totpSecret": key.Secret(),
			"userId":     userUUID,
			"exp":        time.Now().Add(time.Hour * 1).Unix(),
		})

	returnToken, err := token.SignedString(secretKey)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	resp := EnableTotp{
		VerificationToken: returnToken,
	}

	return c.JSON(http.StatusOK, resp)
}

// (PATCH /totp/enable/verify)
func (s *Server) PatchTotpEnableVerify(c echo.Context) error {
	body := new(EnableTotpVerify)

	if err := c.Bind(&body); err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	JWTkey := os.Getenv("JWT_KEY")
	secretKey := []byte(JWTkey)
	token, err := jwt.Parse(body.VerificationToken, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil || !token.Valid {
		fmt.Println(err)
		return c.JSON(http.StatusUnauthorized, nil)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.JSON(http.StatusInternalServerError, nil)
	}

	totpSecret, _ := claims["totpSecret"].(string)
	userId, _ := claims["userId"].(string)

	userIdStr := c.Get("userId").(string)
	userUUID, err := uuid.Parse(userIdStr)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	if userId != userUUID.String() {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	isValid := totp.Validate(strconv.Itoa(body.Code), totpSecret)

	if !isValid {
		return c.JSON(http.StatusInternalServerError, nil)
	}

	ctx := c.Request().Context()
	s.queries.EnableTotp(ctx, sqlc.EnableTotpParams{
		Totp: totpSecret,
		ID:   userUUID,
	})

	return c.JSON(http.StatusOK, true)
}

// (PATCH /totp/disable)
func (s *Server) PatchTotpDisable(c echo.Context) error {
	body := new(Totp)
	if err := c.Bind(&body); err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	userIdStr := c.Get("userId").(string)
	userUUID, err := uuid.Parse(userIdStr)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, nil)
	}
	ctx := c.Request().Context()
	queryResp, err := s.queries.GetTotpSecret(ctx, userUUID)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, false)
	}
	isValid := totp.Validate(strconv.Itoa(body.Code), queryResp.Totp)

	if !isValid {
		return c.JSON(http.StatusInternalServerError, false)
	}
	s.queries.DisableTotp(ctx, userUUID)
	return c.JSON(http.StatusOK, true)
}

// (POST /totp/authenticate)
func (s *Server) PostTotpAuthenticate(c echo.Context) error {
	body := new(EnableTotpVerify)
	if err := c.Bind(&body); err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	JWTkey := os.Getenv("JWT_KEY")
	secretKey := []byte(JWTkey)
	token, err := jwt.Parse(body.VerificationToken, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil || !token.Valid {
		fmt.Println("err", err)
		return c.JSON(http.StatusUnauthorized, nil)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.JSON(http.StatusInternalServerError, nil)
	}

	userIdStr, _ := claims["userId"].(string)

	userUUID, err := uuid.Parse(userIdStr)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	ctx := c.Request().Context()
	queryResp, err := s.queries.GetTotpSecret(ctx, userUUID)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, false)
	}

	isValid := totp.Validate(strconv.Itoa(body.Code), queryResp.Totp)

	if !isValid {
		return c.JSON(http.StatusInternalServerError, nil)
	}

	jwtToken, err := utils.GenerateJWT(queryResp.Name, queryResp.ID.String())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}

	resp := UserCredentials{
		Name: queryResp.Name,
		JWT:  jwtToken,
	}

	return c.JSON(http.StatusOK, resp)
}
