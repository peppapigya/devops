package util

import (
	"reflect"
)

// Struct2MapList 使用反射进行类型转换
func Struct2MapList[T any, M any](list []T) ([]M, error) {
	result := make([]M, len(list))

	for i, item := range list {
		converted, err := convertType[T, M](item)
		if err != nil {
			return nil, err
		}
		result[i] = converted
	}

	return result, nil
}

func convertType[T any, M any](item T) (M, error) {
	var zero M

	tType := reflect.TypeOf((*T)(nil)).Elem()
	mType := reflect.TypeOf((*M)(nil)).Elem()

	// 如果类型相同或兼容，直接转换
	if tType == mType {
		return any(item).(M), nil
	}

	// 使用反射进行字段映射
	tValue := reflect.ValueOf(item)
	mValue := reflect.New(mType).Elem()

	if tValue.Kind() == reflect.Ptr {
		tValue = tValue.Elem()
	}
	if mValue.Kind() == reflect.Ptr {
		mValue = mValue.Elem()
	}

	// 结构体字段映射
	if tValue.Kind() == reflect.Struct && mValue.Kind() == reflect.Struct {
		for i := 0; i < tValue.NumField(); i++ {
			tField := tValue.Type().Field(i)
			if !tField.IsExported() {
				continue
			}

			mField := mValue.FieldByName(tField.Name)
			if mField.IsValid() && mField.CanSet() {
				if mField.Type() == tField.Type {
					mField.Set(tValue.Field(i))
				}
			}
		}
		return mValue.Interface().(M), nil
	}

	return zero, nil
}
