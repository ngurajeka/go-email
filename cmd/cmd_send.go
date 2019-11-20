package cmd

import (
	"encoding/csv"
	"github.com/mitchellh/cli"
	"github.com/ngurajeka/go-email"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"io/ioutil"
	"net/mail"
	"os"
	"strings"
)

// SendCmd command to run the cli version of go-email
type SendCmd struct {
	ui     cli.Ui
	logger *zap.Logger
}

// NewSendCmd create new SendCmd command
func NewSendCmd(ui cli.Ui, logger *zap.Logger) *SendCmd {
	return &SendCmd{ui, logger}
}

// Help return help text
func (c *SendCmd) Help() string {
	helpText := `
			Usage: go-email send [options] -f target.csv
			Sending email to targets.
	`
	return strings.TrimSpace(helpText)
}

// Synopsis return synopsis text
func (c *SendCmd) Synopsis() string {
	return "Running the sending email process"
}

// Run running the sending email process
func (c *SendCmd) Run(args []string) int {
	subject, ok := getKey(args, "subject")
	if !ok {
		c.logger.Error("invalid subject")
		return 1
	}
	targetPath, ok := getKey(args, "f")
	if !ok {
		c.logger.Error("invalid target file")
		return 1
	}
	htmlTemplatePath, ok := getKey(args, "html")
	if !ok {
		c.logger.Error("invalid html template")
		return 1
	}

	account := email.NewAccount(viper.GetString("email.sender_name"), viper.GetString("email.sender_email"))
	account.SetCredential(
		viper.GetString("email.host"),
		viper.GetString("email.username"),
		viper.GetString("email.password"),
		viper.GetInt("email.port"),
	)

	htmlTemplate, err := ioutil.ReadFile(htmlTemplatePath)
	if err != nil {
		c.logger.Error("reading html template failed", zap.String("path", htmlTemplatePath), zap.Error(err))
		return 1
	}

	f, err := os.Open(targetPath)
	if err != nil {
		c.logger.Error("reading target file failed", zap.Error(err))
		return 1
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	rows, err := csvReader.ReadAll()
	if err != nil {
		c.logger.Error("reading target file failed", zap.Error(err))
		return 1
	}

	for _, row := range rows {
		message := email.Default()
		message.SetFrom(viper.GetString("email.sender_name"), viper.GetString("email.sender_email"))
		message.SetSubject(subject)
		message.SetTo(mail.Address{Name: row[1], Address: row[0]})
		message.SetHTMLBody(htmlTemplate)
		message.AddParam("fullName", row[1])
		if err := message.Send(account); err != nil {
			c.logger.Error("sending email failed", zap.Errors("errors", err))
			continue
		} else {
			c.logger.Info("sending email succeed", zap.String("email", row[1]))
		}
	}

	return 0
}

func getKey(args []string, key string) (string, bool) {
	var (
		v     string
		exist bool
	)

	for _, value := range args {
		s := strings.Split(value, "=")
		if len(s) != 2 {
			continue
		}
		if s[0] == "-"+key {
			exist = true
			v = s[1]
		}
	}

	return v, exist
}
