import fs from "node:fs/promises";
import path from "node:path";
export async function init(projectName) {
    const projectPath = path.join(process.cwd(), projectName);
    try {
        await fs.mkdir(projectPath);
        console.log(`Created project directory: ${projectPath}`);
    }
    catch (error) {
        console.error(`Error creating project directory: ${error}`);
        process.exit(1);
    }
    console.log(`Creating local-first project ${projectName}...`);
    const templatePath = new URL("../../templates/default", import.meta.url);
    await fs.cp(templatePath, projectPath, { recursive: true });
    console.log(`Project ${projectName} initialized successfully.`);
    console.log(`\n cd ${projectName}`);
}
