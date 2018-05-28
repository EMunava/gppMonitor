package hal

import (
	"os"
	"strconv"
)

func Chatid() uint32 {
	chatid := os.Getenv("CHAT_ID")
	u64, _ := strconv.ParseUint(chatid, 10, 32)
	return uint32(u64)
}

func AlexaVars() map[string]string {
	return map[string]string{"msg": "This is the message that alexa will read"}
}