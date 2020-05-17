package iotmaker_db_mssql_util

import (
	"errors"
	"regexp"
	"strings"
)

func NameRules(name string) (error, string) {
	var regx *regexp.Regexp
	var err error
	var list []string

	regx, err = regexp.Compile(`(^\d+)(.*)`)
	if err != nil {
		return err, ""
	}
	name = regx.ReplaceAllString(name, "$2")

	regx, err = regexp.Compile(`([^\w\d])`)
	if err != nil {
		return err, ""
	}

	name = regx.ReplaceAllString(name, "_")
	list = strings.Split(name, "_")

	for k := range list {
		list[k] = strings.ToLower(list[k])
		list[k] = strings.Title(list[k])
	}

	name = strings.Join(list, "")
	if name == "" {
		return errors.New("name is empty"), ""
	}

	return nil, name
}
