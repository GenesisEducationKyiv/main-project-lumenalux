package port

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type StubStorage struct {
	data []map[string]string
	err  error
}

func (s *StubStorage) Append(record map[string]string) error {
	if s.err != nil {
		return s.err
	}
	s.data = append(s.data, record)
	return nil
}

func (s *StubStorage) AllRecords() ([]map[string]string, error) {
	if s.err != nil {
		return nil, s.err
	}
	return s.data, nil
}

func TestAdd(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		existingData []map[string]string
		emailToAdd   string
		expectedErr  error
	}{
		{
			name:         "Add new user successfully",
			existingData: []map[string]string{{"email": "existingEmail"}},
			emailToAdd:   "newEmail",
			expectedErr:  nil,
		},
		{
			name:         "Add an existing user",
			existingData: []map[string]string{{"email": "existingEmail"}},
			emailToAdd:   "existingEmail",
			expectedErr:  ErrAlreadyAdded,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			stubStorage := &StubStorage{data: tt.existingData}
			userRepository := NewUserRepository(stubStorage)

			err := userRepository.Add(&User{Email: tt.emailToAdd})

			require.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestFindByEmail(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		existingData []map[string]string
		emailToFind  string
		expectedErr  error
	}{
		{
			name:         "Find user successfully",
			existingData: []map[string]string{{"email": "existingEmail"}},
			emailToFind:  "existingEmail",
			expectedErr:  nil,
		},
		{
			name:         "User not found",
			existingData: []map[string]string{{"email": "existingEmail"}},
			emailToFind:  "nonExistingEmail",
			expectedErr:  ErrCannotFindByEmail,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			stubStorage := &StubStorage{data: tt.existingData}
			userRepository := NewUserRepository(stubStorage)

			_, err := userRepository.FindByEmail(tt.emailToFind)

			require.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestAll(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		existingData  []map[string]string
		storageError  error
		expectedErr   error
		expectedCount int
	}{
		{
			name: "Retrieve all users successfully",
			existingData: []map[string]string{
				{"email": "user1"}, {"email": "user2"},
			},
			storageError:  nil,
			expectedErr:   nil,
			expectedCount: 2,
		},
		{
			name:          "Error retrieving users",
			existingData:  nil,
			storageError:  ErrCannotLoadUsers,
			expectedErr:   ErrCannotLoadUsers,
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			stubStorage := &StubStorage{data: tt.existingData, err: tt.storageError}
			userRepository := NewUserRepository(stubStorage)

			users, err := userRepository.All()

			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tt.expectedCount, len(users))
		})
	}
}
