package wisdombook

import "math/rand"

type WisdomBook struct {
	quotes []string
}

func NewWisdomBook(quotes []string) *WisdomBook {
	return &WisdomBook{
		quotes: quotes,
	}
}

// GetRandomQuote returns a random quote from the wisdom book.
func (wb *WisdomBook) GetRandomQuote() string {
	return wb.quotes[rand.Intn(len(wb.quotes))]
}
