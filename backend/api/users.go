package api

import (
	"fmt"
	"net/http"

	"github.com/eeritvan/calendar/internal/sqlc"
	"github.com/eeritvan/calendar/internal/utils"
	"github.com/labstack/echo/v4"
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

	resp := PostLogin{
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

	jwtToken, err := utils.GenerateJWT(queryResp.Name, queryResp.ID.String())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}

	resp := PostLogin{
		Name: queryResp.Name,
		JWT:  jwtToken,
	}

	return c.JSON(http.StatusOK, resp)
}
