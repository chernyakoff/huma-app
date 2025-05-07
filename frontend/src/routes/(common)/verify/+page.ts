
import { goto } from '$app/navigation';
import { verifyEmail } from '$lib/client';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ url }) => {

    const email = localStorage.getItem("verify")
    if (!email) {
        goto('/login')
    }

    const token = url.searchParams.get('token') || null
    if (token) {

        const { error } = await verifyEmail({ query: { token } })
        if (!error) {
            goto('/login')
        }

    }


}