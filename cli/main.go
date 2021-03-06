package main

import (
	"github.com/mitchellh/cli"
	"github.com/ngurajeka/go-email"
	"github.com/ngurajeka/go-email/cmd"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"log"
	"os"
)

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	_ = viper.ReadInConfig()
}

var ui cli.Ui

func main() {
	zapConfig := zap.NewProductionConfig()
	zapConfig.OutputPaths = []string{"stdout"}
	zapConfig.ErrorOutputPaths = []string{"stderr"}
	zapConfig.DisableStacktrace = true
	logger, _ := zapConfig.Build()
	defer logger.Sync()

	smtpConfig := email.SMTPConfig{
		Host:     viper.GetString("smtp.host"),
		Username: viper.GetString("smtp.username"),
		Password: viper.GetString("smtp.password"),
		Port:     viper.GetInt("smtp.port"),
	}
	sender := email.Target{
		Name:  viper.GetString("sender.name"),
		Email: viper.GetString("sender.email"),
	}
	emailSvc := email.NewSMTPEmailService(smtpConfig, sender)

	ui = &cli.BasicUi{Writer: os.Stdout}

	commands := &cli.CLI{
		Args: os.Args[1:],
		Commands: map[string]cli.CommandFactory{
			"send": func() (cli.Command, error) {
				return cmd.NewSendCmd(ui, logger, emailSvc), nil
			},
		},
		HelpFunc: cli.BasicHelpFunc("go-email"),
		Version:  "1.0.0",
	}

	exitCode, err := commands.Run()
	if err != nil {
		log.Print("error executing go-email", zap.Strings("args", os.Args), zap.Error(err))
		os.Exit(1)
	}

	os.Exit(exitCode)
}
