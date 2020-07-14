package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/apex/gateway"
	"github.com/gin-gonic/gin"
)

var (
	// aes key provided from surveycake
	key string
)

func helloHandler(c *gin.Context) {
	name := c.Param("name")
	c.String(http.StatusOK, "Hello %s", name)
}

func welcomeHandler(c *gin.Context) {
	c.String(http.StatusOK, "Hello World from Go")
}

func rootHandler(c *gin.Context) {
	svid := c.Query("svid")
	hash := c.Query("hash")
	path := fmt.Sprintf("https://www.surveycake.com/webhook/v0/%s/%s", svid, hash)

	body, err := requestToBody(path)
	if err != nil {
		c.String(http.StatusServiceUnavailable, "Failed to send request to surveycake. Error msg: %s", err.Error())
	}
	questionnarie, err := Decrypt(body, []byte(key))
	if err != nil {
		c.String(http.StatusServiceUnavailable, "Failed to decode. Error msg: %s", err.Error())
	}

	c.String(http.StatusOK, "body: %s", questionnarie)
}

func requestToBody(path string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return "", err
	}
	client := &http.Client{Timeout: time.Minute}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	} else if !(200 <= resp.StatusCode && resp.StatusCode < 300) {
		return "", fmt.Errorf("unexpected resp code %d", resp.StatusCode)
	}

	// read body
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func routerEngine() *gin.Engine {
	// set server mode
	gin.SetMode(gin.DebugMode)

	r := gin.New()

	// Global middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/welcome", welcomeHandler)
	r.GET("/user/:name", helloHandler)
	r.GET("/frankie", rootHandler)

	return r
}

func main() {
	log.Fatal(gateway.ListenAndServe("", routerEngine()))
}

/*package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
)

var ginLambda *ginadapter.GinLambda

func init() {
	// stdout and stderr are sent to AWS CloudWatch Logs
	log.Printf("Gin cold start")
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "root",
		})
	})

	ginLambda = ginadapter.New(r)
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	b, _ := json.Marshal(req)
	fmt.Println("----------------------", string(b))
	// If no name is provided in the HTTP request body, throw an error
	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(Handler)
}*/
