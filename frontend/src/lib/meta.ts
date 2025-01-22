import { page } from '$app/state';


const titles: Record<string, string> = {
    '/': "Home",
    '/login': "Login",
    "/register": "Register",
    "/app": "Dashboard",
    "/app/users": "App users"
}


export const getTitle = (): string => {

    if (typeof titles[page.url.pathname] != 'undefined') {
        return titles[page.url.pathname]
    }
    return 'No title'
}