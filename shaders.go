package main

var VertexShader string = `
#version 430

uniform mat4 Model;
uniform mat4 View;
uniform mat4 Projection;

in vec3 VertexPosition;
in vec2 TextureCoord;

out vec2 TexCoord;

void main() {
        TexCoord = TextureCoord;
        gl_Position = Projection * (View * (Model * vec4(VertexPosition, 1.0)));
}
`

var FragmentShader string = `
#version 430

uniform sampler2D Texture;
uniform vec4 Foreground;
uniform vec4 Background;

in vec2 TexCoord;
out vec4 FragColor;

void main() {
        float gamma = 1.43;

        vec4 fg = pow(Foreground, vec4(1.0/gamma));
        vec4 bg = pow(Background, vec4(1.0/gamma));

        /*
        vec4 fg = Foreground;
        vec4 bg = Background;
        */

        vec4 current = texture(Texture, TexCoord);

        float r = current.r * fg.r + (1 - current.r) * bg.r;
        float g = current.g * fg.g + (1 - current.g) * bg.g;
        float b = current.b * fg.b + (1 - current.b) * bg.b;

        FragColor = pow(vec4(r, g, b, current.a), vec4(gamma));
        //FragColor = vec4(r, g, b, 1);
}
`
