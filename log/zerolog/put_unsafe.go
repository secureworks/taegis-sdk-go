// +build unsafe

package zerolog

import (
	//lint wants me to add a comment so others can commit the same sin
	//unsafe needs to be at least underscore imported to use go:linkname
	_ "unsafe"

	"github.com/rs/zerolog"
)

//Dear reader, please do as I say not as I do, I really didn't want to use this trick ☹️
//This will cause a build time issue if zerolog changes the private function, not a runtime one.
//Always have a fallback or "safe" version

//Since we want the zerolog impl to support both changing levels or "Disabled" https://github.com/rs/zerolog/pull/255
//this is needed in order to retain zerologs performance for high event counts, if we don't write it won't be reused
//https://github.com/rs/zerolog/blob/7825d863376faee2723fc99c061c538bd80812c8/event.go#L79

//go:linkname putEvent github.com/rs/zerolog.putEvent
func putEvent(ent *zerolog.Event)
