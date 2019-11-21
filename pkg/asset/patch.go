package asset

import (
	"errors"
	"fmt"
	"io/ioutil"

	jsonpatch "github.com/evanphx/json-patch"
	"github.com/ghodss/yaml"
)

type PatchJSON6902 struct {
	// these fields specify the patch target resource
	Group      string `yaml:"group"`
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	// Name and Namespace are optional
	// NOTE: technically name is required now, but we default it elsewhere
	// Third party users of this type / library would need to set it.
	Metadata metadata `yaml:"metadata"`
	// Patch should contain the contents of the json patch as a string
	Patch     string `yaml:"patch"`
	JsonPatch jsonpatch.Patch
	Matcher   patchMatch
}

type patchMatch struct {
	Kind       string `json:"kind,omitempty"`
	APIVersion string `json:"apiVersion,omitempty"`
	Metadata   metadata
}
type metadata struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace,omitempty"`
}

func getMatchInfo(data []byte) (patchMatch, error) {
	r := patchMatch{}
	if err := yaml.Unmarshal(data, &r); err != nil {
		return patchMatch{}, errors.New("Failed to parse metadata")
	}
	return r, nil
}

func (a *Asset) GenResourceMatchInfo() {
	m, _ := getMatchInfo(a.Data)
	a.Matcher = m
}

func ReadPatchDir(dir string) ([]PatchJSON6902, error) {
	patches := []PatchJSON6902{}
	patchfiles, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, f := range patchfiles {
		file, _ := ioutil.ReadFile(dir + "/" + f.Name())
		p, err := ReadPatch(file)
		if err != nil {
			return nil, err
		}
		patches = append(patches, p)
	}
	return patches, nil
}

func ReadPatch(file []byte) (PatchJSON6902, error) {
	p := PatchJSON6902{}
	m, err := getMatchInfo(file)
	if err = yaml.Unmarshal(file, &p); err != nil {
		return PatchJSON6902{}, errors.New("Unable to parse patch metadata")
	}
	p.Matcher = m
	p, err = convertJSON6902Patch(p)
	if err != nil {
		return PatchJSON6902{}, err
	}
	return p, nil
}

func (a *Asset) patch6902(patch PatchJSON6902) (match bool, err error) {
	if !a.matches(patch.Matcher) {
		return false, nil
	}
	aj, err := yaml.YAMLToJSON(a.Data)
	if err != nil {
		return true, err
	}
	patched, err := patch.JsonPatch.Apply(aj)
	if err != nil {
		fmt.Println(err)
		return true, err
	}
	ay, err := yaml.JSONToYAML(patched)
	if err != nil {
		return true, err
	}
	a.Data = ay
	return true, nil
}

// TBM to assets.go
func (a Asset) matches(pm patchMatch) bool {
	am := &a.Matcher
	return am.Kind == pm.Kind && am.APIVersion == pm.APIVersion && am.Metadata.Name == pm.Metadata.Name
}

func convertJSON6902Patch(patch PatchJSON6902) (PatchJSON6902, error) {
	patchJSON, err := yaml.YAMLToJSON([]byte(patch.Patch))
	if err != nil {
		return PatchJSON6902{}, err
	}
	jpatch, err := jsonpatch.DecodePatch(patchJSON)
	if err != nil {
		return PatchJSON6902{}, err
	}
	patch.JsonPatch = jpatch
	return patch, nil
}

// TBM to assets.go
func (as Assets) Patch(patches []PatchJSON6902) error {
	for i, _ := range as {
		as[i].GenResourceMatchInfo()
		for _, p := range patches {
			if _, err := as[i].patch6902(p); err != nil {
				return err
			}
		}
	}
	return nil
}
