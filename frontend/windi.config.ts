/*
 * windi.config.ts Copyright (c) 2023 z0ne.
 * All Rights Reserved.
 * Licensed under the EUPL 1.2 License.
 * See LICENSE the project root for license information.
 *
 * SPDX-License-Identifier: EUPL-1.2
 */

const {defineConfig} = require("windicss/helpers")

module.exports = defineConfig({
    darkMode: "class",
    shortcuts: {},
    theme: {},
    plugins: [],
    preflight: true,
    extract: {
        include: ["../**/*.jet"],
        exclude: ["../**/*.windi.jet"]
    }
})