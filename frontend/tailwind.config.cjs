/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./src/**/*.{js,jsx,ts,tsx}",
  ],
  theme: {
    extend: {
      keyframes: {
        refresh_rotate_scaled: {
          "0%": { transform: "rotate(0deg) scale(0.75)" },
          "100%": { transform: "rotate(180deg) scale(0.75)" },
        }
      },
      animation: {
        refresh_rotate_scaled: "refresh_rotate_scaled 500ms ease-in-out"
      }
    },
  },
  plugins: [],
}
