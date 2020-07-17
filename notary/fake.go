package notary

import "github.com/theupdateframework/notary/client"

//type MyRepository interface {
//GetAllTargetMetadataByName(name string) ([]client.TargetSignedStruct, error)
//}

// Fake Repository for testing purposes
type fake struct {
	targetsMetadata map[string][]client.TargetSignedStruct
}

func (f fake) GetAllTargetMetadataByName(tag string) ([]client.TargetSignedStruct, error) {
	return f.targetsMetadata[tag], nil
}
