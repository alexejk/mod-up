package upgrade

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"golang.org/x/mod/modfile"
	"golang.org/x/mod/module"
)

const modFile = "go.mod"

type Upgrader struct {

	workDir string
	modFile *modfile.File
}

func New(dir string) (*Upgrader, error) {

	dir, err := homedir.Expand(dir)
	if err != nil {
		return nil, err
	}

	modFilePath := path.Join(dir, modFile)
	workDir := filepath.Dir(modFilePath)

	f, err := os.Open(modFilePath)
	if err != nil {
		return nil, err
	}

	// Read mod file
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	mf, err := modfile.Parse(modFilePath, data, nil)
	if err != nil {
		return nil, err
	}

	return &Upgrader{
		workDir: workDir,
		modFile: mf,
	}, nil
}

// Upgrade runs module update on all dependencies.
// If module is overridden with replace() directive - it will be skipped.
// If indirect is false - only direct dependencies are included
func (u *Upgrader) Upgrade(indirect bool) error {

	updateList := u.updateCandidates(indirect)
	for _, mod := range updateList {
		_ = u.UpdateModule(mod)
	}

	return nil
}

// Tidy up go modules by running "go mod tidy" in the directory
func (u *Upgrader) Tidy() error {

	fmt.Printf("# Tidying up go.mod")

	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = u.workDir
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	return cmd.Run()
}

// UpdateModule runs "go get -u <mod>" to update it to latest version
func (u *Upgrader) UpdateModule(mod module.Version) error {

	fmt.Printf("# Updating module:\t%s\n", mod.Path)
	defer fmt.Println("# -----")

	cmd := exec.Command("go", "get", "-u", mod.Path)
	cmd.Dir = u.workDir
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	err := cmd.Run()
	if err != nil {
		fmt.Printf("\t! Update error: %s\n", err)
		return err
	}

	return nil
}

func (u *Upgrader) updateCandidates(indirect bool) []module.Version {

	var updateList []module.Version
	for _, r := range u.modFile.Require {

		// Skip Indirect modules
		if !indirect && r.Indirect {
			continue
		}

		// Skip modules we override
		if u.isModuleOverridden(r.Mod) {
			continue
		}
		updateList = append(updateList, r.Mod)
	}
	return updateList
}


func (u *Upgrader) isModuleOverridden(mod module.Version) bool {

	for _, r := range u.modFile.Replace {
		if r.Old.Path == mod.Path {
			return true
		}
	}
	return false
}