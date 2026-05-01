package cipher

import (
	"strings"
	"unicode/utf8"
)

var Masker _Masker

type _Masker struct{}

func (m *_Masker) Mask(plaintext string, category Category) string {
	switch category {
	case CategoryName:
		return m.maskName(plaintext)
	case CategoryIdNum:
		return m.maskIdNum(plaintext)
	case CategoryMobile:
		return m.maskMobile(plaintext)
	case CategoryEmail:
		return m.maskEmail(plaintext)
	case CategoryText:
		return m.maskText(plaintext)
	}

	return plaintext
}

func (m *_Masker) maskName(name string) string {
	if name == "" {
		return ""
	}

	l := utf8.RuneCountInString(name)

	if l < 2 {
		return MaskSymbol
	}

	return strings.Repeat(MaskSymbol, l-1) + string([]rune(name)[l-1:])
}

func (m *_Masker) maskIdNum(idNum string) string {
	if idNum == "" {
		return ""
	}

	l := utf8.RuneCountInString(idNum)

	if l < 16 {
		return strings.Repeat(MaskSymbol, l)
	}

	return string([]rune(idNum)[:1]) + strings.Repeat(MaskSymbol, l-2) + string([]rune(idNum)[l-1:])
}

func (m *_Masker) maskMobile(mobile string) string {
	if mobile == "" {
		return ""
	}
	l := utf8.RuneCountInString(mobile)
	if l < 3 {
		return strings.Repeat(MaskSymbol, l)
	}

	if l < 11 {
		return string([]rune(mobile)[:1]) + strings.Repeat(MaskSymbol, l-2) + string([]rune(mobile)[l-1:])
	}

	return string([]rune(mobile)[:3]) + strings.Repeat(MaskSymbol, l-7) + string([]rune(mobile)[l-4:])
}

func (m *_Masker) maskEmail(email string) string {
	if email == "" {
		return ""
	}

	if strings.Contains(email, EmailSeparator) {
		res := strings.Split(email, EmailSeparator)
		if len(res[0]) < 3 {
			return strings.Repeat(MaskSymbol, len(res[0])) + res[1]
		}

		return string([]rune(email)[:3]) + strings.Repeat(MaskSymbol, len(res[0])-3) + EmailSeparator + res[1]
	}

	return strings.Repeat(MaskSymbol, len(email))
}

func (m *_Masker) maskText(text string) string {
	if text == "" {
		return ""
	}
	l := utf8.RuneCountInString(text)
	if l < 5 {
		return strings.Repeat(MaskSymbol, l)
	}

	return string([]rune(text)[:1]) + strings.Repeat(MaskSymbol, 3) + string([]rune(text)[l-1:])
}
