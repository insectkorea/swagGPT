package api

import (
	"fmt"
	"strings"
)

const systemPrompt = `
You are a helpful assistant for generating Swagger annotation comments for Go handler functions.
`

const userPromptTemplate = `
Here is an example of a Go handler function:
// ListAccounts lists all existing accounts
//
//  @Summary      List accounts
//  @Description  get accounts
//  @Tags         accounts
//  @Accept       json
//  @Produce      json
//  @Param        q    query     string  false  "name search by q"  Format(email)
//  @Success      200  {array}   model.Account
//  @Failure      400  {object}  httputil.HTTPError
//  @Failure      404  {object}  httputil.HTTPError
//  @Failure      500  {object}  httputil.HTTPError
//  @Router       /accounts [get]
Generate Swagger comments for the following function:
%s
Do not add any explanation and just return the Swagger comments. Do not wrap it in Markdown. Do not return the function itself.
If it is not a handler function(e.g. a function in a test file or a helper function), return nothing.
`

// EstimateTokens estimates the number of tokens for generating Swagger comments.
func EstimateTokens(functionContent string) int {
	userPrompt := fmt.Sprintf(userPromptTemplate,
		functionContent,
	)
	prompt := userPrompt + systemPrompt

	// Estimate tokens based on number of words
	words := len(strings.Split(prompt, " "))
	// Average tokens per word is approximately 1.33 for English text
	// But we need to be conservative, so we multiply by 1.5
	return int(float64(words) * 1.5)
}
