package httphandler

import (
	"authservice/internal/domain"
	"authservice/internal/service"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func AdminGetUserInfo(resp http.ResponseWriter, req *http.Request) {

	respBody := &HTTPResponse{}
	defer func() {
		resp.Write(respBody.Marshall())
	}()

	id := req.URL.Query().Get("user_id")
	userID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		respBody.SetError(errors.New("invalid input"))
		return
	}

	info, err := service.GetUserFullInfo(userID)
	if err != nil {
		resp.WriteHeader(http.StatusNotFound)
		respBody.SetError(err)
	}

	respBody.SetData(info)
}

func AdminBlockedUser(resp http.ResponseWriter, req *http.Request) {
	respBody := &HTTPResponse{}
	defer func() {
		resp.Write(respBody.Marshall())
	}()

	var input domain.SetBlockUser

	if err := readBody(req, &input); err != nil {
		resp.WriteHeader(http.StatusUnprocessableEntity)
		respBody.SetError(err)
		return
	}

	if !input.IsValid() {
		resp.WriteHeader(http.StatusBadRequest)
		respBody.SetError(errors.New("invalid input"))
		return
	}

	err := service.BlockUser(&input)

	if err != nil {
		resp.WriteHeader(http.StatusConflict)
		respBody.SetError(err)
		return
	}
}

func AdminUnblockedUser(resp http.ResponseWriter, req *http.Request) {
	respBody := &HTTPResponse{}
	defer func() {
		resp.Write(respBody.Marshall())
	}()

	var input domain.SetBlockUser

	if err := readBody(req, &input); err != nil {
		resp.WriteHeader(http.StatusUnprocessableEntity)
		respBody.SetError(err)
		return
	}

	if !input.IsValid() {
		resp.WriteHeader(http.StatusBadRequest)
		respBody.SetError(errors.New("invalid input"))
		return
	}

	err := service.UnblockUser(&input)

	if err != nil {
		resp.WriteHeader(http.StatusConflict)
		respBody.SetError(err)
		return
	}
}

func AdminSetRoleUser(resp http.ResponseWriter, req *http.Request) {

	respBody := &HTTPResponse{}
	defer func() {
		resp.Write(respBody.Marshall())
	}()

	var input domain.UserRole

	if err := readBody(req, &input); err != nil {
		resp.WriteHeader(http.StatusUnprocessableEntity)
		respBody.SetError(err)
		return
	}

	if !input.IsValid() {
		resp.WriteHeader(http.StatusBadRequest)
		respBody.SetError(errors.New("invalid input"))
		return
	}

	err := service.SetRoleUser(&input)

	if err != nil {
		resp.WriteHeader(http.StatusConflict)
		respBody.SetError(err)
		return
	}
}

func AdminChangePsw(resp http.ResponseWriter, req *http.Request) {
	respBody := &HTTPResponse{}
	defer func() {
		resp.Write(respBody.Marshall())
	}()

	var input domain.UserPassword

	if err := readBody(req, &input); err != nil {
		resp.WriteHeader(http.StatusUnprocessableEntity)
		respBody.SetError(err)
		return
	}

	if !input.IsValid() {
		resp.WriteHeader(http.StatusBadRequest)
		respBody.SetError(errors.New("invalid input"))
		return
	}

	err := service.ChangePsw(&input)
	if err != nil {
		resp.WriteHeader(http.StatusNotFound)
		respBody.SetError(err)
		return
	}

	info, err := service.GetUserFullInfo(input.ID)
	newPassword := input.Password
	userName := info.Name
	userEmail := info.Email

	subject := "Your New Password"
	body := fmt.Sprintf("Hello, %s.\nYour new password is: %s.\nPlease change it after logging in.", userName, newPassword)

	if err := sendEmail(userEmail, subject, body); err != nil {
		http.Error(resp, "Failed to send email", http.StatusInternalServerError)
		return
	}
}

func sendEmail(to, subject, body string) error {
	fmt.Printf("\nEmail: %s, Subject: %s\n\n%s\n", to, subject, body)
	return nil
}
