package uuid

import gouuid "github.com/satori/go.uuid"

func UUID() string {
	id := gouuid.NewV4()
	return id.String()
}
