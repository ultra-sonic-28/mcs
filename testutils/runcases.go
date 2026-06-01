package testutils

import (
	"reflect"
	"testing"
)

func RunCases[T any](t *testing.T, cases []T, fn func(*testing.T, T)) {
	t.Helper()

	for _, tc := range cases {
		tc := tc

		name := extractName(tc)

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			fn(t, tc)
		})
	}
}

func extractName(v any) string {

	val := reflect.ValueOf(v)

	if val.Kind() == reflect.Struct {

		field := val.FieldByName("name")

		if field.IsValid() && field.Kind() == reflect.String {
			return field.String()
		}
	}

	return "case"
}
