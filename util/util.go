package util

func IsDigit(chr rune) bool {
    return chr >= '0' && chr <= '9'
}

func IsAlpha(chr rune) bool {
    return (chr >= 'a' && chr <= 'z') ||
           (chr >= 'A' && chr <= 'Z')
}

func IsIdentStart(chr rune) bool {
    return IsAlpha(chr) || chr == '_'
}

func IsIdentChar(chr rune) bool {
    return IsAlpha(chr) || IsDigit(chr) || chr == '_'
}
