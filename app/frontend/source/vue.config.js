const path = require("path")

console.log(__dirname)
console.log(path.resolve(__dirname, "../volume/target"))

module.exports = {
  outputDir: path.resolve(__dirname, "../volume/target"),
  transpileDependencies: [
    'vuetify'
  ]
}
