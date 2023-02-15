package domain

import "context"

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE
type HashKeeper interface {
	Store(context.Context, Hash) error
}

type HashGetter interface {
	Get(context.Context) (Hash, error)
}
