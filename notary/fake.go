package notary

import "github.com/theupdateframework/notary/client"

//type MyRepository interface {
//GetAllTargetMetadataByName(name string) ([]client.TargetSignedStruct, error)
//}

type Fake struct {
	targetsMetadata map[string][]client.TargetSignedStruct
}

func (f *Fake) GetAllTargetMetadataByName(tag string) ([]client.TargetSignedStruct, error) {
	return f.targetsMetadata[tag], nil
}
