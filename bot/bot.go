package bot

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
)

var TicketCountBugReport int = 0
var TicketCountFeatureRequest int = 0

func checkNilErr(e error) {
	if e != nil {
		log.Fatal("Error message:", e)
	}
}

func Run(token string) {
	discord, err := discordgo.New("Bot " + token)
	checkNilErr(err)

	// Add handlers
	discord.AddHandler(interactionHandler) // Handle slash commands

	// Open session
	err = discord.Open()
	checkNilErr(err)
	defer discord.Close()

	// Register slash commands
	registerCommands(discord)

	fmt.Println("Bot running....")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	// Cleanup commands on exit
	unregisterCommands(discord)
}

func registerCommands(discord *discordgo.Session) {
	commands := []*discordgo.ApplicationCommand{
		{
			Name:        "bug",
			Description: "Submit a bug report.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "name",
					Description: "Name of the bug report.",
					Type:        discordgo.ApplicationCommandOptionString,
					Required:    true,
				},
			},
		},
		{
			Name:        "feature",
			Description: "Submit a feature request.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "name",
					Description: "Name of the feature request.",
					Type:        discordgo.ApplicationCommandOptionString,
					Required:    true,
				},
			},
		},
		{
			Name:        "complete",
			Description: "Mark a bug report or feature request as complete.",
		},
		{
			Name:        "notdoing",
			Description: "Mark a feature request as not going to be implemented.",
		},
	}

	for _, cmd := range commands {
		_, err := discord.ApplicationCommandCreate(discord.State.User.ID, "", cmd) // "" for global commands
		if err != nil {
			log.Fatalf("Cannot create '%v' command: %v", cmd.Name, err)
		}
	}
}

func unregisterCommands(discord *discordgo.Session) {
	commands, err := discord.ApplicationCommands(discord.State.User.ID, "")
	checkNilErr(err)

	for _, cmd := range commands {
		err := discord.ApplicationCommandDelete(discord.State.User.ID, "", cmd.ID)
		checkNilErr(err)
	}
}

func interactionHandler(discord *discordgo.Session, interaction *discordgo.InteractionCreate) {
	if interaction.Type != discordgo.InteractionApplicationCommand {
		return
	}

	switch interaction.ApplicationCommandData().Name {
	case "bug":
		handleBugCommand(discord, interaction)
	case "feature":
		handleFeatureCommand(discord, interaction)
	case "complete":
		handleCompleteCommand(discord, interaction)
	case "notdoing":
		handleNotDoingCommand(discord, interaction)
	}
}

func handleBugCommand(discord *discordgo.Session, interaction *discordgo.InteractionCreate) {
	// Retrieve the name argument
	name := interaction.ApplicationCommandData().Options[0].StringValue()

	// Retrieve the channel (thread) information
	channel, err := discord.Channel(interaction.ChannelID)
	checkNilErr(err)

	// Explicitly join the thread
	err = discord.ThreadMemberAdd(channel.ID, discord.State.User.ID)
	if err != nil {
		fmt.Println("Failed to add bot to thread:", err)
		return
	}

	// Check if the current channel is a thread
	if channel.Type != discordgo.ChannelTypeGuildPublicThread {
		err := discord.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "This command can only be used in a thread.",
			},
		})
		checkNilErr(err)
		return
	}

	// Retrieve the first message in the thread (original post)
	threadMessages, err := discord.ChannelMessages(channel.ID, 1, "", "", "")
	checkNilErr(err)

	var originalPosterMention string
	if len(threadMessages) > 0 {
		originalPosterMention = threadMessages[0].Author.Mention()
	} else {
		originalPosterMention = "unknown user"
	}

	// Increment the bug ticket count
	TicketCountBugReport++

	// Update the thread title
	newTitle := fmt.Sprintf("[VICE-BUG-%d] %s", TicketCountBugReport, name)
	_, err = discord.ChannelEdit(channel.ID, &discordgo.ChannelEdit{
		Name: newTitle,
	})
	checkNilErr(err)

	// Respond to the interaction, mentioning the original poster
	err = discord.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf(
				"Good news, %s! Your bug report was accepted by the Vice Development Team. We will notify you when it’s complete! Your ticket number is %d. Please do not reply to this message.\n\nIf you have any questions, please contact the Vice Development Team on Discord.\n\nThank you for your patience!",
				originalPosterMention,
				TicketCountBugReport,
			),
		},
	})
	checkNilErr(err)

	fmt.Println("Bug report accepted for:", originalPosterMention)
}

func handleFeatureCommand(discord *discordgo.Session, interaction *discordgo.InteractionCreate) {
	// Retrieve the name argument
	name := interaction.ApplicationCommandData().Options[0].StringValue()

	// Retrieve the channel (thread) information
	channel, err := discord.Channel(interaction.ChannelID)
	checkNilErr(err)

	// Explicitly join the thread
	err = discord.ThreadMemberAdd(channel.ID, discord.State.User.ID)
	if err != nil {
		fmt.Println("Failed to add bot to thread:", err)
		return
	}

	// Check if the current channel is a thread
	if channel.Type != discordgo.ChannelTypeGuildPublicThread {
		err := discord.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "This command can only be used in a thread.",
			},
		})
		checkNilErr(err)
		return
	}

	// Retrieve the first message in the thread (original post)
	threadMessages, err := discord.ChannelMessages(channel.ID, 1, "", "", "")
	checkNilErr(err)

	var originalPosterMention string
	if len(threadMessages) > 0 {
		originalPosterMention = threadMessages[0].Author.Mention()
	} else {
		originalPosterMention = "unknown user"
	}

	// Increment the feature request ticket count
	TicketCountFeatureRequest++

	// Update the thread title
	newTitle := fmt.Sprintf("[VICE-FEAT-%d] %s", TicketCountFeatureRequest, name)
	_, err = discord.ChannelEdit(channel.ID, &discordgo.ChannelEdit{
		Name: newTitle,
	})
	checkNilErr(err)

	// Respond to the interaction, mentioning the original poster
	err = discord.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf(
				"Good news, %s! Your feature request was accepted by the Vice Development Team. We will notify you when it’s complete! Your ticket number is %d. Please do not reply to this message.\n\nIf you have any questions, please contact the Vice Development Team on Discord.\n\nThank you for your patience!",
				originalPosterMention,
				TicketCountFeatureRequest,
			),
		},
	})
	checkNilErr(err)

	fmt.Println("Feature request accepted for:", originalPosterMention)
}

func handleCompleteCommand(discord *discordgo.Session, interaction *discordgo.InteractionCreate) {
	// Retrieve the channel (thread) information
	channel, err := discord.Channel(interaction.ChannelID)
	checkNilErr(err)

	// Check if the current channel is a thread
	if channel.Type != discordgo.ChannelTypeGuildPublicThread {
		err := discord.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "This command can only be used in a thread.",
			},
		})
		checkNilErr(err)
		return
	}

	// Instead of fetching only 1 message (which might be the bot’s),
	// fetch a larger number of messages (e.g., 50).
	// These are returned in descending order (newest -> oldest).
	threadMessages, err := discord.ChannelMessages(channel.ID, 50, "", "", "")
	checkNilErr(err)

	var originalPosterMention string
	if len(threadMessages) > 0 {
		// The last element in the slice is the oldest message in the fetch
		oldestMessage := threadMessages[len(threadMessages)-1]
		originalPosterMention = oldestMessage.Author.Mention()
	} else {
		originalPosterMention = "unknown user"
	}

	// Update the thread title to mark it as complete
	newTitle := fmt.Sprintf("✅ %s", channel.Name)
	_, err = discord.ChannelEdit(channel.ID, &discordgo.ChannelEdit{
		Name: newTitle,
	})
	checkNilErr(err)

	// Notify the original poster
	err = discord.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf(
				"Good news, %s! Your bug report or feature request has been marked as complete. Thank you for your patience! ✅",
				originalPosterMention,
			),
		},
	})
	checkNilErr(err)

	fmt.Println("Bug report or feature request marked as complete for:", originalPosterMention)
}

func handleNotDoingCommand(discord *discordgo.Session, interaction *discordgo.InteractionCreate) {
	// Retrieve the channel (thread) information
	channel, err := discord.Channel(interaction.ChannelID)
	checkNilErr(err)

	// Check if the current channel is a thread
	if channel.Type != discordgo.ChannelTypeGuildPublicThread {
		err := discord.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "This command can only be used in a thread.",
			},
		})
		checkNilErr(err)
		return
	}

	// Retrieve the first message in the thread (original post)
	threadMessages, err := discord.ChannelMessages(channel.ID, 1, "", "", "")
	checkNilErr(err)

	var originalPosterMention string
	if len(threadMessages) > 0 {
		originalPosterMention = threadMessages[0].Author.Mention()
	} else {
		originalPosterMention = "unknown user"
	}

	// Update the thread title with an "X" suffix
	newTitle := fmt.Sprintf("%s ❌", channel.Name)
	_, err = discord.ChannelEdit(channel.ID, &discordgo.ChannelEdit{
		Name: newTitle,
	})
	checkNilErr(err)

	// Notify the original poster
	err = discord.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf(
				"Unfortunately, %s, your feature request will not be implemented at this time. We appreciate your suggestion and encourage you to keep submitting ideas in the future.",
				originalPosterMention,
			),
		},
	})
	checkNilErr(err)

	fmt.Println("Feature request marked as not doing for:", originalPosterMention)
}
