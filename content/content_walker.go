package content

import (
	"os"
	"path/filepath"
)

type FilterJob interface {
	Filter(path string, info os.FileInfo) ([]byte, error)

	Process(jBytes []byte, uri string) error

	OnComplete() error

	LimitReached() bool
}

func FilterAndProcess(dir string, job FilterJob) error {
	walkFunc := func(path string, info os.FileInfo, err error) error {
		if job.LimitReached() {
			return LimitReached
		}

		jBytes, err := job.Filter(path, info)
		if err != nil {
			return err
		}

		if jBytes == nil {
			return nil
		}
		return job.Process(jBytes, path)
	}

	if err := filepath.Walk(dir, walkFunc); err != nil {
		return err
	}

	return job.OnComplete()
}
