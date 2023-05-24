package code

type Language string

func (l Language) IsEmpty() bool {
	return l == ""
}

const (
	LChinese  Language = "zh"
	LEnglish  Language = "en"
	LSpanish  Language = "es"
	LJapanese Language = "jp"
	LKorean   Language = "ko"
)
