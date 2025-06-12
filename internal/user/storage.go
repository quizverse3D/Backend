package user

import "sync"

type Storage struct {
	users map[string]User
	mu    sync.Mutex
}

func NewStorage() *Storage {
	return &Storage{users: make(map[string]User)}
}

func (s *Storage) CreateUser(u User) error {
	s.mu.Lock()         // синхронная блокировка единственного хранилища
	defer s.mu.Unlock() // разблокировка хранилища после исполнения всего кода метода (за счёт defer)

	if _, exists := s.users[u.Username]; exists { // _ содержит значение ключа, но нам оно неинтересно
		return ErrUserExists
	}

	s.users[u.Username] = u
	return nil
}

func (s *Storage) GetUser(username string) (User, bool) {
	s.mu.Lock() // да, на чтение тоже блокировка-разблокировка
	defer s.mu.Unlock()

	u, ok := s.users[username]
	return u, ok
}
