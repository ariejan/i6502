package i6502

/*
ACIA 6551 Serial IO
*/
type Acia6551 struct {
	rx chan byte // Reading (Acia Input) line
	tx chan byte // Transmitting (Acia Output) line
}

func NewAcia6551(rx chan byte, tx chan byte) (*Acia6551, error) {
	return &Acia6551{tx: tx, rx: rx}, nil
}

func (r *Acia6551) Size() uint16 {
	// We have a only 4 addresses, RX, TX, Command and Control
	return 0x04
}

/*
func (r *Rom) Read(address uint16) byte {
	return r.data[address]
}

func (r *Rom) Write(address uint16, data byte) {
	panic(fmt.Errorf("Trying to write to ROM at 0x%04X", address))
}
*/
