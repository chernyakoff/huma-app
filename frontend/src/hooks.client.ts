import { client } from '$lib/client/sdk.gen';

//TODO baseUrl from env
client.setConfig({
    baseUrl: 'http://localhost:8888',
    credentials: 'include'

});
