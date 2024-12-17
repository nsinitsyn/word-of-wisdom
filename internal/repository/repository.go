package repository

import (
	"bufio"
	"bytes"
	_ "embed"
	"math/rand"
	"word-of-wisdom/internal/server"
)

var _ server.Repository = (*fileRepository)(nil)

//go:embed source.txt
var source []byte

type fileRepository struct {
	quotes []string
}

func NewFileRepository() fileRepository {
	f := fileRepository{quotes: []string{}}
	f.init()
	return f
}

func (f *fileRepository) init() {
	reader := bytes.NewReader(source)
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		f.quotes = append(f.quotes, scanner.Text())
	}
}

func (f *fileRepository) GetQuote() string {
	return f.quotes[rand.Intn(len(f.quotes))]
}
