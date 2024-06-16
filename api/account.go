package api

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/Country7/backend-captaincode-mysql/token"
	"github.com/go-sql-driver/mysql"

	db "github.com/Country7/backend-captaincode-mysql/db/sqlc"
	"github.com/gin-gonic/gin"
	// "github.com/go-sql-driver/mysql"
)

type createAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,currency"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	var req createAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil { // Должен быть привязан JSON
		ctx.JSON(http.StatusBadRequest, errorResponse(err)) // 400
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload) // полезная нагрузка auth

	arg := db.CreateAccountParams{
		Owner:    authPayload.Username,
		Balance:  0,
		Currency: req.Currency,
	}

	err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) {
			switch mysqlErr.Number {
			case 1452: // ER_NO_REFERENCED_ROW_2 (foreign key violation) нарушение внешнего ключа
				// ошибки на сервере при создании аккаунта без юзера
				ctx.JSON(http.StatusForbidden, errorResponse(err)) // 403
				return
			case 1062: // ER_DUP_ENTRY (unique violation) нарушение уникального ограничения
				// и создании аккаунта с одинаковой валютой счета
				ctx.JSON(http.StatusForbidden, errorResponse(err)) // 403
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err)) // 500
		return
	}
	account, err := server.store.GetLastAccount(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getAccount(ctx *gin.Context) {
	var req getAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil { // Должен быть привязан Uri идентификатор
		ctx.JSON(http.StatusBadRequest, errorResponse(err)) // 400
		return
	}

	account, err := server.store.GetAccount(ctx, req.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err)) // 404
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err)) // 500
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload) // полезная нагрузка auth
	if account.Owner != authPayload.Username {
		// учетная запись не принадлежит прошедшему проверку подлинности пользователю
		err := errors.New("account doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err)) // 401
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type listAccountRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listAccount(ctx *gin.Context) {
	var req listAccountRequest
	if err := ctx.ShouldBindQuery(&req); err != nil { // Должен быть привязан запрос
		ctx.JSON(http.StatusBadRequest, errorResponse(err)) // 400
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload) // полезная нагрузка auth
	arg := db.ListAccountsParams{
		Owner:  authPayload.Username,
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}
	accounts, err := server.store.ListAccounts(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err)) // 500
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}

func (server *Server) deleteAccount(ctx *gin.Context) {
	var req getAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil { // Должен быть привязан Uri идентификатор
		ctx.JSON(http.StatusBadRequest, errorResponse(err)) // 400
		return
	}

	err := server.store.DeleteAccount(ctx, req.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err)) // 404
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err)) // 500
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
