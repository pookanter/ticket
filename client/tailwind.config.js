/** @type {import('tailwindcss').Config} */
export default {
  content: ['./src/**/*.{html,js,svelte,ts}',  "./node_modules/flowbite/**/*.js"],
  theme: {
    extend: {
      colors: {
        'primary': {
          '50': '#eff4ff',
          '100': '#dbe6fe',
          '200': '#bfd3fe',
          '300': '#93b4fd',
          '400': '#6090fa',
          '500': '#3b76f6',
          '600': '#2563eb',
          '700': '#1d58d8',
          '800': '#1e4baf',
          '900': '#1e408a',
          '950': '#172a54',
        },
      },
    },
  },
  plugins: [
    require('flowbite/plugin'),
  ],
}

