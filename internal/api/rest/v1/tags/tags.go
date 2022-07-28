package tags

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

type TagDTO struct {
	Id    int
	Name  string
	State string
}

type TagListDTO struct {
	Count  int
	Offset int
	Limit  int
	Data   []TagDTO
}

type TagEditDTO struct {
	Name  string `json:"name" binding:"required"`
	State string `json:"state" binding:"required"`
}

type TagCreateDTO struct {
	Name  string `json:"name" binding:"required"`
	State string `json:"state" binding:"required"`
}

func convertTags(tags []entities.Tag) []TagDTO {
	if tags == nil {
		return make([]TagDTO, 0)
	}
	var result []TagDTO
	for _, tag := range tags {
		result = append(result, convertTag(tag))
	}
	return result
}

func convertTag(tag entities.Tag) TagDTO {
	return TagDTO{Id: tag.Id, Name: tag.Name, State: tag.State}
}

func GetTags(c *gin.Context) {
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

	data, err := db.Tx(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) (any, error) {
		tags, err := queries.GetTags(tx, ctx, limit, offset)
		return tags, err
	})()

	if err != nil {
		c.JSON(http.StatusInternalServerError, "Unable to get tags")
		log.Printf("Unable to get to tags : %s", err)
		return
	}

	tags, ok := data.([]entities.Tag)
	if !ok {
		c.JSON(http.StatusInternalServerError, "Unable to get tags")
		log.Printf("Unable to get to tags : %s", api.ERROR_ASSERT_RESULT_TYPE)
		return
	}

	result := &TagListDTO{Data: convertTags(tags), Count: len(tags), Offset: offset, Limit: limit}
	c.JSON(http.StatusOK, result)
}

func GetTag(c *gin.Context) {
	tagIdStr := c.Param("id")

	if tagIdStr == "" {
		c.JSON(http.StatusBadRequest, "Missed ID")
		return
	}

	var tagId int
	var parseErr error
	if tagId, parseErr = strconv.Atoi(tagIdStr); parseErr != nil {
		c.JSON(http.StatusBadRequest, api.ERROR_ID_WRONG_FORMAT)
		return
	}

	data, err := db.Tx(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) (any, error) {
		tag, err := queries.GetTag(tx, ctx, tagId)
		return tag, err
	})()

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, api.PAGE_NOT_FOUND)
		} else {
			c.JSON(http.StatusInternalServerError, "Unable to get tag")
			log.Printf("Unable to get to tag : %s", err)
		}
		return
	}

	tag, ok := data.(entities.Tag)
	if !ok {
		c.JSON(http.StatusInternalServerError, "Unable to get tag")
		log.Printf("Unable to get to tag : %s", api.ERROR_ASSERT_RESULT_TYPE)
		return
	}

	c.JSON(http.StatusOK, convertTag(tag))
}

func CreateTag(c *gin.Context) {
	var tag TagCreateDTO

	if err := c.ShouldBindJSON(&tag); err != nil {
		validation.ProcessAndSendValidationErrorMessage(c, err)
		return
	}

	possibleTagStates := entities.GetPossibleTagStates()
	if !utils.Contains(possibleTagStates, tag.State) {
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Unable to create tag. Wrong 'State' value. Possible values: %v", possibleTagStates))
		return
	}

	if tag.State == entities.TAG_STATE_DELETED {
		c.JSON(http.StatusBadRequest, api.DELETE_VIA_POST_REQUEST_IS_FODBIDDEN)
		return
	}

	data, err := db.Tx(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) (any, error) {
		result, err := queries.CreateTag(tx, ctx, tag.Name, tag.State)
		return result, err
	})()

	if err != nil || data == -1 {
		if err.Error() == db.ErrorTagDuplicateKey.Error() {
			c.JSON(http.StatusBadRequest, api.DUPLICATE_FOUND)
		} else {
			c.JSON(http.StatusInternalServerError, "Unable to create tag")
			log.Printf("Unable to create tag : %s", err)
		}
		return

	}
	c.JSON(http.StatusCreated, data)
}

func UpdateTag(c *gin.Context) {
	tagIdStr := c.Param("id")

	if tagIdStr == "" {
		c.JSON(http.StatusBadRequest, "Missed ID")
		return
	}

	var tagId int
	var parseErr error
	if tagId, parseErr = strconv.Atoi(tagIdStr); parseErr != nil {
		c.JSON(http.StatusBadRequest, api.ERROR_ID_WRONG_FORMAT)
		return
	}

	var tag TagEditDTO

	if err := c.ShouldBindJSON(&tag); err != nil {
		validation.ProcessAndSendValidationErrorMessage(c, err)
		return
	}

	if tag.State == entities.TAG_STATE_DELETED {
		c.JSON(http.StatusBadRequest, api.DELETE_VIA_PUT_REQUEST_IS_FODBIDDEN)
		return
	}

	possibleTagStates := entities.GetPossibleTagStates()
	if !utils.Contains(possibleTagStates, tag.State) {
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Unable to update tag. Wrong 'State' value. Possible values: %v", possibleTagStates))
		return
	}

	err := db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
		err := queries.UpdateTag(tx, ctx, tagId, tag.Name, tag.State)
		return err
	})()

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, api.PAGE_NOT_FOUND)
		} else if err.Error() == db.ErrorTagDuplicateKey.Error() {
			c.JSON(http.StatusBadRequest, api.DUPLICATE_FOUND)
		} else {
			c.JSON(http.StatusInternalServerError, "Unable to update tag")
			log.Printf("Unable to update tag : %s", err)
		}
		return
	}

	c.JSON(http.StatusOK, api.DONE)
}

func DeleteTag(c *gin.Context) {
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

	err := db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
		err := queries.DeleteTag(tx, ctx, id)
		return err
	})()

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, api.PAGE_NOT_FOUND)
		} else {
			c.JSON(http.StatusInternalServerError, "Unable to delete tag")
			log.Printf("Unable to delete tag: %s", err)
		}
		return
	}

	c.JSON(http.StatusOK, api.DONE)
}
