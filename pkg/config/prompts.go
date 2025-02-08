package config

func wrapBlockCode(t string, code string) string {
	return "```" + t + "\n" + code + "\n```"
}

var COMMIT_PROMPT = wrapBlockCode("diff", "{{diff}}") + "\n\n" + "Write a concise and informative commit message for the change with commitizen convention. If multiple files are changed, provide a summary of the changes without being too specific per-file changes. Ensure the message is readable and clearly conveys the purpose of the changes. Make sure the title has maximum 50 characters and message is wrapped at 72 characters. DON'T WRAP IN CODE BLOCK."
