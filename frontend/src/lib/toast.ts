import { toast } from '@zerodevx/svelte-toast'

export const success = (m: string) => toast.push(m, {
    theme: {
        '--toastBackground': 'green',
        '--toastColor': 'white',
        '--toastBarBackground': 'olive'
    }
})


export const error = (m: string | undefined) => {
    if (!m) m = 'Error'
    toast.push(m, {
        dismissable: false,
        duration: 2000,
        theme: {
            '--toastBackground': 'oklch(var(--er))',
            '--toastColor': 'white',
            '--toastBarHeight': 0
        }
    })
}


