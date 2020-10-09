package plugin

import (
	"context"
	"strings"

	"github.com/sirupsen/logrus"
)

// consider validates if a `drone.yml` was allowed using p.considerFile
func (p *Plugin) consider(file string, req *request) bool {
	// all files are allowed if p.considerFile was not specified
	if p.considerFile == "" {
		return true
	}

	for _, considerListEntry := range req.ConsiderList {
		if file == considerListEntry {
			return true
		}
	}
	return false
}

// getConsiderFile returns the 'drone.yml' entries in a consider file as a string slice
func (p *Plugin) getConsiderFile(ctx context.Context, req *request) ([]string, error) {
	files := make([]string, 0)
	if p.considerFile == "" {
		return files, nil
	}

	// download considerFile from github
	fc, err := p.getScmFile(ctx, req, p.considerFile)
	if err != nil {
		logrus.Errorf("%s skipping: %s is not present: %v", req.UUID, p.considerFile, err)
		return files, err
	}

	// collect entries
	for _, v := range strings.Split(fc, "\n") {
		// skip empty lines and comments
		if strings.TrimSpace(v) == "" || strings.HasPrefix(v, "#") {
			continue
		}
		// skip lines which do not contain a 'drone.yml' reference
		if !strings.HasSuffix(v, req.Repo.Config) {
			logrus.Warnf("%s skipping invalid reference to %s in %s", req.UUID, v, p.considerFile)
			continue
		}
		files = append(files, v)
	}

	return files, nil
}
