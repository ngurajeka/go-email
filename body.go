package email

import (
	"gopkg.in/flosch/pongo2.v3"
)

func ParsingBody(body []byte, params map[string]interface{}) (string, error) {
	var (
		tpl    *pongo2.Template
		result string
		err    error
	)

	tpl, err = pongo2.FromString(string(body))
	if err != nil {
		return result, err
	}

	result, err = tpl.Execute(params)
	if err != nil {
		return result, err
	}

	return result, nil
}
