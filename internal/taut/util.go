package taut

var (
	mockConnection *MockSlackConnection
	mockDiver      *Driver
)

func initMockDriver() *Driver {
	mockConnection = NewMockSlackConnection()
	mockDiver = New(mockConnection)

	return mockDiver
}
