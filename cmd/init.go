package cmd

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"github.com/manifoldco/promptui"
)

func init() {
	RootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a secure, local-first application workspace",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to the LIAB wizard! Let's get started with setting up your local-first application.")
		scanner := bufio.NewScanner(os.Stdin)
		validPattern := "^[a-zA-Z0-9]+$"
		var projectName string

		for {
			fmt.Print("What is your project name? ")
			scanner.Scan()
			input := strings.TrimSpace(scanner.Text())

			isValid, _ := regexp.MatchString(validPattern, input)

			if input == "" {
				fmt.Println("No project name provided. Please enter a valid project name.")
				continue
			}
			if !isValid {
				fmt.Println("Invalid project name. Please use only alphanumeric characters.")
				continue
			}

			projectName = input
			break
		}

			templates := &promptui.SelectTemplates{
				Label:    "{{ . }}?",                               
				Active:   "➡️  \033[36m{{ . }}\033[0m",              
				Inactive: "   {{ . }}",                             
				Selected: "✔️  Selected: \033[32m{{ . }}\033[0m",
			}

			frameworks := []string{"React", "Vue", "Angular", "TypeScript", "JavaScript"}

			prompt := promptui.Select{
				Label: "Great name choice! Now select a framework for your project",
				Items: frameworks,
				Templates: templates,
			}

			_, framework, err := prompt.Run()
			if err != nil {
				fmt.Printf("Error occurred while selecting a framework: %v\n", err)
				return
			}
			fmt.Printf("You have selected %s as your framework.\n", framework)
			fmt.Println("Now select a database for your project (if applicable). If you don't want to use a database, you can skip this step.")

			databases := []string{"RXDB", "SQLite", "IndexedDB", "PouchDB", "None"}

			promptDB := promptui.Select{
				Label: "Select a database for your project",
				Items: databases,
				Templates: templates,
			}

			_, database, err := promptDB.Run()
			if err != nil {
				fmt.Printf("Error occurred while selecting a database: %v\n", err)
				return
			}
			if database == "None" {
				fmt.Println("You have chosen not to use a database for your project.")
			} else {
				fmt.Printf("You have selected %s as your database.\n", database)
			}

			fmt.Printf("Project \"%s\" has been initialized successfully with %s and %s.\n", projectName, framework, database)
	},
}
