import { getTaskById } from "@/api/data/task";
import { notFound } from "next/navigation";
import React from "react";
import Task from "./_components/Task";

export const generateMetadata = async ({
  params,
}: {
  params: Promise<{ id: string }>;
}) => {
  const { id } = await params;
  const task = await getTaskById(id);
  if (!task) {
    return {
      title: "Task not found",
      description: "Task not found",
    };
  }

  return {
    title: task.title,
    description: task.description,
  };
};

const TaskPage = async ({ params }: { params: Promise<{ id: string }> }) => {
  const { id } = await params;
  const task = await getTaskById(id);
  if (!task) {
    notFound();
  }

  return <Task task={task} />;
};

export default TaskPage;
