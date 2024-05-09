/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./template/*.html"],
  theme: {
    screens: {
      sm: '480px',
      md: '768px',
      lg: '976px',
      xl: '1440px',
    },
    fontFamily: {
      sans: ['Graphik', 'sans-serif'],
      serif: ['Merriweather', 'serif'],
    },
    colors: {
      'black': '#000000',
      'white': '#ffffff',
      'red': '#b51d14',
      'blue': '#3669ba',
      'gray': '#8492a6',
      'gray-dark': '#273444',
      'gray-light': '#d3dce6',
    },
    extend: {},
  },
  plugins: [
    require('@tailwindcss/forms'),
    require('@tailwindcss/typography'),
  ],
}

