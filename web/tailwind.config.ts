import type { Config } from 'tailwindcss';

const config: Config = {
  content: [
    './pages/**/*.{js,ts,jsx,tsx,mdx}',
    './components/**/*.{js,ts,jsx,tsx,mdx}',
    './app/**/*.{js,ts,jsx,tsx,mdx}',
  ],
  theme: {
    extend: {
      colors: {
        navbarColor: '#B6B7B2',
        dash: '#F0F5EF',
        card: '#EFF1E6',
        bg: '#E0E2D7',
        btn: '#D2D2C8',
      },
      backgroundImage: {
        'gradient-radial': 'radial-gradient(var(--tw-gradient-stops))',
        'gradient-conic': 'conic-gradient(from 180deg at 50% 50%, var(--tw-gradient-stops))',
      },
      margin: {
        '1p': '1%',
        '2p': '2%',
        '3p': '3%',
        '4p': '4%',
        '5p': '5%',
        '10p': '10%',
        '15p': '15%',
      },
    },
  },
  plugins: [],
};

export default config;
