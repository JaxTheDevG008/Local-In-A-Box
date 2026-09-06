#!/usr/bin/env node
import { Command } from "commander";
import { init } from "./commands/init.js";
const program = new Command();
program
    .name("LIAB")
    .description("Local-first development toolkit")
    .version("1.0.0");
program
    .command("init <project-name>")
    .description("Initialize a new LIAB project")
    .action(init);
program.parse(process.argv);
