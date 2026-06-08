import { NextResponse, type NextRequest } from "next/server";

type InitStatusResponse = {
  code?: number;
  data?: {
    initStatus?: boolean;
  };
};

const INIT_CHECK_PATH = "/configs/check-initialization";

function getServerHost() {
  return (process.env.SERVER_HOST || "http://localhost:8080").replace(/\/$/, "");
}

function getAdminHost(request: NextRequest) {
  const adminHost = process.env.NEXT_PUBLIC_ADMIN_HOST || process.env.ADMIN_HOST;

  if (adminHost) {
    return adminHost;
  }

  const fallbackUrl = request.nextUrl.clone();
  fallbackUrl.port = "8081";
  fallbackUrl.pathname = "/";
  fallbackUrl.search = "";
  return fallbackUrl.toString();
}

export async function middleware(request: NextRequest) {
  try {
    const res = await fetch(`${getServerHost()}${INIT_CHECK_PATH}`, {
      cache: "no-store",
    });

    if (!res.ok) {
      return NextResponse.next();
    }

    const body = (await res.json()) as InitStatusResponse;

    if (body.code === 0 && body.data?.initStatus === false) {
      return NextResponse.redirect(getAdminHost(request), 307);
    }
  } catch {
    return NextResponse.next();
  }

  return NextResponse.next();
}

export const config = {
  matcher: [
    /*
     * Only guard page requests. Static assets, Next internals, and API rewrites
     * should pass through unchanged.
     */
    "/((?!api|_next/static|_next/image|favicon.ico|robots.txt|sitemap.xml|.*\\..*).*)",
  ],
};
