package email

import (
	"log"

	"gopkg.in/flosch/pongo2.v3"
)

func ParsingBody(body []byte, params map[string]interface{}) string {

	var (
		tpl    *pongo2.Template
		result string
		err    error
	)

	tpl, err = pongo2.FromString(string(body))
	if err != nil {
		log.Println(err.Error())
		return result
	}

	result, err = tpl.Execute(params)
	if err != nil {
		log.Println(err.Error())
	}

	return result
}
