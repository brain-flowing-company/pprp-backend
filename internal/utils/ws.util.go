package utils

import (
	"encoding/json"
	"time"

	"github.com/brain-flowing-company/pprp-backend/apperror"
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"github.com/gofiber/contrib/websocket"
)

func WebsocketError(conn *websocket.Conn, err interface{}) {
	r := struct {
		Error models.ErrorResponses `json:"error"`
	}{}

	switch apperr := err.(type) {
	case *apperror.AppError:
		r.Error = models.ErrorResponses{
			Code:    apperr.Code(),
			Name:    apperr.Name(),
			Message: apperr.Error(),
		}

	case *apperror.AppErrorType:
		r.Error = models.ErrorResponses{
			Code: apperr.Code,
			Name: apperr.Name,
		}
	}

	reason, _ := json.Marshal(r)
	data := websocket.FormatCloseMessage(websocket.CloseMessage, string(reason))
	conn.WriteControl(websocket.CloseMessage,
		data,
		time.Now().Add(5*time.Second))
}
