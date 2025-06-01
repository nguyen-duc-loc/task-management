"use client";

import { createTask } from "@/api/actions/task";
import TaskForm from "@/components/forms/TaskForm";
import { TaskData } from "@/validation";
import { addDays } from "date-fns";
import React from "react";

const CreateTaskForm = () => {
  return (
    <TaskForm
      defaultValues={{
        title: "",
        description: "",
        deadline: addDays(new Date(), 1),
      }}
      onSubmit={(data: TaskData) => createTask(data)}
      formType="CREATE"
    />
  );
};

export default CreateTaskForm;
