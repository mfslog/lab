package registry


type Service struct{
	Name string `json:"name"`
	version string `json:"version"`
	MetaData map[string]string `json:"meta_data"`
	Nodes []*Node `json:"nodes"`
}

type Node struct{
	ID   int32 `json:"id"`
	Address string `json:"address"`
	MetaData map[string]string `json:"meta_data"`
}