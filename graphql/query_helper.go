package graphql

import "fmt"

// AddVarNamesToQuery accepts a query or mutation definition and will populate the definition with the variables passed.
// The definitions use the notation: $%[1]s, $%[2]s, etc... to avoid passing the same variables multiple times in the varNames field(s).
// When wanting to reference a variable several times inside a format string,
// the variables can be reference by position using %[n] where n is the index of the parameter (1 based).
func AddVarNamesToQuery(query string, varNames ...interface{}) string {
	return fmt.Sprintf(query, varNames...)
}
