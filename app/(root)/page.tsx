import Heading from "@/components/heading";
import { IconLayoutDashboardFilled } from "@tabler/icons-react";
import React from "react";
import TaskContainer from "../_components/TaskContainer";

const RootPage = ({
  searchParams,
}: {
  searchParams: Promise<{ [key: string]: string }>;
}) => {
  return (
    <div className="w-5xl">
      <Heading
        heading="Dashboard"
        Icon={IconLayoutDashboardFilled}
        className="mb-0!"
      />
      <p className="text-muted-foreground mb-12">
        Welcome to the task management
      </p>
      <TaskContainer searchParams={searchParams} />
    </div>
  );
};

export default RootPage;
