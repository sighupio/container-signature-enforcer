package config

// custom errors needed
type ErrNoRepositoryMatched struct{}

func (e ErrNoRepositoryMatched) Error() string { return "No repository matched" }
