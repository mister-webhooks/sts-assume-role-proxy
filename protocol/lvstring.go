package protocol

import "encoding/binary"

type LVString string

func (lvs LVString) MarshalBinary() ([]byte, error) {
	encoded := make([]byte, 0, 4+len(lvs))
	encoded = binary.BigEndian.AppendUint32(encoded, uint32(len(lvs)))
	encoded = append(encoded, []byte(lvs)...)

	return encoded, nil
}

func (lvs *LVString) UnmarshalBinary(data []byte) error {
	strlen := binary.BigEndian.Uint32(data[0:4])

	*lvs = LVString((data[4:(4 + strlen)]))

	return nil
}

func (lvs LVString) EncodedLength() int {
	return 4 + len(lvs)
}
