package files

import (
	"bufio"
	"os"
)

type Indexer struct{
    file *os.File
    writer *bufio.Writer
}

func NewIndexer(path string) (*Indexer, error) {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644) // overwrite if exists
	if err != nil {
        logger.ERROR.Printf("Error when opening file: %v\n", err)
		return nil, err
	}

    return &Indexer{
        file: file,
        writer: bufio.NewWriter(file),
    }, nil
}

func (idx *Indexer) WriteLine(line string) error {
    _, err := idx.writer.WriteString(line + "\n")
    if err != nil {
        logger.ERROR.Printf("Error when writing line: %v\n", err)
    }
    return err
}

func (idx *Indexer) Close() {
    idx.writer.Flush()
    idx.file.Close()
}

