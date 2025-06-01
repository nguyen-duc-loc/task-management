import { cookies } from "next/headers";
import { NextRequest, NextResponse } from "next/server";

import ROUTES from "./constants/routes";
import { AUTH_TOKEN_KEY } from "./lib/cookies";

export async function middleware(request: NextRequest) {
  const authTokenValue = (await cookies()).get(AUTH_TOKEN_KEY)?.value;
  if (!authTokenValue) {
    return NextResponse.redirect(new URL(ROUTES.signin, request.url));
  }

  return NextResponse.next();
}

export const config = {
  matcher: [
    "/((?!api|_next/static|_next/image|favicon.ico|sitemap.xml|robots.txt|signin|signup).*)",
  ],
};
