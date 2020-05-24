package cmd

import (
	"encoding/csv"
	"github.com/mitchellh/cli"
	"github.com/ngurajeka/go-email"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"
)

// SendCmd command to run the cli version of go-email
type SendCmd struct {
	ui       cli.Ui
	logger   *zap.Logger
	emailSvc email.Service
}

// NewSendCmd create new SendCmd command
func NewSendCmd(ui cli.Ui, logger *zap.Logger, emailSvc email.Service) *SendCmd {
	return &SendCmd{ui: ui, logger: logger, emailSvc: emailSvc}
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
	start := time.Now()
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
	htmlTemplate, err := ioutil.ReadFile(htmlTemplatePath)
	if err != nil {
		c.logger.Error("file not found")
		return 1
	}
	targetFile, err := os.Open(targetPath)
	if err != nil {
		c.logger.Error("reading target file failed", zap.Error(err))
		return 1
	}
	defer targetFile.Close()

	csvReader := csv.NewReader(targetFile)
	rows, err := csvReader.ReadAll()
	if err != nil {
		c.logger.Error("reading target file failed", zap.Error(err))
		return 1
	}

	headers := rows[0]

	var wg sync.WaitGroup

	for _, row := range rows[1:] {
		wg.Add(1)

		var params = make(map[string]interface{})
		for i, column := range headers {
			params[column] = row[i]
		}

		target := email.Target{
			Name:  row[1],
			Email: row[0],
		}

		htmlMessage, err := email.ParseTemplate(htmlTemplate, params)
		if err != nil {
			c.logger.Error("parsing html message failed", zap.Error(err))
			continue
		}

		message := email.Message{HTML: htmlMessage}

		go func(target email.Target, message email.Message) {
			defer wg.Done()
			if _, err := c.emailSvc.Send(subject, target, nil, nil, message, nil); err != nil {
				c.logger.Error("sending email failed", zap.String("target", target.String()), zap.Error(err))
			} else {
				c.logger.Info("sending email succeed", zap.String("target", target.String()))
			}
		}(target, message)
	}

	wg.Wait()

	end := time.Now()

	c.logger.Info("task has been executed", zap.Time("start", start), zap.Time("end", end))

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
