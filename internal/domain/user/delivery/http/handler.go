package http

import (
	"net/http"

	"github.com/canonflow/canonflow-go-ddd/internal/domain/user/dto"
	"github.com/canonflow/canonflow-go-ddd/internal/domain/user/usecase"
	"github.com/canonflow/canonflow-go-ddd/pkg/response"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	UserUsecase usecase.UserUsecase
}

func NewUserHandler(userUsecase usecase.UserUsecase) *UserHandler {
	return &UserHandler{
		UserUsecase: userUsecase,
	}
}

func (h *UserHandler) SignUp(ctx *gin.Context) {
	// TODO: Validate
	var request dto.CreateUserRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, response.BaseErrorResponse{
			Code:   http.StatusBadRequest,
			Status: "Bad Request",
			Error:  err.Error(),
		})
		ctx.Abort()
		return
	}

	// TODO: Find the username
	_, err := h.UserUsecase.FindByUsername(request.Username)

	if err == nil {
		ctx.JSON(http.StatusBadRequest, response.BaseErrorResponse{
			Code:   http.StatusBadRequest,
			Status: "Bad Request",
			Error:  "Username already taken",
		})
		ctx.Abort()
		return
	}

	// TODO: Create user
	user, err := h.UserUsecase.Create(ctx, request.Username, request.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.BaseErrorResponse{
			Code:   http.StatusInternalServerError,
			Status: "Internal Server Error",
			Error:  "Failed to create new user",
		})
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, response.BaseSuccessResponse{
		Code:   http.StatusOK,
		Status: "New user created successfully",
		Data:   user,
	})
}
