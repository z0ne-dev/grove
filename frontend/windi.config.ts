const {defineConfig} = require("windicss/helpers")

module.exports = defineConfig({
    darkMode: "class",
    shortcuts: {},
    theme: {},
    plugins: [],
    preflight: true,
    prefix: "grove-",
    exclude: [/node_modules\/\*\*\/\*/],
    extract: {
        include: ["**/*.jet"],
    }
})