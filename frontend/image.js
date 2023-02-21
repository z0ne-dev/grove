const fs = require("fs")
const yaml = require("yaml")
const sharp = require("sharp")
const path = require("path")

const mode = process.argv[2] || "watch"
if (mode !== "watch" && mode !== "build") {
    console.error("Invalid mode", mode)
    process.exit(1)
}
console.log("Running in", mode, "mode")

generateImages().catch(console.error)
if (mode === "watch") {
    let watchTimeout = null
    fs.watchFile("images.yaml", () => {
        if (watchTimeout) {
            clearTimeout(watchTimeout)
            watchTimeout = null
        }
        watchTimeout = setTimeout(() => {
            generateImages().catch(console.error)
        }, 300)
    })
}

async function generateImages() {
    if (!fs.existsSync("build/images")) {
        fs.mkdirSync("build/images", {recursive: true})
    }

    const content = fs.readFileSync("images.yaml", "utf8")

    /** @type {{input: string, output: string, quality: number, resize: string}[]} */
    const images = yaml.parse(content, {strict: true})

    for (const image of images) {
        console.log(image.input, "=>", image.output, "@", image.resize, "%", image.quality)

        let [w, h] = image.resize.split("x")
        if (w === "") {
            w = null
        } else {
            w = +w
        }
        if (h === "") {
            h = null
        } else {
            h = +h
        }

        await sharp(path.join("images", image.input), {animated: true})
            .resize(w, h, {fit: "outside"})
            .avif({quality: image.quality, lossless: true, effort: 9})
            .toFile(path.join("build/images", image.output))

        // output file size difference
        const inputSize = fs.statSync(path.join("images", image.input)).size
        const outputSize = fs.statSync(path.join("build/images", image.output)).size
        const sizeDiff = Math.round((inputSize - outputSize) / 1024)
        console.log("   Δ", sizeDiff, "KB")
    }
}