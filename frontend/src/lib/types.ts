
export enum ROLES {
    ADMIN = 'admin',
    USER = 'user',
    EDITOR = 'editor'
}

export type RoleType = keyof typeof ROLES;