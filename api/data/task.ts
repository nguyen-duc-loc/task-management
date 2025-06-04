import { createAuthHeader } from "@/lib/auth-header";
import { API_BASE_URL } from "@/lib/backend-url";
import { fetchHandler } from "../fetch";
import TAGS from "@/constants/tag";
import { addDays, formatISO, subMilliseconds } from "date-fns";

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

export const getTasks = async (options: {
  limit?: number;
  page?: number;
  completed?: boolean;
  fromDate?: Date;
  toDate?: Date;
  search: string;
}) => {
  const { limit = 12, page = 1, completed, fromDate, toDate, search } = options;

  const urlSearchParams: Record<string, string> = {};
  urlSearchParams["page"] = String(page);
  urlSearchParams["limit"] = String(limit);
  urlSearchParams["title"] = String(search);
  urlSearchParams["description"] = String(search);
  if (completed !== undefined) {
    urlSearchParams["completed"] = String(completed);
  }
  if (fromDate) {
    urlSearchParams["start_deadline"] = formatISO(fromDate);
  }
  if (toDate) {
    urlSearchParams["end_deadline"] = formatISO(
      subMilliseconds(addDays(toDate, 1), 1)
    );
  }

  const searchParams = new URLSearchParams(urlSearchParams);

  const response = await fetchHandler(
    `${API_BASE_URL}/tasks?${searchParams.toString()}`,
    {
      headers: {
        Authorization: await createAuthHeader(),
      },
      next: {
        tags: [TAGS.tasks],
      },
    }
  );
  return (
    response.success ? response.data : { total: 0, tasks: [] }
  ) as GetTasksResponseData;
};
