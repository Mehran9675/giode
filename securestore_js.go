//go:build js

package giode

func secureSet(service, name string, value []byte) error { return ErrNotSupported }

func secureGet(service, name string) ([]byte, error) { return nil, ErrNotSupported }

func secureDelete(service, name string) error { return ErrNotSupported }
