import Heading from "@/components/heading";
import { IconSquareRoundedPlusFilled } from "@tabler/icons-react";
import { Metadata } from "next";
import React from "react";
import CreateTaskForm from "./_components/CreateTaskForm";

export const metadata: Metadata = {
  title: "Create task",
  description: "Create a new task",
};

const TaskPage = () => {
  return (
    <>
      <Heading heading="Create a new task" Icon={IconSquareRoundedPlusFilled} />
      <CreateTaskForm />
    </>
  );
};

export default TaskPage;
