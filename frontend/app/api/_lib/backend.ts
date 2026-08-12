import { NextRequest, NextResponse } from "next/server";

const BACKEND_BASE_URL = process.env.BACKEND_URL ?? "http://localhost:8080";
const HOP_BY_HOP_HEADERS = [
  "connection",
  "keep-alive",
  "proxy-authenticate",
  "proxy-authorization",
  "te",
  "trailer",
  "transfer-encoding",
  "upgrade",
];

function buildForwardHeaders(request: NextRequest) {
  const headers = new Headers();

  request.headers.forEach((value, key) => {
    const normalizedKey = key.toLowerCase();
    if (
      normalizedKey === "host" ||
      HOP_BY_HOP_HEADERS.includes(normalizedKey)
    ) {
      return;
    }
    headers.set(key, value);
  });

  return headers;
}

function buildResponseHeaders(source: Headers) {
  const headers = new Headers(source);
  HOP_BY_HOP_HEADERS.forEach((header) => headers.delete(header));
  return headers;
}

export async function proxyBackendRequest(
  request: NextRequest,
  path: string,
): Promise<NextResponse> {
  const url = `${BACKEND_BASE_URL}${path}`;
  const headers = buildForwardHeaders(request);
  const init: RequestInit = {
    method: request.method,
    headers,
    redirect: "manual",
  };

  if (request.method !== "GET" && request.method !== "HEAD") {
    init.body = await request.text();
  }

  const backendResponse = await fetch(url, init);
  const body = await backendResponse.text();
  const responseHeaders = buildResponseHeaders(backendResponse.headers);

  console.log(backendResponse);

  if (backendResponse.status === 401) {
    return NextResponse.redirect(new URL("/login", request.url));
  }

  return new NextResponse(body, {
    status: backendResponse.status,
    headers: responseHeaders,
  });
}
