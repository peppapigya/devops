package common

func ConvertList[R any, T any](sources []*T, convertFunc func(*T) *R) ([]*R, error) {
	if len(sources) == 0 {
		return []*R{}, nil
	}

	results := make([]*R, len(sources))
	for i, source := range sources {
		results[i] = convertFunc(source)
	}
	return results, nil
}

// Convert 单个对象的转化
func Convert[R any, T any](source *T, convertFunc func(*T) *R) (*R, error) {
	if source == nil {
		return nil, nil
	}
	return convertFunc(source), nil
}
