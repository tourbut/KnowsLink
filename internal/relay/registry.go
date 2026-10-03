// Registry exposes the central frozen seed; adding an intent never creates an execution handler.
package relay

import (
	"crypto/sha256"
	_ "embed"
	"encoding/hex"
)

//go:embed registry.json
var registry []byte

const RegistrySHA256 = "b9759a1ec4035281c704d23f35f78921d8b54d232fa4cd0f27002b183e5eafda"

func RegistryValid() bool {
	sum := sha256.Sum256(registry)
	return hex.EncodeToString(sum[:]) == RegistrySHA256
}
