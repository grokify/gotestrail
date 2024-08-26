package gotestrail

type Metadatas []Metadata

type Metadata struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
