package email

import "github.com/flosch/pongo2"

func ParseTemplate(template []byte, params map[string]interface{}) (string, error) {
	tpl, err := pongo2.FromString(string(template))
	if err != nil {
		return "", err
	}

	return parseTemplate(tpl, params)
}

func ParseTemplateFile(path string, params map[string]interface{}) (string, error) {
	tpl, err := pongo2.FromFile(path)
	if err != nil {
		return "", err
	}

	return parseTemplate(tpl, params)
}

func parseTemplate(tpl *pongo2.Template, params map[string]interface{}) (string, error) {
	result, err := tpl.Execute(params)
	if err != nil {
		return result, err
	}

	return result, nil
}
