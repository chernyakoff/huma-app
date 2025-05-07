import daisyui from 'daisyui';
import type { Config } from 'tailwindcss';

import typography from '@tailwindcss/typography';

const config = {
    // darkMode: ['class', '[data-theme="dark"]'],
    daisyui: {
        themes: [
            'dark', 'nord'
        ]
    },
    content: [
        './src/**/*.{html,js,svelte,ts}',


    ],
    theme: {
        extend: {

        },
    },
    plugins: [
        typography,
        daisyui,
    ]
} satisfies Config;

export default config;