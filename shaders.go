package main

var VertexShader string = `
#version 430

in vec3 VertexPosition;
in vec3 VertexColor;

out vec3 Color;

void main() {
        Color = VertexColor;
        gl_Position = vec4(VertexPosition, 1.0);
}
`

var FragmentShader string = `
#version 430

in vec3 Color;
out vec4 FragColor;

void main() {
        FragColor = vec4(Color, 1.0);
}
`
