package utils

type FlagDef struct {
	Name      string
	Shorthand string
	Usage     string
}

var (
	QuantumSafeFlag = FlagDef{
		Name:      "quantum-safe",
		Shorthand: "q",
		Usage:     "Use post-quantum cryptography (Kyber1024)",
	}
)
