package tests

// TestContext holds the context for running tests.
type TestContext struct {
	Services   *Services
	Token      string
	ObjectName string
}

// NewTestContext creates a new TestContext with default values.
func NewTestContext(objectName string, services *Services) *TestContext {
	token := getExternalToken(services)

	return &TestContext{
		Services:   services,
		Token:      token,
		ObjectName: objectName,
	}
}
