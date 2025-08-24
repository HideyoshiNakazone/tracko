package utils

import "testing"

func Test_CheckModelHasTag(t *testing.T) {
	type TestStruct struct {
		Field1 string `mapstructure:"field1" restricted:"true"`
		Field2 string `mapstructure:"field2"`
	}

	if !CheckModelHasTag(TestStruct{}, "field1", "restricted", "true") {
		t.Error("Expected Field1 to have restricted tag with value 'true'")
	}

	if CheckModelHasTag(TestStruct{}, "field2", "restricted", "true") {
		t.Error("Expected Field2 to not have restricted tag with value 'true'")
	}
}
