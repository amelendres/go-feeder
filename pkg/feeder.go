package feed

import (
	"errors"
	"regexp"
)

const (
	fs_provider = "fs"
	gd_provider = "gd"
)

var ErrUnknownFeed = errors.New("unknown feed")

type Feeder interface {
	Feeds(path string) (*ParsedItems, error)
	Destination(d *Destination)
}

type feeder struct {
	fileProviders map[string]FileProvider
	parser        Parser
	feeds         []Item
}

func NewFeeder(p Parser, providers []FileProvider) Feeder {

	feeder := &feeder{
		parser: p,
	}

	feeder.fileProviders = make(map[string]FileProvider)
	for _, pro := range providers {
		feeder.AddProvider(pro)
	}

	return feeder
}

func (s *feeder) Destination(d *Destination) {
	s.parser.Destination(d)
}

func (s *feeder) Feeds(path string) (*ParsedItems, error) {
	prov := fs_provider
	if isGoogleDrive(path) {
		prov = gd_provider
	}
	f, err := s.fileProviders[prov].File(path)

	if err != nil {
		return nil, err
	}
	return s.parser.Parse(f)
}

func (s *feeder) AddProvider(p FileProvider) {
	s.fileProviders[p.Name()] = p
}

func isGoogleDrive(path string) bool {
	re := regexp.MustCompile("docs.google.com")
	return len(re.FindAllString(path, 1)) == 1
}
