package server

func NewMockStore() MapStore {
	return &mockStore{
		ips: make([]string, 0),
	}
}

type mockStore struct {
	ips []string
}

var _ MapStore = &mockStore{}

// GetAll implements MapStore
func (s *mockStore) GetAll() ([]string, error) {
	return s.ips, nil
}

// Put implements MapStore
func (s *mockStore) Put(ipStr string) error {
	s.ips = append(s.ips, ipStr)
	return nil
}

// Delete implements MapStore
func (s *mockStore) Delete(ipStr string) error {
	for i, ip := range s.ips {
		if ip == ipStr {
			s.ips = append(s.ips[:i], s.ips[i+1:]...)
			break
		}
	}

	return nil
}
