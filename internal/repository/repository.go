package repository

//go:generate go tool mockgen -destination ./mock/mock.go -package mock . Repository
type Repository interface {
}
