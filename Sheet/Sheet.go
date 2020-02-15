package Sheet

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/semyon-dev/RusSeated/Config"
	"github.com/semyon-dev/RusSeated/Model"
	"golang.org/x/oauth2/google"
	"gopkg.in/Iwark/spreadsheet.v2"
	"io/ioutil"
)

func Connect(){
	data, err := ioutil.ReadFile("credentials.json")
	checkError(err)

	conf, err := google.JWTConfigFromJSON(data, spreadsheet.Scope)
	checkError(err)

	client := conf.Client(context.TODO())
	service := spreadsheet.NewServiceWithClient(client)

	spreadsheetID := Config.SpreadsheetID
	sheet, err := service.FetchSpreadsheet(spreadsheetID)
	checkError(err)

	fmt.Println(sheet.ID)
}

func NewForm(c *gin.Context){
	var form Model.Form
	err := c.ShouldBind(&form)
	if err != nil{
		fmt.Println(err.Error())
	} else {
		// TODO: добавляем форму в конец таблицы
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err.Error())
	}
}
