package notes

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

type NoteDTO struct {
	Id     int
	Text   string
	Topic  string
	TagId  int
	UserId int
	State  string
}

type NoteListDTO struct {
	Count  int
	Offset int
	Limit  int
	Data   []NoteDTO
}

type NoteEditDTO struct {
	Text   string `json:"text" binding:"required"`
	Topic  string `json:"topic" binding:"required"`
	TagId  int    `json:"tagId" binding:"required"`
	UserId int    `json:"userId" binding:"required"`
	State  string `json:"state" binding:"required"`
}

type NoteCreateDTO struct {
	Text   string `json:"text" binding:"required"`
	Topic  string `json:"topic" binding:"required"`
	TagId  int    `json:"tagId" binding:"required"`
	UserId int    `json:"userId" binding:"required"`
	State  string `json:"state" binding:"required"`
}

func convertNotes(notes []entities.Note) []NoteDTO {
	if notes == nil {
		return make([]NoteDTO, 0)
	}
	var result []NoteDTO
	for _, note := range notes {
		result = append(result, convertNote(note))
	}
	return result
}

func convertNote(note entities.Note) NoteDTO {
	return NoteDTO{Id: note.Id, Text: note.Text, Topic: note.Topic, TagId: note.TagId, UserId: note.UserId, State: note.State}
}

func GetNotes(c *gin.Context) {
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
		notes, err := queries.GetNotes(tx, ctx, limit, offset)
		return notes, err
	})()

	if err != nil {
		c.JSON(http.StatusInternalServerError, "Unable to get notes")
		log.Printf("Unable to get to notes : %s", err)
		return
	}

	notes, ok := data.([]entities.Note)
	if !ok {
		c.JSON(http.StatusInternalServerError, "Unable to get notes")
		log.Printf("Unable to get to notes : %s", api.ERROR_ASSERT_RESULT_TYPE)
		return
	}

	result := &NoteListDTO{Data: convertNotes(notes), Count: len(notes), Offset: offset, Limit: limit}
	c.JSON(http.StatusOK, result)
}

func GetNote(c *gin.Context) {
	noteIdStr := c.Param("id")

	if noteIdStr == "" {
		c.JSON(http.StatusBadRequest, "Missed ID")
		return
	}

	var noteId int
	var parseErr error
	if noteId, parseErr = strconv.Atoi(noteIdStr); parseErr != nil {
		c.JSON(http.StatusBadRequest, api.ERROR_ID_WRONG_FORMAT)
		return
	}

	data, err := db.Tx(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) (any, error) {
		note, err := queries.GetNote(tx, ctx, noteId)
		return note, err
	})()

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, api.PAGE_NOT_FOUND)
		} else {
			c.JSON(http.StatusInternalServerError, "Unable to get note")
			log.Printf("Unable to get to note : %s", err)
		}
		return
	}

	note, ok := data.(entities.Note)
	if !ok {
		c.JSON(http.StatusInternalServerError, "Unable to get note")
		log.Printf("Unable to get to note : %s", api.ERROR_ASSERT_RESULT_TYPE)
		return
	}

	c.JSON(http.StatusOK, convertNote(note))
}

func CreateNote(c *gin.Context) {
	var note NoteCreateDTO

	if err := c.ShouldBindJSON(&note); err != nil {
		validation.ProcessAndSendValidationErrorMessage(c, err)
		return
	}

	possibleNoteStates := entities.GetPossibleNoteStates()
	if !utils.Contains(possibleNoteStates, note.State) {
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Unable to create note. Wrong 'State' value. Possible values: %v", possibleNoteStates))
		return
	}

	if note.State == entities.NOTE_STATE_DELETED {
		c.JSON(http.StatusBadRequest, api.DELETE_VIA_POST_REQUEST_IS_FODBIDDEN)
		return
	}

	data, err := db.Tx(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) (any, error) {
		result, err := queries.CreateNote(tx, ctx, note.Text, note.Topic, note.TagId, note.UserId, note.State)
		return result, err
	})()

	if err != nil || data == -1 {
		c.JSON(http.StatusInternalServerError, "Unable to create note")
		log.Printf("Unable to create note : %s", err)
		return
	}

	c.JSON(http.StatusCreated, data)
}

// TODO: add optional field updating (field is not reqired and missed -> do not update it)
func UpdateNote(c *gin.Context) {
	noteIdStr := c.Param("id")

	if noteIdStr == "" {
		c.JSON(http.StatusBadRequest, "Missed ID")
		return
	}

	var noteId int
	var parseErr error
	if noteId, parseErr = strconv.Atoi(noteIdStr); parseErr != nil {
		c.JSON(http.StatusBadRequest, api.ERROR_ID_WRONG_FORMAT)
		return
	}

	var note NoteEditDTO

	if err := c.ShouldBindJSON(&note); err != nil {
		validation.ProcessAndSendValidationErrorMessage(c, err)
		return
	}

	if note.State == entities.TASK_STATE_DELETED {
		c.JSON(http.StatusBadRequest, api.DELETE_VIA_PUT_REQUEST_IS_FODBIDDEN)
		return
	}

	possibleNoteStates := entities.GetPossibleNoteStates()
	if !utils.Contains(possibleNoteStates, note.State) {
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Unable to update note. Wrong 'State' value. Possible values: %v", possibleNoteStates))
		return
	}

	err := db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
		err := queries.UpdateNote(tx, ctx, noteId, note.Text, note.Topic, note.TagId, note.UserId, note.State)
		return err
	})()

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, api.PAGE_NOT_FOUND)
		} else {
			c.JSON(http.StatusInternalServerError, "Unable to update note")
			log.Printf("Unable to update note : %s", err)
		}
		return
	}

	c.JSON(http.StatusOK, api.DONE)
}

func DeleteNote(c *gin.Context) {
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
		err := queries.DeleteNote(tx, ctx, id)
		return err
	})()

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, api.PAGE_NOT_FOUND)
		} else {
			c.JSON(http.StatusInternalServerError, "Unable to delete note")
			log.Printf("Unable to delete note: %s", err)
		}
		return
	}

	c.JSON(http.StatusOK, api.DONE)
}
