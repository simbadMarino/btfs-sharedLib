//go:build js || wasip1
// +build js || wasip1

package path

func volumes() ([]*volume, error) {
	var vs []*volume
	vs = append(vs, &volume{
		Name:       "Root",
		MountPoint: "/",
	})
	return vs, nil
}
