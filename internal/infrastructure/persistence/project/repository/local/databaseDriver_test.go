package local

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/apenella/ransidble/internal/domain/core/entity"
	"github.com/apenella/ransidble/internal/infrastructure/logger"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestFind(t *testing.T) {

	databasePath := filepath.Join("fixtures", "persistence-project-database-local", "read")

	fs := afero.NewCopyOnWriteFs(
		afero.NewReadOnlyFs(
			afero.NewBasePathFs(afero.NewOsFs(), "../../../../../../test"),
		),
		afero.NewMemMapFs(),
	)

	recordGoodReference := "57b44cfa498dee3ead824dd01a895d68807974dbe63db347e4bfda89732ebdfc"

	tests := []struct {
		desc     string
		id       string
		driver   *DatabaseDriver
		expected *entity.Project
		err      error
	}{
		{
			desc:   "Testing finding a valid project from the database",
			id:     recordGoodReference,
			driver: NewDatabaseDriver(fs, databasePath, logger.NewFakeLogger()),
			expected: &entity.Project{
				Format:    "plain",
				Name:      "project-1",
				Reference: "test/projects/project-1",
				Storage:   "local",
			},
			err: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			result, err := test.driver.Find(test.id)
			assert.Equal(t, test.expected, result)
			assert.Equal(t, test.err, err)
		})
	}
}

func TestFindAll(t *testing.T) {
	databasePath := filepath.Join("fixtures", "persistence-project-database-local", "read-all", "valid")
	fs := afero.NewCopyOnWriteFs(
		afero.NewReadOnlyFs(
			afero.NewBasePathFs(afero.NewOsFs(), "../../../../../../test"),
		),
		afero.NewMemMapFs(),
	)

	tests := []struct {
		desc     string
		driver   *DatabaseDriver
		expected []*entity.Project
		err      error
	}{
		{
			desc:   "Testing reading all projects from the database",
			driver: NewDatabaseDriver(fs, databasePath, logger.NewFakeLogger()),
			expected: []*entity.Project{
				{
					Format:    "plain",
					Name:      "project-1",
					Reference: "project-1",
					Storage:   "local",
				},
				{
					Format:    "plain",
					Name:      "project-2",
					Reference: "project-2",
					Storage:   "local",
				},
				{
					Format:    "plain",
					Name:      "project-3",
					Reference: "project-3",
					Storage:   "local",
				},
			},
			err: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			result, err := test.driver.FindAll()

			assert.Equal(t, len(test.expected), len(result))
			assert.Equal(t, test.err, err)
		})
	}
}

func TestStore(t *testing.T) {
	databasePath := filepath.Join("fixtures", "persistence-project-database-local")
	fs := afero.NewCopyOnWriteFs(
		afero.NewReadOnlyFs(
			afero.NewBasePathFs(afero.NewOsFs(), "../../../../../../test"),
		),
		afero.NewMemMapFs(),
	)

	tests := []struct {
		desc       string
		id         string
		data       *entity.Project
		driver     *DatabaseDriver
		assertFunc func(*testing.T, *DatabaseDriver)
		err        error
	}{
		{
			desc: "Testing store a valid project",
			id:   "project-2",
			data: &entity.Project{
				Name:      "project-2",
				Reference: "project-2",
				Format:    "plain",
				Storage:   "local",
			},
			driver: NewDatabaseDriver(fs, databasePath, logger.NewFakeLogger()),
			assertFunc: func(t *testing.T, driver *DatabaseDriver) {
				_, err := driver.fs.Stat(filepath.Join(databasePath, "project-2"))
				assert.Nil(t, err, "unexpected error when checking stored file")
				expected := &entity.Project{
					Name:      "project-2",
					Reference: "project-2",
					Format:    "plain",
					Storage:   "local",
				}
				project, err := driver.read("project-2")
				assert.Nil(t, err, "unexpected error when reading project")
				assert.Equal(t, expected, project)
			},
			err: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)
			err := test.driver.Store(test.id, test.data)
			if err != nil && test.err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.Nil(t, err, "unexpected error occurred")
				assert.Nil(t, test.err, "an expected error did not occur")
				if test.assertFunc != nil {
					test.assertFunc(t, test.driver)
				}
			}
		})
	}
}

func TestSafeStore(t *testing.T) {
	databasePath := filepath.Join("fixtures", "persistence-project-database-local", "exists")

	fs := afero.NewCopyOnWriteFs(
		afero.NewReadOnlyFs(
			afero.NewBasePathFs(afero.NewOsFs(), "../../../../../../test"),
		),
		afero.NewMemMapFs(),
	)

	tests := []struct {
		desc       string
		id         string
		data       *entity.Project
		driver     *DatabaseDriver
		assertFunc func(*testing.T, *DatabaseDriver)
		err        error
	}{
		{
			desc: "Testing storing a valid project (does not exist)",
			id:   "project-2",
			data: &entity.Project{
				Name:      "project-2",
				Reference: "project-2",
				Format:    "plain",
				Storage:   "local",
			},
			driver: NewDatabaseDriver(fs, databasePath, logger.NewFakeLogger()),
			assertFunc: func(t *testing.T, driver *DatabaseDriver) {
				project, err := driver.read("project-2")
				assert.Nil(t, err, "unexpected error when reading project")
				expected := &entity.Project{
					Name:      "project-2",
					Reference: "project-2",
					Format:    "plain",
					Storage:   "local",
				}
				assert.Equal(t, expected, project)
			},
			err: nil,
		},
		{
			desc:       "Testing error when storing an existing project (already exists)",
			id:         "project-1",
			data:       &entity.Project{},
			driver:     NewDatabaseDriver(fs, databasePath, logger.NewFakeLogger()),
			assertFunc: nil,
			err:        fmt.Errorf("%s: %s %s", ErrStoringProject, "project-1", ErrProjectExists),
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			err := test.driver.SafeStore(test.id, test.data)
			if test.err != nil {
				assert.NotNil(t, err, "expected error but got nil")
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.Nil(t, err, "unexpected error occurred")
				if test.assertFunc != nil {
					test.assertFunc(t, test.driver)
				}
			}
		})
	}
}

func TestDelete(t *testing.T) {
	databasePath := filepath.Join("fixtures", "persistence-project-database-local", "remove")

	fs := afero.NewCopyOnWriteFs(
		afero.NewReadOnlyFs(
			afero.NewBasePathFs(afero.NewOsFs(), "../../../../../../test"),
		),
		afero.NewMemMapFs(),
	)

	tests := []struct {
		desc        string
		id          string
		driver      *DatabaseDriver
		arrangeFunc func(*DatabaseDriver)
		assertFunc  func(*testing.T, *DatabaseDriver)
		err         error
	}{
		{
			desc:   "Testing removing a record from the database",
			id:     "project-1",
			driver: NewDatabaseDriver(fs, databasePath, logger.NewFakeLogger()),
			arrangeFunc: func(d *DatabaseDriver) {
				d.fs.OpenFile(filepath.Join(databasePath, "project-1"), os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
			},
			assertFunc: func(t *testing.T, driver *DatabaseDriver) {

				ok := false

				_, err := driver.fs.Stat(filepath.Join(databasePath, "project-1"))
				if err != nil && os.IsNotExist(err) {
					ok = true
				}

				assert.True(t, ok, "the record was not removed from the database")
			},
			err: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {

			if test.arrangeFunc != nil {
				test.arrangeFunc(test.driver)
			}

			err := test.driver.Delete(test.id)
			if test.err != nil {
				assert.NotNil(t, err, "expected error but got nil")
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.Nil(t, err, "unexpected error occurred")
				if test.assertFunc != nil {
					test.assertFunc(t, test.driver)
				}
			}
		})
	}
}

func TestRead(t *testing.T) {

	// Arranging the database with the records needed for the test

	databasePath := filepath.Join("fixtures", "persistence-project-database-local", "read")

	fs := afero.NewCopyOnWriteFs(
		afero.NewReadOnlyFs(
			afero.NewBasePathFs(afero.NewOsFs(), "../../../../../../test"),
		),
		afero.NewMemMapFs(),
	)

	// Arranging a record that can be read and returns a project. recordGoodReference has been generated using the generaterecordReference function. It value is the sha256 hash of "recordGood"
	recordGoodReference := "57b44cfa498dee3ead824dd01a895d68807974dbe63db347e4bfda89732ebdfc"

	// Arranging a record that can be read but returns an error because the hash is invalid. recordGoodReference has been generated using the generaterecordReference function. It value is the sha256 hash of "recordInvalidHash"
	recordInvalidHashReference := "4341e480935b1f271c71005026c54a4d012f2c3babf39c81b0b6139c0d5eacf0"

	// Arranging a record that has a wrong JSON formated. recordInvalidJSONReference has been generated using the generaterecordReference function. It value is the sha256 hash of "recordInvalidJSON"
	recordInvalidJSONReference := "7e5e5388886deab9e1a7bf3f98a518dbd84b62017f540702e090bec7bc011946"

	// Arranging a record that has a wrong JSON formated project. recordInvalidProjectJSONReference has been generated using the generaterecordReference function. It value is the sha256 hash of "recordInvalidProjectJSON"
	recordInvalidProjectJSONReference := "ea54a981bccdab7286b6c2222d6e3da7bc9fd09d181827b0b4a521d0c1455e3c"

	// Arranging a record that is a directory. recordInvalidRecordDirectoryReference has been generated using the generaterecordReference function. It value is the sha256 hash of "invalidRecordDirectory"
	recordInvalidRecordDirectoryReference := "5cbe38cf0619b09b0c46f1aee4cb9c3d3ef7e19a2e2265d4c5571f50fde23203"

	tests := []struct {
		desc     string
		id       string
		driver   *DatabaseDriver
		expected *entity.Project
		err      error
	}{
		{
			desc:   "Testing reading a valid record from the database that returns a project",
			id:     recordGoodReference,
			driver: NewDatabaseDriver(fs, databasePath, logger.NewFakeLogger()),
			expected: &entity.Project{
				Format:    "plain",
				Name:      "project-1",
				Reference: "test/projects/project-1",
				Storage:   "local",
			},
			err: nil,
		},
		{
			desc:     "Testing error reading a record from the database that returns an error because the hash is invalid",
			id:       recordInvalidHashReference,
			driver:   NewDatabaseDriver(fs, databasePath, logger.NewFakeLogger()),
			expected: nil,
			err:      fmt.Errorf("%s: %s", ErrVerifyingRecord, ErrVerifyingRecordInvalidHash),
		},
		{
			desc:     "Testing error reading a record from the database when id is not provided",
			id:       "",
			driver:   NewDatabaseDriver(fs, databasePath, logger.NewFakeLogger()),
			expected: nil,
			err:      fmt.Errorf("%s", ErrIDIsNotProvided),
		},
		{
			desc:     "Testing error reading a record from the database when filesystem is not initialized",
			id:       recordGoodReference,
			driver:   NewDatabaseDriver(nil, databasePath, logger.NewFakeLogger()),
			expected: nil,
			err:      fmt.Errorf("%s", ErrFilesystemNotInitialized),
		},
		{
			desc:     "Testing error reading a record from the database when database path is not initialized",
			id:       recordGoodReference,
			driver:   NewDatabaseDriver(fs, "", logger.NewFakeLogger()),
			expected: nil,
			err:      fmt.Errorf("%s", ErrDatabasePathNotInitialized),
		},
		{
			desc:     "Testing error reading a record from the database when the record is not found",
			id:       "notfound",
			driver:   NewDatabaseDriver(fs, databasePath, logger.NewFakeLogger()),
			expected: nil,
			err:      fmt.Errorf("error reading record: record not found: stat ../../../../../../test/fixtures/persistence-project-database-local/read/notfound: no such file or directory"),
		},
		{
			desc:     "Testing error reading a record from the database when the record file is a directory",
			id:       recordInvalidRecordDirectoryReference,
			driver:   NewDatabaseDriver(fs, databasePath, logger.NewFakeLogger()),
			expected: nil,
			err:      fmt.Errorf("%s: %s", ErrInvalidRecordFormat, ErrInvalidRecordFormatIsDir),
		},
		{
			desc:     "Testing error reading a record from the database when the record file is not a valid JSON",
			id:       recordInvalidJSONReference,
			driver:   NewDatabaseDriver(fs, databasePath, logger.NewFakeLogger()),
			expected: nil,
			// this test case deponds on the error return json.Unmarshal function
			err: fmt.Errorf("%s: %s", ErrExtractingRecordFromDatabase, "unexpected end of JSON input"),
		},
		{
			desc:     "Testing error reading a record from the database when the record contains a invalid project JSON",
			id:       recordInvalidProjectJSONReference,
			driver:   NewDatabaseDriver(fs, databasePath, logger.NewFakeLogger()),
			expected: nil,
			// this test case deponds on the error return json.Unmarshal function
			err: fmt.Errorf("%s: %s", ErrExtractingProjectFromRecord, "json: cannot unmarshal string into Go value of type entity.Project"),
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			project, err := test.driver.read(test.id)
			if err != nil && test.err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.Nil(t, err, "occured an unexpected error")
				assert.Nil(t, test.err, "an expected error not occurred")
				assert.Equal(t, test.expected, project)
			}
		})
	}
}

func TestReadAll(t *testing.T) {
	// Arranging the database with the records needed for the test

	databasePathValid := filepath.Join("fixtures", "persistence-project-database-local", "read-all", "valid")
	databasePathInvalidHash := filepath.Join("fixtures", "persistence-project-database-local", "read-all", "invalid-hash")

	fs := afero.NewCopyOnWriteFs(
		afero.NewReadOnlyFs(
			afero.NewBasePathFs(afero.NewOsFs(), "../../../../../../test"),
		),
		afero.NewMemMapFs(),
	)

	tests := []struct {
		desc     string
		driver   *DatabaseDriver
		expected []*entity.Project
		err      error
	}{
		{
			desc:   "Testing reading all records from the database",
			driver: NewDatabaseDriver(fs, databasePathValid, logger.NewFakeLogger()),
			expected: []*entity.Project{
				{
					Name:      "project-1",
					Reference: "project-1",
					Format:    "plain",
					Storage:   "local",
				},
				{
					Name:      "project-2",
					Reference: "project-2",
					Format:    "plain",
					Storage:   "local",
				},
				{
					Name:      "project-3",
					Reference: "project-3",
					Format:    "plain",
					Storage:   "local",
				},
			},
			err: nil,
		},
		{
			desc:     "Testing error reading all records from the database when the filesystem is not initialized",
			driver:   NewDatabaseDriver(nil, databasePathValid, logger.NewFakeLogger()),
			expected: nil,
			err:      fmt.Errorf("%s", ErrFilesystemNotInitialized),
		},
		{
			desc:     "Testing error reading all records from the database when the database path is not initialized",
			driver:   NewDatabaseDriver(fs, "", logger.NewFakeLogger()),
			expected: nil,
			err:      fmt.Errorf("%s", ErrDatabasePathNotInitialized),
		},
		{
			desc:   "Testing error reading all records from the database when the records have an invalid hash",
			driver: NewDatabaseDriver(fs, databasePathInvalidHash, logger.NewFakeLogger()),
			expected: []*entity.Project{
				{
					Name:      "project-1",
					Reference: "project-1",
					Format:    "plain",
					Storage:   "local",
				},
			},
			err: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			projectList, err := test.driver.readAll()
			if err != nil && test.err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.Nil(t, err, "occured an unexpected error")
				assert.Nil(t, test.err, "an expected error not occurred")
				assert.Equal(t, len(test.expected), len(projectList))

				switch len(test.expected) {
				case 1:
					assert.Contains(t, projectList, test.expected[0])
				case 2:
					assert.Contains(t, projectList, test.expected[0])
					assert.Contains(t, projectList, test.expected[1])
				case 3:
					assert.Contains(t, projectList, test.expected[0])
					assert.Contains(t, projectList, test.expected[1])
					assert.Contains(t, projectList, test.expected[2])
				}
			}
		})
	}
}

func TestWrite(t *testing.T) {
	// Arranging the database with the records needed for the test

	databasePath := filepath.Join("fixtures", "persistence-project-database-local")
	databasePathNotUnexist := filepath.Join("fixtures", "persistence-project-database-local", "unexist")

	fs := afero.NewCopyOnWriteFs(
		afero.NewReadOnlyFs(
			afero.NewBasePathFs(afero.NewOsFs(), "../../../../../../test"),
		),
		afero.NewMemMapFs(),
	)

	tests := []struct {
		desc       string
		id         string
		data       *entity.Project
		driver     *DatabaseDriver
		assertFunc func(*testing.T, *DatabaseDriver)
		err        error
	}{
		{
			desc: "Testing write a valid record to the database",
			id:   "project-1",
			data: &entity.Project{
				Name:      "project-1",
				Reference: "project-1",
				Format:    "plain",
				Storage:   "local",
			},
			driver: NewDatabaseDriver(fs, databasePath, logger.NewFakeLogger()),
			assertFunc: func(t *testing.T, driver *DatabaseDriver) {
				_, err := driver.fs.Stat(filepath.Join(databasePath, "project-1"))
				assert.Nil(t, err, "occured an unexpected error")

				expected := &entity.Project{
					Name:      "project-1",
					Reference: "project-1",
					Format:    "plain",
					Storage:   "local",
				}

				project, err := driver.read("project-1")
				assert.Nil(t, err, "occured an unexpected error")
				assert.Equal(t, expected, project)
			},
			err: nil,
		},
		{
			desc: "Testing error writing a record to the database when the database path does not exist",
			id:   "project-1",
			data: &entity.Project{
				Name:      "project-1",
				Reference: "project-1",
				Format:    "plain",
				Storage:   "local",
			},
			driver:     NewDatabaseDriver(fs, databasePathNotUnexist, logger.NewFakeLogger()),
			assertFunc: nil,
			err:        fmt.Errorf("%s: %s", ErrDatabasePathMustExist, "stat ../../../../../../test/fixtures/persistence-project-database-local/unexist: no such file or directory"),
		},
		{
			desc: "Testing error writing a record to the database when id is not provided",
			id:   "",
			data: &entity.Project{
				Name:      "project-1",
				Reference: "project-1",
				Format:    "plain",
				Storage:   "local",
			},
			driver:     NewDatabaseDriver(fs, databasePath, logger.NewFakeLogger()),
			assertFunc: nil,
			err:        fmt.Errorf("%s", ErrIDIsNotProvided),
		},
		{
			desc:       "Testing error writing a record to the database when data is not provided",
			id:         "project-1",
			data:       nil,
			driver:     NewDatabaseDriver(fs, databasePath, logger.NewFakeLogger()),
			assertFunc: nil,
			err:        fmt.Errorf("%s", ErrDataToWriteIsNotProvided),
		},
		{
			desc:       "Testing error writing a record to the database when filesystem is not initialized",
			id:         "project-1",
			data:       &entity.Project{},
			driver:     NewDatabaseDriver(nil, databasePath, logger.NewFakeLogger()),
			assertFunc: nil,
			err:        fmt.Errorf("%s", ErrFilesystemNotInitialized),
		},
		{
			desc:       "Testing error writing a record to the database when database path is not initialized",
			id:         "project-1",
			data:       &entity.Project{},
			driver:     NewDatabaseDriver(fs, "", logger.NewFakeLogger()),
			assertFunc: nil,
			err:        fmt.Errorf("%s", ErrDatabasePathNotInitialized),
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			err := test.driver.write(test.id, test.data)
			if err != nil && test.err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.Nil(t, err, "occured an unexpected error")
				assert.Nil(t, test.err, "an expected error not occurred")

				if test.assertFunc != nil {
					test.assertFunc(t, test.driver)
				}
			}
		})
	}
}

func TestRemove(t *testing.T) {
	// Arranging the database with the records needed for the test

	databasePath := filepath.Join("fixtures", "persistence-project-database-local", "remove")
	fs := afero.NewCopyOnWriteFs(
		afero.NewReadOnlyFs(
			afero.NewBasePathFs(afero.NewOsFs(), "../../../../../../test"),
		),
		afero.NewMemMapFs(),
	)

	tests := []struct {
		desc        string
		id          string
		driver      *DatabaseDriver
		arrangeFunc func(*DatabaseDriver)
		assertFunc  func(*testing.T, *DatabaseDriver)
		err         error
	}{
		{
			desc:        "Testing error removing a record from the database when the id is not provided",
			id:          "",
			driver:      NewDatabaseDriver(fs, databasePath, logger.NewFakeLogger()),
			arrangeFunc: nil,
			assertFunc:  nil,
			err:         fmt.Errorf("%s", ErrIDIsNotProvided),
		},
		{
			desc:        "Testing error removing a record from the database when the filesystem is not initialized",
			id:          "project-1",
			driver:      NewDatabaseDriver(nil, databasePath, logger.NewFakeLogger()),
			arrangeFunc: nil,
			assertFunc:  nil,
			err:         fmt.Errorf("%s", ErrFilesystemNotInitialized),
		},
		{
			desc:        "Testing error removing a record from the database when the database path is not initialized",
			id:          "project-1",
			driver:      NewDatabaseDriver(fs, "", logger.NewFakeLogger()),
			arrangeFunc: nil,
			assertFunc:  nil,
			err:         fmt.Errorf("%s", ErrDatabasePathNotInitialized),
		},
		{
			desc:        "Testing removing a record from the database when the record is not found",
			id:          "notfound",
			driver:      NewDatabaseDriver(fs, databasePath, logger.NewFakeLogger()),
			arrangeFunc: nil,
			assertFunc:  nil,
			err:         nil,
		},
		{
			// This test case is coupled to the afero NewCopyOnWriteFs file system. That does not allow to remove a file when the file is created on the base file system.
			desc:        "Testing error removing a record from the database",
			id:          "abc",
			driver:      NewDatabaseDriver(fs, databasePath, logger.NewFakeLogger()),
			arrangeFunc: func(d *DatabaseDriver) {},
			assertFunc:  func(t *testing.T, driver *DatabaseDriver) {},
			err:         fmt.Errorf("%s: %s", ErrRemovingRecord, "remove fixtures/persistence-project-database-local/remove/abc: file does not exist"),
		},
		{
			desc:   "Testing removing a record from the database",
			id:     "project-1",
			driver: NewDatabaseDriver(fs, databasePath, logger.NewFakeLogger()),
			arrangeFunc: func(d *DatabaseDriver) {
				d.fs.OpenFile(filepath.Join(databasePath, "project-1"), os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
			},
			assertFunc: func(t *testing.T, driver *DatabaseDriver) {

				ok := false

				_, err := driver.fs.Stat(filepath.Join(databasePath, "project-1"))
				if err != nil && os.IsNotExist(err) {
					ok = true
				}

				assert.True(t, ok, "the record was not removed from the database")
			},
			err: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			if test.arrangeFunc != nil {
				test.arrangeFunc(test.driver)
			}

			err := test.driver.remove(test.id)
			if err != nil && test.err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.Nil(t, err, "occured an unexpected error")
				assert.Nil(t, test.err, "an expected error not occurred")

				if test.assertFunc != nil {
					test.assertFunc(t, test.driver)
				}
			}
		})
	}
}

func TestExists(t *testing.T) {
	// Arranging the database with the records needed for the test

	databasePath := filepath.Join("fixtures", "persistence-project-database-local", "exists")

	fs := afero.NewCopyOnWriteFs(
		afero.NewReadOnlyFs(
			afero.NewBasePathFs(afero.NewOsFs(), "../../../../../../test"),
		),
		afero.NewMemMapFs(),
	)

	tests := []struct {
		desc     string
		id       string
		driver   *DatabaseDriver
		expected bool
		err      error
	}{
		{
			desc:     "Testing error checking if a record exists in the database when the id is not provided",
			id:       "",
			driver:   NewDatabaseDriver(fs, databasePath, logger.NewFakeLogger()),
			expected: false,
			err:      fmt.Errorf("%s", ErrIDIsNotProvided),
		},
		{
			desc:     "Testing error checking if a record exists in the database when the filesystem is not initialized",
			id:       "project-1",
			driver:   NewDatabaseDriver(nil, databasePath, logger.NewFakeLogger()),
			expected: false,
			err:      fmt.Errorf("%s", ErrFilesystemNotInitialized),
		},
		{
			desc:     "Testing error checking if a record exists in the database when the database path is not initialized",
			id:       "project-1",
			driver:   NewDatabaseDriver(fs, "", logger.NewFakeLogger()),
			expected: false,
			err:      fmt.Errorf("%s", ErrDatabasePathNotInitialized),
		},
		{
			desc:     "Testing checking if a record exists in the database",
			id:       "project-1",
			driver:   NewDatabaseDriver(fs, databasePath, logger.NewFakeLogger()),
			expected: true,
			err:      nil,
		},
		{
			desc:     "Testing checking if a record does not exist in the database",
			id:       "notfound",
			driver:   NewDatabaseDriver(fs, databasePath, logger.NewFakeLogger()),
			expected: false,
			err:      nil,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			exists, err := test.driver.exists(test.id)
			if err != nil && test.err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.Nil(t, err, "occured an unexpected error")
				assert.Nil(t, test.err, "an expected error not occurred")
				assert.Equal(t, test.expected, exists)
			}
		})
	}
}
