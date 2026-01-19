package cli

import "path/filepath"

func OpenURL(url string) error {
	err := openBrowser(url)

	if err != nil {
		Error("Unable to start your browser. try manually.")
		Error(url)
	}

	return nil
}

func OpenFile(name string) error {
	path, err := filepath.Abs(name)

	if err != nil {
		return err
	}

	if err := openBrowser("file://" + path); err != nil {
		Error("Unable to open file. try manually")
		Error(name)
	}

	return nil
}
