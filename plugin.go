package main

import (
	"fmt"
	"strings"

	"github.com/drone/drone-template-lib/template"
	"github.com/slack-go/slack"
)

type (
	Repo struct {
		Owner string
		Name  string
	}

	Build struct {
		Tag      string
		Event    string
		Number   int
		Commit   string
		Ref      string
		Branch   string
		Author   Author
		Pull     string
		Message  Message
		DeployTo string
		Status   string
		Link     string
		Started  int64
		Created  int64
	}

	Author struct {
		Username string
		Name     string
		Email    string
		Avatar   string
	}

	Message struct {
		msg   string
		Title string
		Body  string
	}

	Config struct {
		//Webhook   string
		Token     string
		Channel   string
		Recipient string
		Username  string
		Template  string
		Fallback  string
		ImageURL  string
		IconURL   string
		IconEmoji string
		Color     string
		LinkNames bool
	}

	Job struct {
		Started int64
	}

	Plugin struct {
		Repo   Repo
		Build  Build
		Config Config
		Job    Job
	}
)

func (a Author) String() string {
	return a.Username
}

func newCommitMessage(m string) Message {
	// not checking the length here
	// as split will always return at least one element
	// check it if using more than the first element
	splitMsg := strings.Split(m, "\n")

	return Message{
		msg:   m,
		Title: strings.TrimSpace(splitMsg[0]),
		Body:  strings.TrimSpace(strings.Join(splitMsg[1:], "\n")),
	}
}
func (m Message) String() string {
	return m.msg
}

func (p Plugin) Exec() error {
	attachment := slack.Attachment{
		Color:      p.Config.Color,
		ImageURL:   p.Config.ImageURL,
		MarkdownIn: []string{"text", "fallback"},
	}
	if p.Config.Color == "" {
		attachment.Color = color(p.Build)
	}
	if p.Config.Fallback != "" {
		f, err := templateMessage(p.Config.Fallback, p)
		if err != nil {
			return err
		}
		attachment.Fallback = f
	} else {
		attachment.Fallback = fallback(p.Repo, p.Build)
	}
	api := slack.New(p.Config.Token)
	// TODO: Convert this over to the new library
	// Username := p.Config.Username
	// IconUrl := p.Config.IconURL
	// IconEmoji := p.Config.IconEmoji

	var channel string
	if p.Config.Recipient != "" {
		// TODO: Check if Recipient is a email address, and only then use the GetUserByEmail
		// channel = prepend("@", p.Config.Recipient)

		// find slack id of commit author
		user, err := api.GetUserByEmail(p.Config.Recipient)

		if err != nil {
			fmt.Printf("%s\n", err)
			return err
		}

		channel = user.ID

	} else if p.Config.Channel != "" {
		channel = prepend("#", p.Config.Channel)
	}
	// TODO: Convert this over to the new library
	// if p.Config.LinkNames == true {
	// 	LinkNames := "1"
	// }
	if p.Config.Template != "" {
		var err error
		attachment.Text, err = templateMessage(p.Config.Template, p)
		if err != nil {
			return err
		}
	} else {
		attachment.Text = message(p.Repo, p.Build)
	}

	channelID, timestamp, err := api.PostMessage(
		channel,
		slack.MsgOptionAttachments(attachment),
	)

	if err != nil {
		fmt.Printf("%s\n", err)
		return err
	}
	fmt.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
	return err
}

// TODO: Revise webhook varient to use the new Slack library

// func (p Plugin) Exec() error {
// 	attachment := slack.Attachment{
// 		Color:      p.Config.Color,
// 		ImageURL:   p.Config.ImageURL,
// 		MarkdownIn: []string{"text", "fallback"},
// 	}
// 	if p.Config.Color == "" {
// 		attachment.Color = color(p.Build)
// 	}
// 	if p.Config.Fallback != "" {
// 		f, err := templateMessage(p.Config.Fallback, p)
// 		if err != nil {
// 			return err
// 		}
// 		attachment.Fallback = f
// 	} else {
// 		attachment.Fallback = fallback(p.Repo, p.Build)
// 	}

// 	payload := slack.WebHookPostPayload{}
// 	payload.Username = p.Config.Username
// 	payload.Attachments = []*slack.Attachment{&attachment}
// 	payload.IconUrl = p.Config.IconURL
// 	payload.IconEmoji = p.Config.IconEmoji

// 	if p.Config.Recipient != "" {
// 		payload.Channel = prepend("@", p.Config.Recipient)
// 	} else if p.Config.Channel != "" {
// 		payload.Channel = prepend("#", p.Config.Channel)
// 	}
// 	if p.Config.LinkNames == true {
// 		payload.LinkNames = "1"
// 	}
// 	if p.Config.Template != "" {
// 		var err error
// 		attachment.Text, err = templateMessage(p.Config.Template, p)
// 		if err != nil {
// 			return err
// 		}
// 	} else {
// 		attachment.Text = message(p.Repo, p.Build)
// 	}

// 	client := slack.NewWebHook(p.Config.Webhook)
// 	return client.PostMessage(&payload)
// }

func detectRef(build Build) string {
	if build.Commit != "" {
		return build.Commit[:8]
	}

	return build.Tag
}

func templateMessage(t string, plugin Plugin) (string, error) {
	return template.RenderTrim(t, plugin)
}

func message(repo Repo, build Build) string {
	return fmt.Sprintf("*%s* <%s|%s/%s#%s> (%s) by %s",
		build.Status,
		build.Link,
		repo.Owner,
		repo.Name,
		detectRef(build),
		build.Branch,
		build.Author,
	)
}

func fallback(repo Repo, build Build) string {
	return fmt.Sprintf("%s %s/%s#%s (%s) by %s",
		build.Status,
		repo.Owner,
		repo.Name,
		detectRef(build),
		build.Branch,
		build.Author,
	)
}

func color(build Build) string {
	switch build.Status {
	case "success":
		return "good"
	case "failure", "error", "killed":
		return "danger"
	default:
		return "warning"
	}
}

func prepend(prefix, s string) string {
	if !strings.HasPrefix(s, prefix) {
		return prefix + s
	}

	return s
}
