package mock

type MockContext struct {

}

func (mc *MockContext) LoadContract(key []byte) ([]byte, error) {
	return nil, nil
}
