"use server";

import { cookies } from "next/headers";
import { redirect } from "next/navigation";

import ROUTES from "@/constants/routes";
import { SigninData, SignupData } from "@/validation";
import { AUTH_TOKEN_KEY } from "@/lib/cookies";
import { API_BASE_URL } from "@/lib/url";
import { fetchHandler } from "./fetch";

export const signup = async (data: SignupData) => {
  const response = await fetchHandler(`${API_BASE_URL}/users`, {
    method: "POST",
    body: JSON.stringify({
      username: data.username,
      password: data.password,
    }),
  });
  if (response.success) {
    await signin({
      username: data.username,
      password: data.password,
    });
  }
  return response;
};

export const signin = async (data: SigninData) => {
  const response = await fetchHandler(`${API_BASE_URL}/users/login`, {
    method: "POST",
    body: JSON.stringify(data),
  });
  if (response.success) {
    const { access_token, access_token_expire_at } =
      response.data as SigninResponseData;
    (await cookies()).set({
      name: AUTH_TOKEN_KEY,
      value: access_token,
      path: "/",
      httpOnly: true,
      expires: new Date(access_token_expire_at),
    });
    redirect(ROUTES.dashboard);
  }
  return response;
};
