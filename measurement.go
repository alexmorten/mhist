package mhist

//Measurement represents a single meassured value in time
type Measurement struct {
	Ts    int64
	Value float64
}

//Reset resets the Measurement to its zero value
func (m *Measurement) Reset() {
	m.Ts = 0
	m.Value = 0
}
