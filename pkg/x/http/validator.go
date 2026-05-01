package http

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/locales"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
	"github.com/pkg/errors"
)

var (
	customTrans = map[string]string{
		"mobile": "手机号码有误",
		"id_num": "身份证号码有误",
	}
	locale = "zh"
)

type Validator struct {
	v *validator.Validate
	t ut.Translator
}

func NewValidator() *Validator {
	validate := validator.New()
	translator, err := newLocalTranslator(locale)
	if err != nil {
		panic(err)
	}
	err = register(validate, translator)
	if err != nil {
		panic(err)
	}
	return &Validator{
		v: validate,
		t: translator,
	}
}

func (v *Validator) Validate(r *http.Request, data any) error {
	err := v.v.Struct(data)
	if err != nil {
		var newErr validator.ValidationErrors
		switch {
		case errors.As(err, &newErr):
			for _, e := range newErr {
				return errors.New(e.Translate(v.t))
			}
		default:
			return err
		}
	}

	return nil
}

func newLocalTranslator(locale string) (ut.Translator, error) {
	translators := []locales.Translator{zh.New(), en.New()}
	uni := ut.New(translators[0], translators...)
	translator, found := uni.GetTranslator(locale)
	if !found {
		return nil, errors.Errorf("%s locale t not found", locale)
	}
	return translator, nil
}

func register(v *validator.Validate, t ut.Translator) error {
	err := registerZhTranslations(v, t)
	if err != nil {
		return err
	}
	return registerCustomValidations(v)
}

func registerZhTranslations(v *validator.Validate, trans ut.Translator) error {
	err := zhTranslations.RegisterDefaultTranslations(v, trans)
	if err != nil {
		return err
	}
	for tag, msg := range customTrans {
		registerFn := func(ut ut.Translator) error {
			return ut.Add(tag, msg, false)
		}

		transFn := func(ut ut.Translator, fe validator.FieldError) string {
			t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
			if err != nil {
				return fe.(error).Error()
			}
			return t
		}

		err = v.RegisterTranslation(tag, trans, registerFn, transFn)
		if err != nil {
			return err
		}
	}
	return nil
}

func registerCustomValidations(v *validator.Validate) error {
	err := v.RegisterValidation("mobile", ValidateMobile)
	if err != nil {
		return err
	}
	err = v.RegisterValidation("id_num", ValidateIdNum)
	if err != nil {
		return err
	}
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	return nil
}
