# Goflakes
Generate snowflakes with go.
## Functions
NewSnowflakeGenerator()
does what it sounds like

Generate()
Generates a new snowflakes

GenerateMultiple()
Generate multiple snowflakes and put them all into a slice

AsyncGenerate()
Similar to GenerateMultiple, but sends snowflakes over a channel.

# CAUTION
Currently generating there is only protection against generating repeat snowflakes in the context of a single call, there is *__currently__* no protection against generating repeats when running multiple generating function in paralllel.