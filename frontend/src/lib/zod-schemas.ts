import { z } from 'zod';

const passwordValidation = new RegExp(
    /^(?=.*?[A-Z])(?=.*?[a-z])(?=.*?[0-9])(?=.*?[#?!@$%^&*-]).{6,}$/
);


export const userSchema = z.object({
    email: z
        .string({ required_error: 'Email is required' })
        .email({ message: 'Please enter a valid email address' }),
    password: z
        .string({ required_error: 'Password is required' })
        .min(6, { message: 'Password must be at least 6 characters' })
        .regex(passwordValidation, {
            message: 'Your password is not valid',
        })
        .trim(),
    confirmPassword: z
        .string({ required_error: 'Password confirm is required' })
        .trim(),
    //terms: z.boolean({ required_error: 'You must accept the terms and privacy policy' }),
    role: z
        .enum(['user', 'editor', 'admin'], { required_error: 'You must have a role' })
        .default('admin'),
    verified: z.boolean().default(false),
    terms: z.literal<boolean>(true, {
        errorMap: () => ({ message: "You must accept the terms & privacy policy" }),
    }),
    token: z.string().optional(),
    receiveEmail: z.boolean().default(true),
    createdAt: z.date().optional(),
    updatedAt: z.date().optional()
});

export type UserSchema = typeof userSchema;


export const registerSchema = userSchema.pick({ email: true, password: true, confirmPassword: true }).superRefine(({ confirmPassword, password }, ctx) => {
    if (confirmPassword !== password) {
        ctx.addIssue({
            code: 'custom',
            message: 'Passwords must match',
            path: ['confirmPassword']
        });
    }
})
export type RegisterSchema = typeof registerSchema;

export const loginSchema = userSchema.pick({ email: true, password: true });
export type LoginSchema = typeof loginSchema;


export const userUpdatePasswordSchema = userSchema
    .pick({ password: true, confirmPassword: true })
    .superRefine(({ confirmPassword, password }, ctx) => {
        if (confirmPassword !== password) {
            ctx.addIssue({
                code: 'custom',
                message: 'Password and Confirm Password must match',
                path: ['password']
            });
            ctx.addIssue({
                code: 'custom',
                message: 'Password and Confirm Password must match',
                path: ['confirmPassword']
            });
        }
    });

export type UserUpdatePasswordSchema = typeof userUpdatePasswordSchema;