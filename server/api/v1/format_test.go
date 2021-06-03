package v1

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/txfs19260817/scopelens/server/utils/response"
	"github.com/txfs19260817/scopelens/server/utils/testsuite"
)

func TestGetFormats(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	router.GET("/format", GetFormats)

	w := testsuite.PerformRequest(router, "GET", "/format")
	assert.Equal(t, http.StatusOK, w.Code)

	var r response.Response
	err := json.Unmarshal([]byte(w.Body.String()), &r)
	assert.NoError(t, err)
	assert.Equal(t, response.SUCCESS, r.Code)
	assert.NotEmpty(t, r.Data)
}
