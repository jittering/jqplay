const tailwindcss = require('tailwindcss');
const autoprefixer = require('autoprefixer');

console.log("* loading postcss.config.cjs");

const config = {
  plugins: [
    //Some plugins, like tailwindcss/nesting, need to run before Tailwind,
    require("postcss-import")(),
    require("postcss-url")(),
    require("postcss-input-range")(),

    require("tailwindcss/nesting")(),
    tailwindcss,

    //But others, like autoprefixer, need to run after,
    autoprefixer
  ],
  extract: true,
};

module.exports = config;
