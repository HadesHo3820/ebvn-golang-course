package repository

import (
	"testing"

	"github.com/HadesHo3820/ebvn-golang-course/internal/model"
	"github.com/HadesHo3820/ebvn-golang-course/internal/test/fixture"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// TestUser_CreateUser tests the CreateUser method of the User repository.
// It uses table-driven tests with the UserCommonTestDB fixture to verify:
//   - Successful user creation with valid input
//   - Error handling when attempting to create a user with duplicate unique fields (email)
//
// Each test case includes a verifyFunc to confirm the user was correctly persisted in the database.
func TestUser_CreateUser(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name              string
		setupDB           func(t *testing.T) *gorm.DB
		inputUser         *model.User
		expectedErrString string
		expectedOutput    *model.User
		verifyFunc        func(db *gorm.DB, user *model.User)
	}{
		{
			name: "normal case",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixture.NewFixture(t, &fixture.UserCommonTestDB{})
			},
			inputUser: &model.User{
				ID:          "229ac10b-58cc-4372-a567-0e02b2c3d479",
				DisplayName: "Johnny Ho",
				Username:    "johnny.ho1",
				Email:       "johnny.ho1@example.com",
				Password:    "$2a$$2a$10$wfpS7JvQgcHvvHLk86eFs.jhKCIucgr9fhPkyBLVQntSH0nB05106$wfpS23sf",
			},
			expectedOutput: &model.User{
				ID:          "229ac10b-58cc-4372-a567-0e02b2c3d479",
				DisplayName: "Johnny Ho",
				Username:    "johnny.ho1",
				Email:       "johnny.ho1@example.com",
				Password:    "$2a$$2a$10$wfpS7JvQgcHvvHLk86eFs.jhKCIucgr9fhPkyBLVQntSH0nB05106$wfpS23sf",
			},
			verifyFunc: func(db *gorm.DB, user *model.User) {
				checkUser := &model.User{}
				err := db.Where("username = ?", user.Username).First(checkUser).Error
				assert.Nil(t, err)
				assert.Equal(t, checkUser, user)
			},
		},
		{
			name: "err case - username already exists",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixture.NewFixture(t, &fixture.UserCommonTestDB{})
			},
			inputUser: &model.User{
				ID:          "f47ac10b-58cc-4372-a567-0e02b2c3d479",
				DisplayName: "Johnny Ho",
				Username:    "johnny.ho",
				Email:       "johnny.ho@example.com",
				Password:    "$2a$$2a$10$wfpS7JvQgcHvvHLk86eFs.jhKCIucgr9fhPkyBLVQntSH0nB05106$wfpS23sf",
			},
			expectedErrString: "UNIQUE constraint failed: users.email",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()

			// Setup
			db := tc.setupDB(t)
			userRepo := NewUser(db)

			// Execute
			output, err := userRepo.CreateUser(ctx, tc.inputUser)

			// Assert
			if err != nil {
				assert.ErrorContains(t, err, tc.expectedErrString)
			}

			assert.Equal(t, tc.expectedOutput, output)

			if err == nil {
				tc.verifyFunc(db, tc.expectedOutput)
			}
		})
	}

}
