package schema

func SchemaOK(kind string) (*struct{}, error) {
	_, err := schemaFor(kind)
	if err != nil {
		return nil, err
	}
	return &struct{}{}, nil
}
