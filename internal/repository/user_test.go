package repository

import (
	"testing"

	"github.com/HadesHo3820/ebvn-golang-course/internal/model"
	"github.com/HadesHo3820/ebvn-golang-course/internal/test/fixture"
	"github.com/HadesHo3820/ebvn-golang-course/pkg/dbutils"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// whereIDClause is the SQL WHERE clause for filtering by ID.
const whereIDClause = "id = ?"

// TestUser_CreateUser tests the CreateUser method of the User repository.
// It uses table-driven tests with the UserCommonTestDB fixture to verify:
//   - Successful user creation with valid input
//   - Error handling when attempting to create a user with duplicate unique fields (email)
//
// Each test case includes a verifyFunc to confirm the user was correctly persisted in the database.
func TestUser_CreateUser(t *testing.T) {
	t.Parallel()

	updatedAt := fixture.FixtureTimestamp
	createdAt := fixture.FixtureTimestamp

	testCases := []struct {
		name           string
		setupDB        func(t *testing.T) *gorm.DB
		inputUser      *model.User
		expectedErr    error
		expectedOutput *model.User
		verifyFunc     func(db *gorm.DB, user *model.User)
	}{
		{
			name: "normal case",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixture.NewFixture(t, &fixture.UserCommonTestDB{})
			},
			inputUser: &model.User{
				ID:          "229ac10b-58cc-4372-a567-0e02b2c3d479",
				DisplayName: fixture.FixtureUserOneDisplayName,
				Username:    "johnny.ho1",
				Email:       "johnny.ho1@example.com",
				Password:    fixture.FixtureUserPassword,
				CreatedAt:   createdAt,
				UpdatedAt:   updatedAt,
			},
			expectedOutput: &model.User{
				ID:          "229ac10b-58cc-4372-a567-0e02b2c3d479",
				DisplayName: fixture.FixtureUserOneDisplayName,
				Username:    "johnny.ho1",
				Email:       "johnny.ho1@example.com",
				Password:    fixture.FixtureUserPassword,
				CreatedAt:   createdAt,
				UpdatedAt:   updatedAt,
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
				ID:          fixture.FixtureUserOneID,
				DisplayName: fixture.FixtureUserOneDisplayName,
				Username:    fixture.FixtureUserOneUsername,
				Email:       fixture.FixtureUserOneEmail,
				Password:    fixture.FixtureUserPassword,
			},
			expectedErr: dbutils.ErrDuplicationType,
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
				assert.ErrorIs(t, err, tc.expectedErr)
			}

			assert.Equal(t, tc.expectedOutput, output)

			if err == nil {
				tc.verifyFunc(db, tc.expectedOutput)
			}
		})
	}

}

// TestUser_GetUserByUsername tests the GetUserByUsername method of the User repository.
// It uses table-driven tests with the UserCommonTestDB fixture to verify:
//   - Successful retrieval of an existing user by username
//   - Error handling when attempting to retrieve a non-existent user
func TestUser_GetUserByUsername(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		setupDB        func(t *testing.T) *gorm.DB
		inputUsername  string
		expectedErr    error
		expectedOutput *model.User
	}{
		{
			name: "success - user exists",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixture.NewFixture(t, &fixture.UserCommonTestDB{})
			},
			inputUsername: fixture.FixtureUserOneUsername,
			expectedOutput: &model.User{
				ID:          fixture.FixtureUserOneID,
				DisplayName: fixture.FixtureUserOneDisplayName,
				Username:    fixture.FixtureUserOneUsername,
				Email:       fixture.FixtureUserOneEmail,
				Password:    fixture.FixtureUserPassword,
				CreatedAt:   fixture.FixtureTimestamp,
				UpdatedAt:   fixture.FixtureTimestamp,
			},
		},
		{
			name: "error - user not found",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixture.NewFixture(t, &fixture.UserCommonTestDB{})
			},
			inputUsername: "nonexistent.user",
			expectedErr:   dbutils.ErrNotFoundType,
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
			output, err := userRepo.GetUserByUsername(ctx, tc.inputUsername)

			// Assert
			if tc.expectedErr != nil {
				assert.ErrorIs(t, err, tc.expectedErr)
				assert.Nil(t, output)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedOutput, output)
		})
	}
}

// TestUser_GetUserById tests the GetUserById method of the User repository.
// It uses table-driven tests with the UserCommonTestDB fixture to verify:
//   - Successful retrieval of an existing user by ID
//   - Error handling when attempting to retrieve a non-existent user
func TestUser_GetUserById(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		setupDB        func(t *testing.T) *gorm.DB
		inputUserID    string
		expectedErr    error
		expectedOutput *model.User
	}{
		{
			name: "success - user exists",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixture.NewFixture(t, &fixture.UserCommonTestDB{})
			},
			inputUserID: "322ac10b-58cc-4372-a567-0e02b2c3d479",
			expectedOutput: &model.User{
				ID:          "322ac10b-58cc-4372-a567-0e02b2c3d479",
				DisplayName: "Huy Ho",
				Username:    "huy.ho",
				Email:       "huy.ho@example.com",
				Password:    fixture.FixtureUserPassword,
				CreatedAt:   fixture.FixtureTimestamp,
				UpdatedAt:   fixture.FixtureTimestamp,
			},
		},
		{
			name: "error - user not found",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixture.NewFixture(t, &fixture.UserCommonTestDB{})
			},
			inputUserID: "00000000-0000-0000-0000-000000000000",
			expectedErr: dbutils.ErrNotFoundType,
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
			output, err := userRepo.GetUserById(ctx, tc.inputUserID)

			// Assert
			if tc.expectedErr != nil {
				assert.ErrorIs(t, err, tc.expectedErr)
				assert.Nil(t, output)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedOutput, output)
		})
	}
}

// TestUser_UpdateUser tests the UpdateUser method of the User repository.
// It uses table-driven tests with the UserCommonTestDB fixture to verify:
//   - Successful update of display_name only
//   - Successful update of email only
//   - Successful update of both display_name and email
//   - Behavior when updating a non-existent user (no error returned by GORM for 0 rows affected)
func TestUser_UpdateUser(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name             string
		setupDB          func(t *testing.T) *gorm.DB
		inputUserID      string
		inputDisplayName string
		inputEmail       string
		expectedErr      error
		verifyFunc       func(t *testing.T, db *gorm.DB, userID string)
	}{
		{
			name: "success - update display_name only",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixture.NewFixture(t, &fixture.UserCommonTestDB{})
			},
			inputUserID:      fixture.FixtureUserOneID,
			inputDisplayName: "Johnny Ho Updated",
			inputEmail:       "",
			verifyFunc: func(t *testing.T, db *gorm.DB, userID string) {
				var user model.User
				err := db.Where(whereIDClause, userID).First(&user).Error
				assert.NoError(t, err)
				assert.Equal(t, "Johnny Ho Updated", user.DisplayName)
				assert.Equal(t, fixture.FixtureUserOneEmail, user.Email) // Email unchanged
			},
		},
		{
			name: "success - update email only",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixture.NewFixture(t, &fixture.UserCommonTestDB{})
			},
			inputUserID:      fixture.FixtureUserOneID,
			inputDisplayName: "",
			inputEmail:       "johnny.ho.updated@example.com",
			verifyFunc: func(t *testing.T, db *gorm.DB, userID string) {
				var user model.User
				err := db.Where(whereIDClause, userID).First(&user).Error
				assert.NoError(t, err)
				assert.Equal(t, fixture.FixtureUserOneDisplayName, user.DisplayName) // DisplayName unchanged
				assert.Equal(t, "johnny.ho.updated@example.com", user.Email)
			},
		},
		{
			name: "success - update both display_name and email",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixture.NewFixture(t, &fixture.UserCommonTestDB{})
			},
			inputUserID:      fixture.FixtureUserTwoID,
			inputDisplayName: "Huy Ho Updated",
			inputEmail:       "huy.ho.updated@example.com",
			verifyFunc: func(t *testing.T, db *gorm.DB, userID string) {
				var user model.User
				err := db.Where(whereIDClause, userID).First(&user).Error
				assert.NoError(t, err)
				assert.Equal(t, "Huy Ho Updated", user.DisplayName)
				assert.Equal(t, "huy.ho.updated@example.com", user.Email)
			},
		},
		{
			name: "success - update non-existent user (no error, 0 rows affected)",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixture.NewFixture(t, &fixture.UserCommonTestDB{})
			},
			inputUserID:      "00000000-0000-0000-0000-000000000000",
			inputDisplayName: "Ghost User",
			inputEmail:       "ghost@example.com",
			verifyFunc: func(t *testing.T, db *gorm.DB, userID string) {
				// Verify user was not created
				var user model.User
				err := db.Where(whereIDClause, userID).First(&user).Error
				assert.ErrorIs(t, err, gorm.ErrRecordNotFound) // Should not find the user
			},
		},
		{
			name: "error - duplicate email",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixture.NewFixture(t, &fixture.UserCommonTestDB{})
			},
			inputUserID:      fixture.FixtureUserOneID,
			inputDisplayName: "",
			inputEmail:       fixture.FixtureUserTwoEmail, // Already taken by another user
			expectedErr:      dbutils.ErrDuplicationType,
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
			err := userRepo.UpdateUser(ctx, tc.inputUserID, tc.inputDisplayName, tc.inputEmail)

			// Assert
			if tc.expectedErr != nil {
				assert.ErrorIs(t, err, tc.expectedErr)
				return
			}

			assert.NoError(t, err)
			if tc.verifyFunc != nil {
				tc.verifyFunc(t, db, tc.inputUserID)
			}
		})
	}
}
