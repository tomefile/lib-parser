package libparser

type FormatModifier func(string) string

func GetModifier(name string, args []string) FormatModifier {
	return nil
}
