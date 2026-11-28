// Byte budget of the volatile Redis message copy (SAR-REDIS-CAP-001): a bound on each record's stored JSON.
package relay

// Byte budget of the Redis message copy (SAR-REDIS-CAP-001). Redis runs maxmemory 128 MiB with noeviction and holds two
// copies during a commit (writeEphemeral, then retire), so a 32 MiB copy keeps the peak at half the limit. One sending
// owner holds at most 6 MiB, so at least five owners fit before the shared budget refuses everyone, like ownerQueueLimit.
const (
	messageBytesLimit      = 32 << 20
	ownerMessageBytesLimit = 6 << 20
	// messageOverhead bounds one record's fixed metadata, idempotency entry and remote Inbox map key in the Redis JSON
	// (about 1.4 KiB at the longest keys and IDs); TestMessageBytesBoundStoredJSON checks it. 20000 records fit in 30 MiB.
	messageOverhead = 1536
)

// storedLen bounds the JSON size of b in the Redis copy: encoding/json writes <, >, & and control bytes as \u00XX,
// U+2028 and U+2029 (lead byte 0xE2) as a 6-byte escape, and one more byte for quotes and backslashes inside strings.
func storedLen[T ~string | ~[]byte](b T) int {
	if len(b) == 0 {
		return 0 // omitted (omitempty) or nil
	}
	n := len(b) + 2
	for i := 0; i < len(b); i++ {
		switch c := b[i]; {
		case c < 0x20 || c == '<' || c == '>' || c == '&':
			n += 5
		case c == '"' || c == '\\':
			n++
		case c == 0xE2:
			n += 3
		}
	}
	return n
}

// messageBytes bounds what m holds in Redis now plus the Inbox copy a later persist or remote receive adds.
func messageBytes(m *Message) int {
	n := messageOverhead + storedLen(m.Envelope) + storedLen(m.Inbox)
	if len(m.Envelope) > 0 && len(m.Inbox) == 0 {
		n += storedLen(m.Envelope)
	}
	if c := m.Conversation; c != nil {
		n += storedLen(c.Text) + storedLen(c.PersonA) + storedLen(c.PersonB) + storedLen(c.AgentA) + storedLen(c.AgentB)
	}
	return n
}

// retainedBytes sums messageBytes and remote Inbox copies, in total and for the records sent by owner's agents.
func (st *State) retainedBytes(owner string) (total, mine int) {
	charge := func(from string, n int) {
		total += n
		if a := st.Agents[from]; owner != "" && a != nil && a.Owner == owner {
			mine += n
		}
	}
	for _, m := range st.Messages {
		charge(m.Receipt.From, messageBytes(m))
	}
	if st.Remote != nil {
		for _, i := range st.Remote.Installs {
			for id, raw := range i.Inbox {
				from := ""
				if m := st.Messages[id]; m != nil {
					from = m.Receipt.From
				}
				charge(from, storedLen(raw))
			}
		}
	}
	return total, mine
}
