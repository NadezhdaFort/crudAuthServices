package httphandler

import (
	"authservice/internal/service"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func isBlocked(next http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {

		id := req.Header.Get(HeaderUserID)
		userID, _ := primitive.ObjectIDFromHex(id)

		userInfo, err := service.GetUserFullInfo(userID)
		if err != nil {
			resp.WriteHeader(http.StatusInternalServerError)

			respBody := &HTTPResponse{}
			respBody.SetError(errors.New("can not find user by id"))
			resp.Write(respBody.Marshall())

			return
		}

		if userInfo.IsBlocked != false {
			resp.WriteHeader(http.StatusForbidden)

			respBody := &HTTPResponse{}
			respBody.SetError(errors.New("user is blocked"))
			resp.Write(respBody.Marshall())

			return
		}

		next.ServeHTTP(resp, req)
	})
}
