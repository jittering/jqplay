const smelteTailwind = require('smelte/tailwind.config');

console.log("* loading tailwind.config.cjs");

const config = {
  content: [
    './src/**/*.{html,js,svelte,ts}',
    './node_modules/smelte/src/**/*.{html,js.svelte,ts}'
  ],

  safelist: ["nav", "hidden", {
    pattern: /((dark(-(hover|focus))?):)?((bg|text)-)?(primary|secondary|error|success|alert|blue|dark)(-.*?)?/
  }, ],

  ...smelteTailwind({
    // default colors from smelte/tailwind.config
    colors: {
      primary: "#b027b0",
      secondary: "#009688",
      error: "#f44336",
      success: "#4caf50",
      alert: "#ff9800",
      blue: "#2196f3",
      dark: "#212121"
    },

    // whether or not dark mode classes should be generated
    // see smelte docs
    darkMode: true,
  }),

};

module.exports = config;
