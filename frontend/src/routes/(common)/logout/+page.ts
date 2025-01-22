
import type { PageLoad } from './$types';
import { logout } from '$lib/client';
import { goto } from '$app/navigation';
import { authState } from '$lib/states.svelte';


export const load: PageLoad = async () => {
    await logout()
    authState.clear()
    goto('/')

}