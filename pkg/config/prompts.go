package config

func wrapBlockCode(t string, code string) string {
	return "```" + t + "\n" + code + "\n```"
}

var COMMIT_PROMPT = wrapBlockCode("diff", "{{diff}}") + "\n\n" + "Write commit message for the change with commitizen convention. Make sure the title has maximum 50 characters and message is wrapped at 72 characters. Don't wrap the message in code block."
