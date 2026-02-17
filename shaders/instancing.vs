#version 330

// Per-vertex attributes
layout(location = 0) in vec3 vertexPosition;
layout(location = 1) in vec2 vertexTexCoord;

// Per-instance attributes (divisor=1): vec3(posX, posY, spriteIndex)
layout(location = 2) in vec3 instanceData;

uniform mat4 mvp;
uniform float cellSize;

out vec2 fragTexCoord;

void main()
{
    vec2 instancePosition = instanceData.xy;
    float spriteIndex = instanceData.z;

    // Compute sprite grid position (4x4 grid)
    float spriteX = mod(spriteIndex, 4.0);
    float spriteY = floor(spriteIndex / 4.0);

    // Compute UV sub-rectangle for this sprite
    fragTexCoord = vertexTexCoord / 4.0 + vec2(spriteX / 4.0, spriteY / 4.0);

    // Scale quad by cellSize and translate to world position
    vec3 worldPos = vertexPosition * cellSize
        + vec3(instancePosition.x + cellSize * 0.5, 0.0, instancePosition.y + cellSize * 0.5);

    gl_Position = mvp * vec4(worldPos, 1.0);
}
