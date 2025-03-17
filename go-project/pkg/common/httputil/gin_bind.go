package httputil

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	validator "github.com/go-playground/validator/v10"
)

func inspectType(c *gin.Context, v interface{}) (interface{}, error) {
	// 获取反射类型
	vType := reflect.TypeOf(v)
	vValue := reflect.ValueOf(v)
	// 检查是否为指针
	if vType.Kind() != reflect.Ptr {
		return nil, errors.New("err: input must be a pointer")
	}

	// 检查指针指向的类型是否为结构体
	if vType.Elem().Kind() != reflect.Struct {
		return nil, errors.New("err: input must be a pointer to a struct")
	}

	// 获取指针指向的实际值
	vValue = vValue.Elem()
	vType = vValue.Type()

	// 遍历结构体字段
	for i := 0; i < vValue.NumField(); i++ {
		fieldValue := vValue.Field(i)
		fieldType := vType.Field(i)

		// 检查字段是否可导出
		if !fieldValue.CanInterface() {
			continue
		}

		// 检查是否有 form 标签
		formTag, hasFormTag := fieldType.Tag.Lookup("form")
		if !hasFormTag {
			continue
		}
		// 检查字段类型是否为切片
		if fieldValue.Kind() == reflect.Slice {
			// 获取切片元素的类型
			elemType := fieldValue.Type().Elem().Kind()
			// 从请求中获取表单数据
			res := c.PostFormArray(formTag)
			if len(res) == 0 {
				continue
			}

			// 创建目标类型的切片
			sliceValue := reflect.MakeSlice(fieldValue.Type(), len(res), len(res))

			// 根据元素类型进行转换
			switch elemType {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				for i, v := range res {
					val, err := strconv.ParseInt(v, 10, 64)
					if err != nil {
						return nil, fmt.Errorf("failed to convert '%s' to %v: %v", v, elemType, err)
					}
					sliceValue.Index(i).SetInt(val)
				}

			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				for i, v := range res {
					val, err := strconv.ParseUint(v, 10, 64)
					if err != nil {
						return nil, fmt.Errorf("failed to convert '%s' to %v: %v", v, elemType, err)
					}
					sliceValue.Index(i).SetUint(val)
				}

			case reflect.String:
				for i, v := range res {
					sliceValue.Index(i).SetString(v)
				}

			case reflect.Float32, reflect.Float64:
				for i, v := range res {
					val, err := strconv.ParseFloat(v, 64)
					if err != nil {
						return nil, fmt.Errorf("failed to convert '%s' to %v: %v", v, elemType, err)
					}
					sliceValue.Index(i).SetFloat(val)
				}

			default:
				return nil, fmt.Errorf("unsupported slice element type: %v", elemType)
			}

			// 将转换后的切片设置到字段中
			fieldValue.Set(sliceValue)
		}
	}

	return v, nil
}

// 注意,不支持
func ShouldBindFormSlice(c *gin.Context, obj any) error {
	obj, err := inspectType(c, obj)
	if err != nil {
		return err
	}
	return validate(obj)
}

func validate(data any) error {
	v := binding.Validator.Engine().(*validator.Validate)
	return v.Struct(data)
}
