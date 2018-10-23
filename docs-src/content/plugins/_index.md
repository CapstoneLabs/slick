---
title: "Plugins"
weight: 10
---

Creating a new plugin is as simple as:

1. Create a new go package (no `main()` function)
2. Include a `slick.Listener`

```go
// listenDeploy was hooked into a plugin elsewhere..
func listenDeploy() {
	keywords := []string{"project1", "project2", "project3"}
	bot.Listen(&slick.Listener{
		Matches:        regexp.MustCompile("(can you|could you|please|plz|c'mon|icanhaz) deploy (" + strings.Join(keywords, "|") + ") (with|using)( revision| commit)? `?([a-z0-9]{4,42})`?"),
		MentionsMeOnly: true,
		MessageHandlerFunc: func(listen *slick.Listener, msg *slick.Message) {

			projectName := msg.Match[2]
			revision := msg.Match[5]

			go func() {
				go msg.AddReaction("work_hard")
				defer msg.RemoveReaction("work_hard")

				// Do the deployment with projectName and revision...

			}()
		},
	})
}
```