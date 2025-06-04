"use server";

import { API_BASE_URL } from "@/lib/backend-url";
import { TaskData } from "@/validation";
import { fetchHandler } from "../fetch";
import { formatISO } from "date-fns";
import { redirect } from "next/navigation";
import ROUTES from "@/constants/routes";
import { createAuthHeader } from "@/lib/auth-header";
import { revalidateTag } from "next/cache";
import TAGS from "@/constants/tag";

export const createTask = async (data: TaskData) => {
  const response = await fetchHandler(`${API_BASE_URL}/tasks`, {
    method: "POST",
    headers: {
      Authorization: await createAuthHeader(),
    },
    body: JSON.stringify({
      ...data,
      deadline: formatISO(data.deadline),
    }),
  });
  if (response.success) {
    const task = response.data as CreateTaskResponseData;
    redirect(ROUTES.task(task.id));
  }
  return response;
};

export const updateTask = async (id: string, data: TaskData) => {
  const response = await fetchHandler(`${API_BASE_URL}/tasks/${id}`, {
    method: "PUT",
    headers: {
      Authorization: await createAuthHeader(),
    },
    body: JSON.stringify({
      ...data,
      deadline: formatISO(data.deadline),
    }),
  });
  if (response.success) {
    const task = response.data as UpdateTaskResponseData;
    revalidateTag(TAGS.tasks);
    redirect(ROUTES.task(task.id));
  }
  return response;
};

export const markTaskAsDone = async (id: string) => {
  const response = await fetchHandler(`${API_BASE_URL}/tasks/${id}`, {
    method: "PUT",
    headers: {
      Authorization: await createAuthHeader(),
    },
    body: JSON.stringify({
      completed: true,
    }),
  });
  if (response.success) {
    const task = response.data as UpdateTaskResponseData;
    revalidateTag(TAGS.tasks);
  }
  return response;
};
