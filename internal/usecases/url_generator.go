package usecases

type UrlGenerator interface {
	Generate(l int) string
}
