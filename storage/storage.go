package storage

import (
	"TelegramBot/lib/e"
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"strconv"
)

type Storage interface {
	Save(ctx context.Context, p *Page) error
	PickRandom(ctx context.Context, userName string, idUser int) (*Page, error)
	Remove(ctx context.Context, p *Page) error
	IsExists(ctx context.Context, p *Page) (bool, error)
}

var ErrNoSavedPages = errors.New("no saved pages")

type Page struct {
	URL      string
	UserName string
	IdUser   int
}

func (p Page) Hash() (string, error) {
	h := sha1.New()

	if _, err := io.WriteString(h, p.URL); err != nil {
		return "", e.Wrap("can't calculate hash", err)
	}

	if _, err := io.WriteString(h, p.UserName); err != nil {
		return "", e.Wrap("can't calculate hash", err)
	}
	if _, err := io.WriteString(h, strconv.Itoa(p.IdUser)); err != nil {
		return "", e.Wrap("can't calculate hash", err)
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
