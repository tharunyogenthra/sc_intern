package folder

import (
	"errors"
	// "fmt"
	"strings"
)

func (f *driver) MoveFolder(name string, dst string) ([]Folder, error) {
	nameExists := f.CheckFolderExists(name)
	dstExists := f.CheckFolderExists(dst)
	
	if !nameExists {
		return nil, errors.New("Error: Source folder does not exist")
	}
	
	if !dstExists {
		return nil, errors.New("Error: Destination folder does not exist")
	}

	if (name == dst) {
		return nil, errors.New("Error: Cannot move a folder to itself")
	}

	nameOrg := f.GetFolderOrgID(name)
	dstOrg := f.GetFolderOrgID(dst)

	if (nameOrg != dstOrg) {
		return nil, errors.New("Error: Cannot move a folder to a different organization")
	}

	// nameOrg and dstOrg can be used interchangeably now
	nameChildFolders, err := f.GetAllChildFolders(nameOrg, name)
	if (err != nil) {
		return nil, err
	}
	
	// iterating over nameChildFolders slice to avoid circular dependency
	// This will work for both immediate connections but also deep connections
	for _, folder := range nameChildFolders {
		if (folder.Name == dst) {
			return nil, errors.New("Error: Cannot move a folder to a child of itself")
		}
	}

	// Finished Error handling

	// This is really trivial as we are stated in the spec to not persist state
	folders := f.folders
	dstPath := ""
	for _, folder := range folders {
		if (folder.Name == dst) {
			dstPath = folder.Paths
		}
	}

	// rewrite prefix
	prefix := dstPath + "." + name


	for i := range folders {
		if (folders[i].Name == name) {
			folders[i].Paths = dstPath + "." + name
		} else if (isInChildren(folders[i].Name, nameChildFolders)) {
			folders[i].Paths = concatPaths(folders[i].Paths, prefix)
		} 
	}

	return folders, nil
}

func concatPaths(prefix, suffix string) string {
	prefixParts := strings.Split(prefix, ".")
	suffixParts := strings.Split(suffix, ".")

	return strings.Join(append(suffixParts, prefixParts[2:]...), ".")
}

func isInChildren(folderName string, childrenFolder []Folder) bool {
	for _, child := range childrenFolder {
		if child.Name == folderName {
			return true
		}
	}
	return false
}