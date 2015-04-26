package obj

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/go-gl/mathgl/mgl32"
)

type Data struct {
	Vertex []float32
	UV     []float32
	Normal []float32
}

func LoadFile(filename string) (*Data, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return Load(bufio.NewReader(file))
}

func Load(r io.Reader) (*Data, error) {
	var err error

	data := &Data{}

	var vertices []mgl32.Vec3
	var uvs []mgl32.Vec2
	var normals []mgl32.Vec3

	line := 0
	for err == nil {
		var hdr string
		line += 1

		_, err = fmt.Fscanf(r, "%s", &hdr)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("Error at %d: %v", line, err)
		}

		switch hdr {
		case "v":
			var v mgl32.Vec3
			_, err = fmt.Fscanf(r, "%f %f %f\n", &v[0], &v[1], &v[2])
			vertices = append(vertices, v)
		case "vt":
			var uv mgl32.Vec2
			_, err = fmt.Fscanf(r, "%f %f\n", &uv[0], &uv[1])
			uv[1] = -uv[1] // for DDS!!!
			uvs = append(uvs, uv)
		case "vn":
			var v mgl32.Vec3
			_, err = fmt.Fscanf(r, "%f %f %f\n", &v[0], &v[1], &v[2])
			normals = append(normals, v)
		case "f":
			var vi, uvi, ni [3]int
			_, err = fmt.Fscanf(r, "%d/%d/%d %d/%d/%d %d/%d/%d\n",
				&vi[0], &uvi[0], &ni[0],
				&vi[1], &uvi[1], &ni[1],
				&vi[2], &uvi[2], &ni[2])

			data.Vertex = append(data.Vertex, vertices[vi[0]-1][:]...)
			data.Vertex = append(data.Vertex, vertices[vi[1]-1][:]...)
			data.Vertex = append(data.Vertex, vertices[vi[2]-1][:]...)

			data.UV = append(data.UV, uvs[uvi[0]-1][:]...)
			data.UV = append(data.UV, uvs[uvi[1]-1][:]...)
			data.UV = append(data.UV, uvs[uvi[2]-1][:]...)

			data.Normal = append(data.Normal, normals[ni[0]-1][:]...)
			data.Normal = append(data.Normal, normals[ni[1]-1][:]...)
			data.Normal = append(data.Normal, normals[ni[2]-1][:]...)
		default:
			var s string
			_, err = fmt.Fscanf(r, "%s\n", &s)
		}
	}
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("Error at %d: %v", line, err)
	}

	return data, nil
}
