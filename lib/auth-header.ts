import { cookies } from "next/headers";
import { AUTH_TOKEN_KEY } from "./cookies";

export const createAuthHeader = async () => {
  return `Bearer ${(await cookies()).get(AUTH_TOKEN_KEY)?.value}`;
};
