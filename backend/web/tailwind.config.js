/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./internal/templ/**/*.{templ,go}",
    "../internal/templ/**/*.{templ,go}",
    "./static/**/*.{html,js}",
  ],
  theme: {
    extend: {},
  },
  plugins: [],
};
