package handler

import (
	"crud/internal/domain"
	"crud/internal/pkg/authclient"
	"crud/internal/service"
	"encoding/json"
	"github.com/valyala/fasthttp"
	"log"
)

func ServerHandler(ctx *fasthttp.RequestCtx) {

	ctx.Response.Header.Set(fasthttp.HeaderAccessControlAllowOrigin, "*")
	ctx.Response.Header.Set(fasthttp.HeaderAccessControlAllowMethods, fasthttp.MethodPost)
	ctx.Response.Header.Set(fasthttp.HeaderAccessControlAllowMethods, fasthttp.MethodGet)
	ctx.Response.Header.Set(fasthttp.HeaderAccessControlAllowMethods, fasthttp.MethodDelete)
	ctx.Response.Header.Set(fasthttp.HeaderAccessControlAllowHeaders, fasthttp.HeaderContentType)
	ctx.Response.Header.Set(fasthttp.HeaderAccessControlAllowHeaders, fasthttp.HeaderAuthorization)

	if ctx.IsOptions() {
		return
	}

	switch {
	case ctx.IsGet():
		GetHandler(ctx)
	case ctx.IsDelete():
		DeleteHandler(ctx)
	case ctx.IsPost():
		PostHandler(ctx)
	}

}

func GetHandler(ctx *fasthttp.RequestCtx) {
	id := ctx.QueryArgs().Peek("id")
	if len(id) == 0 {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}
	rec, err := service.Get(string(id))
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return
	}

	marshal, err := json.Marshal(rec)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	if _, err = ctx.Write(marshal); err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
}

func DeleteHandler(ctx *fasthttp.RequestCtx) {

	id := ctx.QueryArgs().Peek("id")
	if len(id) == 0 {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}
	recipeById, err := service.Get(string(id))
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return
	}

	authorId := recipeById.AuthorID

	token := ctx.Request.Header.Peek(fasthttp.HeaderAuthorization)
	userInfo, err := authclient.ValidateToken(string(token))
	//log.Println(string(token) == "", err != nil, (*userInfo).UserRole == "admin", string(id) == (*userInfo).UserId, string(token) == "" || err != nil)
	if err != nil || string(token) == "" {
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		log.Println("Get request", string(ctx.Method()), string(token), "error", fasthttp.StatusUnauthorized)
		return
	}

	//log.Println("(*userInfo).UserRole = ", (*userInfo).UserRole, "(*userInfo).UserId = ", (*userInfo).UserId)
	//log.Println(string(token) == "", (*userInfo).UserRole == "admin", authorId == (*userInfo).UserId)
	if ((*userInfo).UserRole == "admin") || (authorId == (*userInfo).UserId) {

		err := service.Delete(string(id))
		if err != nil {
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			log.Println("Delete request", string(ctx.Method()), string(token), "error")
			return
		}

	} else {
		ctx.SetStatusCode(fasthttp.StatusForbidden)
		log.Println("Get request", string(ctx.Method()), string(token), "error", fasthttp.StatusForbidden)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)

}

func PostHandler(ctx *fasthttp.RequestCtx) {

	token := ctx.Request.Header.Peek(fasthttp.HeaderAuthorization)
	userInfo, err := authclient.ValidateToken(string(token))
	//log.Println(string(token) == "", err != nil, string(token) == "" || err != nil)
	if err != nil || string(token) == "" {
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		//log.Println("Get request", string(ctx.Method()), string(token), "error", fasthttp.StatusUnauthorized)
		return
	}

	var rec domain.Recipe

	if err := json.Unmarshal(ctx.PostBody(), &rec); err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	rec.AuthorID = (*userInfo).UserId

	if err := service.AddOrUpd(&rec); err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	resp := IdResponse{ID: rec.ID}

	marshal, err := json.Marshal(resp)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	if _, err := ctx.Write(marshal); err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
}
