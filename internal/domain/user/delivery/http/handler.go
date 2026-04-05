package http

import (
	"context"
	"log"
	"net/http"

	"github.com/canonflow/canonflow-go-ddd/internal/domain/user/dto"
	"github.com/canonflow/canonflow-go-ddd/internal/domain/user/model"
	userQueue "github.com/canonflow/canonflow-go-ddd/internal/domain/user/queue"
	"github.com/canonflow/canonflow-go-ddd/internal/domain/user/usecase"
	"github.com/canonflow/canonflow-go-ddd/pkg/jwt"
	queuePkg "github.com/canonflow/canonflow-go-ddd/pkg/queue"
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

// SignUp godoc
// @Summary      Create an account
// @Description  create a new account with the provided information
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      dto.CreateUserRequest  true  "User information"
// @Success      200  {object}  response.BaseSuccessResponse
// @Failure      400  {object}  response.BaseErrorResponse
// @Failure      404  {object}  response.BaseErrorResponse
// @Failure      500  {object}  response.BaseErrorResponse
// @Router       /api/v1/auth/signup [post]
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

// SignIn godoc
// @Summary      Sign in to an existing account
// @Description  sign in to an existing account with the provided information
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      dto.UserLoginRequest  true  "User credentials"
// @Success      200  {object}  response.BaseSuccessResponse
// @Failure      400  {object}  response.BaseErrorResponse
// @Failure      404  {object}  response.BaseErrorResponse
// @Failure      500  {object}  response.BaseErrorResponse
// @Router       /api/v1/auth/signin [post]
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

// Me godoc
// @Summary      Get user information
// @Description  get user information for the authenticated user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.BaseSuccessResponse
// @Failure      400  {object}  response.BaseErrorResponse
// @Failure      404  {object}  response.BaseErrorResponse
// @Failure      500  {object}  response.BaseErrorResponse
// @Router       /api/v1/auth/me [get]
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

	//* Test Queue
	queuePayload := map[string]interface{}{
		"user_id":    user.ID,
		"username":   user.Username,
		"created_at": user.CreatedAt,
	}
	log.Printf("[User Handler - Me] Queue Payload: %v", queuePayload)

	err := queuePkg.Dispatch(
		context.Background(),
		userQueue.NewUserQueue(userQueue.QUEUE_NAME),
		queuePkg.QueueMessage{
			UniqueID: utils.GenerateUUID(),
			Payload:  queuePayload,
		},
	)
	if err != nil {
		log.Printf("Failed to dispatch message to queue: %s", err)
	}

	ctx.JSON(http.StatusOK, response.BaseSuccessResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   user,
	})
}

// SignOut godoc
// @Summary      Sign out of the current account
// @Description  sign out of the current account and invalidate the access token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.BaseSuccessResponse
// @Failure      400  {object}  response.BaseErrorResponse
// @Failure      404  {object}  response.BaseErrorResponse
// @Failure      500  {object}  response.BaseErrorResponse
// @Router       /api/v1/auth/signout [post]
func (h *UserHandler) SignOut(ctx *gin.Context) {
	jwt.DeleteToken(ctx)

	ctx.JSON(http.StatusOK, response.BaseSuccessResponse{
		Status: "OK",
		Code:   http.StatusOK,
		Data:   "Logged out successfully",
	})
}
