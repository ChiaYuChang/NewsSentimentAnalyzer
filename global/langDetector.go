package global

import (
	"sync"

	"github.com/pemistahl/lingua-go"
)

var languageDetector detectorSingleton

type detectorSingleton struct {
	lingua.LanguageDetector
	sync.Once
}

func LanguageDetector() lingua.LanguageDetector {
	languageDetector.Do(func() {
		languageDetector.LanguageDetector = lingua.NewLanguageDetectorBuilder().
			FromLanguages(lingua.AllLanguages()...).
			Build()
	})
	return languageDetector.LanguageDetector
}

func SetLanguageDetector(lang ...lingua.Language) lingua.LanguageDetector {
	languageDetector.Do(func() {
		languageDetector.LanguageDetector = lingua.NewLanguageDetectorBuilder().
			FromLanguages(lang...).
			Build()
	})
	return languageDetector.LanguageDetector
}
