package goconv

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var (
	ErrorTypeConvert    = errors.New("convert typ failed")
	ErrorTypeUnSoupport = errors.New("typ un soupport!")
)

var (
	DefaultDateLayout         = "2006-01-02"
	DefaultDateTimeLayout     = "2006-01-02 15:04:05"
	DefaultDateTimeZoneLayout = "2006-01-02 15:04:05 -0700 MST"
)

type DTOConverter struct {
	td *TypeDirect

	pc *ParseConfig
}

type ParseConfig struct {
	arrayPrefix, arraySuffix string
	arraySplitFlag           string

	timeLayout string
	dateLayout string
}

func NewDTOConverter(opts ...Option) *DTOConverter {
	converter := &DTOConverter{
		td: DefaultTypeDirect,
		pc: &ParseConfig{
			arrayPrefix:    "[",
			arraySuffix:    "]",
			arraySplitFlag: ",",
			timeLayout:     DefaultDateTimeLayout,
			dateLayout:     DefaultDateLayout,
		},
	}

	for _, o := range opts {
		o(converter)
	}

	return converter
}

type Option func(*DTOConverter)

func WithArraySplitConfig(arrayPrefix, arraySuffix, arraySplitFlag string) Option {
	return func(d *DTOConverter) {
		d.pc.arrayPrefix = arrayPrefix
		d.pc.arraySuffix = arraySuffix
		d.pc.arraySplitFlag = arraySplitFlag
	}
}

func WithDateTimeLayout(layout string) Option {
	return func(d *DTOConverter) {
		d.pc.timeLayout = layout
	}
}

func WithDateLayout(layout string) Option {
	return func(d *DTOConverter) {
		d.pc.dateLayout = layout
	}
}

func WithTypeDirect(td *TypeDirect) Option {
	return func(d *DTOConverter) {
		d.td = td
	}
}

func (dc DTOConverter) Convert(typ, value string) (result interface{}, err error) {
	t := dc.td.DirectTyp(typ)
	if t == UnKnow {
		return nil, fmt.Errorf("type %s UnKnow, alias miss match", typ)
	}
	if dc.pc == nil {
		return nil, errors.New("DTOconverter need ParseConfig != nil")
	}
	parseConfig := dc.pc

	switch t {
	case Int:
		var v int64
		v, err = strconv.ParseInt(value, 0, 0)
		result = int(v)
	case Int32:
		var v int64
		v, err = strconv.ParseInt(value, 0, 0)
		result = int32(v)
	case Int64:
		result, err = strconv.ParseInt(value, 0, 0)
	case Float32:
		var v float64
		v, err = strconv.ParseFloat(value, 32)
		result = float32(v)
	case Float64:
		result, err = strconv.ParseFloat(value, 64)
	case String:
		result = value
	case Boolean:
		result, err = strconv.ParseBool(value)
	case Date:
		result, err = time.ParseInLocation(parseConfig.dateLayout, value, time.Local)
	case DateTime:
		result, err = time.ParseInLocation(parseConfig.timeLayout, value, time.Local)
	case Ints:
		vs := strings.Split(value, parseConfig.arraySplitFlag)
		ints := make([]int, 0, len(vs))
		for _, v := range vs {
			i, err := strconv.ParseInt(v, 0, 0)
			if err != nil {
				return nil, ErrorTypeConvert
			}
			ints = append(ints, int(i))
		}
		result = ints
	case Float32s:
		vs := strings.Split(value, parseConfig.arraySplitFlag)
		float32s := make([]float32, 0, len(vs))
		for _, v := range vs {
			f, err := strconv.ParseFloat(v, 32)
			if err != nil {
				return nil, ErrorTypeConvert
			}
			float32s = append(float32s, float32(f))
		}
		result = float32s
	case Float64s:
		vs := strings.Split(value, parseConfig.arraySplitFlag)
		float64s := make([]float64, 0, len(vs))
		for _, v := range vs {
			f, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return nil, ErrorTypeConvert
			}
			float64s = append(float64s, f)
		}
		result = float64s
	case Strings:
		result = strings.Split(value, parseConfig.arraySplitFlag)
	case Booleans:
		vs := strings.Split(value, parseConfig.arraySplitFlag)
		bs := make([]bool, 0, len(vs))
		for _, v := range vs {
			b, err := strconv.ParseBool(v)
			if err != nil {
				return nil, ErrorTypeConvert
			}
			bs = append(bs, b)
		}
		result = bs
	case Dates:
		vs := strings.Split(value, parseConfig.arraySplitFlag)
		ts := make([]time.Time, 0, len(vs))
		for _, v := range vs {
			t, err := time.ParseInLocation(parseConfig.dateLayout, v, time.Local)
			if err != nil {
				return nil, ErrorTypeConvert
			}
			ts = append(ts, t)
		}
		result = ts
	case DateTimes:
		vs := strings.Split(value, parseConfig.arraySplitFlag)
		ts := make([]time.Time, 0, len(vs))
		for _, v := range vs {
			t, err := time.ParseInLocation(parseConfig.timeLayout, v, time.Local)
			if err != nil {
				return nil, ErrorTypeConvert
			}
			ts = append(ts, t)
		}
		result = ts
	default:
		return nil, ErrorTypeUnSoupport
	}

	if err != nil {
		return nil, ErrorTypeConvert
	}

	return
}

func (dc DTOConverter) ConvertReflect(typ, value string, result interface{}) (err error) {
	v, err := dc.Convert(typ, value)
	if err != nil {
		return err
	}
	return setValueReflect(v, result)
}

// func (dc DTOConverter) JsonMarshal() []byte   { return nil }
// func (dc DTOConverter) JsonUnMarshal() map[string]interface{} { return nil }

func setValueReflect(origin, result interface{}) error {
	outVal := reflect.ValueOf(result)
	if outVal.Kind() != reflect.Ptr {
		return errors.New("output need ptr")
	}

	if outVal.IsNil() {
		outVal.Set(reflect.New(outVal.Type().Elem()))
	}

	outValElem := outVal.Elem()
	originValue := reflect.ValueOf(origin)
	if originValue.Kind() != outValElem.Kind() {
		return fmt.Errorf("got kind %s != output kind %s", originValue.Kind().String(), outValElem.Kind().String())
	}
	if !outValElem.CanSet() {
		return fmt.Errorf("output cant set value, %s", outValElem.Kind().String())
	}
	outValElem.Set(originValue)
	return nil
}
