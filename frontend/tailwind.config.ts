import type { Config } from 'tailwindcss';
import daisyui from 'daisyui';

import typography from '@tailwindcss/typography';

const config = {

    daisyui: {
        themes: [
            'nord', 'dark', 'cupcake'
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