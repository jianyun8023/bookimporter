package epub

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
)

// Regular expression for 13-digit ISBN
var isbnRegexp = regexp.MustCompile(`97[89]\d{10}`)

// FindIsbn Find ISBNs in epub files
// This ISBN will not be read from the EPUB metadata.
func FindIsbn(epubPath string) (string, error) {

	// Open EPUB file
	r, err := zip.OpenReader(epubPath)
	if err != nil {
		return "", err
	}
	defer func(r *zip.ReadCloser) {
		err := r.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(r)

	for _, f := range r.File {

		// Skip the OPF format file.
		if strings.HasSuffix(f.Name, ".opf") {
			continue
		}

		//Search for ISBN in the file
		isbn, _ := extractISBN(f)
		if isbn != "" {
			return isbn, nil
		}
	}
	return "", nil
}

// extractISBN Extract ISBN from file
func extractISBN(f *zip.File) (string, error) {
	rc, err := f.Open()
	if err != nil {
		return "", err
	}
	defer func(rc io.ReadCloser) {
		_ = rc.Close()
	}(rc)

	// Read the text content line by line and match ISBN.
	scanner := bufio.NewScanner(rc)
	for scanner.Scan() {
		line := scanner.Text()

		// There is a string 'isbn' in the line.
		if !strings.Contains(strings.ToLower(line), "isbn") {
			continue
		}
		// Remove hyphens from the string
		line = strings.ReplaceAll(line, "-", "")

		isbn := isbnRegexp.FindString(line)
		if isbn != "" {
			//Print the current file name and path.
			fmt.Println(f.Name)
			return isbn, nil
		}
	}
	return "", nil
}
