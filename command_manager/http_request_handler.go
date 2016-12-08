package command_manager

import (
	"fmt"
	"gopkg.in/vmihailenco/msgpack.v2"
	"gopkg.in/gin-gonic/gin.v1"
)

type Response struct {
	Message string `msg:"message"`
	Status int `msg:"status"`
}

func HandleHTTMRequests()  {
	webApp := gin.Default()
	webApp.GET("/run_command/:command_uid", requestRunCommand)
	webApp.GET("/get_info/:command_uid", requestGetCommand)
	webApp.GET("/delete/:command_uid", requestDeleteCommand)
	webApp.Run(":8080")
}

func requestRunCommand(c *gin.Context) {
	c.Header("Content-Type", "application/x-msgpack")

	commandUID := c.Param("command_uid")

	if command := CreateCommand("cat /dev/random", commandUID); command != nil {
		c.String(200, msgResponse("Created", 1))
	} else {
		c.String(201, msgResponse("Something wrong", 1))
	}
}

func requestGetCommand(c *gin.Context) {
	commandUID := c.Param("command_uid")

	result := GetCommandInfo(commandUID)

	c.String(200, result)
}

func requestDeleteCommand(c *gin.Context) {
	commandUID := c.Param("command_uid")

	result := AbortCommand(commandUID)

	c.String(200, result)

}


func msgResponse(messageName string, status int) string {
	response := Response{Message: messageName, Status: status}
	decodedResponse, err := msgpack.Marshal(response)
	PanicError(err)

	return fmt.Sprintf("%s", decodedResponse)
}