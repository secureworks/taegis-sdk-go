package graphql

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/websocket"

	"github.com/stretchr/testify/require"
)

func NewMockSubServer(t *testing.T, expectedQuery string, expectedVars map[string]interface{}, outputs ...interface{}) *httptest.Server {
	var wsUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := wsUpgrader.Upgrade(w, r, nil)
		require.NoError(t, err)

		msg := &operationMessage{}
		err = conn.ReadJSON(msg)
		require.NoError(t, err)
		err = conn.WriteJSON(&operationMessage{Type: connectionAckMsg})
		require.NoError(t, err)
		err = conn.WriteJSON(&operationMessage{Type: connectionKaMsg})
		require.NoError(t, err)

		err = conn.ReadJSON(msg)
		require.NoError(t, err)
		req := &Request{}
		err = json.Unmarshal(msg.Payload, req)
		require.NoError(t, err)
		require.Equal(t, expectedQuery, req.Query)
		require.Equal(t, expectedVars, req.Variables)

		resp := &Response{}
		result := &operationMessage{Type: dataMsg}

		for _, output := range outputs {
			switch v := output.(type) {
			case error:
				resp.Error = append(resp.Error, Error{Message: v.Error()})
			default:
				outData, err := json.Marshal(output)
				require.NoError(t, err)
				resp.Data = json.RawMessage(outData)
			}

			respData, err := json.Marshal(resp)
			require.NoError(t, err)

			result.Payload = respData

			err = conn.WriteJSON(result)
			require.NoError(t, err)
		}
	}))
}
