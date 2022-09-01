package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
)

const DataSeed = "ZEROPS_RECIPE_DATA_SEED"
const DropTable = "ZEROPS_RECIPE_DROP_TABLE"
const DbUrl = "DB_URL"

func main() {
	ctx := context.Background()

	dbUrl, ok := os.LookupEnv(DbUrl)
	if !ok {
		panic("database url missing set " + DbUrl + " env")
	}

	conn, err := pgxpool.Connect(ctx, dbUrl)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	seeds, err := getSeeds()
	if err != nil {
		panic(err)
	}

	dropTable, err := getDropTable()
	if err != nil {
		panic(err)
	}

	model := TodoRepository{conn}
	err = model.PrepareDatabase(ctx, dropTable, seeds)
	if err != nil {
		panic(err)
	}

	handler := todoHandler{model}

	r := gin.Default()
	r.Use(cors.Default())
	r.Use(func(c *gin.Context) {
		c.Header("content-type", "application/json")
	})
	r.Use(gin.ErrorLoggerT(gin.ErrorTypePublic | gin.ErrorTypeBind))
	r.RedirectTrailingSlash = true

	r.GET("", handler.getTodos)

	g := r.Group("todos")

	g.GET("", handler.getTodos)
	g.GET("/:id", handler.getTodo)
	g.POST("", handler.createTodo)
	g.PATCH("/:id", handler.editTodo)
	g.DELETE("/:id", handler.deleteTodo)

	log.Fatal(r.Run(":3000"))
}

func getSeeds() ([]string, error) {
	dbSeed, ok := os.LookupEnv(DataSeed)
	if !ok {
		dbSeed = "[]"
	}
	var seeds []string
	err := json.Unmarshal([]byte(dbSeed), &seeds)
	return seeds, err
}

func getDropTable() (bool, error) {
	dropTable, ok := os.LookupEnv(DropTable)
	if !ok {
		return false, nil
	}
	return strconv.ParseBool(dropTable)
}
