#version 330 core
uniform mat4 Projection;
uniform mat4 Camera;
uniform mat4 Model;

in vec3 vertex;
in vec2 vertexUV;

out vec2 UV;

void main(){
	gl_Position = Projection * Camera * Model * vec4(vertex, 1);
	UV = vertexUV;
}