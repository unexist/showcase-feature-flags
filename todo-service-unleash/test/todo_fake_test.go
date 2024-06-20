//
// @package Showcase-Feature-Flags
//
// @file Todo tests for fake repository
// @copyright 2023-present Christoph Kappel <christoph@unexist.dev>
// @version $Id$
//
// This program can be distributed under the terms of the Apache License v2.0.
// See the file LICENSE for details.
//

package test

import (
	"github.com/Unleash/unleash-client-go/v3"
	unleashContext "github.com/Unleash/unleash-client-go/v3/context"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"time"

	"os"
	"testing"

	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"

	"github.com/unexist/showcase-feature-flags/adapter"
	"github.com/unexist/showcase-feature-flags/domain"
	"github.com/unexist/showcase-feature-flags/infrastructure"
)

/* Test globals */
var engine *gin.Engine
var todoRepository *infrastructure.TodoFakeRepository

func init() {
	unleash.Initialize(
		unleash.WithListener(&unleash.DebugListener{}),
		unleash.WithAppName("todo-service-unleash"),
		unleash.WithUrl(os.Getenv("API_URL")),
		unleash.WithCustomHeaders(http.Header{"Authorization": {os.Getenv("API_TOKEN")}}),
		unleash.WithRefreshInterval(1*time.Second),
		unleash.WithMetricsInterval(1*time.Second),
	)
}

func TestMain(m *testing.M) {
	/* Create business stuff */
	todoRepository = infrastructure.NewTodoFakeRepository()
	todoService := domain.NewTodoService(todoRepository)
	todoResource := adapter.NewTodoResource(todoService)

	/* Finally start Gin */
	engine = gin.Default()

	todoResource.RegisterRoutes(engine)

	code := m.Run()

	os.Exit(code)
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()

	engine.ServeHTTP(recorder, req)

	return recorder
}

func checkResponseCode(t *testing.T, expected, actual int) {
	assert.Equal(t, expected, actual, "Expected different response code")
}

func TestEmptyTable(t *testing.T) {
	todoRepository.Clear()

	req, _ := http.NewRequest("GET", "/todo", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); "[]" != body {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func TestGetNonExistentTodo(t *testing.T) {
	todoRepository.Clear()

	req, _ := http.NewRequest("GET", "/todo/11", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)

	assert.Equal(t, "Todo not found", m["error"],
		"Expected the 'error' key of the response to be set to 'Todo not found'")
}

func TestCreateTodo(t *testing.T) {
	todoRepository.Clear()

	var jsonStr = []byte(`{"title":"string", "description": "string"}`)

	req, _ := http.NewRequest("POST", "/todo", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	assert.Equal(t, 1.0, m["id"], "Expected todo ID to be '1'")
	assert.Equal(t, "string", m["title"], "Expected todo title to be 'string'")
	assert.Equal(t, "string", m["description"], "Expected todo description to be 'string'")
}

func TestCreateTodoBadword(t *testing.T) {
	todoRepository.Clear()

	ctx := unleashContext.Context{
		UserId:        "1",
		SessionId:     "some-session-id",
		RemoteAddress: "test",
	}

	if unleash.IsEnabled("feat.CheckBadwords", unleash.WithContext(ctx)) {
		var jsonStr = []byte(`{"title":"string crap", "description": "string"}`)

		req, _ := http.NewRequest("POST", "/todo", bytes.NewBuffer(jsonStr))
		req.Header.Set("Content-Type", "application/json")

		response := executeRequest(req)
		checkResponseCode(t, http.StatusExpectationFailed, response.Code)

		var m map[string]interface{}
		json.Unmarshal(response.Body.Bytes(), &m)

		assert.Equal(t, "Title contains badword", m["error"],
			"Title contains badword")
	}
}

func TestGetTodo(t *testing.T) {
	todoRepository.Clear()
	addTodos(1)

	req, _ := http.NewRequest("GET", "/todo/1", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func addTodos(count int) {
	if 1 > count {
		count = 1
	}

	todo := domain.Todo{}

	for i := 0; i < count; i++ {
		todo.ID = i
		todo.Title = "Todo " + strconv.Itoa(i)
		todo.Description = "string"

		todoRepository.CreateTodo(&todo)
	}
}

func TestUpdateTodo(t *testing.T) {
	todoRepository.Clear()
	addTodos(1)

	req, _ := http.NewRequest("GET", "/todo/1", nil)
	response := executeRequest(req)

	var origTodo map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &origTodo)

	var jsonStr = []byte(`{"title":"new string", "description": "new string"}`)

	req, _ = http.NewRequest("PUT", "/todo/1", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	assert.Equal(t, origTodo["id"], m["id"], "Expected the id to remain the same")
	assert.NotEqual(t, origTodo["title"], m["title"], "Expected the title to change")
	assert.NotEqual(t, origTodo["description"], m["description"], "Expected the description to change")
}

func TestDeleteTodo(t *testing.T) {
	todoRepository.Clear()
	addTodos(1)

	req, _ := http.NewRequest("GET", "/todo/1", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/todo/1", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNoContent, response.Code)

	req, _ = http.NewRequest("GET", "/todo/1", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}
