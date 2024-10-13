package folder

import (
	"errors"
	"github.com/gofrs/uuid"
)

func GetAllFolders() []Folder {
	return GetSampleData()
}

func (f *driver) GetFoldersByOrgID(orgID uuid.UUID) []Folder {
	folders := f.folders

	res := []Folder{}
	for _, f := range folders {
		if f.OrgId == orgID {
			res = append(res, f)
		}
	}

	return res

}

func (f *driver) GetAllChildFolders(orgID uuid.UUID, name string) ([]Folder, error) {

	exists := f.CheckFolderExists(name)
	if !exists {
		return nil, errors.New("Error: Folder does not exist")
	}

	existsOrg := f.CheckFolderExistsWithinOrg(orgID, name)
	if !existsOrg {
		return nil, errors.New("Error: Folder does not exist in the specified organization")
	}

	rootPath := ""
	folderOrgID := f.GetFoldersByOrgID(orgID)

	for _, folder := range folderOrgID {
		if folder.Name == name {
			rootPath = folder.Paths
			break
		}
	}

	if rootPath == "" {
		return []Folder{}, nil
	}

	children := []Folder{}
	for _, folder := range f.folders {
		if IsChildFolder(folder, rootPath) {
			children = append(children, folder)
		}
	}

	return children, nil
}

// Checks whether a folder exists regardless of org
func (f *driver) CheckFolderExists(name string) bool {
	for _, folder := range f.folders {
		if folder.Name == name {
			return true
		}
	}
	return false
}

// Checks whether a folder exists within a specific org
func (f *driver) CheckFolderExistsWithinOrg(orgID uuid.UUID, name string) bool {
	folders := f.GetFoldersByOrgID(orgID)

	for _, folder := range folders {
		if folder.Name == name {
			return true
		}
	}
	return false
}

// Checks whether a folder is a child using string manipulation
func IsChildFolder(folder Folder, rootPath string) bool {
	// as its a subfolder this condition must apply
	if len(rootPath) >= len(folder.Paths) {
		return false
	}
	// :len(rootPath) uses string splicing very similar to python
	// x := "bravo.charlie"
	// fmt.Println(x[:len(root)]) -> bravo
	if folder.Paths != rootPath && folder.Paths[:len(rootPath)] == rootPath {
		return true
	}
	return false
}

// Returns a folder orgID as uuid
func (f *driver) GetFolderOrgID(name string) uuid.UUID {
	for _, folder := range f.folders {
		if (folder.Name == name) {
			return folder.OrgId
		}
	}
	return uuid.UUID{}
}

