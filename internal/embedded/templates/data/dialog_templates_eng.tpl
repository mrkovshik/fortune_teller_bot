{{define "help"}}
There’s an old tradition of fortune-telling with books: you ask a question, pick a random page and line, and the book gives you an answer.
Here it works almost the same way :) You can type in your question and let the bot use it to generate a “random” prediction, or you can just ask for a random quote from a book.
What would you like to do?
{{end}}

{{define "start"}}
What would you like to do?
{{end}}

{{define "invalid_button"}}
That option doesn’t work. Please use the last menu or start over with /start
{{end}}

{{define "select_book_for_question"}}
Which book should we use to answer your question?
{{end}}

{{define "select_book_for_random"}}
Which book should we use to give you a random quote?
{{end}}

{{define "type_question"}}
Type your question, and the book will give you an answer.
{{end}}

{{define "list_books"}}
Which books would you like to choose from?
{{end}}

{{define "back"}}
Going back…
{{end}}

{{define "error"}}
⚠️ Something went wrong. Please try again later
{{end}}

{{define "quote"}}
{{.Text}}
    {{.Title}}
{{end}}

{{define "change_language"}}
Bot language has been changed to {{.}}
{{end}}