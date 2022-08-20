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

	EmojiFatal    = "💀"
	EmojiError    = "❌"
	EmojiWarn     = "🟡"
	EmojiOkay     = "✅"
	EmojiInfo     = "💬"
	EmojiProgress = "⌛️"
	EmojiDebug    = "🐛"
)

var Emojis = map[Level]string{
	LevelFatal:    EmojiFatal,
	LevelError:    EmojiError,
	LevelWarn:     EmojiWarn,
	LevelOkay:     EmojiOkay,
	LevelInfo:     EmojiInfo,
	LevelProgress: EmojiProgress,
	LevelDebug:    EmojiDebug,
}

func GetLevelEmoji(l Level) string {
	return Emojis[l]
}

/*
⭕ ✅ ❌ ❎
🔴 🟠 🟡 🟢 🔵 🟣 🟤 ⚫ ⚪
🟥 🟧 🟨 🟩 🟦 🟪 🟫 ⬛ ⬜ ◾ ◽
🔶 🔷 🔸 🔹 🔺 🔻 💠 🔘 🔳 🔲
*/
