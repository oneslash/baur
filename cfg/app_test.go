package cfg

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ExampleApp_IsValid(t *testing.T) {
	a := ExampleApp("shop")
	if err := a.Validate(); err != nil {
		t.Error("example app conf fails validation: ", err)
	}
}

func Test_ExampleApp_WrittenAndReadCfgIsValid(t *testing.T) {
	tmpfileFD, err := ioutil.TempFile("", "baur")
	if err != nil {
		t.Fatal("opening tmpfile failed: ", err)
	}

	tmpfileName := tmpfileFD.Name()
	tmpfileFD.Close()
	os.Remove(tmpfileName)

	a := ExampleApp("shop")
	if err := a.Validate(); err != nil {
		t.Error("example conf fails validation: ", err)
	}

	if err := a.ToFile(tmpfileName); err != nil {
		t.Fatal("writing conf to file failed: ", err)
	}

	rRead, err := AppFromFile(tmpfileName)
	if err != nil {
		t.Fatal("reading conf from file failed: ", err)
	}

	if err := rRead.Validate(); err != nil {
		t.Errorf("validating conf from file failed: %s\nFile Content: %+v", err, rRead)
	}
}

func Test_AppHasOneTaskDefinition(t *testing.T) {
	app := App{
		Name: "testapp",
	}

	err := app.Validate()
	assert.EqualError(t, err, "The Tasks section must define exactly 1 Task")

	app = App{
		Name: "testapp",
		Tasks: []*Task{
			&Task{},
			&Task{},
		},
	}
	err = app.Validate()
	assert.EqualError(t, err, "The Tasks section must define exactly 1 Task")
}

func Test_OnlyBuildTaskAllowed(t *testing.T) {
	testcases := []struct {
		taskName   string
		shouldFail bool
	}{
		{
			taskName:   "check",
			shouldFail: true,
		},
		{
			taskName:   "test",
			shouldFail: true,
		},
		{
			taskName:   "",
			shouldFail: true,
		},
		{
			taskName:   "build",
			shouldFail: false,
		},
	}

	for _, testcase := range testcases {
		t.Run(fmt.Sprintf("taskname %s", testcase.taskName), func(t *testing.T) {
			app := App{
				Name: "testapp",
				Tasks: []*Task{
					&Task{
						Name:    testcase.taskName,
						Command: "check",
						Input: &Input{
							Files: FileInputs{
								Paths: []string{"*.txt"},
							},
						},
						Output: &Output{
							File: []*FileOutput{
								{
									Path: "test.tar",
									FileCopy: FileCopy{
										Path: "/tmp/",
									},
								},
							},
						},
					},
				},
			}

			err := app.Validate()
			if testcase.shouldFail {
				require.Error(t, err)
				require.Contains(t, err.Error(), fmt.Sprintf("invalid task name"))
				return
			}

			require.NoError(t, err)
		})
	}
}
