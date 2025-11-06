const shared = require('../shared/tailwind.theme.js');
const daisyui = require('daisyui');

/** @type {import('tailwindcss').Config} */
module.exports = {
	...shared,
	content: ['./src/**/*.{svelte,js,ts}'],
	plugins: [
		daisyui
	],
	daisyui: {
		themes: [
			{
				barracuda: {
					primary: '#8ec07c',
					'primary-focus': '#a0d28c',
					neutral: '#3c3836',
					'base-100': '#3c3836',
					'base-content': '#FFFFFF',
					info: '#458588',
					success: '#8ec07c',
					warning: '#d79921',
					error: '#cc241d'
				}
			}
		],
		darkTheme: 'barracuda',
		base: true,
		styled: true,
		utils: true,
		prefix: '',
		logs: true,
		themeRoot: ':root'
	}
};

