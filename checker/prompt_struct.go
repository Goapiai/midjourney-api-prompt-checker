package checker

import (
	"errors"
	"strings"
)

const (
	ErrPromptEmpty           = "Prompt is empty, please enter a prompt"
	ErrPromptEmptyWithParams = "Prompt is empty, not allowed to start with -- params"
	ErrPromptTooLong         = "Prompt is too long, limited to 6000 characters"
)

func CheckPromptLength(prompt string) error {
	if len(prompt) > 6000 {
		return errors.New(ErrPromptTooLong)
	}
	return nil
}

func CheckPromptEmpty(prompt string) error {
	if len(prompt) == 0 {
		return errors.New(ErrPromptEmpty)
	}
	return nil
}

func CheckPromptEmptyWithParams(prompt string) error {
	if strings.HasPrefix(prompt, "--") {
		return errors.New(ErrPromptEmptyWithParams)
	}
	return nil
}

func CheckPromptStruct(prompt string, allowEmpty bool) error {
	if !allowEmpty {
		if err := CheckPromptEmpty(prompt); err != nil {
			return err
		}
	}
	if err := CheckPromptEmptyWithParams(prompt); err != nil {
		return err
	}
	if err := CheckPromptLength(prompt); err != nil {
		return err
	}
	return nil
}
