package users

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/ArtemVoronov/indefinite-studies-api/internal/api"
	"github.com/ArtemVoronov/indefinite-studies-api/internal/api/validation"
	"github.com/ArtemVoronov/indefinite-studies-api/internal/app/utils"
	"github.com/ArtemVoronov/indefinite-studies-api/internal/db"
	"github.com/ArtemVoronov/indefinite-studies-api/internal/db/entities"
	"github.com/ArtemVoronov/indefinite-studies-api/internal/db/queries"
	"github.com/gin-gonic/gin"
)

type UserDTO struct {
	Id    int
	Login string
	Email string
	Role  string
	State string
}

type UserListDTO struct {
	Count  int
	Offset int
	Limit  int
	Data   []UserDTO
}

type UserEditDTO struct {
	Login    string `json:"login" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role" binding:"required"`
	State    string `json:"state" binding:"required"`
}

type UserCreateDTO struct {
	Login    string `json:"login" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role" binding:"required"`
	State    string `json:"state" binding:"required"`
}

func convertUsers(users []entities.User) []UserDTO {
	if users == nil {
		return make([]UserDTO, 0)
	}
	var result []UserDTO
	for _, user := range users {
		result = append(result, convertUser(user))
	}
	return result
}

func convertUser(user entities.User) UserDTO {
	return UserDTO{Id: user.Id, Login: user.Login, Email: user.Email, Role: user.Role, State: user.State}
}

func GetUsers(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "50")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 50
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}

	db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
		users, err := queries.GetUsers(tx, ctx, limit, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, "Unable to get users")
			log.Printf("Unable to get to users : %s", err)
			return err
		}
		result := &UserListDTO{Data: convertUsers(users), Count: len(users), Offset: offset, Limit: limit}
		c.JSON(http.StatusOK, result)
		return err
	})()
}

func GetUser(c *gin.Context) {
	userIdStr := c.Param("id")

	if userIdStr == "" {
		c.JSON(http.StatusBadRequest, "Missed ID")
		return
	}

	var userId int
	var parseErr error
	if userId, parseErr = strconv.Atoi(userIdStr); parseErr != nil {
		c.JSON(http.StatusBadRequest, api.ERROR_ID_WRONG_FORMAT)
		return
	}

	db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
		user, err := queries.GetUser(tx, ctx, userId)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, api.PAGE_NOT_FOUND)
			} else {
				c.JSON(http.StatusInternalServerError, "Unable to get user")
				log.Printf("Unable to get to user : %s", err)
			}
			return err
		}
		c.JSON(http.StatusOK, convertUser(user))
		return err
	})()
}

func CreateUser(c *gin.Context) {
	var user UserCreateDTO

	if err := c.ShouldBindJSON(&user); err != nil {
		validation.ProcessAndSendValidationErrorMessage(c, err)
		return
	}

	possibleUserRoles := entities.GetPossibleUserRoles()
	if !utils.Contains(possibleUserRoles, user.Role) {
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Unable to create user. Wrong 'Role' value. Possible values: %v", possibleUserRoles))
		return
	}

	possibleUserStates := entities.GetPossibleUserStates()
	if !utils.Contains(possibleUserStates, user.State) {
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Unable to create user. Wrong 'State' value. Possible values: %v", possibleUserStates))
		return
	}

	if user.State == entities.USER_STATE_DELETED {
		c.JSON(http.StatusBadRequest, api.DELETE_VIA_POST_REQUEST_IS_FODBIDDEN)
		return
	}

	db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
		result, err := queries.CreateUser(tx, ctx, user.Login, user.Email, utils.CreateSHA512HashHexEncoded(user.Password), user.Role, user.State)
		if err != nil || result == -1 {
			if err.Error() == db.ErrorUserDuplicateKey.Error() {
				c.JSON(http.StatusBadRequest, api.DUPLICATE_FOUND)
			} else {
				c.JSON(http.StatusInternalServerError, "Unable to create user")
				log.Printf("Unable to create user : %s", err)
			}
			return err

		}
		c.JSON(http.StatusCreated, result)
		return err
	})()
}

// TODO: add optional field updating (field is not reqired and missed -> do not update it)
func UpdateUser(c *gin.Context) {
	userIdStr := c.Param("id")

	if userIdStr == "" {
		c.JSON(http.StatusBadRequest, "Missed ID")
		return
	}

	var userId int
	var parseErr error
	if userId, parseErr = strconv.Atoi(userIdStr); parseErr != nil {
		c.JSON(http.StatusBadRequest, api.ERROR_ID_WRONG_FORMAT)
		return
	}

	var user UserEditDTO

	if err := c.ShouldBindJSON(&user); err != nil {
		validation.ProcessAndSendValidationErrorMessage(c, err)
		return
	}

	if user.State == entities.USER_STATE_DELETED {
		c.JSON(http.StatusBadRequest, api.DELETE_VIA_PUT_REQUEST_IS_FODBIDDEN)
		return
	}

	possibleUserRoles := entities.GetPossibleUserRoles()
	if !utils.Contains(possibleUserRoles, user.Role) {
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Unable to update user. Wrong 'Role' value. Possible values: %v", possibleUserRoles))
		return
	}

	possibleUserStates := entities.GetPossibleUserStates()
	if !utils.Contains(possibleUserStates, user.State) {
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Unable to update user. Wrong 'State' value. Possible values: %v", possibleUserStates))
		return
	}

	// TODO: check password hash
	// TODO: add route for changing password
	// TODO: add route for restoring password
	db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
		err := queries.UpdateUser(tx, ctx, userId, user.Login, user.Email, utils.CreateSHA512HashHexEncoded(user.Password), user.Role, user.State)

		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, api.PAGE_NOT_FOUND)
			} else if err.Error() == db.ErrorUserDuplicateKey.Error() {
				c.JSON(http.StatusBadRequest, api.DUPLICATE_FOUND)
			} else {
				c.JSON(http.StatusInternalServerError, "Unable to update user")
				log.Printf("Unable to update user : %s", err)
			}
			return err
		}
		c.JSON(http.StatusOK, api.DONE)
		return err
	})()
}

func DeleteUser(c *gin.Context) {
	idStr := c.Param("id")

	if idStr == "" {
		c.JSON(http.StatusBadRequest, "Missed ID")
		return
	}

	var id int
	var parseErr error
	if id, parseErr = strconv.Atoi(idStr); parseErr != nil {
		c.JSON(http.StatusBadRequest, api.ERROR_ID_WRONG_FORMAT)
		return
	}

	db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
		err := queries.DeleteUser(tx, ctx, id)

		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, api.PAGE_NOT_FOUND)
			} else {
				c.JSON(http.StatusInternalServerError, "Unable to delete user")
				log.Printf("Unable to delete user: %s", err)
			}
			return err
		}
		c.JSON(http.StatusOK, api.DONE)
		return err
	})()
}
