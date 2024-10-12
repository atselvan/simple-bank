package api

import (
	"database/sql"
	"errors"
	"fmt"
	db "github.com/atselvan/simple-bank/db/sqlc"
	"github.com/gin-gonic/gin"
	"net/http"
)

type (
	listAccountsRequest struct {
		PageId   int32 `form:"page" binding:"required,min=1"`
		PageSize int32 `form:"pageSize" binding:"required,min=5,max=10"`
	}

	getAccountRequest struct {
		Id int64 `uri:"id" binding:"required,min=1"`
	}

	createAccountRequest struct {
		Owner    string `json:"owner" binding:"required"`
		Currency string `json:"currency" binding:"required,oneof=USD EUR"`
	}
)

func (server *Server) listAccounts(ctx *gin.Context) {
	var req listAccountsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, server.errResponse(err))
		return
	}

	params := db.ListAccountsParams{
		Limit:  req.PageSize,
		Offset: (req.PageId - 1) * req.PageSize,
	}
	accounts, err := server.store.ListAccounts(ctx, params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, server.errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}

func (server *Server) getAccount(ctx *gin.Context) {
	var req getAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, server.errResponse(err))
		return
	}

	account, err := server.store.GetAccount(ctx, req.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err := fmt.Errorf("account with id (%d) was not found", req.Id)
			ctx.JSON(http.StatusNotFound, server.errResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, server.errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

func (server *Server) createAccount(ctx *gin.Context) {
	var req createAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, server.errResponse(err))
		return
	}

	params := db.CreateAccountParams{
		Owner:    req.Owner,
		Currency: req.Currency,
		Balance:  0,
	}

	account, err := server.store.CreateAccount(ctx, params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, server.errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}
