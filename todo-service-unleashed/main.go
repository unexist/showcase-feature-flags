//
// @package Showcase-Feature-Flags
//
// @file Todo main
// @copyright 2023-present Christoph Kappel <christoph@unexist.dev>
// @version $Id$
//
// This program can be distributed under the terms of the Apache License v2.0.
// See the file LICENSE for details.
//

package main

import (
	"net/http"

	"github.com/Unleash/unleash-client-go/v3"
	"github.com/gin-gonic/gin"

	"github.com/unexist/showcase-feature-flags/adapter"
	"github.com/unexist/showcase-feature-flags/domain"
	"github.com/unexist/showcase-feature-flags/infrastructure"

	"log"
)

func init() {
	unleash.Initialize(
		unleash.WithListener(&unleash.DebugListener{}),
		unleash.WithAppName("todo-service-unleash"),
		unleash.WithUrl("http://localhost:8080/todo"),
		unleash.WithCustomHeaders(http.Header{"Authorization": {"<API token>"}}),
	)
}

func main() {
	/* Create business stuff */
	var todoRepository *infrastructure.TodoFakeRepository

	todoRepository = infrastructure.NewTodoFakeRepository()

	defer todoRepository.Close()

	todoService := domain.NewTodoService(todoRepository)
	todoResource := adapter.NewTodoResource(todoService)

	/* Finally start Gin */
	engine := gin.Default()

	todoResource.RegisterRoutes(engine)

	log.Fatal(http.ListenAndServe(":8080", engine))
}
