package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(initCmd)
}

func copyDir(src, dst string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := src + "/" + entry.Name()
		dstPath := dst + "/" + entry.Name()

		if entry.Name() == "package.json.txt" || entry.Name() == "node_modules" { 
			continue 
		}

		if entry.IsDir() {
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			fileData, err := os.ReadFile(srcPath)
			if err != nil {
				return err
			}
			if err := os.WriteFile(dstPath, fileData, 0644); err != nil {
				return err
			}
		}
	}
	return nil
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

		if err := scanner.Err(); err != nil {
			fmt.Printf("Error reading input: %v\n", err)
			return
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

        projectRootDir := projectName

		if _, err := os.Stat(projectRootDir); !os.IsNotExist(err) {
			fmt.Printf("Error: The target path %s already exists. Please choose a different project name or remove the existing directory.\n", projectRootDir)
			return
		}

		targetStoragePath := fmt.Sprintf("%s/src/storage.ts", projectRootDir)

		if err := os.MkdirAll(fmt.Sprintf("%s/src", projectRootDir), 0755); err != nil {
			fmt.Printf("Error creating project directory: %v\n", err)
			return
		}

		var sourceDatabaseFile string

		pkgBytes, err := os.ReadFile("templates/package.json.txt")
		if err != nil {
			fmt.Printf("Error reading package.json template: %v\n", err)
			return
		}

		pkgString := string(pkgBytes)

		var scripts, coreDeps, devDeps, extraDeps string

		switch framework {
		case "React":
			scripts = "    \"dev\": \"vite\",\n    \"build\": \"tsc && vite build\",\n    \"preview\": \"vite preview\""
			coreDeps = "    \"react\": \"^19.0.0\",\n    \"react-dom\": \"^19.0.0\""
			devDeps = "    \"@types/react\": \"^19.0.0\",\n    \"@types/react-dom\": \"^19.0.0\",\n    \"@vitejs/plugin-react\": \"^4.3.0\",\n    \"typescript\": \"^5.0.0\",\n    \"vite\": \"^5.0.0\""

			err = copyDir("templates/react", projectRootDir)
			if err != nil {
				fmt.Printf("Error scaffolding application from template: %v\n", err)
				return
			}
		default: 
			scripts = "    \"dev\": \"vite\",\n    \"build\": \"tsc && vite build\""
			coreDeps = "    \"canvas-confetti\": \"^1.9.0\""
			devDeps = "    \"@types/canvas-confetti\": \"^1.6.0\",\n    \"typescript\": \"^5.0.0\",\n    \"vite\": \"^5.0.0\""
		}
	
		switch database {
		case "IndexedDB":
			sourceDatabaseFile = "templates/databases/indexeddb.txt"
			extraDeps = ""
		case "RXDB":
			sourceDatabaseFile = "templates/databases/rxdb.txt"
			extraDeps = ",\n    \"rxdb\": \"^15.0.0\",\n    \"dexie\": \"^4.0.0\""
		default:
			sourceDatabaseFile = ""
			extraDeps = ""
		}

		pkgString = strings.ReplaceAll(pkgString, "{{projectName}}", projectName)
		pkgString = strings.ReplaceAll(pkgString, "  \"scripts\": \"{{scripts}}\"", "  \"scripts\": {\n"+scripts+"\n  }")
		pkgString = strings.ReplaceAll(pkgString, "  \"dependencies\": \"{{coreDependencies}}{{extraDependencies}}\"", "  \"dependencies\": {\n"+coreDeps+extraDeps+"\n  }")
		pkgString = strings.ReplaceAll(pkgString, "  \"devDependencies\": \"{{devDependencies}}\"", "  \"devDependencies\": {\n"+devDeps+"\n  }")

		targetPkgPath := fmt.Sprintf("%s/package.json", projectRootDir)
		err = os.WriteFile(targetPkgPath, []byte(pkgString), 0644)
		if err != nil {
			fmt.Printf("Error writing package.json to %s: %v\n", targetPkgPath, err)
			return
		}
		fmt.Printf("package.json has been written to %s successfully!\n", targetPkgPath)

		if sourceDatabaseFile != "" {
			fmt.Printf("Database layer requested. Reading boilerplate from %s...\n", sourceDatabaseFile)
			dbCode, err := os.ReadFile(sourceDatabaseFile)
			if err != nil {
				fmt.Printf("Error reading database template: %v\n", err)
				return
			}

			err = os.WriteFile(targetStoragePath, dbCode, 0644)
			if err != nil {
				fmt.Printf("Error writing database code to %s: %v\n", targetStoragePath, err)
				return
			}
			fmt.Printf("Database code has been written to %s successfully!\n", targetStoragePath)
		} else {
			fmt.Println("Skipping database layer setup as no database was selected.")
		}

		templateBytes, err := os.ReadFile(sourceDatabaseFile)
		if err != nil {
			fmt.Printf("Error reading database template: %v\n", err)
			return
		}

		templateString := string(templateBytes)

		finalString := strings.ReplaceAll(templateString, "{{projectName}}", projectName)
		finalBytes := []byte(finalString)
		err = os.WriteFile(targetStoragePath, finalBytes, 0644)
		if err != nil {
			fmt.Printf("Error writing final code to %s: %v\n", targetStoragePath, err)
			return
		}

		fmt.Println("Installing npm dependencies...this may take a moment...")
		npmCmd := exec.Command("npm", "install")
		npmCmd.Dir = projectRootDir
		npmCmd.Stdout = os.Stdout
		npmCmd.Stderr = os.Stderr
		err = npmCmd.Run()
		if err != nil {
			fmt.Printf("Error occurred while installing npm dependencies: %v\n", err)
			fmt.Println("Please ensure you have npm installed and try running \"npm install\" manually in the project directory.")
			return
		} else {
			fmt.Println("npm dependencies installed successfully!")
		}

		targetHtmlPath := fmt.Sprintf("%s/index.html", projectRootDir)
		htmlBytes, err := os.ReadFile(targetHtmlPath)
		if err == nil {
			finalHtmlString := strings.ReplaceAll(string(htmlBytes), "{{projectName}}", projectName)
			err = os.WriteFile(targetHtmlPath, []byte(finalHtmlString), 0644)
			if err != nil {
				fmt.Printf("Error writing final HTML to %s: %v\n", targetHtmlPath, err)
				return
			}
		}

		fmt.Printf("Project directory %s has been created successfully.\n", projectRootDir)

		fmt.Printf("Project \"%s\" has been initialized successfully with %s and %s.\n", projectName, framework, database)
	},
}
