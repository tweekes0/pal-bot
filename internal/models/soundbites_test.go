package models

import (
	"testing"
	"time"

	test "github.com/tweekes0/pal-bot/internal/testing"
)

var (
	s1 = &Soundbite{
		ID:       1,
		Name:     "test1",
		Username: "test_username_1",
		UserID:   "111111",
		FilePath: "/path/to/file/1",
		FileHash: "sha256:111111",
	}
	s2 = &Soundbite{
		ID:       2,
		Name:     "test2",
		Username: "test_username_2",
		UserID:   "222222",
		FilePath: "/path/to/file/2",
		FileHash: "sha256:222222",
	}
	s3 = &Soundbite{
		ID:       3,
		Name:     "test3",
		Username: "test_username_3",
		UserID:   "333333",
		FilePath: "/path/to/file/3",
		FileHash: "sha256:333333",
	}
)

func mockInsert(m SoundbiteModel, s *Soundbite) (int, error) {
	return m.Insert(s.Name, s.Username, s.UserID, s.FilePath, s.FileHash)
}

func TestInsert(t *testing.T) {
	t.Run("insert single soundbite", func(t *testing.T) {
		m, teardown := modelsTestSetup(t)
		defer teardown()

		id, err := mockInsert(m, s1)

		test.AssertType(t, id, 1)
		test.AssertError(t, err, nil)
	})

	t.Run("insert multiple soundbites", func(t *testing.T) {
		m, teardown := modelsTestSetup(t)
		defer teardown()

		id, err := mockInsert(m, s1)
		test.AssertType(t, id, 1)
		test.AssertError(t, err, nil)

		id, err = mockInsert(m, s2)
		test.AssertType(t, id, 2)
		test.AssertError(t, err, nil)
	})

	t.Run("insert duplicate soundbites", func(t *testing.T) {
		m, teardown := modelsTestSetup(t)
		defer teardown()

		id, err := mockInsert(m, s1)
		test.AssertType(t, id, 1)
		test.AssertError(t, err, nil)

		id, err = mockInsert(m, s1)
		test.AssertType(t, id, 0)
		test.AssertError(t, err, ErrUniqueConstraint)
	})
}

func TestGet(t *testing.T) {
	t.Run("get single record", func(t *testing.T) {
		m, teardown := modelsTestSetup(t)
		defer teardown()
		_, _ = mockInsert(m, s1)
		s, err := m.Get(s1.Name)
		s.Created = time.Time{}

		test.AssertError(t, err, nil)
		test.AssertType(t, s, s1)
	})

	t.Run("get multiple records", func(t *testing.T) {
		m, teardown := modelsTestSetup(t)
		defer teardown()

		_, err := mockInsert(m, s1)
		test.AssertError(t, err, nil)

		_, err = mockInsert(m, s2)
		test.AssertError(t, err, nil)

		b1, err := m.Get(s1.Name)
		b1.Created = time.Time{}
		test.AssertError(t, err, nil)
		test.AssertType(t, b1, s1)

		b2, err := m.Get(s2.Name)
		b2.Created = time.Time{}
		test.AssertError(t, err, nil)
		test.AssertType(t, b2, s2)
	})

	t.Run("get not-existent record", func(t *testing.T) {
		m, teardown := modelsTestSetup(t)
		defer teardown()

		_, err := mockInsert(m, s1)
		test.AssertError(t, err, nil)

		_, err = m.Get(s2.Name)
		test.AssertError(t, err, ErrDoesNotExist)
	})
}

func TestGetAll(t *testing.T) {
	t.Run("get all from empty table", func(t *testing.T) {
		m, teardown := modelsTestSetup(t)
		defer teardown()

		sounds, err := m.GetAll()
		test.AssertError(t, err, ErrNoRecords)

		if len(sounds) != 0 {
			t.Fatalf("got: %v, expected: %v", len(sounds), 0)
		}
	})

	t.Run("get all from populated table", func(t *testing.T) {
		m, teardown := modelsTestSetup(t)
		defer teardown()

		expected := []*Soundbite{s1, s2, s3}

		_, _ = mockInsert(m, s1)
		_, _ = mockInsert(m, s2)
		_, _ = mockInsert(m, s3)

		sounds, err := m.GetAll()
		for _, s := range sounds {
			s.Created = time.Time{}
		}

		test.AssertError(t, err, nil)
		test.AssertType(t, sounds, expected)
	})
}

func TestExists(t *testing.T) {
	t.Run("exists in an empty table", func(t *testing.T) {
		m, teardown := modelsTestSetup(t)
		defer teardown()

		b, err := m.Exists(s1.Name, s1.FileHash)
		test.AssertError(t, err, nil)
		test.AssertType(t, b, false)
	})

	t.Run("exists in a table", func(t *testing.T) {
		m, teardown := modelsTestSetup(t)
		defer teardown()

		_, _ = mockInsert(m, s1)

		b, err := m.Exists(s1.Name, s1.FileHash)
		test.AssertError(t, err, nil)
		test.AssertType(t, b, true)
	})

	t.Run("exists but incorrect name and filehash", func(t *testing.T) {
		m, teardown := modelsTestSetup(t)
		defer teardown()

		_, _ = mockInsert(m, s1)

		b, err := m.Exists(s2.Name, s2.FileHash)
		test.AssertError(t, err, nil)
		test.AssertType(t, b, false)

		b, err = m.Exists(s1.Name, s1.FileHash)
		test.AssertError(t, err, nil)
		test.AssertType(t, b, true)
	})
}

func TestDelete(t *testing.T) {
	t.Run("delete from empty table", func(t *testing.T) {
		m, teardown := modelsTestSetup(t)
		defer teardown()

		err := m.Delete(s1.Name, s1.UserID)
		test.AssertError(t, err, ErrDoesNotExist)
	})

	t.Run("delete from table", func(t *testing.T) {
		m, teardown := modelsTestSetup(t)
		defer teardown()

		_, _ = mockInsert(m, s1)

		err := m.Delete(s1.Name, s1.UserID)
		test.AssertError(t, err, nil)

		sounds, err := m.GetAll()
		test.AssertError(t, err, ErrNoRecords)
		if len(sounds) != 0 {
			t.Fatalf("got: %v, expected: %v", len(sounds), 0)
		}
	})
}

func TestUserCreatedSound(t *testing.T) {
	t.Run("test user created sound", func(t *testing.T) {
		m, teardown := modelsTestSetup(t)
		defer teardown()

		_, _ = mockInsert(m, s1)

		err := m.userCreatedSound(s1.Name, s1.UserID)
		test.AssertError(t, err, nil)

		err = m.userCreatedSound(s1.Name, s2.UserID)
		test.AssertError(t, err, ErrCommandOwnership)
	})
}

func updateNameTestFunc(t *testing.T, oldName, newName string, expectedErr error) {
	m, teardown := modelsTestSetup(t)
	defer teardown()

	_, _ = mockInsert(m, s1)
	_, _ = mockInsert(m, s2)

	err := m.UpdateName(oldName, newName)

	test.AssertError(t, err, expectedErr)
}

func TestUpdateName(t *testing.T) {
	tt := []struct {
		description string
		oldName     string
		newName     string
		expectedErr error
		testFunc    func(*testing.T, string, string, error)
	}{
		{
			description: "update single soundbite",
			oldName:     s1.Name,
			newName:     s3.Name,
			expectedErr: nil,
			testFunc:    updateNameTestFunc,
		},
		{
			description: "update soundbite with existing name",
			oldName:     s1.Name,
			newName:     s2.Name,
			expectedErr: ErrUniqueConstraint,
			testFunc:    updateNameTestFunc,
		},
		{
			description: "update non-existent soundbite",
			oldName:     s3.Name,
			newName:     s2.Name,
			expectedErr: ErrUniqueConstraint,
			testFunc:    updateNameTestFunc,
		},
	}

	for _, tc := range tt {
		t.Run(tc.description, func(t *testing.T) {
			tc.testFunc(t, tc.oldName, tc.newName, tc.expectedErr)
		})
	}
}
