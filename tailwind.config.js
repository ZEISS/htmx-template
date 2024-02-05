/** @type {import('tailwindcss').Config} */
module.exports = {
    content: ['./public/**/*.html'],
    purge: ['./public/**/*.html'],
    theme: {
        extend: {},
    },
    plugins: [
        require("@tailwindcss/typography"),
        require("daisyui")
    ],
    darkMode: ["class", '[data-theme="dark"]'],
    daisyui: {
        themes: ["light", "dark"],
    }
}