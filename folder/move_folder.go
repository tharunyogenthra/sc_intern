package folder

import (
	"errors"
	"fmt"
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

	if name == dst {
		return nil, errors.New("Error: Cannot move a folder to itself")
	}

	nameOrg := f.GetFolderOrgID(name)
	dstOrg := f.GetFolderOrgID(dst)

	if nameOrg != dstOrg {
		return nil, errors.New("Error: Cannot move a folder to a different organization")
	}

	// nameOrg and dstOrg can be used interchangeably now
	nameChildFolders, err := f.GetAllChildFolders(nameOrg, name)
	if err != nil {
		return nil, err
	}

	// iterating over nameChildFolders slice to avoid circular dependency
	// This will work for both immediate connections but also deep connections
	for _, folder := range nameChildFolders {
		if folder.Name == dst {
			return nil, errors.New("Error: Cannot move a folder to a child of itself")
		}
	}

	// Finished Error handling

	// This is really trivial as we are stated in the spec to not persist state
	folders := f.folders
	dstPath := ""

	for _, folder := range folders {
		if folder.Name == dst {
			dstPath = folder.Paths
		}
	}

	// rewrite prefix
	prefix := dstPath + "." + name

	for i := range folders {
		if folders[i].Name == name {
			folders[i].Paths = dstPath + "." + name
		} else if isInChildren(folders[i].Name, nameChildFolders) {
			folders[i].Paths = concatPaths(folders[i].Paths, prefix)
		}
	}

	return folders, nil
}

// This function is the main logic of this component
// It concats two file paths to make the ordering of children work
// To showcase this I am going to show this through running how test 3 works "We move a subfolder to a diff folder in the same org" (b -> g)
/*
str 		= alpha.bravo.charlie
prefix 		= golf.bravo
strSplit 	= [alpha bravo charlie]
prefixSplit = [golf bravo]

The end goal is the make charlie path look like golf.bravo.charlie to reflect it being moved

we use a loop to find where the bravo is in str

we then concat prefixSplit with everythin after bravo in str
*/
func concatPaths(str string, prefix string) string {
	strSplit := strings.Split(str, ".")
	prefixSplit := strings.Split(prefix, ".")

	stoppingIndex := 0

	for i, name := range strSplit {
		if name == prefixSplit[len(prefixSplit)-1] {
			stoppingIndex = i
		}
	}
	fmt.Println(str, strSplit, prefix, prefixSplit)

	result := strings.Join(prefixSplit, ".") + "." + strings.Join(strSplit[stoppingIndex+1:], ".")
	return result
}

// Just checks if the name of the folder is within the children folder
func isInChildren(folderName string, childrenFolder []Folder) bool {
	for _, child := range childrenFolder {
		if child.Name == folderName {
			return true
		}
	}
	return false
}
