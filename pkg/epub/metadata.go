package epub

import (
	"fmt"
	"github.com/kapmahc/epub"
	"strings"
)

var (
	IsbnKey = "ISBN" // IsbnKey defines the identification of ISBN.
)

// ReadMetadata Extract metadata from an EPUB file
func ReadMetadata(epubPath string) (*Metadata, error) {
	book, err := epub.Open(epubPath)
	if err != nil {
		return nil, err
	}
	if book == nil || len(book.Opf.Metadata.Title) == 0 {
		return nil, fmt.Errorf("Unable to obtain book title")
	}
	fmt.Printf("book.Opf.Metadata: %v\n", book.Opf.Metadata)

	var title string
	if len(book.Opf.Metadata.Title) > 0 {
		title = book.Opf.Metadata.Title[0]
	}
	var author string
	if len(book.Opf.Metadata.Creator) > 0 {
		author = book.Opf.Metadata.Creator[0].Data
	}
	var description string
	if len(book.Opf.Metadata.Description) > 0 {
		description = book.Opf.Metadata.Description[0]
	}

	var publisher string
	if len(book.Opf.Metadata.Publisher) > 0 {
		publisher = book.Opf.Metadata.Publisher[0]
	}
	var date string
	if len(book.Opf.Metadata.Date) > 0 {
		date = book.Opf.Metadata.Date[0].Data
	}
	var language string
	if len(book.Opf.Metadata.Language) > 0 {
		language = book.Opf.Metadata.Language[0]
	}

	var identifier []Identifier
	var isbn string
	for _, d := range book.Opf.Metadata.Identifier {
		if strings.EqualFold(d.Scheme, IsbnKey) {
			isbn = d.Data
		}
		identifier = append(identifier, Identifier{
			ID:     d.ID,
			Scheme: d.Scheme,
			Value:  d.Data,
		})
	}

	return &Metadata{
		Title:       title,
		Author:      author,
		Description: description,
		Publisher:   publisher,
		Date:        date,
		Identifier:  identifier,
		Language:    language,
		Isbn:        isbn,
	}, nil

}
