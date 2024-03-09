package csv_encoder

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"reflect"

	"backend/internal/domain"
)

const key = "json"

func Encode(w http.ResponseWriter, data []domain.TodoListModel) error {
	writer := csv.NewWriter(w)
	defer writer.Flush()

	t := reflect.TypeOf(data[0])

	headers := make([]string, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		headers[i] = t.Field(i).Tag.Get(key)
	}

	err := writer.Write(headers)
	if err != nil {
		return err
	}

	for _, item := range data {
		values := make([]string, t.NumField())
		for i := 0; i < t.NumField(); i++ {
			field := reflect.ValueOf(item).Field(i)
			switch field.Kind() {
			case reflect.Slice:
				slice := field.Interface().([]string)
				values[i] = fmt.Sprintf("%v", slice)
			default:
				values[i] = fmt.Sprintf("%v", field)
			}
		}

		err = writer.Write(values)
		if err != nil {
			return err
		}
	}

	return nil
}
