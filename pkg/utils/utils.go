package utils

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"
	"unicode"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stringx"
)

// 判断字符串是否包含中文
func IsChinese(str string) bool {
	var flag bool
	for _, v := range str {
		if unicode.Is(unicode.Han, v) {
			flag = true
			break
		}
	}
	return flag
}

// 从email中解析出用户名。使用场景如：用于三方登录时用户名的解析
func GetUsernameFromEmail(email string) string {
	tokens := strings.Split(email, "@")
	if len(tokens) < 2 {
		stringx.Seed(time.Now().UnixNano())
		return stringx.Rand()
	} else {
		return tokens[0]
	}
}

// 使用反射获取对象的属性
func GetByReflection[T any](obj any, name string) (*T, error) {

	valueOfExpire := reflect.ValueOf(obj).Elem().FieldByName(name)
	if valueOfExpire.IsValid() {
		rst, ok := valueOfExpire.Interface().(T)
		if !ok {
			return nil, errors.New("type error")
		}
		return &rst, nil
	}
	return nil, errors.New("not found")
}

// 使用反射设置对象的属性，请谨慎使用，可能会导致panic
func SetByReflection[T any](obj any, target any, targetName string) error {

	targetObj := reflect.ValueOf(obj).Elem().FieldByName(targetName)
	if targetObj.IsValid() {
		targetObj.Set(reflect.ValueOf(target))
		return nil
	}
	return errors.New("not found")
}

// 通过输入的生日，计算年龄
// 生日格式为：2000-01-01
func GetAgeByBirthday(birthday string) int {
	birthdayTime, err := time.Parse("2006-01-02", birthday)
	if err != nil {
		return 0
	}
	now := time.Now()
	age := now.Year() - birthdayTime.Year()
	if now.Month() < birthdayTime.Month() || (now.Month() == birthdayTime.Month() && now.Day() < birthdayTime.Day()) {
		age--
	}
	return age
}

// 比较两个时间是否为同一天
func IsSameDay(t1, t2 time.Time) bool {
	return t1.Year() == t2.Year() && t1.Month() == t2.Month() && t1.Day() == t2.Day()
}

// 等待锁释放
func AcquireAndWaitRedisLock(ctx context.Context, lock *redis.RedisLock, timeout int) error {
	lock.SetExpire(5) // 设置锁 5s 超时
	t := time.Now()
	for !time.Now().After(t.Add(time.Duration(timeout) * time.Second)) {
		ok, err := lock.AcquireCtx(ctx)
		if err != nil {
			return err
		}
		if ok {
			return nil
		}
		time.Sleep(time.Millisecond * 20)
	}
	return errors.New("timeout")
}

// 使用动态数据合并到模板，非零字段覆盖模板中的对应字段
// template必须是指向结构体的指针，dynamic可以是结构体或结构体指针
func MergeTemplateWithDynamic(template interface{}, dynamic interface{}) error {
	templateValue := reflect.ValueOf(template)
	if templateValue.Kind() != reflect.Ptr || templateValue.IsNil() {
		return fmt.Errorf("template必须是非空指针")
	}
	templateValue = templateValue.Elem()
	if templateValue.Kind() != reflect.Struct {
		return fmt.Errorf("template必须指向结构体")
	}

	dynamicValue := reflect.ValueOf(dynamic)
	if dynamicValue.Kind() == reflect.Ptr {
		dynamicValue = dynamicValue.Elem()
	}
	if dynamicValue.Kind() != reflect.Struct {
		return fmt.Errorf("dynamic必须是结构体或结构体指针")
	}

	dynamicType := dynamicValue.Type()

	for i := 0; i < dynamicValue.NumField(); i++ {
		fieldName := dynamicType.Field(i).Name
		dynamicField := dynamicValue.Field(i)
		templateField := templateValue.FieldByName(fieldName)

		if !templateField.IsValid() || !templateField.CanSet() {
			continue
		}

		if isZero(dynamicField) {
			continue
		}

		if templateField.Type() != dynamicField.Type() {
			return fmt.Errorf("字段类型不匹配 %s: template类型 %v, dynamic类型 %v",
				fieldName, templateField.Type(), dynamicField.Type())
		}

		switch dynamicField.Kind() {
		case reflect.Ptr:
			if dynamicField.IsNil() {
				continue
			}

			if dynamicField.Elem().Kind() == reflect.Struct {
				if templateField.IsNil() {
					templateField.Set(reflect.New(dynamicField.Elem().Type()))
				}
				if err := MergeTemplateWithDynamic(templateField.Interface(), dynamicField.Interface()); err != nil {
					return err
				}
			} else {
				templateField.Set(dynamicField)
			}

		case reflect.Struct:
			if err := MergeTemplateWithDynamic(templateField.Addr().Interface(), dynamicField.Interface()); err != nil {
				return err
			}

		case reflect.Slice:
			// 处理数组类型，按顺序合并元素
			if dynamicField.Len() > 0 {
				// 如果模板字段是nil，初始化一个相同类型的新切片
				if templateField.IsNil() {
					templateField.Set(reflect.MakeSlice(dynamicField.Type(), dynamicField.Len(), dynamicField.Cap()))
				} else if templateField.Len() < dynamicField.Len() {
					// 如果模板切片长度小于动态切片，扩展模板切片
					newSlice := reflect.MakeSlice(templateField.Type(), dynamicField.Len(), dynamicField.Cap())
					reflect.Copy(newSlice, templateField)
					templateField.Set(newSlice)
				}

				// 按顺序合并元素
				for j := 0; j < dynamicField.Len(); j++ {
					dynamicElem := dynamicField.Index(j)
					// 只处理非零值元素
					if !isZero(dynamicElem) {
						templateElem := templateField.Index(j)

						// 根据元素类型统一处理
						switch dynamicElem.Kind() {
						case reflect.Ptr:
							if dynamicElem.IsNil() {
								continue
							}
							if dynamicElem.Elem().Kind() == reflect.Struct {
								if templateElem.IsNil() {
									templateElem.Set(reflect.New(dynamicElem.Elem().Type()))
								}
								if err := MergeTemplateWithDynamic(templateElem.Interface(), dynamicElem.Interface()); err != nil {
									return err
								}
							} else {
								templateElem.Set(dynamicElem)
							}
						case reflect.Struct:
							if err := MergeTemplateWithDynamic(templateElem.Addr().Interface(), dynamicElem.Interface()); err != nil {
								return err
							}
						default:
							templateElem.Set(dynamicElem)
						}
					}
				}
			}
		default:
			templateField.Set(dynamicField)
		}
	}

	return nil
}

// 更完善的零值判断方法
func isZero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Ptr, reflect.Interface:
		return v.IsNil()
	case reflect.Slice, reflect.Map, reflect.Chan:
		return v.IsNil() || v.Len() == 0
	case reflect.Struct:
		zero := reflect.Zero(v.Type()).Interface()
		return reflect.DeepEqual(v.Interface(), zero)
	default:
		return reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
	}
}
