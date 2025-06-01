"use server";

import { API_BASE_URL } from "@/lib/url";
import { TaskData } from "@/validation";
import { fetchHandler } from "../fetch";
import { formatISO } from "date-fns";
import { redirect } from "next/navigation";
import ROUTES from "@/constants/routes";
import { createAuthHeader } from "@/lib/auth-header";

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
    console.log(response);
    const task = response.data as CreateTaskResponseData;
    redirect(ROUTES.task(task.id));
  }
  return response;
};
