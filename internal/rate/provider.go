package rate

type Provider interface {
	ExchangeRate() (float32, error)
}
