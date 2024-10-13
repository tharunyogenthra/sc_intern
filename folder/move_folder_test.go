package folder_test

import (
	"errors"
	"github.com/georgechieng-sc/interns-2022/folder"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func GetTestingSampleData2() []folder.Folder {
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
			Paths: "alpha.delta.echo",
		},
		{
			Name:  "foxtrot",
			OrgId: uuid.FromStringOrNil("c1556e17-b7c0-45a3-a6ae-9546248fb17a"),
			Paths: "foxtrot",
		},
		{
			Name:  "golf",
			OrgId: uuid.FromStringOrNil("38b9879b-f73b-4b0e-b9d9-4fc4c23643a7"),
			Paths: "golf",
		},
	}

}

func Test_folder_MoveFolder(t *testing.T) {
	t.Parallel()
	orgID := uuid.FromStringOrNil("38b9879b-f73b-4b0e-b9d9-4fc4c23643a7")
	tests := [...]struct {
		name_of_test string
		name         string
		dst          string
		folders      []folder.Folder
		want         []folder.Folder
		wantErr      error
	}{
		{
			name_of_test: "We move the folder to child of same level",
			name:         "bravo",
			dst:          "delta",
			folders:      GetTestingSampleData2(),
			want: []folder.Folder{
				{Name: "alpha", Paths: "alpha", OrgId: orgID},
				{Name: "bravo", Paths: "alpha.delta.bravo", OrgId: orgID},
				{Name: "charlie", Paths: "alpha.delta.bravo.charlie", OrgId: orgID},
				{Name: "delta", Paths: "alpha.delta", OrgId: orgID},
				{Name: "echo", Paths: "alpha.delta.echo", OrgId: orgID},
				{Name: "foxtrot", Paths: "foxtrot", OrgId: uuid.FromStringOrNil("c1556e17-b7c0-45a3-a6ae-9546248fb17a")},
				{Name: "golf", Paths: "golf", OrgId: orgID},
			},
			wantErr: nil,
		},
		{
			name_of_test: "We move a subfolder to a diff folder in the same org",
			name:         "bravo",
			dst:          "golf",
			folders:      GetTestingSampleData2(),
			want: []folder.Folder{
				{Name: "alpha", Paths: "alpha", OrgId: orgID},
				{Name: "bravo", Paths: "golf.bravo", OrgId: orgID},
				{Name: "charlie", Paths: "golf.bravo.charlie", OrgId: orgID},
				{Name: "delta", Paths: "alpha.delta", OrgId: orgID},
				{Name: "echo", Paths: "alpha.delta.echo", OrgId: orgID},
				{Name: "foxtrot", Paths: "foxtrot", OrgId: uuid.FromStringOrNil("c1556e17-b7c0-45a3-a6ae-9546248fb17a")},
				{Name: "golf", Paths: "golf", OrgId: orgID},
			},
			wantErr: nil,
		},
		{
			name_of_test: "Name string is empty",
			name:         "",
			dst:          "charlie",
			folders:      GetTestingSampleData2(),
			want:         []folder.Folder{},
			wantErr:      errors.New("Error: Source folder does not exist"),
		},
		{
			name_of_test: "Dst string is empty",
			name:         "bravo",
			dst:          "",
			folders:      GetTestingSampleData2(),
			want:         []folder.Folder{},
			wantErr:      errors.New("Error: Destination folder does not exist"),
		},
		{
			name_of_test: "Cannot make a circular dependency for immediate child connections",
			name:         "bravo",
			dst:          "charlie",
			folders:      GetTestingSampleData2(),
			want:         []folder.Folder{},
			wantErr:      errors.New("Error: Cannot move a folder to a child of itself"),
		},
		{
			name_of_test: "Cannot make a circular dependency but the file hierarchy is deep",
			name:         "alpha",
			dst:          "charlie",
			folders:      GetTestingSampleData2(),
			want:         []folder.Folder{},
			wantErr:      errors.New("Error: Cannot move a folder to a child of itself"),
		},
		{
			name_of_test: "Cannot make a circular dependency by having a folder as a child of itself",
			name:         "bravo",
			dst:          "bravo",
			folders:      GetTestingSampleData2(),
			want:         []folder.Folder{},
			wantErr:      errors.New("Error: Cannot move a folder to itself"),
		},
		{
			name_of_test: "Cant move a folder to diff org",
			name:         "bravo",
			dst:          "foxtrot",
			folders:      GetTestingSampleData2(),
			want:         []folder.Folder{},
			wantErr:      errors.New("Error: Cannot move a folder to a different organization"),
		},
		{
			name_of_test: "Src doesnt exist",
			name:         "invalid_folder",
			dst:          "delta",
			folders:      GetTestingSampleData2(),
			want:         []folder.Folder{},
			wantErr:      errors.New("Error: Source folder does not exist"),
		},
		{
			name_of_test: "Dest doesnt exist.",
			name:         "bravo",
			dst:          "invalid_folder",
			folders:      GetTestingSampleData2(),
			want:         []folder.Folder{},
			wantErr:      errors.New("Error: Destination folder does not exist"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name_of_test, func(t *testing.T) {
			f := folder.NewDriver(tt.folders)
			got, err := f.MoveFolder(tt.name, tt.dst)

			if err != nil && tt.wantErr != nil {
				assert.Equal(t, tt.wantErr.Error(), err.Error(), "Error message is wrong")
				return
			}

			if err != nil {
				t.Errorf("MoveFolder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.want, got, "The expected output doesn't match")
		})
	}
}
