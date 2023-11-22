package usecases


//go:generate mockery --name UrlGenerator
type UrlGenerator interface {
	Generate(l int) string
}
