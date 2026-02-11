#version 330

// Input vertex attributes
in vec3 vertexPosition;
in vec2 vertexTexCoord;
in vec4 vertexColor;

// Input instance attributes
in mat4 instanceTransform;

// Input uniform values
uniform mat4 mvp;

// Output vertex attributes (to fragment shader)
out vec2 fragTexCoord;
out vec4 fragColor;

void main()
{
    // Extract sprite index from instanceTransform[1][1] (M5)
    // This slot has no visual effect on a plane with Y=0 vertices
    float spriteIndex = instanceTransform[1][1];

    // Compute sprite grid position (4x4 grid)
    float spriteX = mod(spriteIndex, 4.0);
    float spriteY = floor(spriteIndex / 4.0);

    // Compute UV sub-rectangle for this sprite
    fragTexCoord = vertexTexCoord / 4.0 + vec2(spriteX / 4.0, spriteY / 4.0);
    fragColor = vertexColor;

    // Reset [1][1] to 1.0 before position calculation
    mat4 cleanTransform = instanceTransform;
    cleanTransform[1][1] = 1.0;

    gl_Position = mvp * cleanTransform * vec4(vertexPosition, 1.0);
}
