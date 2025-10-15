package springinitializr

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func (i Instance) InitializeProject(options InitializrOptions) error {
	url, err := i.getDownloadURL(options)
	if err != nil {
		return err
	}

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("GET request to spring initializr url failed with status code %d", resp.StatusCode)
	}

	path := options.ArtifactID
	if options.Name != "" {
		path = strings.ReplaceAll(strings.ToLower(options.Name), " ", "-")
	}

	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		err = os.Mkdir(path, 0o755)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	zipFilename := fmt.Sprintf("%s.zip", path)
	zipFile, err := os.Create(zipFilename)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	_, err = io.Copy(zipFile, resp.Body)
	if err != nil {
		return err
	}

	err = unzipProject(zipFilename, path)

	if err != nil {
		return err
	}

	err = os.Remove(zipFilename)

	return err
}

func (i Instance) getDownloadURL(options InitializrOptions) (string, error) {
	u, err := url.Parse(i.BaseURL)
	u.Path = "/starter.zip"
	if err != nil {
		return "", errors.New("error: could not parse base url")
	}

	if options.ArtifactID == "" {
		return "", errors.New("error: artifact id should be set")
	}
	if options.GroupID == "" {
		return "", errors.New("error: group id should be set")
	}

	q := u.Query()

	q.Add("artifactId", options.ArtifactID)
	q.Add("groupId", options.GroupID)

	if options.JavaVersion != "" {
		q.Add("javaVersion", options.JavaVersion)
	}
	if options.Packaging != "" {
		q.Add("packaging", string(options.Packaging))
	}
	if options.Language != "" {
		q.Add("language", string(options.Language))
	}
	if options.ProjectType != "" {
		q.Add("type", string(options.ProjectType))
	}
	if options.Name != "" {
		q.Add("name", options.Name)
	}
	if options.Description != "" {
		q.Add("description", options.Description)
	}
	if options.BootVersion != "" {
		q.Add("bootVersion", options.BootVersion)
	}
	if options.Version != "" {
		q.Add("version", options.Version)
	}
	if options.PackageName != "" {
		q.Add("packageName", options.PackageName)
	}

	if len(options.Dependencies) > 1 {
		depsList := make([]string, len(options.Dependencies))
		for _, dependency := range options.Dependencies {
			depsList = append(depsList, dependency.ID)
		}
		deps := strings.Join(depsList, ",")
		q.Add("dependencies", deps)
	}

	if len(options.Dependencies) == 1 {
		deps := options.Dependencies[0].ID
		q.Add("dependencies", deps)
	}

	u.RawQuery = q.Encode()
	return u.String(), nil
}
