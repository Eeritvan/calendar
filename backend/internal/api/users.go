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

	jwtToken, err := utils.GenerateJWT(queryResp.ID.String())
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

	jwtToken, err := utils.GenerateJWT(queryResp.ID.String())
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
// TODO: verify that totp is not enabled already
func (s *Server) PostTotpEnable(c echo.Context) error {
	userId := c.Get("userId").(uuid.UUID)

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "void",
		AccountName: "#TODO",
	})
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
			"userId":     userId.String(),
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

	userUUID := c.Get("userId").(uuid.UUID)

	if userId != userUUID.String() {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	isValid := totp.Validate(strconv.Itoa(body.Code), totpSecret)

	if !isValid {
		return c.JSON(http.StatusInternalServerError, nil)
	}

	codes := make([]string, 5)
	for i := range len(codes) {
		code := utils.GenerateRecoveryCode()
		codes[i] = code
	}

	hashedCodes, err := utils.HashRecoveryCodes(codes)

	ctx := c.Request().Context()
	s.queries.ClearRecoveryCodes(ctx, userUUID)
	s.queries.InsertRecoveryCodes(ctx, sqlc.InsertRecoveryCodesParams{
		UserID:  userUUID,
		Column2: hashedCodes,
	})
	s.queries.EnableTotp(ctx, sqlc.EnableTotpParams{
		Totp: totpSecret,
		ID:   userUUID,
	})

	resp := RecoveryCodes{
		RecoveryCodes: codes,
	}

	return c.JSON(http.StatusOK, resp)
}

// (PATCH /totp/disable)
func (s *Server) PatchTotpDisable(c echo.Context) error {
	body := new(Totp)
	if err := c.Bind(&body); err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	userId := c.Get("userId").(uuid.UUID)

	ctx := c.Request().Context()
	queryResp, err := s.queries.GetTotpSecret(ctx, userId)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, false)
	}

	isValid := totp.Validate(strconv.Itoa(body.Code), queryResp.Totp)

	if !isValid {
		return c.JSON(http.StatusInternalServerError, false)
	}
	// todo: transactions
	s.queries.DisableTotp(ctx, userId)
	s.queries.ClearRecoveryCodes(ctx, userId)
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

	jwtToken, err := utils.GenerateJWT(queryResp.ID.String())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}

	resp := UserCredentials{
		Name: queryResp.Name,
		JWT:  jwtToken,
	}

	return c.JSON(http.StatusOK, resp)
}

// (POST /totp/recovery)
func (s *Server) PostTotpRecovery(c echo.Context) error {
	body := new(RecoveryCode)
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
	queryResp, err := s.queries.GetUnusedRecoveryCodes(ctx, userUUID)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	matchedId, _ := utils.VerifyRecoveryCode(body.RecoveryCode, queryResp)

	if matchedId == 0 {
		fmt.Println("invalid token")
		return c.JSON(http.StatusInternalServerError, nil)
	}

	s.queries.SetRecoveryCodeAsUsed(ctx, matchedId)
	loginQueryResp, err := s.queries.GetTotpSecret(ctx, userUUID)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	jwtToken, err := utils.GenerateJWT(loginQueryResp.ID.String())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	resp := UserCredentials{
		Name: loginQueryResp.Name,
		JWT:  jwtToken,
	}

	return c.JSON(http.StatusOK, resp)
}
