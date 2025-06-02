import { getTaskById } from "@/api/data/task";
import Heading from "@/components/heading";
import { IconEdit } from "@tabler/icons-react";
import { notFound } from "next/navigation";
import React from "react";
import EditTaskForm from "./_components/EditTaskForm";

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

const EditTaskPage = async ({
  params,
}: {
  params: Promise<{ id: string }>;
}) => {
  const { id } = await params;
  const task = await getTaskById(id);
  if (!task) {
    notFound();
  }

  return (
    <>
      <Heading heading="Edit task" Icon={IconEdit} />
      <EditTaskForm task={task} />
    </>
  );
};

export default EditTaskPage;
