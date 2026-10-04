// Registry exposes the central frozen seed; adding an intent never creates an execution handler.
package relay

import (
	"crypto/sha256"
	_ "embed"
	"encoding/hex"
)

//go:embed registry.json
var registry []byte

const RegistrySHA256 = "2b25c6b58fb6b6973de1c9bac19162834cfc9b678a0416ad39e2cffab3cde1e2"

func RegistryValid() bool {
	sum := sha256.Sum256(registry)
	return hex.EncodeToString(sum[:]) == RegistrySHA256
}
