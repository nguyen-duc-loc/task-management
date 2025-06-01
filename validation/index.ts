import z from "zod";

export const SigninSchema = z.object({
  username: z
    .string()
    .min(1, { message: "Username is required." })
    .regex(/^[a-zA-Z0-9]+$/, {
      message: "Username must contains only alphanumeric characters",
    }),

  password: z
    .string()
    .min(6, { message: "Password must be at least 6 characters." }),
});
export type SigninData = z.infer<typeof SigninSchema>;

export const SignupSchema = SigninSchema.extend({
  confirmPassword: z.string(),
}).refine((data) => data.password === data.confirmPassword, {
  path: ["confirmPassword"],
  message: "Passwords do not match",
});
export type SignupData = z.infer<typeof SignupSchema>;
