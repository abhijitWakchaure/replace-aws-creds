package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"regexp"
	"time"

	"github.com/abhijitWakchaure/replace-aws-creds/logger"
	"github.com/atotto/clipboard"
	"github.com/spf13/viper"
)

// App ...
type App struct {
	Logger          logger.Logger
	Daemon, Verbose bool
}

// const patternAWSCreds = "\\[(?P<sectionName>\\d*_[a-z]{3,4})\\]\naws_access_key_id=(?P<accessKeyID>[A-Z0-9]*)\naws_secret_access_key=(?P<secretAccessKey>[A-Za-z0-9+\\/]*)\naws_session_token=(?P<sessionToken>[\\w\\/+=]*)"
const patternAWSCreds = `\[(?P<sectionName>\d*_[a-z]{3,4})\]\naws_access_key_id=(?P<accessKeyID>[A-Z0-9]*)\naws_secret_access_key=(?P<secretAccessKey>[A-Za-z0-9+\/]*)\naws_session_token=(?P<sessionToken>[\w\/+=]*)`

// ExtractAWSCreds ...
func (app *App) ExtractAWSCreds() {
	r, err := regexp.Compile(patternAWSCreds)
	if err != nil {
		app.Logger.Fatal("failed to parse regex due to %s", err)
	}

	c := app.getClipboardContent()
	for {
		if c == app.getClipboardContent() {
			app.Logger.Debug("clipboard contents unchanged, waiting...")
			time.Sleep(2 * time.Second)
			continue
		}
		c = app.getClipboardContent()

		if !r.MatchString(c) {
			app.Logger.Info("clipboard contents don't match the aws creds pattern!")
			time.Sleep(5 * time.Second)
			continue
		}
		a := r.FindAllStringSubmatch(c, -1)
		sectionName := a[0][1]
		accessKeyID := a[0][2]
		secretAccessKey := a[0][3]
		sessionToken := a[0][4]
		creds := make(map[string]interface{})
		creds[fmt.Sprintf("%s.%s", sectionName, "aws_access_key_id")] = accessKeyID
		creds[fmt.Sprintf("%s.%s", sectionName, "aws_secret_access_key")] = secretAccessKey
		creds[fmt.Sprintf("%s.%s", sectionName, "aws_session_token")] = sessionToken
		app.writeAWSCreds(creds)

		if !app.Daemon {
			break
		}
	}

}

func (app *App) getClipboardContent() string {
	c, err := clipboard.ReadAll()
	if err != nil {
		app.Logger.Fatal("failed to read clipboard contents", err)
	}
	return c
}

func (app *App) writeAWSCreds(m map[string]interface{}) {
	viper.MergeConfigMap(m)
	b, _ := json.MarshalIndent(viper.AllSettings(), "", " ")
	app.Logger.Debugf("updated AWS creds file: %s", string(b))
	viper.WriteConfig()
	defer app.Logger.Info("AWS creds updated!")
	d := viper.Get("default")
	if d == nil {
		return
	}
	app.Logger.Debug("default section detected, rewriting old default section")
	b, err := ioutil.ReadFile(viper.ConfigFileUsed())
	if err != nil {
		app.Logger.Fatal("failed to read AWS credentials file", err)
	}
	s := "[default]\n"
	narr := append([]byte(s), b...)
	err = ioutil.WriteFile(viper.ConfigFileUsed(), narr, 0500)
	if err != nil {
		app.Logger.Fatal("failed to write AWS credentials file", err)
	}
}
