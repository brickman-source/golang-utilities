/*
 * Copyright (c) 2019. Brickman Source.
 */

package i18n

import (
	"github.com/BurntSushi/toml"
	"github.com/brickman-source/golang-utilities/log"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"os"
	"path/filepath"
	"strings"
)

var (
	defaultLang string
	i18nBundle  *i18n.Bundle
)

func LoadI18Bundle(defaultLang string, dir string) {
	if defaultLang == "" {
		defaultLang = "zh-CN"
	}
	i18nBundle = i18n.NewBundle(language.Make(defaultLang))

	i18nBundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if strings.Contains(path, ".toml") {
			_, err = i18nBundle.LoadMessageFile(path)
			if err != nil {
				log.Errorf( "load i18n file err:%s", err.Error())
			}
		}
		return nil
	})
	if err != nil {
		log.Errorf( "load i18n bundle err: %s", err.Error())
	}
}

func GetI18nString(lang string, messageID string, template map[string]interface{}, pluralCount float32) string {
	if lang == "" {
		lang = defaultLang
	}
	str, err := i18n.NewLocalizer(i18nBundle, lang).Localize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: template,
		PluralCount:  pluralCount,
	})
	if err != nil {
		log.Errorf( "localize(%s) %v, pluralCount: %f err:%s", lang, template, pluralCount, err.Error())
		return str
	}
	return str
}

func GetSimpleI18nString(lang string, messageID string) string {
	if lang == "" {
		lang = defaultLang
	}
	str, err := i18n.NewLocalizer(i18nBundle, lang).Localize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: nil,
		PluralCount:  0,
	})
	if err != nil {
		return str
	}
	return str
}
