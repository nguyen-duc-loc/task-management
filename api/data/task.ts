import { createAuthHeader } from "@/lib/auth-header";
import { API_BASE_URL } from "@/lib/url";
import { fetchHandler } from "../fetch";
import TAGS from "@/constants/tag";

export const getTaskById = async (id: string) => {
  const response = await fetchHandler(`${API_BASE_URL}/tasks/${id}`, {
    headers: {
      Authorization: await createAuthHeader(),
    },
    next: {
      tags: [TAGS.tasks],
    },
  });
  return response.success ? (response.data as GetTaskByIdResponseData) : null;
};
