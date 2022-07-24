# Mango
A Dota 2 replay parser in Go

Inspired by Smoke and Clarity by Skadistats

There is a Go version made by Dotabuff but I haven't looked at that to make making this more fun :D (and I didn't know it existed when I started this...)

A work in progress - currently can attempt to parse the entire replay and will error when an unknown packet type is reached (not all have been added to the mappings yet). Currently does nothing with the resultant packet data, but can be accessed and/or displayed by looping over the output of the `ReplayParser.ParseReplay` function.

