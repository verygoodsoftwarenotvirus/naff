package utils

import "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

// RequireNoError creates a require call
func RequireNoError(value, message *jen.Statement, formatArgs ...*jen.Statement) jen.Code {
	return buildSingleValueTestifyFunc(r, "NoError")(value, message, formatArgs...)
}

// RequireNotNil creates a require call
func RequireNotNil(value, message *jen.Statement, formatArgs ...*jen.Statement) jen.Code {
	return buildSingleValueTestifyFunc(r, "NotNil")(value, message, formatArgs...)
}

// RequireNil creates a require call
func RequireNil(value, message *jen.Statement, formatArgs ...*jen.Statement) jen.Code {
	return buildSingleValueTestifyFunc(r, "Nil")(value, message, formatArgs...)
}

// AssertTrue calls assert.True
func AssertTrue(value, message *jen.Statement, formatArgs ...*jen.Statement) jen.Code {
	return buildSingleValueTestifyFunc(a, "True")(value, message, formatArgs...)
}

// AssertFalse calls assert.False
func AssertFalse(value, message *jen.Statement, formatArgs ...*jen.Statement) jen.Code {
	return buildSingleValueTestifyFunc(a, "False")(value, message, formatArgs...)
}

// AssertNotNil calls assert.NotNil
func AssertNotNil(value, message *jen.Statement, formatArgs ...*jen.Statement) jen.Code {
	return buildSingleValueTestifyFunc(a, "NotNil")(value, message, formatArgs...)
}

// AssertNil calls assert.Nil
func AssertNil(value, message *jen.Statement, formatArgs ...*jen.Statement) jen.Code {
	return buildSingleValueTestifyFunc(a, "Nil")(value, message, formatArgs...)
}

// AssertError calls assert.Error
func AssertError(value, message *jen.Statement, formatArgs ...*jen.Statement) jen.Code {
	//if message == nil {
	//	message = jen.Lit("error should be returned")
	//}
	return buildSingleValueTestifyFunc(a, "Error")(value, message, formatArgs...)
}

// AssertNoError calls assert.NoError
func AssertNoError(value, message *jen.Statement, formatArgs ...*jen.Statement) jen.Code {
	//if message == nil {
	//	message = jen.Lit("no error should be returned")
	//}
	return buildSingleValueTestifyFunc(a, "NoError")(value, message, formatArgs...)
}

// AssertEmpty calls assert.NotEmpty
func AssertEmpty(value, message *jen.Statement, formatArgs ...*jen.Statement) jen.Code {
	return buildSingleValueTestifyFunc(a, "Empty")(value, message, formatArgs...)
}

// AssertNotEmpty calls assert.NotEmpty
func AssertNotEmpty(value, message *jen.Statement, formatArgs ...*jen.Statement) jen.Code {
	return buildSingleValueTestifyFunc(a, "NotEmpty")(value, message, formatArgs...)
}

// AssertZero calls assert.NotEmpty
func AssertZero(value, message *jen.Statement, formatArgs ...*jen.Statement) jen.Code {
	return buildSingleValueTestifyFunc(a, "Zero")(value, message, formatArgs...)
}

// AssertNotZero calls assert.NotEmpty
func AssertNotZero(value, message *jen.Statement, formatArgs ...*jen.Statement) jen.Code {
	return buildSingleValueTestifyFunc(a, "NotZero")(value, message, formatArgs...)
}

// AssertLength calls assert.NotEmpty
func AssertLength(value, length, message *jen.Statement, formatArgs ...*jen.Statement) jen.Code {
	return buildDoubleValueTestifyFunc(a, "Len")(value, length, message, formatArgs...)
}

// AssertContains calls assert.NotEmpty
func AssertContains(container, value, message *jen.Statement, formatArgs ...*jen.Statement) jen.Code {
	return buildDoubleValueTestifyFunc(a, "Contains")(container, value, message, formatArgs...)
}

// AssertEqual calls assert.Equal
func AssertEqual(expected, actual, message *jen.Statement, formatArgs ...*jen.Statement) jen.Code {
	//if message == nil && len(formatArgs) == 0 {
	//	message = jen.Lit("expected %v to equal %v")
	//	formatArgs = []*jen.Statement{
	//		expected,
	//		actual,
	//	}
	//}
	return buildDoubleValueTestifyFunc(a, "Equal")(expected, actual, message, formatArgs...)
}

// AssertNotEqual calls assert.Equal
func AssertNotEqual(expected, actual, message *jen.Statement, formatArgs ...*jen.Statement) jen.Code {
	//if message == nil && len(formatArgs) == 0 {
	//	message = jen.Lit("expected %v not to equal %v")
	//	formatArgs = []*jen.Statement{
	//		expected,
	//		actual,
	//	}
	//}
	return buildDoubleValueTestifyFunc(a, "NotEqual")(expected, actual, message, formatArgs...)
}
