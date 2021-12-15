package log

// Success, Warning, Error can also be summary items.
// Grn, Ylw, Red are calm B/G indicator lights .
const (
	// STATE INDICATORS
	Red = "🔴"
	Ylw = "🟡"
	Blu = "🔵"
	Grn = "🟢"
	Org = "🟠"
	Pnk = "🟣"
)

var Emojis = map[Level]string{
	LevelFatal:    "💀❌💀",
	LevelError:    "❌",
	LevelWarn:     "🟨",
	LevelOkay:     "🟩",
	LevelInfo:     "💬",
	LevelProgress: "〰️",
	LevelDebug:    "❓",
}

func GetLevelEmoji(l Level) string {
	emoji, _ := Emojis[l]
	return emoji
}

/*
⭕ ✅ ❌ ❎
🔴 🟠 🟡 🟢 🔵 🟣 🟤 ⚫ ⚪
🟥 🟧 🟨 🟩 🟦 🟪 🟫 ⬛ ⬜ ◾ ◽
🔶 🔷 🔸 🔹 🔺 🔻 💠 🔘 🔳 🔲
*/
