package utils

import (
	"encoding/json"
	"time"

	"github.com/brain-flowing-company/pprp-backend/apperror"
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"github.com/gofiber/contrib/websocket"
)

func parseError(err interface{}) models.ErrorResponses {
	var r models.ErrorResponses

	switch apperr := err.(type) {
	case *apperror.AppError:
		r = models.ErrorResponses{
			Code:    apperr.Code(),
			Name:    apperr.Name(),
			Message: apperr.Error(),
		}

	case *apperror.AppErrorType:
		r = models.ErrorResponses{
			Code: apperr.Code,
			Name: apperr.Name,
		}
	}

	return r
}

func WebsocketError(conn *websocket.Conn, err interface{}) error {
	r := struct{ Error models.ErrorResponses }{
		Error: parseError(err),
	}

	return conn.WriteJSON(r)
}

func WebsocketFatal(conn *websocket.Conn, err interface{}) error {
	r := struct{ Error models.ErrorResponses }{
		Error: parseError(err),
	}

	reason, _ := json.Marshal(r)
	data := websocket.FormatCloseMessage(websocket.CloseMessage, string(reason))
	return conn.WriteControl(websocket.CloseMessage, data, time.Now().Add(5*time.Second))
}
