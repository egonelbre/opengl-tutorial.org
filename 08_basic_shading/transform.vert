#version 330 core

uniform mat4 Projection;
uniform mat4 Camera;
uniform mat4 VP;
uniform vec3 LightPosition;

uniform mat4 Model;
in vec3 vertex;
in vec2 vertexUV;
in vec3 vertexNormal;

out vec2 UV;

void main(){
	gl_Position = VP * Model * vec4(vertex, 1);
	UV = vertexUV;
}