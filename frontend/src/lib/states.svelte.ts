
import type { GetUserByIdRow } from '$lib/client';
import type { RoleType } from '$lib/types';


class AuthState {

    private data: GetUserByIdRow | null = $state(null);
    public valid: boolean = $state(false);

    clear () {
        this.valid = false
        this.data = null
    }

    get email (): string | undefined {
        return this.data?.email
    }

    get id (): number | undefined {
        return this.data?.id
    }

    get role (): RoleType | undefined {
        return this.data?.role as RoleType
    }

    set (data: GetUserByIdRow) {
        this.data = data
        this.valid = true

    }
    get (): GetUserByIdRow | null {
        return this.data
    }
}

export const authState = new AuthState();