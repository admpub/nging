package log

// Success, Warning, Error can also be summary items.
// Grn, Ylw, Red are calm B/G indicator lights .
const (
	// STATE INDICATORS
	Red = "🔴"
	Ylw = "🟡"
	Grn = "🟢"
)

func EmojiOfLevel(L Level) string {
	switch L {
	case LevelFatal:
		return "💀❌💀"
	case LevelError:
		return "❌"
	case LevelWarn:
		return "🟨"
	case LevelOkay:
		return "🟩"
	case LevelInfo:
		return "💬"
	case LevelProgress:
		return "〰️"
	case LevelDebug:
		return "❓"
	}
	return ""
}

/*
⭕ ✅ ❌ ❎
🔴 🟠 🟡 🟢 🔵 🟣 🟤 ⚫ ⚪
🟥 🟧 🟨 🟩 🟦 🟪 🟫 ⬛ ⬜ ◾ ◽
🔶 🔷 🔸 🔹 🔺 🔻 💠 🔘 🔳 🔲
*/
