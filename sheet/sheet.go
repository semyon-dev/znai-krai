// Работа с Google Sheet
package sheet

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/semyon-dev/znai-krai/config"
	"github.com/semyon-dev/znai-krai/log"
	"github.com/semyon-dev/znai-krai/model"
	"golang.org/x/oauth2/google"
	"gopkg.in/Iwark/spreadsheet.v2"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"time"
)

var Service *spreadsheet.Service

// Connect to Google Sheets
func Connect() {
	var data []byte
	var err error
	// если нет переменных окружения значит читаем файл credentials.json
	if len(config.Credentials) == 0 {
		data, err = ioutil.ReadFile("credentials.json")
		fmt.Println("Чтение файла с учетными данными (credentials)")
		if err != nil {
			log.HandleErrWithMsg("Критическая ошибка: не удалось импортировать переменные:", err)
		}
	} else {
		var f model.CredentialsFile
		fmt.Println("read credentials from env var")
		f.Type = "service_account"
		f.ProjectID = "zekovnet"
		f.PrivateKeyID = os.Getenv("private_key_id")
		f.PrivateKey = strings.ReplaceAll(os.Getenv("private_key"), "\\n", "\n")
		f.ClientEmail = os.Getenv("client_email")
		f.ClientID = os.Getenv("client_id")
		f.TokenURL = os.Getenv("token_uri")
		data, err = json.Marshal(f)
		checkError(err)
	}

	conf, err := google.JWTConfigFromJSON(data, spreadsheet.Scope)
	checkError(err)

	client := conf.Client(context.TODO())
	Service = spreadsheet.NewServiceWithClient(client)

	formSpreadsheet, err := Service.FetchSpreadsheet(config.SpreadsheetIDForms)
	checkError(err)
	spreadsheetCoronavirus, err := Service.FetchSpreadsheet(config.SpreadsheetCoronavirus)
	checkError(err)
	fsinPlacesSpreadsheet, err := Service.FetchSpreadsheet(config.SpreadsheetIDFsinPlaces)
	checkError(err)
	emailSpreadsheet, err := Service.FetchSpreadsheet(config.SpreadsheetEmail)
	checkError(err)
	fmt.Println("таблица нарушений (форм):", formSpreadsheet.Properties.Title)
	fmt.Println("таблица ФСИН учреждений:", fsinPlacesSpreadsheet.Properties.Title)
	fmt.Println("таблица с информацией по Коронавирусу:", spreadsheetCoronavirus.Properties.Title)
	fmt.Println("таблица с информацией email рассылки:", emailSpreadsheet.Properties.Title)
}

func AddCoronaViolation(coronaForm model.CoronaViolation) error {
	formSpreadsheet, err := Service.FetchSpreadsheet(config.SpreadsheetCoronavirus)
	checkError(err)

	formsSheet, err := formSpreadsheet.SheetByID(0)
	checkError(err)

	row := len(formsSheet.Rows)
	currentTime := time.Now()

	formsSheet.Update(row, 1, currentTime.Format("2006.01.02 15:04:05"))
	formsSheet.Update(row, 2, coronaForm.NameOfFSIN)
	formsSheet.Update(row, 3, coronaForm.Region)
	formsSheet.Update(row, 4, coronaForm.Info)
	formsSheet.Update(row, 5, coronaForm.CommentFSIN)
	formsSheet.Update(row, 6, coronaForm.Source) // источник данных
	formsSheet.Update(row, 7, coronaForm.PlaceIDString)
	err = formsSheet.Synchronize()
	checkError(err)
	return err
}

func AddViolation(form model.Violation) error {
	// We need a pointer so that we can set the value via reflection
	msValuePtr := reflect.ValueOf(&form)
	msValue := msValuePtr.Elem()

	// нужно для синхронизации
	formSpreadsheet, err := Service.FetchSpreadsheet(config.SpreadsheetIDForms)
	checkError(err)

	formsSheet, err := formSpreadsheet.SheetByID(0)
	checkError(err)

	row := len(formsSheet.Rows)
	column := 0
	currentTime := time.Now()

	formsSheet.Update(row, column, currentTime.Format("2006.01.02 15:04:05"))
	column++

	for ; column < msValue.NumField(); column++ {
		field := msValue.Field(column)

		// Ignore fields that don't have the same type as a string
		if field.Type() != reflect.TypeOf("") || "add_files" == reflect.ValueOf(field).Type().Name() {
			continue
		}

		str := field.Interface().(string)
		str = strings.TrimSpace(str)
		field.SetString(str)

		// добавляем в таблицу
		formsSheet.Update(row, column, field.String())
	}
	err = formsSheet.Synchronize()
	checkError(err)
	return err
}

// Добавление пользователя в таблицу для email рассылки
func AddMailing(formMailing model.Mailing) error {
	formSpreadsheet, err := Service.FetchSpreadsheet(config.SpreadsheetEmail)
	checkError(err)
	formsSheet, err := formSpreadsheet.SheetByID(0)
	checkError(err)
	row := len(formsSheet.Rows)
	formsSheet.Update(row, 0, formMailing.Name)
	formsSheet.Update(row, 1, formMailing.Email)
	err = formsSheet.Synchronize()
	checkError(err)
	return err
}

func checkError(err error) {
	if err != nil {
		log.HandleErr(err)
	}
}
