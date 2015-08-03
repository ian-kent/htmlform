package htmlform

import "fmt"

// Map creates a map with the provided arg pairs
func Map(args ...interface{}) (map[string]interface{}, error) {
	if len(args)%2 != 0 {
		return nil, fmt.Errorf("expecting even number of arguments, got %d", len(args))
	}

	m := make(map[string]interface{})
	fn := ""
	for _, v := range args {
		if len(fn) == 0 {
			if s, ok := v.(string); ok {
				fn = s
				continue
			}
			return m, fmt.Errorf("expecting string for odd numbered arguments, got %+v", v)
		}
		m[fn] = v
		fn = ""
	}

	return m, nil
}

// Arr returns a slice for a given argument list
func Arr(args ...interface{}) []interface{} {
	return args
}

// Extend extends the target map with the provided arg pairs
func Extend(target map[string]interface{}, args ...interface{}) (map[string]interface{}, error) {
	if len(args)%2 != 0 {
		return nil, fmt.Errorf("expecting even number of arguments, got %d", len(args))
	}

	fn := ""
	for _, v := range args {
		if len(fn) == 0 {
			if s, ok := v.(string); ok {
				fn = s
				continue
			}
			return target, fmt.Errorf("expecting string for odd numbered arguments, got %+v", v)
		}
		target[fn] = v
		fn = ""
	}

	return target, nil
}

// FirstNotNil returns th first parameter that isn't nil
func FirstNotNil(args ...interface{}) interface{} {
	for _, v := range args {
		if v != nil {
			return v
		}
	}
	return nil
}
