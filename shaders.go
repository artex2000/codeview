package main

var VertexShader string = `
#version 430

uniform mat4 Model;
uniform mat4 View;
uniform mat4 Projection;

in vec3 VertexPosition;
in vec3 VertexColor;

out vec3 Color;

void main() {
        Color = VertexColor;
        gl_Position = Projection * (View * (Model * vec4(VertexPosition, 1.0)));
//        gl_Position = Projection * vec4(Model * VertexPosition, 1);
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
