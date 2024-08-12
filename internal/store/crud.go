package store

func Select[T any](query string, args ...interface{}) ([]T, error) {
	res := make([]T, 0, 10)
	err := db.Select(&res, query, args...)
	if err != nil {
		return nil, err
	}
	return res, nil
}
