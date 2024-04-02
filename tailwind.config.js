const colors = require("tailwindcss/colors");

/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["internal/adapters/primary/web/**/*.templ"],
  theme: {
    extend: {
      colors: {
        primary: "#FFFFFF",
        secondary: "#FDF4F5",
        neutral: colors.gray,
      },
    },
  },
  plugins: [require("@tailwindcss/forms"), require("@tailwindcss/typography")],
};
