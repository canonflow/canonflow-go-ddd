package http

import (
	"net/http"

	"github.com/canonflow/canonflow-go-ddd/internal/domain/user/dto"
	"github.com/canonflow/canonflow-go-ddd/internal/domain/user/model"
	"github.com/canonflow/canonflow-go-ddd/internal/domain/user/usecase"
	"github.com/canonflow/canonflow-go-ddd/pkg/jwt"
	"github.com/canonflow/canonflow-go-ddd/pkg/response"
	"github.com/canonflow/canonflow-go-ddd/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type UserHandler struct {
	UserUsecase usecase.UserUsecase
	Config      *viper.Viper
}

func NewUserHandler(userUsecase usecase.UserUsecase, config *viper.Viper) *UserHandler {
	return &UserHandler{
		UserUsecase: userUsecase,
		Config:      config,
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

func (h *UserHandler) SignIn(ctx *gin.Context) {
	// TODO: Validate
	var request dto.UserLoginRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, response.BaseErrorResponse{
			Code:   http.StatusBadRequest,
			Status: "Bad Request",
			Error:  err.Error(),
		})
		ctx.Abort()
		return
	}

	// TODO: Find the user by username
	user, err := h.UserUsecase.FindByUsername(request.Username)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.BaseErrorResponse{
			Code:   http.StatusBadRequest,
			Status: "Bad Request",
			Error:  "Invalid Credentials",
		})
		ctx.Abort()
		return
	}

	// TODO: Check the password
	if !utils.CheckPassword(request.Password, user.Password) {
		ctx.JSON(http.StatusBadRequest, response.BaseErrorResponse{
			Code:   http.StatusBadRequest,
			Status: "Bad Request",
			Error:  "Invalid Credentials",
		})
		ctx.Abort()
		return
	}

	// TODO: Create access token
	tokenString, err := jwt.CreateToken(user, h.Config.GetString("JWT_SECRET_KEY"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.BaseErrorResponse{
			Code:   http.StatusInternalServerError,
			Status: "Internal Server Error",
			Error:  "Failed to create access token",
		})
		ctx.Abort()
		return
	}

	// TODO: Set HTTP-only Cookie
	ctx.SetSameSite(http.SameSiteNoneMode)
	ctx.SetCookie(jwt.TOKEN_COOKIE, tokenString, int(jwt.TOKEN_DURATION), "/", "", true, true)
	ctx.JSON(http.StatusOK, response.BaseSuccessResponse{
		Code:   http.StatusOK,
		Status: "Sign in successfully",
		Data:   user,
	})
}

func (h *UserHandler) Me(ctx *gin.Context) {
	contextUser, ok := ctx.Get(jwt.USER_KEY)

	if !ok {
		ctx.JSON(http.StatusUnauthorized, response.BaseErrorResponse{
			Code:   http.StatusUnauthorized,
			Status: "Unauthorized",
			Error:  "Unauthorized Access",
		})
		ctx.Abort()
		return
	}

	contextParse := contextUser.(model.User)

	// Get the user
	user, _ := h.UserUsecase.FindById(contextParse.ID)

	ctx.JSON(http.StatusOK, response.BaseSuccessResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   user,
	})
}

func (h *UserHandler) SignOut(ctx *gin.Context) {
	jwt.DeleteToken(ctx)

	ctx.JSON(http.StatusOK, response.BaseSuccessResponse{
		Status: "OK",
		Code:   http.StatusOK,
		Data:   "Logged out successfully",
	})
}
