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
		test.AssertError(t, err, ErrNoRecords)
	})
}


