package springinitializr

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
)

func (i Instance) GetCapabilities() (*Capabilities, error) {
	u, err := url.Parse(i.BaseURL)
	u.Path = "/"
	if err != nil {
		return nil, errors.New("error: could not parse base url")
	}

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, errors.New("error: could not create GET request to base url")
	}
	req.Header.Set("Accept", "application/vnd.initializr.v2.2+json")

	client := http.DefaultClient

	resp, err := client.Do(req)

	if err != nil {
		return nil, errors.New("error: could not sned GET request to base url")
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("error: could not parse body of capabilities response")
	}
	dto := capabilitiesDto{}
	json.Unmarshal(body, &dto)
	javaVersions := make([]string, len(dto.JavaVersions.Values))
	for i, v := range dto.JavaVersions.Values {
		javaVersions[i] = v.ID
	}

	packagingTypes := make([]AppPackaging, len(dto.PackagingTypes.Values))
	for i, v := range dto.PackagingTypes.Values {
		packagingTypes[i] = AppPackaging(v.ID)
	}

	projectTypes := make([]AppProjectType, len(dto.ProjectTypes.Values))
	for i, v := range dto.ProjectTypes.Values {
		projectTypes[i] = AppProjectType(v.ID)
	}

	languages := make([]AppLanguage, len(dto.Languages.Values))
	for i, v := range dto.Languages.Values {
		languages[i] = AppLanguage(v.ID)
	}

	bootVersions := make([]string, len(dto.BootVersions.Values))
	for i, v := range dto.BootVersions.Values {
		bootVersions[i] = v.ID
	}

	// len of "core" deps + len of "other" deps = total len of deps
	dependencies := []Dependency{}
	for _, category := range dto.Dependencies.Values {
		for _, dependency := range category.Values {
			dependencies = append(dependencies, Dependency{
				ID:          dependency.ID,
				Name:        dependency.Name,
				Description: dependency.Description,
				Category:    DependencyCategory(category.Name),
			})
		}
	}

	return &Capabilities{
		JavaVersions:   javaVersions,
		PackagingTypes: packagingTypes,
		ProjectTypes:   projectTypes,
		Languages:      languages,
		BootVersions:   bootVersions,
		Dependencies:   dependencies,
	}, nil
}
