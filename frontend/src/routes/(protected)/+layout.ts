
import type { LayoutLoad } from './$types';
import { me } from '$lib/client';

export const ssr = false;
export const load: LayoutLoad = async () => {

    const { data, error } = await me()

    if (!error) {

        return {
            user: data
        }
    }



}