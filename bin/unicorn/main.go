package main

import (
	"fmt"

	"github.com/jmhobbs/pinentry-client/pinentry"
)

func main() {
	pe := pinentry.New("pinentry-mac").
		SetDescription("What's your favorite mythological animal?").
		SetPrompt("Animal:").
		SetButtonOk("They're the best.").
		SetButonCancel("I'm not telling you.")

	defer pe.Close()

	secret, err := pe.GetPIN()
	if err != nil {
		panic(err)
	}

	pe.SetDescription(fmt.Sprintf("Are you positive it is a %q?", secret)).
		SetButtonOk("Yes, ugh, stop asking!").
		SetButonNotOk("No, it's not").
		SetButonCancel("I don't want to tell you.")

	confirm, err := pe.Confirm()
	if err != nil {
		panic(err)
	}

	if confirm {
		fmt.Printf("Your favorite mythological animal is a %s!\n", secret)
	} else {
		fmt.Printf("That's ok, you don't have to tell me.")
	}
}
