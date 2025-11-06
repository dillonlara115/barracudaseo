/** @type {import('tailwindcss').Config} */
module.exports = {
	theme: {
		extend: {
			fontFamily: {
				heading: ['Sora', 'sans-serif'],
				body: ['DM Sans', 'sans-serif'],
				mono: ['JetBrains Mono', 'monospace']
			},
			colors: {
				primary: {
					DEFAULT: '#8ec07c',
					focus: '#a0d28c'
				},
				neutral: '#3c3836',
				'base-100': '#3c3836',
				'base-content': '#FFFFFF',
				info: '#458588',
				success: '#8ec07c',
				warning: '#d79921',
				error: '#cc241d'
			}
		}
	}
};

