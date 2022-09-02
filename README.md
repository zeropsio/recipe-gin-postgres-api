## Building a simple TODO app with gin-gonic in Zerops

Building APIs is bread and butter for most programmers. Zerops makes this process as easy as pie by taking care of the infrastructure.

In this article we will show you how simple it is to build a sample TODO API written in GO using *gin-gonic* - one of the most popular web frameworks.

### Requirements

* [Go](https://www.golang.org) (tested on 1.18)
* [PostgreSQL 12](https://www.postgresql.org)

Both of these requirements can be fullfilled using Zerops managed services.

### Getting started

First and foremost we need a database system to store and manage our data. In this case we chose **PostgreSQL**.

> In case you prefer **MariaDB**, which is also supported by Zerops, the process is analogical and source code can be found [here](https://github.com/zeropsio/recipe-gin-mariadb-api).

Let's start with running a PostgreSQL server. We can create one at
Zerops for development purposes. <!-- TODO: WIKI / print screen -->

We can connect to the db service locally using zcli vpn functionality.

* Generate access token
* zcli login
* zcli vpn start projectName

Now we should be able to test that db is accessible by running `ping6 db.zerops`.

### Golang
#### Dependencies

We will create the api as a GO module, for easier versioning and reproducibility of builds.
New GO modules is created with.

```
go mod init [api-name]
```

This command creates files `go.mod` and `go.sum`, that contain dependency information.

Following GO packages are used in the example:

* [github.com/georgysavva/scany *v1.1.0*](github.com/georgysavva/scany) 
* [github.com/gin-contrib/cors *v1.4.0*](github.com/gin-contrib/cors)
* [github.com/gin-gonic/gin *v1.8.1*](github.com/gin-gonic/gin)
* [github.com/jackc/pgx/v4 *v4.17.1*](github.com/jackc/pgx/v4)

and they can be installed using

```
go get [package-url]
```

More information on how to use *go modules* can be found [here](https://go.dev/blog/using-go-modules).

#### Folder structure
This being a sample application means that the project structure is very simple. 

```
todo-api/
├── http.go
├── main.go
├── model.go
├── go.mod
├── go.sum
├── schema.sql
└── zerops.yml
```

Source code of the api is contained in files

* `http.go` - regarding the http server
* `model.go` - for communication with the DB
* `main.go` - initialization and wiring of dependencies together

This is the boostrap of the Gin framework that is enought for it to run in zerops.
The following is our implementation of http server.

First of all we need to initialize the server by calling 
```go
r := gin.Default()
```

To run the http server smoothly not only in Zerops, we use several middlewares. 
That includes CORS support, better error logging that logs to Zerops runtime log, 
and `content-type` header addition, which is an example of custom written middleware. 

```go
r.Use(cors.Default())
r.Use(func(c *gin.Context) {
    c.Header("content-type", "application/json")
})
r.Use(gin.ErrorLoggerT(gin.ErrorTypePublic | gin.ErrorTypeBind))
r.RedirectTrailingSlash = true
```
Now that we are done with basic server setup, the only thing left is to register endpoints.
First we create a router group that will consist of routes with the same path prefix.
```go
g := r.Group("todos")
```
This api contains CRUD operations for working with the `todo` resource. We register `uri` path to
a handler, which processes the http request.
```go
g.GET("", handler.getTodos)
g.GET("/:id", handler.getTodo)
g.POST("", handler.createTodo)
g.PATCH("/:id", handler.editTodo)
g.DELETE("/:id", handler.deleteTodo)
```

We have chosen `createTodo` handler as an example in this blog post.

```go
func (t todoHandler) createTodo(c *gin.Context) {
	var todo Todo
	err := c.Bind(&todo)
	if err != nil {
		return
	}
	todo, err = t.model.Create(c.Request.Context(), todo)
	if err != nil {
		_ = c.AbortWithError(500, err)
		return
	}
	c.JSON(http.StatusOK, todo)
}
```

Finally, we can run this server on the port `3000` using the following code.

```go
log.Fatal(r.Run(":3000"))
```

#### Running the api locally

In the `main.go` file there are 3 environment variables used to connect and migrate the database.
We can do that by creating .env file with following content.
```env
ZEROPS_RECIPE_DATA_SEED=["foo", "bar"]
ZEROPS_RECIPE_DROP_TABLE=1
DB_URL=postgres://${db_user}:${db_password}@${db_hostname}:5432/${db_hostname}
```

The information about `db_user`, `db_password` and `db_hostname` could be found in
environment variables in zerops. You need to have zcli vpn ready to proceed here.

This command sets environment variables and runs the api.

```sh
$ source .env && go run main.go http.go model.go
```

#### Runnning the api on Zerops

After we completed the development of the api, it's time to deploy it to Zerops. For that
we need to create a file called `zerops.yml`, which contains steps to build and deploy our app.
For the GO language this file is rather simple and looks like this.

```yaml
api:
  build:
    base: [ go@1 ]
    build:
      - go build -o app main.go model.go http.go
    deploy: [ app ]
  run:
    start: ./app
```

You need to add this to the root of your Gitlab / Github repository and link that to the
service stack in Zerops. If you are not sure how to do that [heres link](TODO). Additionally,
configure environment variables as you did in the local development inside Zerops.

After you enabled the subdomain and accessed it, you should see response with todo entries from the
`ZEROPS_RECIPE_DATA_SEED` variable.

#### Conclusion
Hopefully you managed to work along this article to deploy the api to Zerops successfully.
For further questions visit our Discord channel.