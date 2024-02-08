package repository

type Entity[Json any] interface {
	GetUid() string
	ToJson() Json
}
