
import { goto } from '$app/navigation';
import { me } from '$lib/client';
import type { LayoutLoad } from './$types';

export const ssr = false;
export const load: LayoutLoad = async () => {

    const { data, error } = await me()

    if (!error) {

        return {
            user: data
        }
    } else {
        goto('/login')
    }

}