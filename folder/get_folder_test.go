package folder_test

import (
	"errors"
	"github.com/georgechieng-sc/interns-2022/folder"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

// This data is provided in the spec for component 1
// i changed the orgId field to make it actually compile
func GetTestingSampleData() []folder.Folder {
	return []folder.Folder{
		{
			Name:  "alpha",
			OrgId: uuid.FromStringOrNil("38b9879b-f73b-4b0e-b9d9-4fc4c23643a7"),
			Paths: "alpha",
		},
		{
			Name:  "bravo",
			OrgId: uuid.FromStringOrNil("38b9879b-f73b-4b0e-b9d9-4fc4c23643a7"),
			Paths: "alpha.bravo",
		},
		{
			Name:  "charlie",
			OrgId: uuid.FromStringOrNil("38b9879b-f73b-4b0e-b9d9-4fc4c23643a7"),
			Paths: "alpha.bravo.charlie",
		},
		{
			Name:  "delta",
			OrgId: uuid.FromStringOrNil("38b9879b-f73b-4b0e-b9d9-4fc4c23643a7"),
			Paths: "alpha.delta",
		},
		{
			Name:  "echo",
			OrgId: uuid.FromStringOrNil("38b9879b-f73b-4b0e-b9d9-4fc4c23643a7"),
			Paths: "echo",
		},
		{
			Name:  "foxtrot",
			OrgId: uuid.FromStringOrNil("c1556e17-b7c0-45a3-a6ae-9546248fb17a"),
			Paths: "foxtrot",
		},
	}
}

// feel free to change how the unit test is structured
func Test_folder_GetFoldersByOrgID(t *testing.T) {
	t.Parallel()
	// reuse to save performace
	orgID := uuid.FromStringOrNil("38b9879b-f73b-4b0e-b9d9-4fc4c23643a7")
	tests := [...]struct {
		name    string
		orgID   uuid.UUID
		folders []folder.Folder
		want    []folder.Folder
	}{
		{
			name:    "One org which returns multiple folders",
			orgID:   orgID,
			folders: GetTestingSampleData(),
			want: []folder.Folder{
				{Name: "alpha", OrgId: orgID, Paths: "alpha"},
				{Name: "bravo", OrgId: orgID, Paths: "alpha.bravo"},
				{Name: "charlie", OrgId: orgID, Paths: "alpha.bravo.charlie"},
				{Name: "delta", OrgId: orgID, Paths: "alpha.delta"},
				{Name: "echo", OrgId: orgID, Paths: "echo"},
			},
		},
		{
			name:    "One org which returns one folder",
			orgID:   uuid.FromStringOrNil("c1556e17-b7c0-45a3-a6ae-9546248fb17a"),
			folders: GetTestingSampleData(),
			want: []folder.Folder{
				{Name: "foxtrot", OrgId: uuid.FromStringOrNil("c1556e17-b7c0-45a3-a6ae-9546248fb17a"), Paths: "foxtrot"},
			},
		},
		{
			name:    "Wrong org returns no folders",
			orgID:   uuid.FromStringOrNil("c1wrong7-b7c0-45a3-a6ae-9546248fb17a"),
			folders: GetTestingSampleData(),
			want:    []folder.Folder{},
		},
		{
			name:    "Empty org string returns nothing",
			orgID:   uuid.FromStringOrNil(""),
			folders: GetTestingSampleData(),
			want:    []folder.Folder{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := folder.NewDriver(tt.folders)
			got := f.GetFoldersByOrgID(tt.orgID)

			assert.Equal(t, tt.want, got, "The expected output doesnt match")
		})
	}
}

func Test_folder_GetAllChildFolders(t *testing.T) {
	t.Parallel()
	orgID := uuid.FromStringOrNil("38b9879b-f73b-4b0e-b9d9-4fc4c23643a7")
	tests := [...]struct {
		name_of_test   string
		orgID          uuid.UUID
		name_of_folder string
		folders        []folder.Folder
		want           []folder.Folder
		wantErr        error
	}{
		{
			name_of_test:   "Gets all child folders from org and name",
			orgID:          orgID,
			name_of_folder: "alpha",
			folders:        GetTestingSampleData(),
			want: []folder.Folder{
				{Name: "bravo", OrgId: orgID, Paths: "alpha.bravo"},
				{Name: "charlie", OrgId: orgID, Paths: "alpha.bravo.charlie"},
				{Name: "delta", OrgId: orgID, Paths: "alpha.delta"},
			},
			wantErr: nil,
		},
		{
			name_of_test:   "Gets all child folders of a child folder",
			orgID:          orgID,
			name_of_folder: "bravo",
			folders:        GetTestingSampleData(),
			want: []folder.Folder{
				{
					Name:  "charlie",
					OrgId: orgID,
					Paths: "alpha.bravo.charlie",
				},
			},
			wantErr: nil,
		},
		{
			name_of_test:   "Given folder has no child folders",
			orgID:          orgID,
			name_of_folder: "charlie",
			folders:        GetTestingSampleData(),
			want:           []folder.Folder{},
			wantErr:        nil,
		},
		{
			name_of_test:   "Given folder doesn't exist",
			orgID:          orgID,
			name_of_folder: "invalid_folder",
			folders:        GetTestingSampleData(),
			want:           []folder.Folder{},
			wantErr:        errors.New("Error: Folder does not exist"),
		},
		{
			name_of_test:   "Given folder does not exist in org",
			orgID:          orgID,
			name_of_folder: "echo",
			folders:        GetTestingSampleData(),
			want:           []folder.Folder{},
			wantErr:        errors.New("Error: Folder does not exist in the specified organization"),
		},
		{
			name_of_test:   "Given folder os an empty string",
			orgID:          orgID,
			name_of_folder: "",
			folders:        GetTestingSampleData(),
			want:           []folder.Folder{},
			wantErr:        errors.New("Error: Folder does not exist"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name_of_test, func(t *testing.T) {
			f := folder.NewDriver(tt.folders)
			got, err := f.GetAllChildFolders(tt.orgID, tt.name_of_folder)

			if err != nil && tt.wantErr != nil {
				assert.Equal(t, tt.wantErr.Error(), err.Error(), "Error message is wrong")
				return
			}

			if err != nil {
				t.Errorf("GetAllChildFolders() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.want, got, "The expected output doesn't match")
		})
	}
}
