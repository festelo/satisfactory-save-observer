package domain

import "io"

type SavesRepository interface {
	ListSaves() ([]*Save, error)
	CopySave(name string, writer io.Writer) error
}

type UrlResolverRepository interface {
	SaveUrl(saveName string) (string, error)
	MapUrl(saveName string) (string, error)
}
