import { defaultPlugins } from '@hey-api/openapi-ts'
import { defineConfig } from '@hey-api/openapi-ts';

export default defineConfig({
    client: '@hey-api/client-fetch',
    input: '../openapi.json',
    output: {
        path: 'src/lib/client',
        case: 'snake_case'
    },

    plugins: [
        ...defaultPlugins,
        '@tanstack/svelte-query',
    ],
});