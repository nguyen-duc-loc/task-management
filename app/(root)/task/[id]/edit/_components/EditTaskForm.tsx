"use client";

import { createTask, updateTask } from "@/api/actions/task";
import TaskForm from "@/components/forms/TaskForm";
import { TaskData } from "@/validation";
import React from "react";

interface EditTaskFormProps {
  task: Task;
}

const EditTaskForm = ({ task }: EditTaskFormProps) => {
  const { id, title, description, deadline } = task;

  return (
    <TaskForm
      defaultValues={{
        title,
        description: description || "",
        deadline: new Date(deadline),
      }}
      onSubmit={(data: TaskData) => updateTask(id, data)}
      formType="UPDATE"
    />
  );
};

export default EditTaskForm;
