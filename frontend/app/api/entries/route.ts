import { NextRequest } from "next/server";
import { proxyBackendRequest } from "../_lib/backend";

export async function POST(request: NextRequest) {
  return proxyBackendRequest(request, "/entries/");
}

export async function GET(request: NextRequest) {
  const url = new URL(request.url);
  const userId = url.searchParams.get("user_id");

  if (!userId) {
    return new Response(
      JSON.stringify({ error: "Missing user_id query parameter" }),
      {
        status: 400,
        headers: { "content-type": "application/json" },
      },
    );
  }

  return proxyBackendRequest(request, `/entries/${encodeURIComponent(userId)}`);
}
