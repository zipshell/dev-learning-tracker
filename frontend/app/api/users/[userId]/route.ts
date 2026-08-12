import { NextRequest } from "next/server";
import { proxyBackendRequest } from "../../_lib/backend";

export async function GET(
  request: NextRequest,
  { params }: { params: Promise<{ userId: string }> },
) {
  const { userId } = await params;

  return proxyBackendRequest(request, `/users/${encodeURIComponent(userId)}`);
}
