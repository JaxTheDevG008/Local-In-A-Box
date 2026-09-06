import fs from "node:fs/promises";
import path from "node:path";

export async function init(projectName: string) {
  const projectPath = path.join(process.cwd(), projectName);

  try {
    await fs.mkdir(projectPath);
    console.log(`Created project directory: ${projectPath}`);
  } catch (error) {
    console.error(`Error creating project directory: ${error}`);
    process.exit(1);
  }

  console.log(`Creating local-first project ${projectName}...`);

  const templatePath = new URL("../../templates/default", import.meta.url);
  await fs.cp(templatePath, projectPath, { recursive: true });
  const filesToRename = ["package.json", "vite.config.js", "src/index.html", "src/storage.ts"];
  for (const file of filesToRename) {
    const filePath = path.join(projectPath, file);
    try {
      const content = await fs.readFile(filePath, "utf-8");
      const updatedContent = content.replaceAll("{{PROJECT_NAME}}", projectName);
      await fs.writeFile(filePath, updatedContent, "utf-8");
    } catch (error) {
      console.error(`Error updating ${file}: ${error}`);
      process.exit(1);
    }
  }

  console.log(`Project ${projectName} initialized successfully.`);
  console.log(`\n cd ${projectName}`);
}