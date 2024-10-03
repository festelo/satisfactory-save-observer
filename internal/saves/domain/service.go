package domain

import (
	"errors"
	"io"
	"sort"
)

type SavesService struct {
	savesRepository SavesRepository
	urlResolver     UrlResolverRepository
}

func NewSavesService(savesRepository SavesRepository, urlResolver UrlResolverRepository) *SavesService {
	return &SavesService{savesRepository, urlResolver}
}

func (s *SavesService) ListSaves() ([]*Save, error) {
	saves, err := s.savesRepository.ListSaves()
	if err != nil {
		return nil, err
	}
	for i, element := range saves {
		saveUrl, err := s.urlResolver.SaveUrl(element.Name)
		if err != nil {
			return nil, err
		}
		saves[i].Link = saveUrl
		mapUrl, err := s.urlResolver.MapUrl(element.Name)
		if err != nil {
			return nil, err
		}
		saves[i].MapLink = mapUrl
	}
	sort.Slice(saves, func(i, j int) bool {
		return saves[i].CreatedAt.After(saves[j].CreatedAt)
	})
	return saves, nil
}

func (s *SavesService) CopySave(name string, writer io.Writer) error {
	return s.savesRepository.CopySave(name, writer)
}

func (s *SavesService) CopyLatestSave(writer io.Writer) error {
	saves, err := s.ListSaves()
	if err != nil {
		return err
	}
	if len(saves) == 0 {
		return errors.New("saves are empty")
	}
	return s.CopySave(saves[0].Name, writer)
}
