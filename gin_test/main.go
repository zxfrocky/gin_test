package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"time"
)

func indexHandler(c *gin.Context){
	c.HTML(200,"index.html",nil)
}

func oneLogin(c *gin.Context){
	c.HTML(200,"1/login.html",nil)
}

func twoLogin(c *gin.Context){
	c.HTML(200,"2/login.html",nil)
}

func welcome(c *gin.Context) {
	firstName := c.DefaultQuery("firstname", "Guest")
	lastName := c.Query("lastname")

	message := fmt.Sprintf("Hello %s %s\n", firstName, lastName)

	c.String(http.StatusOK, message)
}

func welcome1(c *gin.Context) {
	firstName := c.DefaultQuery("firstname", "Guest")
	lastName := c.Query("lastname")

	message := fmt.Sprintf("Hello %s %s\n", firstName, lastName)

	c.String(http.StatusOK, message)
}

func root(c *gin.Context) {
	fmt.Println("root")
	remoteAddr := c.Request.RemoteAddr

	val, _ := c.Request.Cookie("test_cookie")
	fmt.Println("cookie is:%s", val)

	message := fmt.Sprintf("your address is:%s\n", remoteAddr)
	c.SetCookie("test_cookie", "haha", 200, "/", "localhost", true, true)

	c.String(http.StatusOK, message)
}

func TlsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		secureMiddleware := secure.New(secure.Options{
			SSLRedirect: true,
			SSLHost:     "localhost:8081",
		})
		err := secureMiddleware.Process(c.Writer, c.Request)

		// If there was an error, do not continue.
		if err != nil {
			return
		}

		c.Next()
	}
}

func router01() http.Handler {
	e := gin.New()
	e.Use(gin.Recovery())
	//e.GET("/welcome", welcome)
	//e.GET("/", root)
	routerGroupInit(e)

	return e
}

func router02() http.Handler {
	e := gin.New()
	e.Use(gin.Recovery())
	routerGroupInit(e)
	return e
}

func routerGroupInit(engine *gin.Engine) {
	//engine.LoadHTMLGlob("/Users/zhouxinfeng/go/src/html/template/static_pages/index.html")
	//engine.LoadHTMLGlob("/Users/zhouxinfeng/go/src/html/template/static_pages")
	//engine.LoadHTMLGlob("/Users/zhouxinfeng/go/src/html/template/static_pages/1/*.html")
	engine.LoadHTMLGlob("views/**/*")
	//engine.LoadHTMLGlob("/Users/zhouxinfeng/go/src/html/template/static_pages1/*")
	engine.RouterGroup.GET("/welcome", welcome)
	engine.RouterGroup.GET("/welcome/1", welcome1)
	engine.RouterGroup.GET("/", indexHandler)
	engine.RouterGroup.GET("/1/login", oneLogin)
	engine.RouterGroup.GET("/2/login", twoLogin)
}

var (
	g errgroup.Group
)

func main() {

	server01 := &http.Server{
		Addr:         ":8000",
		Handler:      router01(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	server02 := &http.Server{
		Addr:         ":8081",
		Handler:      router02(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,}

	g.Go(func() error {
		return server01.ListenAndServe()
	})

	g.Go(func() error {
		return server02.ListenAndServeTLS("/opt/tiger/tools/ca/server.crt", "/opt/tiger/tools/ca/server.key")
	})

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}

	//router := gin.Default()
	//router.Use(TlsHandler())
	//router.GET("/welcome", welcome)
	//router.GET("/", root)

	//router.Run(":8000")
	//router.RunTLS(":8080", "/opt/tiger/tools/ca/server.crt", "/opt/tiger/tools/ca/server.key")
}
