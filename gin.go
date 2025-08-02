package mox

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/samber/lo"
	"github.com/samber/mo"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var (
	ErrNotSupportKind            = errors.New("not support kind")
	ErrNotSupportOptionValueKind = errors.New("not support option value kind")
	ErrOnlyStruct                = errors.New("only struct")
	OptionFormBinding            = &optionFormBinding{}
	//OptionQueryBinding option not allow pointer value or array
	OptionQueryBinding = &optionQueryBinding{}
)

const defaultMemory = 32 << 20

type optionFormBinding struct{}

func (optionFormBinding) Name() string {
	return "OptionForm"
}

func (optionFormBinding) Bind(req *http.Request, obj any) error {
	if err := req.ParseForm(); err != nil {
		return err
	}
	if err := req.ParseMultipartForm(defaultMemory); err != nil && !errors.Is(err, http.ErrNotMultipart) {
		return err
	}
	if err := mapForm(obj, req.Form); err != nil {
		return err
	}
	return validate(obj)
}

func ShouldBindGinUri(c *gin.Context, obj any) error {
	values := make(map[string][]string, len(c.Params))
	for _, v := range c.Params {
		values[v.Key] = []string{v.Value}
	}
	if err := mapForm(obj, values); err != nil {
		return err
	}
	return validate(obj)
}

type optionQueryBinding struct {
}

func (t *optionQueryBinding) Name() string {
	return "OptionQuery"
}

func (t *optionQueryBinding) Bind(req *http.Request, obj any) error {
	values := req.URL.Query()
	if err := mapForm(obj, values); err != nil {
		return err
	}
	return validate(obj)
}

func validate(obj any) error {
	if binding.Validator == nil {
		return nil
	}
	return binding.Validator.ValidateStruct(obj)
}

func mapForm(ptr any, form map[string][]string) error {
	// Check if ptr is a map
	ptrValue := reflect.ValueOf(ptr)
	if ptrValue.Kind() == reflect.Ptr {
		ptrValue = ptrValue.Elem()
	}
	if ptrValue.Kind() != reflect.Struct {
		return fmt.Errorf("%w: kind=%s", ErrOnlyStruct, ptrValue.Kind().String())
	}

	ptrType := ptrValue.Type()
	for i := range ptrValue.NumField() {
		field := ptrType.Field(i)
		tag := field.Tag.Get("form")
		if tag == "-" {
			continue
		}
		tags := lo.Filter(strings.Split(tag, ","), func(item string, index int) bool {
			return strings.TrimSpace(item) != ""
		})
		// ignore default=defaultValue
		var name string
		if len(tags) > 0 {
			name = tags[0]
		}
		if name == "" {
			name = field.Name
		}
		vs := form[name]
		if len(vs) == 0 {
			continue
		}
		fieldValue := ptrValue.Field(i)

		if IsOption(field.Type) {
			if err := setOptionValue(vs, fieldValue, field); err != nil {
				return err
			}
		} else {
			if err := setValue(vs, fieldValue, field); err != nil {
				return err
			}
		}
	}
	return nil
}
func setValue(vs []string, value reflect.Value, field reflect.StructField) error {
	switch field.Type.Kind() {
	case reflect.Slice:
		if err := setSlice(vs, value, field); err != nil {
			return err
		}
	case reflect.Array:
		if err := setArray(vs, value, field); err != nil {
			return err
		}
	default:
		if err := setWithProperType(vs[0], value, field); err != nil {
			return err
		}
	}
	return nil
}
func setOptionValue(vs []string, value reflect.Value, field reflect.StructField) error {
	optionValue := value.FieldByName("value")
	switch optionValue.Kind() {
	case reflect.Slice:
		if err := setOptionSlice(vs, value, field, optionValue); err != nil {
			return err
		}
	default:
		if err := setWithProperOptionType(vs[0], value, field, optionValue); err != nil {
			return err
		}
	}
	return nil
}

func setIntField(val string, bitSize int, field reflect.Value) error {
	if val == "" {
		val = "0"
	}
	intVal, err := strconv.ParseInt(val, 10, bitSize)
	if err == nil {
		field.SetInt(intVal)
	}
	return err
}

func setUintField(val string, bitSize int, field reflect.Value) error {
	if val == "" {
		val = "0"
	}
	uintVal, err := strconv.ParseUint(val, 10, bitSize)
	if err == nil {
		field.SetUint(uintVal)
	}
	return err
}

func setBoolField(val string, field reflect.Value) error {
	if val == "" {
		val = "false"
	}
	boolVal, err := strconv.ParseBool(val)
	if err == nil {
		field.SetBool(boolVal)
	}
	return err
}

func setFloatField(val string, bitSize int, field reflect.Value) error {
	if val == "" {
		val = "0.0"
	}
	floatVal, err := strconv.ParseFloat(val, bitSize)
	if err == nil {
		field.SetFloat(floatVal)
	}
	return err
}

func setTimeDuration(val string, value reflect.Value) error {
	d, err := time.ParseDuration(val)
	if err != nil {
		return err
	}
	value.Set(reflect.ValueOf(d))
	return nil
}

func setArray(vals []string, value reflect.Value, field reflect.StructField) error {
	for i, s := range vals {
		err := setWithProperType(s, value.Index(i), field)
		if err != nil {
			return err
		}
	}
	return nil
}

func setSlice(vals []string, value reflect.Value, field reflect.StructField) error {
	slice := reflect.MakeSlice(value.Type(), len(vals), len(vals))
	err := setArray(vals, slice, field)
	if err != nil {
		return err
	}
	value.Set(slice)
	return nil
}

func setWithProperType(val string, value reflect.Value, field reflect.StructField) error {
	switch value.Kind() {
	default:
		return fmt.Errorf("%w: %s is %s", ErrNotSupportKind, field.Name, value.Kind())
	case reflect.Bool:
		return setBoolField(val, value)
	case reflect.String:
		value.SetString(val)
	case reflect.Ptr:
		if !value.Elem().IsValid() {
			value.Set(reflect.New(value.Type().Elem()))
		}
		return setWithProperType(val, value.Elem(), field)
	case reflect.Int:
		return setIntField(val, 0, value)
	case reflect.Int8:
		return setIntField(val, 8, value)
	case reflect.Int16:
		return setIntField(val, 16, value)
	case reflect.Int32:
		return setIntField(val, 32, value)
	case reflect.Int64:
		switch value.Interface().(type) {
		case time.Duration:
			return setTimeDuration(val, value)
		}
		return setIntField(val, 64, value)
	case reflect.Uint:
		return setUintField(val, 0, value)
	case reflect.Uint8:
		return setUintField(val, 8, value)
	case reflect.Uint16:
		return setUintField(val, 16, value)
	case reflect.Uint32:
		return setUintField(val, 32, value)
	case reflect.Uint64:
		return setUintField(val, 64, value)
	case reflect.Float32:
		return setFloatField(val, 32, value)
	case reflect.Float64:
		return setFloatField(val, 64, value)
	}
	return nil
}

func setOptionSlice(vs []string, value reflect.Value, field reflect.StructField, optionValue reflect.Value) error {
	if optionValue.Type().Elem().Kind() == reflect.String {
		value.Set(reflect.ValueOf(mo.Some[[]string](vs)))
		return nil
	}
	switch optionValue.Type().Elem().Kind() {
	default:
		return fmt.Errorf("%w: %s is %s", ErrNotSupportOptionValueKind, field.Name, optionValue.Kind().String())
	case reflect.Bool:
		return _setOptionSlice[bool](vs, value, field, optionValue, func(str string) (bool, error) {
			return str == "" || str == "true", nil
		})
	case reflect.Int:
		return _setOptionSlice[int](vs, value, field, optionValue, func(str string) (int, error) {
			v, err := strconv.ParseInt(str, 10, 0)
			if err != nil {
				return 0, err
			}
			return int(v), nil
		})
	case reflect.Int8:
		return _setOptionSlice[int8](vs, value, field, optionValue, func(str string) (int8, error) {
			v, err := strconv.ParseInt(str, 10, 8)
			if err != nil {
				return 0, err
			}
			return int8(v), nil
		})
	case reflect.Int16:
		return _setOptionSlice[int16](vs, value, field, optionValue, func(str string) (int16, error) {
			v, err := strconv.ParseInt(str, 10, 16)
			if err != nil {
				return 0, err
			}
			return int16(v), nil
		})
	case reflect.Int32:
		return _setOptionSlice[int32](vs, value, field, optionValue, func(str string) (int32, error) {
			v, err := strconv.ParseInt(str, 10, 32)
			if err != nil {
				return 0, err
			}
			return int32(v), nil
		})
	case reflect.Int64:
		switch optionValue.Type().Elem() {
		case reflect.TypeOf(time.Nanosecond):
			return _setOptionSlice[time.Duration](vs, value, field, optionValue, func(str string) (time.Duration, error) {
				v, err := time.ParseDuration(str)
				if err != nil {
					return 0, err
				}
				return v, nil
			})
		default:
			return _setOptionSlice[int64](vs, value, field, optionValue, func(str string) (int64, error) {
				v, err := strconv.ParseInt(str, 10, 64)
				if err != nil {
					return 0, err
				}
				return v, nil
			})
		}
	case reflect.Uint:
		return _setOptionSlice[uint](vs, value, field, optionValue, func(str string) (uint, error) {
			v, err := strconv.ParseUint(str, 10, 0)
			if err != nil {
				return 0, err
			}
			return uint(v), nil
		})
	case reflect.Uint8:
		return _setOptionSlice[uint8](vs, value, field, optionValue, func(str string) (uint8, error) {
			v, err := strconv.ParseUint(str, 10, 8)
			if err != nil {
				return 0, err
			}
			return uint8(v), nil
		})
	case reflect.Uint16:
		return _setOptionSlice[uint16](vs, value, field, optionValue, func(str string) (uint16, error) {
			v, err := strconv.ParseUint(str, 10, 16)
			if err != nil {
				return 0, err
			}
			return uint16(v), nil
		})
	case reflect.Uint32:
		return _setOptionSlice[uint32](vs, value, field, optionValue, func(str string) (uint32, error) {
			v, err := strconv.ParseUint(str, 10, 32)
			if err != nil {
				return 0, err
			}
			return uint32(v), nil
		})
	case reflect.Uint64:
		return _setOptionSlice[uint64](vs, value, field, optionValue, func(str string) (uint64, error) {
			v, err := strconv.ParseUint(str, 10, 64)
			if err != nil {
				return 0, err
			}
			return v, nil
		})
	case reflect.Float32:
		return _setOptionSlice[float32](vs, value, field, optionValue, func(str string) (float32, error) {
			v, err := strconv.ParseFloat(str, 32)
			if err != nil {
				return 0, err
			}
			return float32(v), nil
		})
	case reflect.Float64:
		return _setOptionSlice[float64](vs, value, field, optionValue, func(str string) (float64, error) {
			v, err := strconv.ParseFloat(str, 64)
			if err != nil {
				return 0, err
			}
			return v, nil
		})
	}
}

func _setOptionSlice[T int | int8 | int16 | int32 | ~int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64 | string | bool](vs []string, value reflect.Value, field reflect.StructField, optionValue reflect.Value, convert func(v string) (T, error)) error {
	values := make([]T, 0, len(vs))
	for _, v := range vs {
		tv, err := convert(v)
		if err != nil {
			return fmt.Errorf("convert %s to %s: %w", field.Name, optionValue.Elem().Type(), err)
		}
		values = append(values, tv)
	}
	value.Set(reflect.ValueOf(mo.Some[[]T](values)))
	return nil
}

// for can not set no export field
func setWithProperOptionType(val string, option reflect.Value, field reflect.StructField, optionValue reflect.Value) error {
	if val == "" && optionValue.Kind() != reflect.String && optionValue.Kind() != reflect.Bool {
		return fmt.Errorf("can use empty string to %s: %s", field.Name, optionValue.Kind())
	}
	switch optionValue.Kind() {
	default:
		return fmt.Errorf("%w: %s is %s", ErrNotSupportOptionValueKind, field.Name, optionValue.Kind().String())
	case reflect.String:
		option.Set(reflect.ValueOf(mo.Some[string](val)))
	case reflect.Bool:
		option.Set(reflect.ValueOf(mo.Some[bool](val == "" || val == "true")))
	case reflect.Int:
		v, err := strconv.ParseInt(val, 10, 0)
		if err != nil {
			return err
		}
		option.Set(reflect.ValueOf(mo.Some[int](int(v))))
	case reflect.Int8:
		v, err := strconv.ParseInt(val, 10, 8)
		if err != nil {
			return err
		}
		option.Set(reflect.ValueOf(mo.Some[int8](int8(v))))
	case reflect.Int16:
		v, err := strconv.ParseInt(val, 10, 16)
		if err != nil {
			return err
		}
		option.Set(reflect.ValueOf(mo.Some[int16](int16(v))))
	case reflect.Int32:
		v, err := strconv.ParseInt(val, 10, 32)
		if err != nil {
			return err
		}
		option.Set(reflect.ValueOf(mo.Some[int32](int32(v))))
	case reflect.Int64:
		switch option.Interface().(type) {
		case time.Duration:
			v, err := time.ParseDuration(val)
			if err != nil {
				return err
			}
			option.Set(reflect.ValueOf(mo.Some(v)))
		default:
			v, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				return err
			}
			option.Set(reflect.ValueOf(mo.Some[int64](v)))
		}
	case reflect.Uint:
		v, err := strconv.ParseUint(val, 10, 0)
		if err != nil {
			return err
		}
		option.Set(reflect.ValueOf(mo.Some[uint](uint(v))))
	case reflect.Uint8:
		v, err := strconv.ParseUint(val, 10, 8)
		if err != nil {
			return err
		}
		option.Set(reflect.ValueOf(mo.Some[uint8](uint8(v))))
	case reflect.Uint16:
		v, err := strconv.ParseUint(val, 10, 16)
		if err != nil {
			return err
		}
		option.Set(reflect.ValueOf(mo.Some[uint16](uint16(v))))
	case reflect.Uint32:
		v, err := strconv.ParseUint(val, 10, 32)
		if err != nil {
			return err
		}
		option.Set(reflect.ValueOf(mo.Some[uint32](uint32(v))))
	case reflect.Uint64:
		v, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			return err
		}
		option.Set(reflect.ValueOf(mo.Some[uint64](v)))
	case reflect.Float32:
		v, err := strconv.ParseFloat(val, 32)
		if err != nil {
			return err
		}
		option.Set(reflect.ValueOf(mo.Some[float32](float32(v))))
	case reflect.Float64:
		v, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return err
		}
		option.Set(reflect.ValueOf(mo.Some[float64](v)))
	}
	return nil
}
