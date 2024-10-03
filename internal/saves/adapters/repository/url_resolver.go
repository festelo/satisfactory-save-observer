package repository

import (
	"fmt"
	"net/url"
)

type SimpleUrlResolverRepository struct {
	host   string
	mapUrl string
}

func NewSimpleUrlResolverRepository(host string, mapUrl string) *SimpleUrlResolverRepository {
	return &SimpleUrlResolverRepository{host, mapUrl}
}

func (r SimpleUrlResolverRepository) SaveUrl(saveName string) (string, error) {
	return url.JoinPath(r.host, url.PathEscape(saveName))
}

func (r SimpleUrlResolverRepository) MapUrl(saveName string) (string, error) {
	saveUrl, err := r.SaveUrl(saveName)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(
		r.mapUrl,
		url.QueryEscape(saveUrl),
	), nil
}
