import { NextRequest } from "next/server";
import { proxyBackendRequest } from "../../_lib/backend";

export async function PATCH(
  request: NextRequest,
  { params }: { params: Promise<{ entry_id: string }> },
) {
  const { entry_id } = await params;

  return proxyBackendRequest(
    request,
    `/entries/${encodeURIComponent(entry_id)}`,
  );
}

export async function DELETE(
  request: NextRequest,
  { params }: { params: Promise<{ entry_id: string }> },
) {
  const { entry_id } = await params;

  return proxyBackendRequest(
    request,
    `/entries/${encodeURIComponent(entry_id)}`,
  );
}
