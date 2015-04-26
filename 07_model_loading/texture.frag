#version 330 core
in vec2 UV;
out vec3 color;
uniform sampler2D sampler;

void main(){
	color = texture(sampler, UV).rgb;
}