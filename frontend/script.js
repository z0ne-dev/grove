/*
 * script.js Copyright (c) 2023 z0ne.
 * All Rights Reserved.
 * Licensed under the EUPL 1.2 License.
 * See LICENSE the project root for license information.
 *
 * SPDX-License-Identifier: EUPL-1.2
 */

const fs = require("fs")
const path = require("path")

const mode = process.argv[2] || "watch"
if (mode !== "watch" && mode !== "build") {
    console.error("Invalid mode", mode)
    process.exit(1)
}
console.log("Running in", mode, "mode")

copyScripts().catch(console.error)
if (mode === "watch") {
    let watchTimeout = null
    fs.watchFile("pnpm-lock.yaml", () => {
        if (watchTimeout) {
            clearTimeout(watchTimeout)
            watchTimeout = null
        }
        watchTimeout = setTimeout(() => {
            copyScripts().catch(console.error)
        }, 300)
    })
}

async function copyScripts() {
    if (!fs.existsSync("build/scripts")) {
        fs.mkdirSync("build/scripts", {recursive: true})
    }

    const content = fs.readFileSync("package.json", "utf8")

    /** @type {string[]} */
    const dependencies = Object.keys(JSON.parse(content).dependencies)

    for (const dependency of dependencies) {
        const packageJson = fs.readFileSync(path.join("node_modules", dependency, "package.json"), "utf8")
        const packageJsonData = JSON.parse(packageJson)
        const main = packageJsonData.module || packageJsonData.main

        if (!main) {
            console.error("No main file for", dependency)
        }

        console.log(dependency, "=>", main)

        fs.copyFileSync(path.join("node_modules", dependency, main), path.join("build/scripts", dependency + ".js"))
    }
}